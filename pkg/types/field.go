package types

import (
	"gopkg.in/yaml.v3"
)

type Field interface {
	yaml.Unmarshaler

	// ResolveFields resolves yaml fields using rendering suffix
	ResolveFields(rc RenderingContext, depth int, fieldName string) error
}
