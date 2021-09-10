package dukkha

import (
	"context"

	"arhat.dev/rs"
	"mvdan.cc/sh/v3/expand"
)

type RenderingContext interface {
	context.Context

	expand.Environ

	GlobalValues
	EnvValues

	// AddValues will merge provided values into existing values
	AddValues(values map[string]interface{}) error

	Env() map[string]string

	Values() map[string]interface{}

	rs.InterfaceTypeHandler
	rs.RenderingHandler
}

// Renderer to handle rendering suffix
type Renderer interface {
	rs.Field

	// Init the renderer and add itself to the context
	Init(ctx ConfigResolvingContext) error

	RenderYaml(rc RenderingContext, rawData interface{}) (result []byte, err error)
}

// RendererManager to manage renderers
type RendererManager interface {
	AllRenderers() map[string]Renderer
	AddRenderer(name string, renderer Renderer)
}

func newContextRendering(
	ctx context.Context,
	ifaceTypeHandler rs.InterfaceTypeHandler,
	globalEnv map[string]string,
) *contextRendering {
	return &contextRendering{
		Context: ctx,

		envValues: newEnvValues(globalEnv),

		ifaceTypeHandler: ifaceTypeHandler,
		renderers:        make(map[string]Renderer),
		values:           make(map[string]interface{}),
	}
}

var (
	_ RendererManager  = (*contextRendering)(nil)
	_ RenderingContext = (*contextRendering)(nil)
)

type contextRendering struct {
	context.Context

	*envValues

	ifaceTypeHandler rs.InterfaceTypeHandler
	renderers        map[string]Renderer

	values map[string]interface{}
}
