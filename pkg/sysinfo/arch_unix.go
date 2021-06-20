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
		return constant.ArchAMD64
	case strings.HasPrefix(rawArch, "armv8"),
		strings.HasPrefix(rawArch, "arm64"),
		strings.HasPrefix(rawArch, "aarch64"):
		return constant.ArchARM64
	case strings.HasPrefix(rawArch, "armv7"):
		return constant.ArchARMv7
	case strings.HasPrefix(rawArch, "armv6"):
		return constant.ArchARMv6
	case strings.HasPrefix(rawArch, "armv5"):
		return constant.ArchARMv5
	case strings.HasPrefix(rawArch, "i686"),
		strings.HasPrefix(rawArch, "i386"),
		rawArch == "i86pc",
		rawArch == "x86pc",
		rawArch == "x86":
		return constant.ArchX86
	case rawArch == "ppc64":
		if littleEndian {
			return constant.ArchPPC64LE
		}

		return constant.ArchPPC64
	case rawArch == "ppc":
		return constant.ArchPPC64
	case rawArch == "mips64":
		if littleEndian {
			return constant.ArchMIPS64LE
		}
		return constant.ArchMIPS64
	default:
		// uncertain, use build arch
		return version.Arch()
	}
}
