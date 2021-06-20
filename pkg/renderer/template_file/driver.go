package template_file

import (
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/template"
)

const DefaultName = "template_file"

var _ renderer.Interface = (*Driver)(nil)

type Driver struct {
	impl *template.Driver
}

func (d *Driver) Name() string {
	return DefaultName
}

func (d *Driver) Render(ctx *field.RenderingContext, path string) (string, error) {
	tplBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to read template file: %w", DefaultName, err)
	}

	result, err := d.impl.Render(ctx, string(tplBytes))
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to render file %q: %w", DefaultName, path, err)
	}

	return result, nil
}
