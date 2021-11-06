package templateutils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
)

func ExpandEnv(rc dukkha.RenderingContext, toExpand string) (string, error) {
	parser := syntax.NewParser(
		syntax.Variant(syntax.LangBash),
	)

	word, err := parser.Document(strings.NewReader(toExpand))
	if err != nil {
		return "", fmt.Errorf(
			"invalid expansion text %q: %w",
			toExpand, err,
		)
	}

	embeddedShellOutput := &bytes.Buffer{}
	runner, err := CreateEmbeddedShellRunner(
		rc.WorkingDir(), rc, nil, embeddedShellOutput, os.Stderr,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create shell runner for env: %w",
			err,
		)
	}

	printer := syntax.NewPrinter(
		syntax.FunctionNextLine(false),
		syntax.Indent(2),
	)

	result, err := expand.Document(&expand.Config{
		Env: rc,
		CmdSubst: func(w io.Writer, cs *syntax.CmdSubst) error {
			buf := &bytes.Buffer{}
			err2 := printer.Print(buf, cs)
			if err2 != nil {
				return fmt.Errorf("failed to get evaluation commands: %w", err2)
			}

			script := string(buf.Bytes()[2 : buf.Len()-1])

			embeddedShellOutput.Reset()
			err2 = RunScriptInEmbeddedShell(rc, runner, parser, script)
			if err2 != nil {
				return err2
			}

			_, err2 = embeddedShellOutput.WriteTo(w)
			if err2 != nil {
				return fmt.Errorf(
					"failed to write embedded shell output to result value: %w", err,
				)
			}

			return nil
		},
		ProcSubst: nil,
		ReadDir: func(s string) ([]os.FileInfo, error) {
			ents, err2 := os.ReadDir(s)
			if err2 != nil {
				return nil, err2
			}

			ret := make([]os.FileInfo, len(ents))
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
	},
		word,
	)

	if err != nil {
		return "", fmt.Errorf(
			"env expansion failed: %w",
			err,
		)
	}

	return result, nil
}

func CreateEmbeddedShellRunner(
	workingDir string,
	rc dukkha.RenderingContext,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (*interp.Runner, error) {
	cmdExecHandler := interp.DefaultExecHandler(0)
	return interp.New(
		interp.Env(rc),
		interp.Dir(workingDir),
		interp.StdIO(stdin, stdout, stderr),
		interp.Params("-e"),
		interp.ExecHandler(func(ctx context.Context, args []string) error {
			hc := interp.HandlerCtx(ctx)

			if !strings.HasPrefix(args[0], "template:") {
				return cmdExecHandler(ctx, args)
			}

			var pipeReader io.Reader
			if hc.Stdin != stdin {
				// piped context
				pipeReader = hc.Stdin
			}

			return ExecCmdAsTemplateFuncCall(
				rc,
				pipeReader,
				hc.Stdout,
				append(
					[]string{strings.TrimPrefix(args[0], "template:")},
					args[1:]...,
				),
			)
		}),
	)
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
