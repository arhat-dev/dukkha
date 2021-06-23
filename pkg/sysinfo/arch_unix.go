// +build darwin linux freebsd openbsd netbsd dragonfly solaris aix
//go:build darwin || linux || freebsd || openbsd || netbsd || dragonfly || solaris || aix

package sysinfo

import (
	"bytes"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/version"
)

var littleEndian bool

func init() {
	var x uint32 = 0x01020304
	littleEndian = *(*byte)(unsafe.Pointer(&x)) == 0x04
}

func Arch() string {
	var uname unix.Utsname
	_ = unix.Uname(&uname)

	buf := make([]byte, len(uname.Machine))
	for i, b := range uname.Machine {
		// nolint:unconvert
		buf[i] = byte(b)
	}

	if i := bytes.Index(buf, []byte{'\x00'}); i != -1 {
		buf = buf[:i]
	}

	// https://en.wikipedia.org/wiki/Uname
	rawArch := string(buf)
	switch {
	case rawArch == "x86_64",
		rawArch == "i686-64":
		return constant.ARCH_AMD64
	case strings.HasPrefix(rawArch, "armv8"),
		strings.HasPrefix(rawArch, "arm64"),
		strings.HasPrefix(rawArch, "aarch64"):
		return constant.ARCH_ARM64
	case strings.HasPrefix(rawArch, "armv7"):
		return constant.ARCH_ARM_V7
	case strings.HasPrefix(rawArch, "armv6"):
		return constant.ARCH_ARM_V6
	case strings.HasPrefix(rawArch, "armv5"):
		return constant.ARCH_ARM_V5
	case strings.HasPrefix(rawArch, "i686"),
		strings.HasPrefix(rawArch, "i386"),
		rawArch == "i86pc",
		rawArch == "x86pc",
		rawArch == "x86":
		return constant.ARCH_X86
	case rawArch == "ppc64":
		if littleEndian {
			return constant.ARCH_PPC64LE
		}

		return constant.ARCH_PPC64
	case rawArch == "ppc":
		return constant.ARCH_PPC64
	case rawArch == "mips64":
		if littleEndian {
			return constant.ARCH_MIPS64_LE
		}

		return constant.ARCH_MIPS64
	default:
		// uncertain, use build arch
		return version.Arch()
	}
}
