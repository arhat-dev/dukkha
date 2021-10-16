package rs

import (
	"reflect"

	"gopkg.in/yaml.v3"
)

// Struct field tag names used by this package
const (
	TagNameRS = "rs"
	TagName   = TagNameRS
)

// Field defines methods required for rendering suffix support
//
// these methods are implemented by BaseField, so usually you should just embed
// BaseField in your struct as the very first field
type Field interface {
	// A Field must implement yaml.Unmarshaler interface to process rendering suffix
	// in yaml key
	yaml.Unmarshaler

	// ResolveFields sets struct field values with data from unmarshaled yaml
	//
	// depth controls to what level this function call goes
	// when depth >= 1, resolve inner fields until reaching depth limit (depth == 0)
	// when depth == 0, do nothing
	// when depth < 0, resolve recursively
	//
	// fieldNames instructs which fields to be resolved. When it's not empty,
	// resolve specified fields only, otherwise, resolve all exported fields
	// in the underlying struct.
	ResolveFields(rc RenderingHandler, depth int, fieldNames ...string) error
}

type (
	// RenderingHandler is used when resolving yaml fields
	RenderingHandler interface {
		// RenderYaml transforms rawData to result through the renderer
		//
		// renderer is the name of your renderer without type hint (`?<hint>`) and patch suffix (`!`)
		//
		// rawData is the input to your renderer, which can be any type, your renderer is responsible
		// for casting it to its desired data type
		//
		// result can be any type as well, but if you returns an custom struct
		// you must make sure all desired value in it can be marshaled using yaml.Marshal
		RenderYaml(renderer string, rawData interface{}) (result interface{}, err error)
	}

	// RenderingHandleFunc is a helper type to wrap your function as RenderingHandler
	RenderingHandleFunc func(renderer string, rawData interface{}) (result interface{}, err error)
)

var _ RenderingHandler = RenderingHandleFunc(nil)

func (f RenderingHandleFunc) RenderYaml(renderer string, rawData interface{}) (result interface{}, err error) {
	return f(renderer, rawData)
}

type (
	// InterfaceTypeHandler is used when setting values for interface{} typed field
	InterfaceTypeHandler interface {
		// Create request interface type using yaml information
		Create(typ reflect.Type, yamlKey string) (interface{}, error)
	}

	// InterfaceTypeHandleFunc is a helper type to wrap your function as InterfaceTypeHandler
	InterfaceTypeHandleFunc func(typ reflect.Type, yamlKey string) (interface{}, error)
)

var _ InterfaceTypeHandler = InterfaceTypeHandleFunc(nil)

func (f InterfaceTypeHandleFunc) Create(typ reflect.Type, yamlKey string) (interface{}, error) {
	return f(typ, yamlKey)
}
