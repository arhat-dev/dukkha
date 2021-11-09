package testhelper

import (
	"io"
)

// NewAlwaysFailReader creates an reader already fail on Read() call
// with provided err
//
// if err is nil, an io.ErrUnexpectedEOF is returned by default
func NewAlwaysFailReader(err error) io.Reader {
	if err == nil {
		err = io.ErrUnexpectedEOF
	}

	return &alwaysFailReader{
		err: err,
	}
}

type alwaysFailReader struct{ err error }

func (r *alwaysFailReader) Read(p []byte) (int, error)  { return 0, r.err }
func (r *alwaysFailReader) Write(p []byte) (int, error) { return 0, r.err }
