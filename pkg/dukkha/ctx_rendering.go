package dukkha

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"github.com/itchyny/gojq"
	"mvdan.cc/sh/v3/expand"
)

type RenderingContext interface {
	context.Context
	expand.Environ
	rs.InterfaceTypeHandler
	rs.RenderingHandler
	EnvValues

	// AddValues will merge provided values into existing values
	AddValues(values map[string]interface{}) error

	Env() map[string]string

	Values() map[string]interface{}
}

type RendererAttribute string

// Renderer to handle rendering suffix
type Renderer interface {
	rs.Field

	// Init the renderer and add itself to the context
	Init(ctx ConfigResolvingContext) error

	RenderYaml(rc RenderingContext, rawData interface{}, attributes []RendererAttribute) (result []byte, err error)
}

// RendererManager to manage renderers
type RendererManager interface {
	AllRenderers() map[string]Renderer
	AddRenderer(name string, renderer Renderer)
}

func newContextRendering(
	ctx *contextStd,
	ifaceTypeHandler rs.InterfaceTypeHandler,
	globalEnv map[string]string,
) *contextRendering {
	return &contextRendering{
		contextStd: ctx,

		envValues: newEnvValues(globalEnv),

		ifaceTypeHandler: ifaceTypeHandler,
		renderers:        make(map[string]Renderer),
		values:           make(map[string]interface{}),
	}
}

var (
	_ RendererManager  = (*contextRendering)(nil)
	_ RenderingContext = (*contextRendering)(nil)
)

type contextRendering struct {
	*contextStd
	*envValues

	ifaceTypeHandler rs.InterfaceTypeHandler
	renderers        map[string]Renderer

	values map[string]interface{}

	// nolint:revive
	_transform_value string
}

func (c *contextRendering) clone(newCtx *contextStd, deepCopy bool) *contextRendering {

	envValues := c.envValues
	if deepCopy {
		envValues = c.envValues.clone()
	}

	return &contextRendering{
		contextStd: newCtx,

		envValues: envValues,
		renderers: c.renderers,

		// values are global scoped, DO NOT deep copy in any case
		values: c.values,
	}
}

func (c *contextRendering) Env() map[string]string {
	for k, v := range c.envValues.globalEnv {
		c.envValues.env[k] = v
	}

	return c.envValues.env
}

func (c *contextRendering) AddValues(values map[string]interface{}) error {
	mergedValues, err := rs.MergeMap(c.values, values, false, false)
	if err != nil {
		return err
	}

	c.values = mergedValues
	return nil
}

func (c *contextRendering) Values() map[string]interface{} {
	return c.values
}

func (c *contextRendering) RenderYaml(renderer string, rawData interface{}) ([]byte, error) {
	var attributes []RendererAttribute
	attrStart := strings.LastIndexByte(renderer, '#')
	if attrStart != -1 {
		for _, attr := range strings.Split(renderer[attrStart+1:], ",") {
			attributes = append(attributes, RendererAttribute(strings.TrimSpace(attr)))
		}

		renderer = renderer[:attrStart]
	}

	v, ok := c.renderers[renderer]
	if !ok {
		return nil, fmt.Errorf("renderer %q not found", renderer)
	}

	return v.RenderYaml(c, rawData, attributes)
}

func (c *contextRendering) Create(typ reflect.Type, yamlKey string) (interface{}, error) {
	return c.ifaceTypeHandler.Create(typ, yamlKey)
}

func (c *contextRendering) AddRenderer(name string, r Renderer) {
	c.renderers[name] = r
}

func (c *contextRendering) AllRenderers() map[string]Renderer {
	return c.renderers
}

// SetVALUE for transform renderer
func (c *contextRendering) SetVALUE(value string) { c._transform_value = value }

// VALUE for transform renderer
func (c *contextRendering) VALUE() string { return c._transform_value }

// Get implements expand.Environ
func (c *contextRendering) Get(name string) expand.Variable {
	v, exists := c.Env()[name]
	if exists {
		return createVariable(v)
	}

	switch name {
	case "IFS":
		v = " \t\n"
	case "OPTIND":
		v = "1"
	case "PWD":
		v = c.WorkingDir()
	case "UID":
		// os.Getenv("UID") usually retruns empty value
		// so we have to call os.Getuid
		v = strconv.FormatInt(int64(os.Getuid()), 10)
	default:
		kind := expand.Unset
		if strings.HasPrefix(name, valuesEnvPrefix) {
			valRef := strings.TrimPrefix(name, valuesEnvPrefix)

			query, err := gojq.Parse("." + valRef)
			if err != nil {
				goto ret
			}

			result, found, err := textquery.RunQuery(query, c.values, nil)
			if err != nil {
				goto ret
			}

			if !found {
				goto ret
			}

			kind = expand.String
			v = textquery.HandleQueryResult(result, json.Marshal)
		}

	ret:
		return expand.Variable{
			Local:    false,
			Exported: true,
			ReadOnly: false,
			Kind:     kind,
			Str:      v,
		}
	}

	return createVariable(v)
}

// Each implements expand.Environ
func (c *contextRendering) Each(do func(name string, vr expand.Variable) bool) {
	for k, v := range c.Env() {
		if !do(k, createVariable(v)) {
			return
		}
	}

	values, _ := genEnvForValues(c.values)
	for k, v := range values {
		if !do(k, v) {
			return
		}
	}
}

const valuesEnvPrefix = "Values."

func genEnvForValues(values map[string]interface{}) (map[string]expand.Variable, error) {
	out := make(map[string]expand.Variable)
	for k, v := range values {
		err := doGenEnvForInterface(valuesEnvPrefix+k, v, &out)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func doGenEnvForInterface(prefix string, v interface{}, out *map[string]expand.Variable) error {
	switch t := v.(type) {
	case map[string]interface{}:
		dataBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}

		(*out)[prefix] = createVariable(string(dataBytes))

		for k, v := range t {
			err = doGenEnvForInterface(prefix+"."+k, v, out)
			if err != nil {
				return err
			}
		}
	case string:
		(*out)[prefix] = createVariable(t)
	case []byte:
		(*out)[prefix] = createVariable(string(t))
	default:
		dataBytes, err := json.Marshal(t)
		if err != nil {
			return err
		}

		(*out)[prefix] = createVariable(string(dataBytes))
	}

	return nil
}

// createVariable for embedded shell, if exists is false, will lookup values for the name
func createVariable(value string) expand.Variable {
	// TODO: set kind for lists
	return expand.Variable{
		Local:    false,
		Exported: true,
		ReadOnly: false,
		Kind:     expand.String,
		Str:      value,
	}
}
