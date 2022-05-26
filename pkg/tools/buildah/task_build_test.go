package buildah

import (
	"context"
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/pkg/textquery"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/file"
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
		func() dukkha.Task { return newTaskBuild("") },
		func() rs.Field { return &Check{} },
		func(t *testing.T, expected, actual rs.Field) {
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
		func() interface{} { return rs.Init(&TestCase{}, nil) },
		func() interface{} { return rs.Init(&CheckSpec{}, nil) },
		func(t *testing.T, spec, exp interface{}) {
			test, check := spec.(*TestCase), exp.(*CheckSpec)

			ctx := dt.NewTestContext(context.TODO())
			ctx.(di.CacheDirSetter).SetCacheDir(t.TempDir())
			ctx.AddRenderer("file", file.NewDefault("file"))

			assert.NoError(t, test.ResolveFields(ctx, -1))
			assert.NoError(t, check.ResolveFields(ctx, -1))

			query := createManifestPlatformQueryForDigest(test.Query.Kernel, test.Query.Arch)

			result, err := textquery.JQ[byte](query, []byte(test.Manifest))
			if check.ExpectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			digests, err := parseManifestOsArchVariantQueryResult(result)
			assert.NoError(t, err)

			if !assert.EqualValues(t, check.Digests, digests) {
				t.Log("Query:", query)
				t.Log("Manifest:", test.Manifest)
			}
		},
	)
}
