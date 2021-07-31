package templateutils

import (
	"context"
	"fmt"
	"io"
	"strings"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
)

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
		return fmt.Errorf(
			"failed to parse shell script:\n\n%s\n\nin embedded shell: %w",
			script,
			err,
		)
	}

	err = runner.Run(ctx, f)
	if err != nil {
		return fmt.Errorf(
			"failed to run script:\n\n%s\n\nin embedded shell: %w",
			script, err,
		)
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
