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
func TestFixtures[TestCase, Expected any](
	t *testing.T,
	dir string,
	createTestSpec func() TestCase,
	createExpected func() Expected,
	check func(t *testing.T, spec TestCase, exp Expected),
) {
	err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return err
		}

		t.Run(d.Name(), func(t *testing.T) {
			var rd bytes.Reader
			specBytes, err := os.ReadFile(filepath.Join(dir, path))
			if !assert.NoError(t, err) {
				return
			}
			rd.Reset(specBytes)

			dec := yaml.NewDecoder(&rd)

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
