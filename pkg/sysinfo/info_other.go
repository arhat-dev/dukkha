//go:build !(windows || darwin || linux || freebsd || openbsd || netbsd || dragonfly || solaris || aix)
// +build !windows,!darwin,!linux,!freebsd,!openbsd,!netbsd,!dragonfly,!solaris,!aix

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
