package field

import (
	"context"

	"gopkg.in/yaml.v3"
)

type RenderingFunc func(ctx context.Context, renderer, rawData string) (string, error)

type Interface interface {
	yaml.Unmarshaler

	// Resolve
	Resolve(ctx context.Context, render RenderingFunc, depth int) error
}
