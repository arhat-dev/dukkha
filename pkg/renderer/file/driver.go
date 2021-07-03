package file

import (
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

// nolint:revive
const (
	DefaultName = "file"
)

var _ renderer.Interface = (*Driver)(nil)

type Driver struct{}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(_ *field.RenderingContext, rawData interface{}) (string, error) {
	path, ok := rawData.(string)
	if !ok {
		return "", fmt.Errorf("renderer.%s: unexpected non-string input %T", DefaultName, rawData)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("renderer.%s: %w", DefaultName, err)
	}

	return string(data), err
}
