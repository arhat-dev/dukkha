package dukkha

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"arhat.dev/rs"
)

func ResolveEnv(t rs.Field, mCtx RenderingContext, envFieldName, envTagName string) error {
	err := t.ResolveFields(mCtx, 1, envTagName)
	if err != nil {
		return fmt.Errorf("failed to get env overview: %w", err)
	}

	fv := reflect.ValueOf(t)
	for fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
	}

	// avoid panic
	if fv.Kind() != reflect.Struct {
		return fmt.Errorf("unexpected non struct target: %T", t)
	}

	env := fv.FieldByName(envFieldName).Interface().(Env)
	for i := range env {
		err = env[i].ResolveFields(mCtx, -1)
		if err != nil {
			return fmt.Errorf("failed to resolve env %q: %w", env[i].Name, err)
		}

		mCtx.AddEnv(true, env[i])
	}

	return nil
}

type Env []*EnvEntry

func (orig Env) Clone() Env {
	ret := make(Env, 0, len(orig))
	for _, entry := range orig {
		ret = append(ret, &EnvEntry{
			Name:  entry.Name,
			Value: entry.Value,
		})
	}

	return ret
}

type EnvEntry struct {
	rs.BaseField `yaml:"-"`

	// Name of the entry (in other words, key)
	Name string `yaml:"name"`

	// Value associated to the name
	Value string `yaml:"value"`
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

func RegisterRenderer(name string, create RendererCreateFunc) {
	if strings.Contains(name, ":") {
		panic(fmt.Sprintf("invalid renderer name %q containing `:`", name))
	}

	globalTypeManager.register(
		name,
		rendererType,
		regexp.MustCompile(fmt.Sprintf(`^%s(:.+){0,1}$`, name)),
		func(subMatches []string) interface{} {
			if len(subMatches) > 1 {
				return create(name + ":" + subMatches[1])
			}

			return create(name)
		},
	)
}

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
			if len(subMatches) > 1 {
				return create(subMatches[1])
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
