package renderer

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"arhat.dev/pkg/exechelper"
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
		interp.Params("-e"),
		interp.ExecHandler(func(ctx context.Context, args []string) error {
			hc := interp.HandlerCtx(ctx)

			env := make(map[string]string)
			hc.Env.Each(func(name string, vr expand.Variable) bool {
				switch vr.Kind {
				case expand.NameRef:
					env[name], _ = vr.Resolve(environ)
				case expand.String:
					env[name] = vr.Str
				case expand.Indexed:
					env[name] = strings.Join(vr.List, " ")
				default:
					env[name] = vr.String()
				}

				return true
			})

			cmd, err := exechelper.Do(exechelper.Spec{
				Context: ctx,
				Env:     env,
				Dir:     hc.Dir,
				Command: args,

				Stdin:  hc.Stdin,
				Stdout: hc.Stdout,
				Stderr: hc.Stderr,
			})

			if err == nil {
				exitCode, err2 := cmd.Wait()
				if err2 != nil {
					return interp.NewExitStatus(uint8(exitCode))
				}

				return nil
			}

			// copied from interp.DefaultExecHandler
			switch x := err.(type) {
			case *exec.ExitError:
				// started, but errored - default to 1 if OS
				// doesn't have exit statuses
				if status, ok := x.Sys().(syscall.WaitStatus); ok {
					if status.Signaled() {
						if ctx.Err() != nil {
							return ctx.Err()
						}
						return interp.NewExitStatus(uint8(128 + status.Signal()))
					}
					return interp.NewExitStatus(uint8(status.ExitStatus()))
				}
				return interp.NewExitStatus(1)
			case *exec.Error:
				// did not start
				fmt.Fprintf(hc.Stderr, "%v\n", err)
				return interp.NewExitStatus(127)
			default:
				return err
			}
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
			"failed to parse shell script:\n\n%s\n\nin embedded shell: %w",
			script,
			err,
		)
	}

	err = runner.Run(ctx, f)
	if err != nil {
		return fmt.Errorf(
			"failed to run command:\n\n%s\n\nin embedded shell: %w",
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
