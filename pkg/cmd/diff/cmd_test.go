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
		func() rs.Field { return &TestCase{} },
		func() rs.Field { return &CheckSpec{} },
		func(t *testing.T, ctx dukkha.Context, ts, cs rs.Field) {
			srcDoc, baseDoc, newDoc := filepath.Join(t.TempDir(), "src.yaml"),
				filepath.Join(t.TempDir(), "base.yaml"),
				filepath.Join(t.TempDir(), "new.yaml")

			test, check := ts.(*TestCase), cs.(*CheckSpec)
			assert.NoError(t, os.WriteFile(srcDoc, []byte(test.Src), 0644))
			assert.NoError(t, os.WriteFile(baseDoc, []byte(test.Base), 0644))
			assert.NoError(t, os.WriteFile(newDoc, []byte(test.New), 0644))

			err := diffFile(ctx, srcDoc, baseDoc, newDoc)
			if check.ExpectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		},
	)
}
