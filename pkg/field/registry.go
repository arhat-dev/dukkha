package field

import (
	"fmt"
	"reflect"
	"regexp"
)

// nolint:revive
type (
	FieldFactoryFunc func(params []string) interface{}

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
		v.factories = append(v.factories,
			&interfaceFieldFactoryImpl{
				exp:         yamlKeyMatch,
				createField: createField,
			},
		)

		return
	}

	supportedInterfaceTypes[key] = &interfaceFieldFactoryValue{
		factories: []*interfaceFieldFactoryImpl{
			{
				exp:         yamlKeyMatch,
				createField: createField,
			},
		},
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
		if !impl.exp.MatchString(yamlKey) {
			continue
		}

		if impl.exp.NumSubexp() == 0 {
			return impl.createField(nil), nil
		}

		return impl.createField(impl.exp.FindStringSubmatch(yamlKey)[1:]), nil
	}

	return nil, fmt.Errorf("yaml field %q not resolved as %q", yamlKey, interfaceType.String())
}
