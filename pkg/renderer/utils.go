package renderer

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"arhat.dev/pkg/exechelper"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
)

func ToYamlBytes(in interface{}) ([]byte, error) {
	switch t := in.(type) {
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	default:
	}

	ret, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func CreateEnvForEmbeddedShell(rc dukkha.RenderingContext) (expand.Environ, *error) {
	var errEnvMissing error
	return expand.FuncEnviron(func(name string) string {
		v, ok := rc.Env()[name]
		if ok {
			return v
		}

		switch name {
		case "IFS":
			return " \t\n"
		case "OPTIND":
			return "1"
		case "PWD":
			return rc.WorkingDir()
		case "UID":
			// os.Getenv("UID") usually retruns empty value
			// so we have to call os.Getuid
			return strconv.FormatInt(int64(os.Getuid()), 10)
		default:
			errEnvMissing = multierr.Append(errEnvMissing,
				fmt.Errorf("env %q not found", name),
			)
		}

		return ""
	}), &errEnvMissing
}

func CreateEmbeddedShellRunner(
	workingDir string,
	environ expand.Environ,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (*interp.Runner, error) {
	return interp.New(
		interp.Env(environ),
		interp.Dir(workingDir),
		interp.StdIO(stdin, stdout, stderr),
		interp.Params("-e", "-u"),
		interp.ExecHandler(func(_ctx context.Context, args []string) error {
			ctx := interp.HandlerCtx(_ctx)

			cmd, err := exechelper.Do(exechelper.Spec{
				Context: _ctx,
				Dir:     ctx.Dir,
				Command: args,

				Stdin:  ctx.Stdin,
				Stdout: ctx.Stdout,
				Stderr: ctx.Stderr,
			})
			if err != nil {
				return err
			}

			_, err = cmd.Wait()
			if err != nil {
				return err
			}

			return nil
		}),
	)
}

func RunShellScriptInEmbeddedShell(
	ctx context.Context,
	runner *interp.Runner,
	parser *syntax.Parser,
	script string,
) error {
	f, err := parser.Parse(strings.NewReader(script), "")
	if err != nil {
		return fmt.Errorf(
			"failed to parse shell command \n\n%s\n\nusing embedded shell: %w",
			script,
			err,
		)
	}

	err = runner.Run(ctx, f)
	if err != nil {
		return fmt.Errorf(
			"failed to evaluate command \n\n%s\n\nusing embedded shell: %w",
			script, err,
		)
	}

	return nil
}

func RunShellScript(
	rc dukkha.RenderingContext,
	script string,
	isFilePath bool,
	stdout io.Writer,
	getExecSpec dukkha.ExecSpecGetFunc,
) error {
	env, cmd, err := getExecSpec([]string{script}, false)
	if err != nil {
		return fmt.Errorf("failed to get exec spec: %w", err)
	}

	rc.AddEnv(env...)

	p, err := exechelper.Do(exechelper.Spec{
		Context: rc,
		Command: cmd,
		Env:     rc.Env(),

		Stdout: stdout,
		Stderr: os.Stderr,
	})

	if err != nil {
		return fmt.Errorf("failed to run script [%s]: %w",
			strings.Join(cmd, " "), err,
		)
	}

	_, err = p.Wait()
	if err != nil {
		return fmt.Errorf("cmd failed: %w", err)
	}

	return nil
}
