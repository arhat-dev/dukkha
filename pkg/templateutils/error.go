package templateutils

import "arhat.dev/pkg/stringhelper"

type errString string

func (s errString) Error() string { return stringhelper.Convert[string, byte](s) }

const (
	errAtLeastOneArgGotZero errString = "at least 1 arg expected, got 0"
)
