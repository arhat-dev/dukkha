package testhelper

import "errors"

var (
	errTest = errors.New("test error")
)

// Error returns an internal error for testing
func Error() error { return errTest }
