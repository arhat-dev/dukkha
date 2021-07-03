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

func (d *Driver) Render(ctx *field.RenderingContext, rawData interface{}) (string, error) {
	path, ok := rawData.(string)
	if !ok {
		return "", fmt.Errorf("renderer.%s: unexpected non string input %T", DefaultName, rawData)
	}

	tplBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to read template file: %w", DefaultName, err)
	}

	result, err := d.impl.Render(ctx, tplBytes)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: failed to render file %q: %w", DefaultName, path, err)
	}

	return result, nil
}
