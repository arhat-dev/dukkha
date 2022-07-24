package diff

import (
	"os"
	"path/filepath"
	"testing"

	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
)

func TestCmd(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		rs.BaseField

		Src  string `yaml:"src"`
		Base string `yaml:"base"`
		New  string `yaml:"new"`
	}

	type CheckSpec struct {
		rs.BaseField

		ExpectErr bool `yaml:"expect_err"`
	}

	dukkha_test.TestFixturesUsingRenderingSuffix(t, "./fixtures",
		map[string]dukkha.Renderer{
			"file": file.NewDefault("file"),
			"env":  env.NewDefault("env"),
		},
		func() *TestCase { return &TestCase{} },
		func() *CheckSpec { return &CheckSpec{} },
		func(t *testing.T, ctx dukkha.Context, spec *TestCase, exp *CheckSpec) {
			srcDoc, baseDoc, newDoc := filepath.Join(t.TempDir(), "src.yaml"),
				filepath.Join(t.TempDir(), "base.yaml"),
				filepath.Join(t.TempDir(), "new.yaml")

			assert.NoError(t, os.WriteFile(srcDoc, []byte(spec.Src), 0644))
			assert.NoError(t, os.WriteFile(baseDoc, []byte(spec.Base), 0644))
			assert.NoError(t, os.WriteFile(newDoc, []byte(spec.New), 0644))

			err := diffFile(ctx, srcDoc, baseDoc, newDoc)
			if exp.ExpectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		},
	)
}
