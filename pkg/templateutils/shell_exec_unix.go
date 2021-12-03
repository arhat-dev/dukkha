//go:build !windows
// +build !windows

package templateutils

import (
	"arhat.dev/pkg/pathhelper"
	"mvdan.cc/sh/v3/expand"
)

func isSlash(c byte) bool { return pathhelper.IsUnixSlash(c) }

func pathExts(env expand.Environ) []string { return nil }
