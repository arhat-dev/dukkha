package tlang

import (
	"testing"

	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
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
		func() rs.Field { return &rs.AnyObjectMap{} },
		func() rs.Field { return &rs.AnyObjectMap{} },
		func(t *testing.T, ctx dukkha.Context, ts, cs rs.Field) {
			actual, expected := ts.(*rs.AnyObjectMap), cs.(*rs.AnyObjectMap)

			assert.EqualValues(t, expected.NormalizedValue(), actual.NormalizedValue())
		},
	)
}
