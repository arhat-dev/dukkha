package conf_test

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/conf"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/file"

	_ "arhat.dev/dukkha/cmd/dukkha/addon"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	testhelper.TestFixtures(t, "./fixtures",
		func() interface{} { return conf.NewConfig() },
		func() interface{} { return conf.NewConfig() },
		func(t *testing.T, spec, exp interface{}) {
			actual := conf.NewConfig()
			assert.NoError(t, actual.Merge(spec.(*conf.Config)))
			expected := exp.(*conf.Config)

			ctx := dukkha_test.NewTestContext(context.TODO(), t.TempDir())
			ctx.AddRenderer("file", file.NewDefault("file"))

			assert.NoError(t, actual.Resolve(ctx, conf.ReadFlag_Full))
			assert.NoError(t, expected.Resolve(ctx, conf.ReadFlag_Full))

			for k, list := range expected.Tools.Tools {
				if !assert.Len(t, actual.Tools.Tools[k], len(list)) {
					continue
				}

				for i, v := range list {
					assert.EqualValues(t, v.Key(), actual.Tools.Tools[k][i].Key())
				}
			}

			assert.Len(t, actual.Tasks, len(expected.Tasks))
			for k, list := range expected.Tasks {
				if !assert.Len(t, actual.Tasks[k], len(list)) {
					continue
				}

				for i, v := range list {
					// assert.EqualValues(t, v.Key(), actual.Tasks[k][i].Key())
					assert.EqualValues(t, v.Key(), actual.Tasks[k][i].Key())
				}
			}
		},
	)
}
