package testhelper

import (
	"arhat.dev/pkg/errhelper"
)

const (
	errTest errhelper.ErrString = "test error"
)

// Error returns an internal error for testing
func Error() error { return errTest }
