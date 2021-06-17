package renderer

import (
	"fmt"
	"reflect"
)

type (
	factoryKey struct {
		configTypeString string
	}

	factoryValue struct {
		create RendererFactoryFunc
	}

	RendererFactoryFunc func(config interface{}) (Interface, error)
)

var (
	supportedRenderers = make(map[factoryKey]factoryValue)
)

func Register(config interface{}, create RendererFactoryFunc) {
	supportedRenderers[createFactoryKey(config)] = factoryValue{
		create: create,
	}
}

func Create(config interface{}) (Interface, error) {
	f, ok := supportedRenderers[createFactoryKey(config)]
	if !ok {
		return nil, fmt.Errorf("renderer for config %T not found", config)
	}

	return f.create(config)
}

func createFactoryKey(config interface{}) factoryKey {
	return factoryKey{
		configTypeString: reflect.TypeOf(config).String(),
	}
}
