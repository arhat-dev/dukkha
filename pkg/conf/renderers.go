package conf

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/rs"
)

// RendererGroup contains a group of renderers can be initialized
// without depending on each other
type RendererGroup struct {
	rs.BaseField

	Renderers map[string]dukkha.Renderer `yaml:",inline"`
}
