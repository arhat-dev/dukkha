package transform

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/tmpl"
)

var _ dukkha.Renderer = (*Driver)(nil)

func TestDriver_RenderYaml(t *testing.T) {
	t.Parallel()

	type TestSpec struct {
		rs.BaseField

		Data string `yaml:"data"`
	}

	type CheckSpec struct {
		rs.BaseField

		Data string `yaml:"data"`
	}

	testhelper.TestFixtures(t, "./fixtures",
		func() any { return rs.InitAny(&TestSpec{}, nil) },
		func() any { return rs.InitAny(&CheckSpec{}, nil) },
		func(t *testing.T, spec, exp any) {
			actual := spec.(*TestSpec)
			expected := exp.(*CheckSpec)

			ctx := dt.NewTestContext(context.TODO(), t.TempDir())

			ctx.AddRenderer("T", NewDefault("T"))
			ctx.AddRenderer("file", file.NewDefault("file"))
			ctx.AddRenderer("tmpl", tmpl.NewDefault("tmpl"))

			assert.NoError(t, actual.ResolveFields(ctx, -1))
			assert.NoError(t, expected.ResolveFields(ctx, -1))

			assert.EqualValues(t, expected.Data, actual.Data)
		},
	)
}
