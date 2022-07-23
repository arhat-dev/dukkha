package dukkha_test

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
)

func TestFixturesUsingRenderingSuffix(
	t *testing.T,
	dir string,
	renderers map[string]dukkha.Renderer,
	newTestSpec func() rs.Field,
	newCheckSpec func() rs.Field,
	check func(t *testing.T, ctx dukkha.Context, ts, cs rs.Field),
) {
	testhelper.TestFixtures(t, dir,
		func() any { return rs.InitAny(newTestSpec(), nil) },
		func() any { return rs.InitAny(newCheckSpec(), nil) },
		func(t *testing.T, spec, exp any) {
			defer t.Cleanup(func() {})
			s, e := spec.(rs.Field), exp.(rs.Field)

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
