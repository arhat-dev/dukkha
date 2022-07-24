package tlang

import (
	"testing"

	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
)

var _ dukkha.Renderer = (*Driver)(nil)

func TestNewDriver(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, NewDefault(""))
}

func TestDriver_RenderYaml(t *testing.T) {
	t.Parallel()

	dt.TestFixturesUsingRenderingSuffix(t, "./fixtures",
		map[string]dukkha.Renderer{
			"tl": NewDefault("tl"),
		},
		func() *rs.AnyObjectMap { return &rs.AnyObjectMap{} },
		func() *rs.AnyObjectMap { return &rs.AnyObjectMap{} },
		func(t *testing.T, ctx dukkha.Context, spec *rs.AnyObjectMap, exp *rs.AnyObjectMap) {
			assert.EqualValues(t, exp.NormalizedValue(), spec.NormalizedValue())
		},
	)
}
