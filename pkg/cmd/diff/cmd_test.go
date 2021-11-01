package diff

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
			err = diffFile(
				filepath.Join(testdataDir, path, "base.yaml"),
				filepath.Join(testdataDir, path, "new.yaml"),
			)

			assert.NoError(t, err)
		})

		return fs.SkipDir
	})
}
