package dukkha

import (
	"context"
	"fmt"

	"arhat.dev/dukkha/pkg/field"
)

type RenderingContext interface {
	context.Context

	ImmutableValues
	MutableValues

	Env() map[string]string

	field.RenderingHandler
}

// Renderer to handle rendering suffix
type Renderer interface {
	field.Field

	// Init the renderer and add itself to the context
	Init(ctx ConfigResolvingContext) error

	RenderYaml(rc RenderingContext, rawData interface{}) (result []byte, err error)
}

// RendererManager to manage renderers
type RendererManager interface {
	AllRenderers() map[string]Renderer
	AddRenderer(name string, renderer Renderer)
}

func newContextRendering(ctx context.Context, globalEnv map[string]string) *contextRendering {
	return &contextRendering{
		Context: ctx,

		immutableValues: newContextImmutableValues(globalEnv),
		mutableValues:   newContextMutableValues(),

		renderers: make(map[string]Renderer),
	}
}

var (
	_ RendererManager  = (*contextRendering)(nil)
	_ RenderingContext = (*contextRendering)(nil)
)

type contextRendering struct {
	context.Context

	*mutableValues
	*immutableValues

	renderers map[string]Renderer
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
	v, ok := c.renderers[renderer]
	if !ok {
		return nil, fmt.Errorf("renderer %q not found", renderer)
	}

	return v.RenderYaml(c, rawData)
}

func (c *contextRendering) AddRenderer(name string, r Renderer) {
	c.renderers[name] = r
}

func (c *contextRendering) AllRenderers() map[string]Renderer {
	return c.renderers
}
