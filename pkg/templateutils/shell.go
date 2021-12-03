package templateutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
)

// ExpandEnv expand environment variable in unix style (`$FOO`, `${BAR}`)
// if enableExec is set to ture, also supports arbitrary command execution
// using `$(do something)`
func ExpandEnv(rc dukkha.RenderingContext, toExpand string, enableExec bool) (string, error) {
	parser := syntax.NewParser(
		syntax.Variant(syntax.LangBash),
	)

	word, err := parser.Document(strings.NewReader(toExpand))
	if err != nil {
		return "", fmt.Errorf("invalid text for env expansion: %w", err)
	}

	printer := syntax.NewPrinter(
		syntax.FunctionNextLine(false),
		syntax.Indent(2),
	)

	config := &expand.Config{
		Env: rc,
		// reassemble back-quoted string and $() by default
		CmdSubst: func(w io.Writer, cs *syntax.CmdSubst) error {
			script, err2 := rebuildShellEvaluation(printer, cs)
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
		stdout := &bytes.Buffer{}
		runner, err := CreateEmbeddedShellRunner(
			rc.WorkDir(), rc, nil, stdout, os.Stderr,
		)
		if err != nil {
			return "", fmt.Errorf("failed to create shell runner for env: %w", err)
		}

		// reassemble back-quoted string but eval $()
		config.CmdSubst = func(w io.Writer, cs *syntax.CmdSubst) error {
			script, err2 := rebuildShellEvaluation(printer, cs)
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
			err2 = RunScriptInEmbeddedShell(rc, runner, parser, script)
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

	return expand.Document(config, word)
}

var errSkipBackquotedCmdSubst = errors.New("skip CmdSubSt")

func rebuildShellEvaluation(printer *syntax.Printer, cs *syntax.CmdSubst) (string, error) {
	buf := &bytes.Buffer{}
	err2 := printer.Print(buf, cs)
	if err2 != nil {
		return "", fmt.Errorf("failed to get evaluation commands: %w", err2)
	}

	rawCmd := string(buf.Next(buf.Len()))

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

func CreateEmbeddedShellRunner(
	workdir string,
	rc dukkha.RenderingContext,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (*interp.Runner, error) {
	runner, err := interp.New(
		interp.Env(rc),
		interp.Dir(workdir),
		interp.StdIO(stdin, stdout, stderr),
		interp.Params("-e"),
		interp.OpenHandler(fileOpenHandler),
		interp.ExecHandler(newExecHandler(rc, stdin)),
	)

	if err != nil {
		return nil, err
	}

	return runner, nil
}

func RunScriptInEmbeddedShell(
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

func ExecCmdAsTemplateFuncCall(
	rc dukkha.RenderingContext,
	stdin io.Reader,
	stdout io.Writer,
	args []string,
) error {
	tpl := `{{- ` + strings.Join(args, " ")

	var values interface{} = rc
	if stdin != nil {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return err
		}

		type valuesWithStdin struct {
			dukkha.RenderingContext
			// nolint:revive
			DUKKHA_TEMPLATE_STDIN string
		}

		values = &valuesWithStdin{
			RenderingContext:      rc,
			DUKKHA_TEMPLATE_STDIN: string(data),
		}

		tpl += ` .DUKKHA_TEMPLATE_STDIN`
	}

	tpl += ` -}}`

	t, err := CreateTemplate(rc).Parse(tpl)
	if err != nil {
		return fmt.Errorf("failed to convert to template call: %w", err)
	}

	return t.Execute(stdout, values)
}
