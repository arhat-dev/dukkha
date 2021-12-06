package templateutils

import (
	"context"
	"fmt"
	"io"
	"strings"

	"mvdan.cc/sh/v3/interp"

	"arhat.dev/dukkha/pkg/dukkha"
)

func newExecHandler(rc dukkha.RenderingContext, stdin io.Reader) interp.ExecHandlerFunc {
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

		hc := interp.HandlerCtx(ctx)

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
				[]string{strings.TrimPrefix(args[0], "tpl:")},
				args[1:]...,
			),
		)
	}
}
