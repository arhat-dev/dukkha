package renderer

import (
	"fmt"
	"sync"

	"arhat.dev/dukkha/pkg/field"
)

func NewManager() *Manager {
	return &Manager{}
}

type Manager struct {
	renderers sync.Map
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

func (m *Manager) Render(ctx *field.RenderingContext, renderer, rawData string) (string, error) {
	r := m.getRenderer(renderer)
	if r == nil {
		return "", fmt.Errorf("renderer.Manager.Render: renderer %q not found", renderer)
	}

	return r.Render(ctx, rawData)
}

func (m *Manager) getRenderer(rendererName string) Interface {
	r, ok := m.renderers.Load(rendererName)
	if ok {
		return r.(Interface)
	}

	return nil
}
