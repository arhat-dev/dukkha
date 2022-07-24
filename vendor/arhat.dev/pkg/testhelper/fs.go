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

	return AlwaysErrFS{err: err}
}

type AlwaysErrFS struct{ err error }

func (fs AlwaysErrFS) Open(name string) (fs.File, error) { return nil, fs.err }
