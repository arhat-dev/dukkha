package renderer

import (
	"arhat.dev/dukkha/pkg/field"
)

type Interface interface {
	Name() string

	Render(ctx *field.RenderingContext, rawData interface{}) (result string, err error)
}

type Config interface{}
