package dukkha_test

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
)

func TestFixturesUsingRenderingSuffix[TestCase, CheckSpec rs.Field](
	t *testing.T,
	dir string,
	renderers map[string]dukkha.Renderer,
	newTestSpec func() TestCase,
	newCheckSpec func() CheckSpec,
	check func(t *testing.T, ctx dukkha.Context, ts TestCase, cs CheckSpec),
) {
	testhelper.TestFixtures(t, dir,
		func() TestCase { return rs.Init(newTestSpec(), nil).(TestCase) },
		func() CheckSpec { return rs.Init(newCheckSpec(), nil).(CheckSpec) },
		func(t *testing.T, s TestCase, e CheckSpec) {
			defer t.Cleanup(func() {})

			ctx := NewTestContext(context.TODO(), t.TempDir())

			for k, v := range renderers {
				assert.NoError(t, v.Init(ctx.RendererCacheFS(k)))

				ctx.AddRenderer(k, v)
			}

			assert.NoError(t, s.ResolveFields(ctx, -1))
			assert.NoError(t, e.ResolveFields(ctx, -1))

			check(t, ctx, s, e)
		},
	)
}
