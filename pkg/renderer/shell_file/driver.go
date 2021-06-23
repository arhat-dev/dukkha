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

func (d *Driver) Render(ctx *field.RenderingContext, scriptPath string) (string, error) {
	buf := &bytes.Buffer{}

	env, cmd, err := d.getExecSpec(scriptPath, true)
	if err != nil {
		return "", fmt.Errorf(
			"renderer.%s: failed to get exec spec: %w",
			DefaultName, err,
		)
	}

	ctx.AddEnv(env...)

	p, err := exechelper.Do(exechelper.Spec{
		Context: ctx.Context(),
		Command: cmd,
		Env:     ctx.Values().Env,

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
