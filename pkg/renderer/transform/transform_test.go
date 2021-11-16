package transform

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/template"
)

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

			ctx := dukkha_test.NewTestContext(context.TODO())
			ctx.SetCacheDir(t.TempDir())

			ctx.AddRenderer("transform", NewDefault("transform"))
			ctx.AddRenderer("file", file.NewDefault("file"))
			ctx.AddRenderer("template", template.NewDefault("template"))

			assert.NoError(t, actual.ResolveFields(ctx, -1))
			assert.NoError(t, expected.ResolveFields(ctx, -1))

			assert.EqualValues(t, expected.Data, actual.Data)
		},
	)
}
