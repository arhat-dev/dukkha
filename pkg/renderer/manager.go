package renderer

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

func WithManager(ctx context.Context, mgr *Manager) context.Context {
	return context.WithValue(ctx, contextKeyManager, mgr)
}

func GetManager(ctx context.Context, name string) *Manager {
	mgr, ok := ctx.Value(contextKeyManager).(*Manager)
	if ok {
		return mgr
	}

	return nil
}

func NewManager() *Manager {
	return &Manager{
		renderers: &sync.Map{},
	}
}

type Manager struct {
	renderers *sync.Map
}

func (m *Manager) Add(config Config, names ...string) error {
	renderer, err := Create(config)
	if err != nil {
		return fmt.Errorf("renderer.Manger.Add: unable to create renderer for config %T: %w", config, err)
	}

	for _, name := range names {
		// TODO: create renderers according to config
		m.renderers.Store(name, renderer)
	}

	return nil
}

func (m *Manager) Render(ctx context.Context, fieldName string, fieldValue, out interface{}) error {
	return m.render(ctx, fieldName, fieldValue, reflect.ValueOf(out))
}

func (m *Manager) getRenderer(rendererName string) Interface {
	r, ok := m.renderers.Load(rendererName)
	if ok {
		return r.(Interface)
	}

	return nil
}

func (m *Manager) render(ctx context.Context, fieldName string, fieldValue interface{}, out reflect.Value) error {
	parts := strings.SplitN(fieldName, "@", 2)
	if len(parts) == 2 {
		// has renderer, expecting string as value
		fieldName, rendererName := parts[0], parts[1]
		renderer := m.getRenderer(rendererName)
		if renderer == nil {
			return fmt.Errorf("renderer.Manager.render: requested renderer %q not found", rendererName)
		}

		switch t := fieldValue.(type) {
		case string:
			result, err := renderer.Render(ctx, t)
			if err != nil {
				return fmt.Errorf("renderer.Manager.render: rendering failed: %w", err)
			}

			_ = result
		default:
			return fmt.Errorf("renderer.Manager.render: unexpected non string value of %q for rendering", fieldName)
		}
	}

	// no renderer, unmarshal directly
	switch out.Kind() {
	case reflect.Struct:
	case reflect.Array:
	case reflect.Slice:
	case reflect.Bool:
	case reflect.Chan, reflect.Func:
		// invalid values
	}
	// switch t := fieldValue.(type) {
	// case map[string]interface{}:

	// }

	return nil

}
