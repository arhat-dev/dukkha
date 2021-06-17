package file

import (
	"context"
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/renderer"
)

const (
	DefaultName = "file"
)

var _ renderer.Interface = (*Driver)(nil)

type Driver struct{}

func (d *Driver) Name() string { return DefaultName }

func (d *Driver) Render(ctx context.Context, rawValue string) (string, error) {
	data, err := os.ReadFile(rawValue)
	if err != nil {
		return "", fmt.Errorf("renderer.file: %w", err)
	}

	return string(data), err
}
