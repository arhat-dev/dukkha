package renderer

import (
	"context"
	"fmt"
	"sync"
)

func WithManager(ctx context.Context, mgr *Manager) context.Context {
	return context.WithValue(ctx, contextKeyManager, mgr)
}

func GetManager(ctx context.Context) *Manager {
	mgr, ok := ctx.Value(contextKeyManager).(*Manager)
	if ok {
		return mgr
	}

	return nil
}

func NewManager() *Manager {
	return &Manager{
		renderers: &sync.Map{},

		values: &RenderingValues{
			Env: make(map[string]string),
		},
	}
}

type Manager struct {
	renderers *sync.Map

	values *RenderingValues
}

func (m *Manager) UpdateEnv(key, value string) {
	m.values.Env[key] = value
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

func (m *Manager) Render(ctx context.Context, renderer, rawData string) (string, error) {
	r := m.getRenderer(renderer)
	if r == nil {
		return "", fmt.Errorf("renderer.Manager.Render: renderer %q not found", renderer)
	}

	return r.Render(ctx, rawData, m.values)
}

func (m *Manager) getRenderer(rendererName string) Interface {
	r, ok := m.renderers.Load(rendererName)
	if ok {
		return r.(Interface)
	}

	return nil
}
