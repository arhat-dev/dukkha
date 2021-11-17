package dukkha

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"arhat.dev/rs"
)

// Resolvable represents a kind of struct that can be resolved at runtime
type Resolvable interface {
	rs.Field

	// DoAfterFieldsResolved is a helper function to ensure no data race
	//
	// The implementation MUST be safe to be used concurrently
	DoAfterFieldsResolved(
		rc RenderingContext,
		depth int,
		resolveEnv bool,
		do func() error,
		tagNames ...string,
	) error
}

type (
	RendererCreateFunc func(name string) Renderer

	ToolCreateFunc func() Tool
	TaskCreateFunc func(toolName string) Task
)

var globalTypeManager = &TypeManager{
	types: make(map[IfaceTypeKey]*IfaceFactory),
}

var GlobalInterfaceTypeHandler rs.InterfaceTypeHandler = globalTypeManager

// type values for interface type registration
var (
	rendererType = reflect.TypeOf((*Renderer)(nil)).Elem()
	toolType     = reflect.TypeOf((*Tool)(nil)).Elem()
	taskType     = reflect.TypeOf((*Task)(nil)).Elem()
)

// RegisterRenderer associates a renderer factory with name
func RegisterRenderer(name string, create RendererCreateFunc) {
	if strings.Contains(name, ":") || strings.Contains(name, "#") {
		panic(fmt.Sprintf("invalid renderer name %q containing `:` or `#`", name))
	}

	globalTypeManager.register(
		name,
		rendererType,
		regexp.MustCompile(fmt.Sprintf(`^%s(:.+){0,1}$`, name)),
		func(subMatches []string) interface{} {
			if len(subMatches) > 1 && len(subMatches[1]) > 1 {
				return create(name + subMatches[1])
			}

			return create(name)
		},
	)
}

// RegisterTool associates a tool factory with tool kind
func RegisterTool(k ToolKind, create ToolCreateFunc) {
	if strings.Contains(string(k), ":") {
		panic(fmt.Sprintf("invalid tool kind %q containing `:`", k))
	}

	globalTypeManager.register(
		string(k),
		toolType,
		regexp.MustCompile(fmt.Sprintf(`^%s$`, string(k))),
		func(subMatches []string) interface{} { return create() },
	)
}

// RegisterTask associates a task factory with task kind
func RegisterTask(k ToolKind, tk TaskKind, create TaskCreateFunc) {
	if strings.Contains(string(k), ":") {
		panic(fmt.Sprintf("invalid tool kind %q containing `:`", k))
	}

	if strings.Contains(string(tk), ":") {
		panic(fmt.Sprintf("invalid task kind %q containing `:`", tk))
	}

	globalTypeManager.register(
		string(k)+":"+string(tk),
		taskType,
		regexp.MustCompile(
			fmt.Sprintf(`^%s(:.+){0,1}:%s$`, string(k), string(tk)),
		),
		func(subMatches []string) interface{} {
			if len(subMatches) > 1 && len(subMatches[1]) > 1 {
				return create(subMatches[1][1:])
			}

			return create("")
		},
	)
}

// nolint:revive
type (
	IfaceTypeKey struct {
		Typ reflect.Type
	}

	IfaceFactoryFunc func(subMatches []string) interface{}

	IfaceFactoryImpl struct {
		// Name is the raw information about what instance we are creating
		// currently only used in json schema generation
		Name string
		exp  *regexp.Regexp

		Create IfaceFactoryFunc
	}

	IfaceFactory struct {
		Factories []*IfaceFactoryImpl
	}
)

var _ rs.InterfaceTypeHandler = (*TypeManager)(nil)

type TypeManager struct {
	types map[IfaceTypeKey]*IfaceFactory
}

func (h *TypeManager) Types() map[IfaceTypeKey]*IfaceFactory {
	return h.types
}

// Create implements rs.InterfaceTypeHandler
func (h *TypeManager) Create(typ reflect.Type, yamlKey string) (interface{}, error) {
	key := IfaceTypeKey{
		Typ: typ,
	}

	v, ok := h.types[key]
	if !ok {
		return nil, fmt.Errorf(
			"interface type %q not registered: %w",
			typ.String(), rs.ErrInterfaceTypeNotHandled,
		)
	}

	for _, impl := range v.Factories {
		if !impl.exp.MatchString(yamlKey) {
			continue
		}

		if impl.exp.NumSubexp() == 0 {
			return impl.Create(nil), nil
		}

		return impl.Create(impl.exp.FindStringSubmatch(yamlKey)), nil
	}

	return nil, fmt.Errorf("yaml field %q not resolved as %q", yamlKey, typ.String())
}

func (h *TypeManager) register(
	name string,
	ifaceType reflect.Type,
	yamlKeyMatch *regexp.Regexp,
	createField IfaceFactoryFunc,
) {
	key := IfaceTypeKey{
		Typ: ifaceType,
	}

	v, ok := h.types[key]
	if ok {
		v.Factories = append(v.Factories,
			&IfaceFactoryImpl{
				Name: name,
				exp:  yamlKeyMatch,

				Create: createField,
			},
		)

		return
	}

	h.types[key] = &IfaceFactory{
		Factories: []*IfaceFactoryImpl{{
			Name: name,
			exp:  yamlKeyMatch,

			Create: createField,
		}},
	}
}
