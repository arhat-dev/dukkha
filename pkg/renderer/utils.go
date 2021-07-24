package renderer

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

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
		interp.ExecHandler(interp.DefaultExecHandler(0)),
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
