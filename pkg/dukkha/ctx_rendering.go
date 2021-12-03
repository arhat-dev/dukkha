package dukkha

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"github.com/itchyny/gojq"
	"mvdan.cc/sh/v3/expand"

	"arhat.dev/dukkha/pkg/utils"
)

type RenderingContext interface {
	context.Context
	expand.Environ
	rs.InterfaceTypeHandler
	rs.RenderingHandler
	EnvValues

	// FS returns the filesystem with cwd set to DUKKHA_WORKDIR
	FS() *fshelper.OSFS

	// AddValues will merge provided values into existing values
	AddValues(values map[string]interface{}) error

	Env() map[string]utils.LazyValue

	Values() map[string]interface{}
}

func newContextRendering(
	ctx *contextStd,
	ifaceTypeHandler rs.InterfaceTypeHandler,
	globalEnv map[string]utils.LazyValue,
) *contextRendering {
	envValues := newEnvValues(globalEnv)
	return &contextRendering{
		contextStd: ctx,

		envValues: envValues,

		ifaceTypeHandler: ifaceTypeHandler,
		renderers:        make(map[string]Renderer),
		values:           make(map[string]interface{}),

		fs: fshelper.NewOSFS(false, func() (string, error) {
			return envValues.WorkDir(), nil
		}),
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
	_VALUE interface{}

	fs *fshelper.OSFS
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

		fs: fshelper.NewOSFS(false, func() (string, error) {
			return envValues.WorkDir(), nil
		}),
	}
}

func (c *contextRendering) FS() *fshelper.OSFS { return c.fs }

// Env returns all environment variables available
// global environment variables are always kept
func (c *contextRendering) Env() map[string]utils.LazyValue {
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

// Get implements expand.Environ
func (c *contextRendering) Get(name string) expand.Variable {
	v, exists := c.globalEnv[name]
	if exists {
		return createVariable(v.Get())
	}

	v, exists = c.env[name]
	if exists {
		return createVariable(v.Get())
	}

	// non existing env

	// TODO: can we remove all these cases? (except the default case)
	switch name {
	case "IFS":
		v = utils.ImmediateString(" \t\n")
	case "OPTIND":
		v = utils.ImmediateString("1")
	case "PWD":
		v = utils.ImmediateString(c.WorkDir())
	case "UID":
		v = utils.ImmediateString(
			strconv.FormatInt(int64(os.Getuid()), 10),
		)
	case "GID":
		v = utils.ImmediateString(
			strconv.FormatInt(int64(os.Getegid()), 10),
		)
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
			v = utils.ImmediateString(textquery.HandleQueryResult(result, json.Marshal))
		}

	ret:
		str := ""
		if v != nil {
			str = v.Get()
		}
		return expand.Variable{
			Local:    false,
			Exported: true,
			ReadOnly: false,
			Kind:     kind,
			Str:      str,
		}
	}

	return createVariable(v.Get())
}

// Each implements expand.Environ
func (c *contextRendering) Each(do func(name string, vr expand.Variable) bool) {
	visited := make(map[string]struct{})

	for k, v := range c.globalEnv {
		visited[k] = struct{}{}

		if !do(k, createVariable(v.Get())) {
			return
		}
	}

	for k, v := range c.env {
		// do not override
		if _, ok := visited[k]; ok {
			continue
		}

		if !do(k, createVariable(v.Get())) {
			return
		}
	}

	values, _ := genEnvFromValues(c.values)
	for k, v := range values {
		if !do(k, v) {
			return
		}
	}
}

const valuesEnvPrefix = "values."

func genEnvFromValues(values map[string]interface{}) (map[string]expand.Variable, error) {
	out := make(map[string]expand.Variable)
	for k, v := range values {
		err := genEnvFromInterface(valuesEnvPrefix+k, v, &out)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func genEnvFromInterface(prefix string, v interface{}, out *map[string]expand.Variable) error {
	switch t := v.(type) {
	case map[string]interface{}:
		dataBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}

		(*out)[prefix] = createVariable(string(dataBytes))

		for k, v := range t {
			err = genEnvFromInterface(prefix+"."+k, v, out)
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
