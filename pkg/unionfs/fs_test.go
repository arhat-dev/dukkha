package unionfs

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestRemappedFS_Open(t *testing.T) {
	tests := []struct {
		name string
		path string

		expectedName string
		expectDir    bool
		expectErr    error
	}{
		{
			name:         "File Exists",
			path:         "foo",
			expectedName: "foo",
		},
		{
			name:      "File Not Exists",
			path:      "bar",
			expectErr: fs.ErrNotExist,
		},
		{
			name:         "Path Is Dir",
			path:         "some",
			expectedName: "some",
			expectDir:    true,
		},
	}

	rmfs := New()
	rmfs.Remap("foo", "bar", fstest.MapFS{
		"bar": &fstest.MapFile{
			Data: []byte("bar"),
		},
	})

	rmfs.Remap("some/data", "some/data", fstest.MapFS{
		"data": &fstest.MapFile{
			Data: []byte("data"),
		},
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := rmfs.Open(test.path)
			if test.expectErr != nil {
				assert.ErrorIs(t, err, test.expectErr)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			fInfo, err := f.Stat()
			assert.Equal(t, test.expectedName, fInfo.Name())
			assert.Equal(t, test.expectDir, fInfo.IsDir())
		})
	}
}
