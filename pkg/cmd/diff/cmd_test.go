package diff

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
)

func TestDiffFile(t *testing.T) {
	const testdataDir = "./testdata"
	fs.WalkDir(os.DirFS(testdataDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "." {
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		t.Run(d.Name(), func(t *testing.T) {
			rc := dukkha_test.NewTestContext(context.TODO())
			rc.AddRenderer("file", file.NewDefault(""))
			rc.AddRenderer("env", env.NewDefault(""))

			err = diffFile(rc,
				filepath.Join(testdataDir, path, "src.yaml"),
				filepath.Join(testdataDir, path, "base.yaml"),
				filepath.Join(testdataDir, path, "new.yaml"),
			)

			assert.NoError(t, err)
		})

		return fs.SkipDir
	})
}
