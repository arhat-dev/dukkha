package buildah

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/tests"
)

func TestTaskBuild(t *testing.T) {
	t.Parallel()

	type Check struct {
		rs.BaseField

		Data map[string]string `yaml:",inline"`
	}

	t.Skip()

	tests.TestTask(t, "./fixtures/build", &Tool{},
		func() *TaskBuild { return tools.NewTask[TaskBuild, *TaskBuild]("").(*TaskBuild) },
		func() *Check { return &Check{} },
		func(t *testing.T, expected, actual *Check) {
			// TODO: check images
		},
	)
}

func TestCreateManifestPlatformQueryForDigest(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		rs.BaseField

		Query struct {
			Kernel string `yaml:"kernel"`
			Arch   string `yaml:"arch"`
		} `yaml:"query"`
		Manifest string `yaml:"manifest"`
	}

	type CheckSpec struct {
		rs.BaseField

		ExpectErr bool     `yaml:"expect_err"`
		Digests   []string `yaml:"digests"`
	}

	testhelper.TestFixtures(t, "./fixtures/manifest-platform-query",
		func() *TestCase { return rs.Init(&TestCase{}, nil).(*TestCase) },
		func() *CheckSpec { return rs.Init(&CheckSpec{}, nil).(*CheckSpec) },
		func(t *testing.T, spec *TestCase, exp *CheckSpec) {
			ctx := dt.NewTestContext(context.TODO(), t.TempDir())
			ctx.AddRenderer("file", file.NewDefault("file"))

			assert.NoError(t, spec.ResolveFields(ctx, -1))
			assert.NoError(t, exp.ResolveFields(ctx, -1))

			query := createManifestPlatformQueryForDigest(spec.Query.Kernel, spec.Query.Arch)

			result, err := textquery.JQ[byte](query, []byte(spec.Manifest))
			if exp.ExpectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			digests, err := parseManifestOsArchVariantQueryResult(result)
			assert.NoError(t, err)

			if !assert.EqualValues(t, exp.Digests, digests) {
				t.Log("Query:", query)
				t.Log("Manifest:", spec.Manifest)
			}
		},
	)
}
