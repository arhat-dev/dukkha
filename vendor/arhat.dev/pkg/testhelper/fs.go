package testhelper

import "io/fs"

// NewAlwaysErrFS creates an fs.FS that always fail at Open() call
// with provided err
//
// if err is nil, fs.ErrInvalid is returned by default
func NewAlwaysErrFS(err error) fs.FS {
	if err == nil {
		err = fs.ErrInvalid
	}

	return &alwaysErrFS{err: err}
}

type alwaysErrFS struct{ err error }

func (fs *alwaysErrFS) Open(name string) (fs.File, error) { return nil, fs.err }
