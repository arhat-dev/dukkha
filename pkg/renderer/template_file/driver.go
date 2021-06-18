package template_file

import (
	"context"
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/renderer"
)

const DefaultName = "template_file"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct {
	impl *Driver
}

func (d *Driver) Name() string {
	return DefaultName
}

func (d *Driver) Render(ctx context.Context, rawValue string, v *renderer.RenderingValues) (string, error) {
	tplBytes, err := os.ReadFile(rawValue)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to read template file: %w", DefaultName, err)
	}

	result, err := d.impl.Render(ctx, string(tplBytes), v)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to render file %q: %w", DefaultName, rawValue, err)
	}

	return result, nil
}
