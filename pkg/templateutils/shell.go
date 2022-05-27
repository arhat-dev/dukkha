package templateutils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"reflect"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/iohelper"
	"arhat.dev/pkg/stringhelper"
)

// ExpandEnv expands unix style environment variable (`$FOO`, `${BAR}`)
// if enableExec is set to ture, also supports arbitrary command execution
// using `$(do something)`
func ExpandEnv(rc dukkha.RenderingContext, toExpand string, enableExec bool) (string, error) {
	var rd strings.Reader
	rd.Reset(toExpand)

	parser := syntax.NewParser(
		syntax.Variant(syntax.LangBash),
	)

	word, err := parser.Document(&rd)
	if err != nil {
		return "", fmt.Errorf("invalid text for env expansion: %w", err)
	}

	printer := syntax.NewPrinter(
		syntax.FunctionNextLine(false),
		syntax.Indent(2),
	)

	config := expand.Config{
		Env: rc,
		// reassemble back-quoted string and $() by default
		CmdSubst: func(w io.Writer, cs *syntax.CmdSubst) error {
			script, err2 := rebuildScript(printer, cs)
			if err2 != nil {
				if err2 != errSkipBackquotedCmdSubst {
					return err2
				}

				script, err2 = ExpandEnv(rc, script, false)
				if err2 != nil {
					return err2
				}

				script = "`" + script + "`"
			} else {
				script = "$(" + script + ")"
			}

			_, err2 = w.Write([]byte(script))
			return err2
		},
		ProcSubst: nil,
		ReadDir: func(s string) ([]fs.FileInfo, error) {
			ents, err2 := rc.FS().ReadDir(s)
			if err2 != nil {
				return nil, err2
			}

			ret := make([]fs.FileInfo, len(ents))
			for i, e := range ents {
				ret[i], err2 = e.Info()
				if err2 != nil {
					return nil, err2
				}
			}

			return ret, nil
		},
		GlobStar: true,
		NullGlob: true,
		NoUnset:  true,
	}

	if enableExec {
		var stdout bytes.Buffer
		runner, err := CreateShellRunner(
			rc.WorkDir(), rc, nil, &stdout, rc.Stderr(),
		)
		if err != nil {
			return "", fmt.Errorf("creating shell runner for env: %w", err)
		}

		// reassemble back-quoted string but eval $()
		config.CmdSubst = func(w io.Writer, cs *syntax.CmdSubst) error {
			script, err2 := rebuildScript(printer, cs)
			if err2 != nil {
				if err2 != errSkipBackquotedCmdSubst {
					return err2
				}

				script, err2 = ExpandEnv(rc, script, false)
				if err2 != nil {
					return err2
				}

				script = "`" + script + "`"
				_, err2 = w.Write([]byte(script))
				return err2
			}

			stdout.Reset()
			err2 = RunScript(rc, runner, parser, script)
			if err2 != nil {
				return err2
			}

			_, err2 = stdout.WriteTo(w)
			if err2 != nil {
				return fmt.Errorf("shell output not written: %w", err)
			}

			return nil
		}
	}

	return expand.Document(&config, word)
}

const errSkipBackquotedCmdSubst errString = "skip CmdSubSt"

func rebuildScript(printer *syntax.Printer, cs *syntax.CmdSubst) (string, error) {
	var buf bytes.Buffer

	err2 := printer.Print(&buf, cs)
	if err2 != nil {
		return "", fmt.Errorf("rebuild evaluation commands: %w", err2)
	}

	rawCmd := stringhelper.Convert[string, byte](buf.Next(buf.Len()))

	// printed rawCmd is always in `$()` format
	switch {
	case !cs.Backquotes:
		return rawCmd[2 : len(rawCmd)-1], nil
	case cs.Backquotes:
		return rawCmd[2 : len(rawCmd)-1], errSkipBackquotedCmdSubst
	default:
		return "", fmt.Errorf("invalid command substution: %q", rawCmd)
	}
}

