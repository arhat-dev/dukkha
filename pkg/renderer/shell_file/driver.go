package shell_file

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "shell_file"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct {
	getExecSpec field.ExecSpecGetFunc
}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx *field.RenderingContext, rawData interface{}) (string, error) {
	scriptPath, ok := rawData.(string)
	if !ok {
		return "", fmt.Errorf("renderer.%s: unexpected non-string input %T", DefaultName, rawData)
	}

	env, cmd, err := d.getExecSpec([]string{scriptPath}, true)
	if err != nil {
		return "", fmt.Errorf(
			"renderer.%s: failed to get exec spec: %w",
			DefaultName, err,
		)
	}

	execCtx := ctx.Clone()
	execCtx.AddEnv(env...)

	buf := &bytes.Buffer{}
	p, err := exechelper.Do(exechelper.Spec{
		Context: execCtx.Context(),
		Command: cmd,
		Env:     execCtx.Values().Env,

		Stdout: buf,
		Stderr: os.Stderr,
	})

	if err != nil {
		return "", fmt.Errorf(
			"renderer.%s: failed to start command [%s]: %w",
			DefaultName, strings.Join(cmd, " "), err,
		)
	}

	_, err = p.Wait()
	if err != nil {
		return "", fmt.Errorf(
			"renderer.%s: cmd failed: %w",
			DefaultName, err,
		)
	}

	return buf.String(), nil
}
