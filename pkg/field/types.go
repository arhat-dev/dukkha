package field

import (
	"reflect"

	"gopkg.in/yaml.v3"
)

type Field interface {
	yaml.Unmarshaler

	// ResolveFields resolves yaml fields using rendering suffix
	// when depth >= 1, resolve inner fields until reaching depth limit
	// when depth == 0, do nothing
	// when depth < 0, resolve recursively
	//
	// when fieldName is not empty, resolve single field
	// when fieldName is empty, resolve all fields in the struct
	ResolveFields(rc RenderingHandler, depth int, fieldName string) error
}

type RenderingHandler interface {
	// RenderYaml using specified renderer
	RenderYaml(renderer string, rawData interface{}) (result []byte, err error)
}

type InterfaceTypeHandler interface {
	Create(typ reflect.Type, yamlKey string) (interface{}, error)
}
