//go:build !windows
// +build !windows

package templateutils

import (
	"mvdan.cc/sh/v3/expand"
)

func pathExts(env expand.Environ) []string { return nil }
