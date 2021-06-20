package shell_file

import (
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/shell"
)

const DefaultName = "shell_file"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct {
	impl *shell.Driver
}

func (d *Driver) Name() string {
	return DefaultName
}

func (d *Driver) Render(ctx *field.RenderingContext, path string) (string, error) {
	script, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to read script: %w", DefaultName, err)
	}

	result, err := d.impl.Render(ctx, string(script))
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to execute script %q: %w", DefaultName, path, err)
	}

	return result, nil
}
