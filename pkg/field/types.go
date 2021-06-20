package field

import (
	"context"
	"os"
	"strings"

	"arhat.dev/pkg/envhelper"
	"gopkg.in/yaml.v3"
)

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
	osEnv := os.Environ()

	envMap := make(map[string]string)

	for _, e := range osEnv {
		parts := strings.SplitN(e, "=", 2)
		key, value := parts[0], ""
		if len(parts) == 2 {
			value = parts[1]
		}

		envMap[key] = value
	}

	for _, e := range env {
		envhelper.Expand(e, func(varName, origin string) string {

			return ""
		})
	}

	return &RenderingContext{
		ctx: ctx,
		values: &RenderingValues{
			Env: envMap,
		},
	}
}

func (c *RenderingContext) Context() context.Context {
	return c.ctx
}

func (c *RenderingContext) Values() *RenderingValues {
	return c.values
}

type Interface interface {
	yaml.Unmarshaler

	// ResolveFields resolves struct fields with rendering suffix
	ResolveFields(ctx *RenderingContext, render RenderingFunc, depth int) error
}
