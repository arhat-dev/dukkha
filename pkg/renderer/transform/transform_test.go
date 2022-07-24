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
		func() *TestSpec { return rs.Init(&TestSpec{}, nil).(*TestSpec) },
		func() *CheckSpec { return rs.Init(&CheckSpec{}, nil).(*CheckSpec) },
		func(t *testing.T, spec *TestSpec, exp *CheckSpec) {
			ctx := dt.NewTestContext(context.TODO(), t.TempDir())

			ctx.AddRenderer("T", NewDefault("T"))
			ctx.AddRenderer("file", file.NewDefault("file"))
			ctx.AddRenderer("tmpl", tmpl.NewDefault("tmpl"))

			assert.NoError(t, spec.ResolveFields(ctx, -1))
			assert.NoError(t, exp.ResolveFields(ctx, -1))

			assert.EqualValues(t, exp.Data, spec.Data)
		},
	)
}
