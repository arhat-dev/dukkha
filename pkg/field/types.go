package field

import (
	"context"
	"reflect"

	"gopkg.in/yaml.v3"
)

type RenderingFunc func(ctx context.Context, renderer, rawData string) (string, error)

type Interface interface {
	Type() reflect.Type

	yaml.Unmarshaler

	// Resolve
	Resolve(ctx context.Context, render RenderingFunc, depth int) error
}
