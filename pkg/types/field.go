package types

import (
	"reflect"

	"gopkg.in/yaml.v3"
)

type Field interface {
	yaml.Unmarshaler

	// ResolveFields resolves yaml fields using rendering suffix
	ResolveFields(rc RenderingContext, depth int, fieldName string) error
}

type InterfaceTypeHandler interface {
	Create(typ reflect.Type, yamlKey string) (interface{}, error)
}
