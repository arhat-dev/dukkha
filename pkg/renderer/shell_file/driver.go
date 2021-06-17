package shell_file

import (
	"context"
	"fmt"
	"os"

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

func (d *Driver) Render(ctx context.Context, path string) (string, error) {
	script, err := os.ReadFile(string(path))
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to read script: %w", DefaultName, err)
	}

	result, err := d.impl.Render(ctx, string(script))
	if err != nil {
		return "", fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	return result, nil
}
