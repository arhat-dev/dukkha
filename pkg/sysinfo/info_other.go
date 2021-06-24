//go:build not (windows || darwin || linux || freebsd || openbsd || netbsd || dragonfly || solaris || aix)
// +build !darwin,!linux,!freebsd,!openbsd,!netbsd,!dragonfly,!solaris,!aix,!windows

package sysinfo

import (
	"arhat.dev/dukkha/pkg/version"
)

func Arch() string {
	return version.Arch()
}

func OSName() string {
	return ""
}

func OSVersion() string {
	// TODO: check os version using syscall
	return ""
}

func KernelVersion() string {
	return ""
}
