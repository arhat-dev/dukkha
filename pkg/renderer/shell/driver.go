package shell

import (
	"bytes"
	"fmt"
	"os"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "shell"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct {
	doExec ExecFunc
}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx *field.RenderingContext, script string) (string, error) {
	stdout := &bytes.Buffer{}

	code, err := d.doExec(script, &exechelper.Spec{
		Context: ctx.Context(),
		Env:     ctx.Values().Env,
		Stdout:  stdout,
		Stderr:  os.Stderr,
	})
	if err != nil {
		return "", fmt.Errorf("renderer.%s: exit code %d: %w", DefaultName, code, err)
	}

	return stdout.String(), nil
}
