package dukkha

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
)

func ResolveEnv(t field.Field, mCtx RenderingContext, envFieldName string) error {
	err := t.ResolveFields(mCtx, 1, envFieldName)
	if err != nil {
		return fmt.Errorf("failed to get env overview: %w", err)
	}

	fv := reflect.ValueOf(t)
	for fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
	}

	env := fv.FieldByName(envFieldName).Interface().(Env)
	for i := range env {
		err = env[i].ResolveFields(mCtx, -1, "")
		if err != nil {
			return fmt.Errorf("failed to resolve env %q: %w", env[i].Name, err)
		}

		mCtx.AddEnv(true, env[i])
	}

	return nil
}

type Env []EnvEntry

func (orig Env) Clone() Env {
	ret := make(Env, 0, len(orig))
	for _, entry := range orig {
		ret = append(ret, EnvEntry{
			Name:  entry.Name,
			Value: entry.Value,
		})
	}

	return ret
}

type EnvEntry struct {
	field.BaseField

	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

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
	if strings.Contains(name, ":") {
		panic(fmt.Sprintf("invalid renderer name %q containing `:`", name))
	}

	globalTypeManager.register(
		rendererType,
		regexp.MustCompile(fmt.Sprintf(`^%s$`, name)),
		func(subMatches []string) interface{} {
			return create()
		},
	)
}

func RegisterTool(k ToolKind, create ToolCreateFunc) {
	if strings.Contains(string(k), ":") {
		panic(fmt.Sprintf("invalid tool kind %q containing `:`", k))
	}

	globalTypeManager.register(
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
