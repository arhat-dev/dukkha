package dukkha_test

import (
	"context"
	"testing"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
)

func TestFixturesUsingRenderingSuffix(
	t *testing.T,
	dir string,
	renderers map[string]dukkha.Renderer,
	newTestSpec func() rs.Field,
	newCheckSpec func() rs.Field,
	check func(t *testing.T, ts, cs rs.Field),
) {
	testhelper.TestFixtures(t, dir,
		func() interface{} { return rs.Init(newTestSpec(), nil) },
		func() interface{} { return rs.Init(newCheckSpec(), nil) },
		func(t *testing.T, spec, exp interface{}) {
			defer t.Cleanup(func() {})
			s, e := spec.(rs.Field), exp.(rs.Field)

			ctx := NewTestContext(context.TODO())
			ctx.(di.CacheDirSetter).SetCacheDir(t.TempDir())

			for k, v := range renderers {
				assert.NoError(t, v.Init(ctx))

				ctx.AddRenderer(k, v)
			}

			assert.NoError(t, s.ResolveFields(ctx, -1))
			assert.NoError(t, e.ResolveFields(ctx, -1))

			check(t, s, e)
		},
	)
}
