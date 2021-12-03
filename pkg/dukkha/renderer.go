package dukkha

import (
	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"
)

type RendererAttribute string

// Renderer to handle rendering suffix
type Renderer interface {
	rs.Field

	// Init the renderer and add itself to the context
	Init(cacheFS *fshelper.OSFS) error

	RenderYaml(rc RenderingContext, rawData interface{}, attributes []RendererAttribute) (result []byte, err error)
}

// RendererManager to manage renderers
type RendererManager interface {
	AllRenderers() map[string]Renderer
	AddRenderer(name string, renderer Renderer)
}