func CreateShellRunner(
	workdir string, // workdir may be different from rc.WorkDir
	rc dukkha.RenderingContext,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (*interp.Runner, error) {
	runner, err := interp.New(
		interp.Env(rc),
		interp.Dir(workdir),
		interp.StdIO(stdin, stdout, stderr),
		interp.Params("-e"),
		interp.OpenHandler(fileOpenHandler),
		interp.ExecHandler(newExecHandler(rc, stdin, stdout)),
	)

	if err != nil {
		return nil, err
	}

	return runner, nil
}

func RunScript(
	ctx context.Context,
	runner *interp.Runner,
	parser *syntax.Parser,
	script string,
) error {
	f, err := parser.Parse(strings.NewReader(script), "")
	if err != nil {
		return fmt.Errorf("invalid script (%v):\n%s", err, script)
	}

	err = runner.Run(ctx, f)
	if err != nil {
		return fmt.Errorf("embedded shell exited with error (%v):\n%s", err, script)
	}

	return nil
}

// ExecCmdAsTemplateFuncCall executes cmd as a template func
// args[0] should be a template func name, see funcs.go for reference
func ExecCmdAsTemplateFuncCall(
	rc dukkha.RenderingContext,
	pipeStdin io.Reader,
	pipeStdout io.Writer,
	args []string,
) (stdout, stderr string, err error) {
	var (
		bufCallArgs [10]reflect.Value
		callArgs    []reflect.Value

		nArgs     = len(args) - 1
		useStdin  = pipeStdin != nil
		useStdout = pipeStdout != nil

		funcs = CreateTemplateFuncs(rc)
	)
	if nArgs < 0 {
		err = errString("invalid 0 arg call")
		return
	}

	const nBUF = len(bufCallArgs)

	funcName := args[0]
	fid := FuncNameToFuncID(funcName)

	fn := funcs.GetByID(fid)
	if !fn.IsValid() {
		err = fmt.Errorf("%q not found", funcName)
		return
	}

	if useStdout {
		// TODO: generate a func to tell whether fn supports the second last arg as writer
		switch fid {
		case FuncID_enc_YAML, FuncID_enc_JSON, FuncID_toYaml, FuncID_toJson,
			FuncID_enc_Base32, FuncID_enc_Base64, FuncID_hex:
		default:
			useStdout = false
		}
	}

	if useStdin && useStdout {
		if nArgs+2 > nBUF {
			callArgs = make([]reflect.Value, nArgs+2)
		} else {
			callArgs = bufCallArgs[:nArgs+2]
		}

		callArgs[nArgs+1] = reflect.ValueOf(pipeStdin)
		callArgs[nArgs] = reflect.ValueOf(pipeStdout)
	} else if useStdin || useStdout {
		if nArgs+1 > nBUF {
			callArgs = make([]reflect.Value, nArgs+1)
		} else {
			callArgs = bufCallArgs[:nArgs+1]
		}

		if useStdin { // input is the last argument
			callArgs[nArgs] = reflect.ValueOf(pipeStdin)
		} else /* useStdout */ { // output is the second last argument by convension
			if nArgs > 0 {
				callArgs[nArgs-1] = reflect.ValueOf(pipeStdout)
			} else {
				callArgs[0] = reflect.ValueOf(pipeStdout)
			}
		}
	} else {
		if nArgs > nBUF {
			callArgs = make([]reflect.Value, nArgs)
		} else {
			callArgs = bufCallArgs[:nArgs]
		}
	}

	j := 0
	for i := range callArgs {
		if callArgs[i].IsValid() {
			continue
		}

		callArgs[i] = reflect.ValueOf(args[j+1])
		j++
	}

	defer func() {
		p := recover()
		if p != nil {
			err = fmt.Errorf("%s %s, called with %v: %v", funcName, fn.Type().String(), callArgs, p)
		}
	}()

	callRet := fn.Call(callArgs)
	switch len(callRet) {
	case 1:
	case 2:
		ret1 := callRet[1]
		if ret1.IsValid() && ret1.CanInterface() {
			err, _ = ret1.Interface().(error)
			if err != nil {
				return
			}
		}
	default:
		panic("invalid return value count, expecting 1 or 2")
	}

	ret0 := callRet[0]
	if ret0.IsValid() && ret0.CanInterface() {
		switch rv := ret0.Interface().(type) {
		case cmdOutput:
			stdout = rv.Stdout
			stderr = rv.Stderr
		default:
			stdout, err = toString(rv)
		}
		return
	}

	return
}

func newExecHandler(
	rc dukkha.RenderingContext,
	origStdin io.Reader,
	origStdout io.Writer,
) interp.ExecHandlerFunc {
	defaultCmdExecHandler := interp.DukkhaExecHandler(0)

	return func(
		ctx context.Context,
		args []string,
	) error {
		if !strings.HasPrefix(args[0], "tpl:") {
			err := defaultCmdExecHandler(ctx, args)
			if err != nil {
				return fmt.Errorf("exec: %q: %w", strings.Join(args, " "), err)
			}

			return nil
		}

		// has `tpl:` prefix, execute as a template func

		hc := interp.HandlerCtx(ctx)

		var (
			pipeReader io.Reader
			pipeWriter io.Writer
		)
		if hc.Stdin != origStdin {
			// piped context
			pipeReader = hc.Stdin
		}

		if hc.Stdout != origStdout {
			pipeWriter = hc.Stdout
		}

		args[0] = strings.TrimPrefix(args[0], "tpl:")

		out, errOut, err := ExecCmdAsTemplateFuncCall(
			rc,
			pipeReader,
			pipeWriter,
			args,
		)

		if len(out) != 0 {
			hc.Stdout.Write(stringhelper.ToBytes[byte, byte](out))
		}

		if len(errOut) != 0 {
			hc.Stderr.Write(stringhelper.ToBytes[byte, byte](errOut))
		}

		return err
	}
}

func fileOpenHandler(
	ctx context.Context,
	path string,
	flag int,
	perm fs.FileMode,
) (io.ReadWriteCloser, error) {
	const devNullPath = "/dev/null"

	if path == devNullPath {
		return iohelper.NewDevNull(), nil
	}

	hc := interp.HandlerCtx(ctx)
	osfs := fshelper.NewOSFS(false, func() (string, error) {
		return hc.Dir, nil
	})

	f, err := osfs.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	return f.(*os.File), nil
}
