package dukkha

import (
	"context"
	"fmt"
	"sync"

	"arhat.dev/dukkha/pkg/types"
)

// Renderer to handle rendering suffix
type Renderer interface {
	RenderYaml(rc types.RenderingContext, rawData interface{}) (result []byte, err error)
}

// RendererManager to manage renderers
type RendererManager interface {
	AddRenderer(renderer Renderer, names ...string) error
}

func newContextRendering(ctx context.Context, globalEnv map[string]string) *contextRendering {
	return &contextRendering{
		Context: ctx,

		immutableValues: newContextImmutableValues(globalEnv),
		mutableValues:   newContextMutableValues(),

		renderers: new(sync.Map),
	}
}

var (
	_ RendererManager        = (*contextRendering)(nil)
	_ types.RenderingContext = (*contextRendering)(nil)
)

type contextRendering struct {
	context.Context

	*mutableValues
	*immutableValues

	renderers *sync.Map
}

func (c *contextRendering) clone(newCtx context.Context) *contextRendering {
	return &contextRendering{
		Context: newCtx,

		immutableValues: c.immutableValues,
		mutableValues:   c.mutableValues.clone(),
		renderers:       c.renderers,
	}
}

func (c *contextRendering) Env() map[string]string {
	for k, v := range c.immutableValues.globalEnv {
		c.mutableValues.env[k] = v
	}

	return c.mutableValues.env
}

func (c *contextRendering) RenderYaml(renderer string, rawData interface{}) ([]byte, error) {
	v, ok := c.renderers.Load(renderer)
	if !ok {
		return nil, fmt.Errorf("renderer %q not found", renderer)
	}

	return v.(Renderer).RenderYaml(c, rawData)
}

func (c *contextRendering) AddRenderer(renderer Renderer, names ...string) error {
	for _, name := range names {
		c.renderers.Store(name, renderer)
	}

	return nil
}
