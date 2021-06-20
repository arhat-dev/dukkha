package field

import (
	"context"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Interface interface {
	yaml.Unmarshaler

	// ResolveFields resolves struct fields with rendering suffix
	ResolveFields(ctx *RenderingContext, render RenderingFunc, depth int) error
}

type (
	RenderingFunc func(ctx *RenderingContext, renderer, rawData string) (string, error)

	RenderingValues struct {
		Env map[string]string
	}

	RenderingContext struct {
		ctx    context.Context
		values *RenderingValues
	}
)

func WithRenderingValues(ctx context.Context, env []string) *RenderingContext {
	ret := &RenderingContext{
		ctx: ctx,
		values: &RenderingValues{
			Env: make(map[string]string),
		},
	}

	for _, e := range append(os.Environ(), env...) {
		ret.SetEnv(e)
	}

	return ret
}

func (c *RenderingContext) Context() context.Context {
	return c.ctx
}

func (c *RenderingContext) SetEnv(entry string) {
	parts := strings.SplitN(entry, "=", 2)
	key, value := parts[0], ""
	if len(parts) == 2 {
		value = parts[1]
	}

	// do not expand environment variables here
	c.values.Env[key] = value
}

func (c *RenderingContext) Values() *RenderingValues {
	return c.values
}
