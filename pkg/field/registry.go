package field

import (
	"fmt"
	"reflect"
	"regexp"
)

// nolint:revive
type (
	FieldFactoryFunc func() interface{}

	interfaceFieldFactoryKey struct {
		typ reflect.Type
	}

	interfaceFieldFactoryImpl struct {
		exp         *regexp.Regexp
		createField FieldFactoryFunc
	}

	interfaceFieldFactoryValue struct {
		factories []*interfaceFieldFactoryImpl
	}
)

var supportedInterfaceTypes = make(map[interfaceFieldFactoryKey]*interfaceFieldFactoryValue)

func RegisterInterfaceField(
	interfaceType reflect.Type,

	yamlKeyMatch *regexp.Regexp,
	createField FieldFactoryFunc,
) {
	key := interfaceFieldFactoryKey{
		typ: interfaceType,
	}

	v, ok := supportedInterfaceTypes[key]
	if ok {
		v.factories = append(v.factories, &interfaceFieldFactoryImpl{
			exp:         yamlKeyMatch,
			createField: createField,
		})
	} else {
		supportedInterfaceTypes[key] = &interfaceFieldFactoryValue{
			factories: []*interfaceFieldFactoryImpl{
				{
					exp:         yamlKeyMatch,
					createField: createField,
				},
			},
		}
	}
}

func CreateInterfaceField(interfaceType reflect.Type, yamlKey string) (interface{}, error) {
	key := interfaceFieldFactoryKey{
		typ: interfaceType,
	}

	v, ok := supportedInterfaceTypes[key]
	if !ok {
		return nil, fmt.Errorf("interface type %q not registered", interfaceType.String())
	}

	for _, impl := range v.factories {
		if impl.exp.MatchString(yamlKey) {
			return impl.createField(), nil
		}
	}

	return nil, fmt.Errorf("yaml field %q not resolved as %q", yamlKey, interfaceType.String())
}
