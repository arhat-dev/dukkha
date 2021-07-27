package dukkha

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"mvdan.cc/sh/v3/expand"

	"arhat.dev/dukkha/pkg/field"
)

type RenderingContext interface {
	context.Context

	expand.Environ

	ImmutableValues
	MutableValues

	Env() map[string]string

	field.InterfaceTypeHandler
	field.RenderingHandler
}

// Renderer to handle rendering suffix
type Renderer interface {
	field.Field

	// Init the renderer and add itself to the context
	Init(ctx ConfigResolvingContext) error

	RenderYaml(rc RenderingContext, rawData interface{}) (result []byte, err error)
}

// RendererManager to manage renderers
type RendererManager interface {
	AllRenderers() map[string]Renderer
	AddRenderer(name string, renderer Renderer)
}

func newContextRendering(
	ctx context.Context,
	globalEnv map[string]string,
	ifaceTypeHandler field.InterfaceTypeHandler,
) *contextRendering {
	return &contextRendering{
		Context: ctx,

		immutableValues: newContextImmutableValues(globalEnv),
		mutableValues:   newContextMutableValues(),

		ifaceTypeHandler: ifaceTypeHandler,
		renderers:        make(map[string]Renderer),
	}
}

var (
	_ RendererManager  = (*contextRendering)(nil)
	_ RenderingContext = (*contextRendering)(nil)
)

type contextRendering struct {
	context.Context

	*mutableValues
	*immutableValues

	ifaceTypeHandler field.InterfaceTypeHandler
	renderers        map[string]Renderer
}

func (c *contextRendering) clone(newCtx context.Context) *contextRendering {
	return &contextRendering{
		Context: newCtx,

		immutableValues: c.immutableValues,
		mutableValues:   c.mutableValues.clone(),
		renderers:       c.renderers,
	}
}

func (c *contextRendering) Env() map[string]string {
	for k, v := range c.immutableValues.globalEnv {
		c.mutableValues.env[k] = v
	}

	return c.mutableValues.env
}

func (c *contextRendering) RenderYaml(renderer string, rawData interface{}) ([]byte, error) {
	v, ok := c.renderers[renderer]
	if !ok {
		return nil, fmt.Errorf("renderer %q not found", renderer)
	}

	return v.RenderYaml(c, rawData)
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

// Get retrieves a variable by its name. To check if the variable is
// set, use Variable.IsSet.
//
// for expand.Environ
func (c *contextRendering) Get(name string) expand.Variable {
	v, ok := c.Env()[name]
	return c.createVariable(name, v, ok)
}

// Each iterates over all the currently set variables, calling the
// supplied function on each variable. Iteration is stopped if the
// function returns false.
//
// The names used in the calls aren't required to be unique or sorted.
// If a variable name appears twice, the latest occurrence takes
// priority.
//
// Each is required to forward exported variables when executing
// programs.
//
// for expand.Environ
func (c *contextRendering) Each(do func(name string, vr expand.Variable) bool) {
	for k, v := range c.Env() {
		if !do(k, c.createVariable(k, v, true)) {
			return
		}
	}
}

func (c *contextRendering) createVariable(name, value string, eixists bool) expand.Variable {
	// TODO: set kind for lists
	kind := expand.String
	if !eixists {
		switch name {
		case "IFS":
			value = " \t\n"
		case "OPTIND":
			value = "1"
		case "PWD":
			value = c.WorkingDir()
		case "UID":
			// os.Getenv("UID") usually retruns empty value
			// so we have to call os.Getuid
			value = strconv.FormatInt(int64(os.Getuid()), 10)
		default:
			kind = expand.Unset
		}
	}

	return expand.Variable{
		Local:    false,
		Exported: true,
		ReadOnly: false,
		Kind:     kind,
		Str:      value,
	}
}
