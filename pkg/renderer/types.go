package renderer

import "context"

type contextKey string

const (
	contextKeyManager contextKey = "renderer.manager"
)

type Interface interface {
	Name() string

	Render(ctx context.Context, rawValue string, v *RenderingValues) (result string, err error)
}

type Config interface{}

type RenderingValues struct {
	// resolved environment variables
	Env map[string]string
}
