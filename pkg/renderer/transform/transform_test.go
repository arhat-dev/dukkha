package transform

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/tpl"
)

var _ dukkha.Renderer = (*Driver)(nil)

func TestDriver_RenderYaml(t *testing.T) {
	type TestSpec struct {
		rs.BaseField

		Data string `yaml:"data"`
	}

	type CheckSpec struct {
		rs.BaseField

		Data string `yaml:"data"`
	}

	testhelper.TestFixtures(t, "./fixtures",
		func() interface{} { return rs.Init(&TestSpec{}, nil) },
		func() interface{} { return rs.Init(&CheckSpec{}, nil) },
		func(t *testing.T, spec, exp interface{}) {
			actual := spec.(*TestSpec)
			expected := exp.(*CheckSpec)

			ctx := dt.NewTestContext(context.TODO())
			ctx.(di.CacheDirSetter).SetCacheDir(t.TempDir())

			ctx.AddRenderer("T", NewDefault("T"))
			ctx.AddRenderer("file", file.NewDefault("file"))
			ctx.AddRenderer("tpl", tpl.NewDefault("tpl"))

			assert.NoError(t, actual.ResolveFields(ctx, -1))
			assert.NoError(t, expected.ResolveFields(ctx, -1))

			assert.EqualValues(t, expected.Data, actual.Data)
		},
	)
}
