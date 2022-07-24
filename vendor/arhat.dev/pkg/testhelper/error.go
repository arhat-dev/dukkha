package testhelper

import (
	"arhat.dev/pkg/errhelper"
)

var (
	errTest errhelper.ErrString = "test error"
)

// Error returns an internal error for testing
func Error() error { return errTest }
