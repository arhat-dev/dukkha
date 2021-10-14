package rshelper

import (
	"fmt"
	"text/template"

	"arhat.dev/rs"

	"arhat.dev/pkg/envhelper"
)

// DefaultRenderingManager creates a RenderingManager with env, file rendering handler
func DefaultRenderingManager(env map[string]string, funcMap template.FuncMap) *RenderingManager {
	m := &RenderingManager{
		m: make(map[string]rs.RenderingHandler),
	}

	m.Add(&envhelper.EnvRenderingHandler{Env: env, AllowNotFound: false}, "env", "env-strict")
	m.Add(&envhelper.EnvRenderingHandler{Env: env, AllowNotFound: true}, "env-loose")

	m.Add(&FileHandler{}, "file")

	m.Add(&TemplateHandler{CreateFuncMap: func() template.FuncMap {
		return funcMap
	}}, "template")

	return m
}

var _ rs.RenderingHandler = (*RenderingManager)(nil)

// RenderingManager is a collection of named rendering handlers
type RenderingManager struct {
	m map[string]rs.RenderingHandler
}

func (r *RenderingManager) Add(h rs.RenderingHandler, names ...string) {
	if r.m == nil {
		r.m = make(map[string]rs.RenderingHandler)
	}

	for _, name := range names {
		r.m[name] = h
	}
}

func (r *RenderingManager) RenderYaml(
	name string, rawData interface{},
) (interface{}, error) {
	h, ok := r.m[name]
	if !ok {
		return nil, fmt.Errorf("rendering handler %q not found", name)
	}

	ret, err := h.RenderYaml(name, rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to render raw data with %q: %w", name, err)
	}

	return ret, nil
}
