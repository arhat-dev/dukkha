package conf

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
)

// RendererGroup contains a group of renderers can be initialized
// without depending on each other
type RendererGroup struct {
	rs.BaseField

	Renderers map[string]dukkha.Renderer `yaml:",inline"`
}
