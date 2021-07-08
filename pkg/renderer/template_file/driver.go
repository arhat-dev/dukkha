package template_file

import (
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/types"
)

const DefaultName = "template_file"

func New() dukkha.Renderer {
	return &driver{impl: template.New()}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct {
	impl dukkha.Renderer
}

func (d *driver) Name() string {
	return DefaultName
}

func (d *driver) RenderYaml(rc types.RenderingContext, rawData interface{}) ([]byte, error) {
	path, ok := rawData.(string)
	if !ok {
		return nil, fmt.Errorf("renderer.%s: unexpected non string input %T", DefaultName, rawData)
	}

	tplBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: failed to read template file: %w", DefaultName, err)
	}

	result, err := d.impl.RenderYaml(rc, tplBytes)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: failed to render file %q: %w", DefaultName, path, err)
	}

	return result, nil
}
