package rs

import (
	"errors"
	"reflect"

	"gopkg.in/yaml.v3"
)

const (
	TagNameRS = "rs"
	TagName   = TagNameRS
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
	ResolveFields(rc RenderingHandler, depth int, fieldNames ...string) error
}

type RenderingHandler interface {
	// RenderYaml using specified renderer
	RenderYaml(renderer string, rawData interface{}) (result []byte, err error)
}

var (
	ErrInterfaceTypeNotHandled = errors.New("interface type not handled")
)

type InterfaceTypeHandler interface {
	// Create request interface type using yaml information
	Create(typ reflect.Type, yamlKey string) (interface{}, error)
}

type InterfaceTypeHandleFunc func(typ reflect.Type, yamlKey string) (interface{}, error)

func (f InterfaceTypeHandleFunc) Create(typ reflect.Type, yamlKey string) (interface{}, error) {
	return f(typ, yamlKey)
}
