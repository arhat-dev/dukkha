//go:build !windows && !js && !illumos && !ios && !plan9
// +build !windows,!js,!illumos,!ios,!plan9

package sysinfo

import (
	"runtime"
	"strings"

	"golang.org/x/sys/unix"

	"arhat.dev/dukkha/pkg/constant"
)

func OSName() string {
	// TODO: check real name using syscall
	switch runtime.GOOS {
	case constant.KERNEL_WINDOWS:
		return "windows"
	case constant.KERNEL_DARWIN:
		return "macos"
	default:
		return runtime.GOOS
	}
}

func OSVersion() string {
	// TODO: check os version using syscall
	return ""
}

func KernelVersion() string {
	var uname unix.Utsname
	_ = unix.Uname(&uname)

	buf := make([]byte, len(uname.Release))
	for i, b := range uname.Release {
		// nolint:unconvert
		buf[i] = byte(b)
	}

	kernelVersion := string(buf)
	if i := strings.Index(kernelVersion, "\x00"); i != -1 {
		kernelVersion = kernelVersion[:i]
	}

	return kernelVersion
}
