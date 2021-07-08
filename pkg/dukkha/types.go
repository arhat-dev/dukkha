package dukkha

import (
	"fmt"
	"reflect"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
)

type (
	RendererCreateFunc func() Renderer

	ToolCreateFunc func() Tool
	TaskCreateFunc func(toolName string) Task
)

var globalTypeManager = &typeManager{
	types: make(map[ifaceTypeKey]*ifaceFactory),
}

var GlobalInterfaceTypeHandler field.InterfaceTypeHandler = globalTypeManager

// type values for interface type registration
var (
	rendererType = reflect.TypeOf((*Renderer)(nil)).Elem()
	toolType     = reflect.TypeOf((*Tool)(nil)).Elem()
	taskType     = reflect.TypeOf((*Task)(nil)).Elem()
)

func RegisterRenderer(name string, create RendererCreateFunc) {
	globalTypeManager.register(
		rendererType,
		regexp.MustCompile(`^%s(:.+)?$`),
		func(subMatches []string) interface{} {
			return create()
		},
	)
}

func RegisterTool(kind ToolKind, create ToolCreateFunc) {
	globalTypeManager.register(
		toolType,
		regexp.MustCompile(fmt.Sprintf(`^%s$`, string(kind))),
		func(subMatches []string) interface{} { return create() },
	)
}

func RegisterTask(k ToolKind, tk TaskKind, create TaskCreateFunc) {
	globalTypeManager.register(
		taskType,
		regexp.MustCompile(
			fmt.Sprintf(`^%s(:.+){0,1}:%s$`, string(k), string(tk)),
		),
		func(subMatches []string) interface{} {
			if len(subMatches) > 1 {
				return create(subMatches[1])
			}

			return create("")
		},
	)
}

// nolint:revive
type (
	ifaceTypeKey struct {
		typ reflect.Type
	}

	ifaceFactoryFunc func(subMatches []string) interface{}

	ifaceFactoryImpl struct {
		exp         *regexp.Regexp
		createField ifaceFactoryFunc
	}

	ifaceFactory struct {
		factories []*ifaceFactoryImpl
	}
)

var _ field.InterfaceTypeHandler = (*typeManager)(nil)

type typeManager struct {
	types map[ifaceTypeKey]*ifaceFactory
}

func (h *typeManager) Create(typ reflect.Type, yamlKey string) (interface{}, error) {
	key := ifaceTypeKey{
		typ: typ,
	}

	v, ok := h.types[key]
	if !ok {
		return nil, fmt.Errorf("interface type %q not registered", typ.String())
	}

	for _, impl := range v.factories {
		if !impl.exp.MatchString(yamlKey) {
			continue
		}

		if impl.exp.NumSubexp() == 0 {
			return impl.createField(nil), nil
		}

		return impl.createField(impl.exp.FindStringSubmatch(yamlKey)), nil
	}

	return nil, fmt.Errorf("yaml field %q not resolved as %q", yamlKey, typ.String())
}

func (h *typeManager) register(
	ifaceType reflect.Type,
	yamlKeyMatch *regexp.Regexp,
	createField ifaceFactoryFunc,
) {
	key := ifaceTypeKey{
		typ: ifaceType,
	}

	v, ok := h.types[key]
	if ok {
		v.factories = append(v.factories,
			&ifaceFactoryImpl{
				exp:         yamlKeyMatch,
				createField: createField,
			},
		)

		return
	}

	h.types[key] = &ifaceFactory{
		factories: []*ifaceFactoryImpl{{
			exp:         yamlKeyMatch,
			createField: createField,
		}},
	}
}
