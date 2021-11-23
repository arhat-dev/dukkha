package af

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func addTarFile(w *tar.Writer, name string, content []byte) error {
	err := w.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,

		Name:   name,
		Format: tar.FormatPAX,
		Size:   int64(len(content)),
	})

	if err != nil {
		return err
	}

	_, err = w.Write(content)
	return err
}

func addTarDir(w *tar.Writer, name string) error {
	return w.WriteHeader(&tar.Header{
		Typeflag: tar.TypeDir,
		Name:     name,
		Format:   tar.FormatPAX,
		Size:     0,
	})
}

func addTarSymLink(w *tar.Writer, name, target string) error {
	return w.WriteHeader(&tar.Header{
		Typeflag: tar.TypeSymlink,

		Name:     name,
		Format:   tar.FormatPAX,
		Linkname: target,
	})
}

func newTarBuf(add func(w *tar.Writer)) io.ReadSeeker {
	buf := &bytes.Buffer{}
	w := tar.NewWriter(buf)

	add(w)
	err := w.Close()
	if err != nil {
		panic(err)
	}

	return io.NewSectionReader(bytes.NewReader(buf.Bytes()), 0, int64(buf.Len()))
}

func TestUntar(t *testing.T) {
	const (
		fileContent = "test-data"
	)
	for _, test := range []struct {
		name      string
		tar       io.ReadSeeker
		target    string
		expectErr bool
	}{
		{
			name: "Only File - Match",
			tar: newTarBuf(func(w *tar.Writer) {
				assert.NoError(t, addTarFile(w, "foo", []byte(fileContent)))
			}),
			target:    "foo",
			expectErr: false,
		},
		{
			name: "Only File - Missing",
			tar: newTarBuf(func(w *tar.Writer) {
				assert.NoError(t, addTarFile(w, "foo", []byte(fileContent)))
			}),
			target:    "bar",
			expectErr: true,
		},
		{
			name: "Dir - Error",
			tar: newTarBuf(func(w *tar.Writer) {
				assert.NoError(t, addTarDir(w, "foo"))
			}),
			target:    "foo",
			expectErr: true,
		},
		{
			name: "Internal Symlink - Redirect",
			tar: newTarBuf(func(w *tar.Writer) {
				assert.NoError(t, addTarFile(w, "foo", []byte(fileContent)))
				assert.NoError(t, addTarSymLink(w, "link", "foo"))
			}),
			target:    "link",
			expectErr: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r, err := untar(test.tar, test.target)
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			data, err := io.ReadAll(r)
			assert.NoError(t, err)
			assert.EqualValues(t, fileContent, string(data))
		})
	}
}
