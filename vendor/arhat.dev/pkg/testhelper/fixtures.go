package testhelper

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// TestFixtures run tests using multi-doc yaml file
//
/*
# first doc is the test spec
some_args: ...
spec:
	a: b
---
# second doc is the expected result
c: d
*/
func TestFixtures(
	t *testing.T,
	dir string,
	createTestSpec func() interface{},
	createExpected func() interface{},
	check func(t *testing.T, spec interface{}, exp interface{}),
) {
	err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return err
		}

		t.Run(d.Name(), func(t *testing.T) {
			specBytes, err := os.ReadFile(filepath.Join(dir, path))
			if !assert.NoError(t, err) {
				return
			}

			dec := yaml.NewDecoder(bytes.NewReader(specBytes))

			spec := createTestSpec()
			if !assert.NoError(t, dec.Decode(spec)) {
				return
			}

			exp := createExpected()
			if !assert.NoError(t, dec.Decode(exp)) {
				return
			}

			check(t, spec, exp)
		})
		return nil
	})

	assert.NoError(t, err)
}
