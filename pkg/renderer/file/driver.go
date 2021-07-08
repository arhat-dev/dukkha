package file

import (
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/types"
)

// nolint:revive
const (
	DefaultName = "file"
)

func New() dukkha.Renderer {
	return &driver{}
}

var _ dukkha.Renderer = (*driver)(nil)

type driver struct{}

func (d *driver) Name() string { return DefaultName }

func (d *driver) RenderYaml(_ types.RenderingContext, rawData interface{}) ([]byte, error) {
	path, ok := rawData.(string)
	if !ok {
		return nil, fmt.Errorf("renderer.%s: unexpected non-string input %T", DefaultName, rawData)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	return data, err
}
