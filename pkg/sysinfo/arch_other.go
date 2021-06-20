//go:build not (windows || darwin || linux || freebsd || openbsd || netbsd || dragonfly || solaris || aix)
// +build !darwin,!linux,!freebsd,!openbsd,!netbsd,!dragonfly,!solaris,!aix,!windows

package sysinfo

import (
	"arhat.dev/dukkha/pkg/version"
)

func Arch() string {
	return version.Arch()
}
