//go:build !windows && !js && !illumos && !ios && !plan9
// +build !windows,!js,!illumos,!ios,!plan9

package cpuhelper

import (
	"bytes"
	"strings"

	"arhat.dev/pkg/archconst"
	"arhat.dev/pkg/versionhelper"
	"golang.org/x/sys/unix"
)

// Arch returns runtime cpu arch value as defined in package arhat.dev/pkg/archconst
func Arch(cpu CPU) string {
	hostArch := ArchByCPUFeatures(cpu)
	if len(hostArch) != 0 {
		return hostArch
	}

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

	// ref: https://en.wikipedia.org/wiki/Uname

	hostArch = string(buf)
	switch {
	case hostArch == "x86_64", hostArch == "i686-64":
		return archconst.ARCH_AMD64
	case strings.HasPrefix(hostArch, "armv8"):
		if Bits() == 64 {
			return archconst.ARCH_ARM64
		}

		return versionhelper.Arch()
	case strings.HasPrefix(hostArch, "arm64"),
		strings.HasPrefix(hostArch, "aarch64"):
		return archconst.ARCH_ARM64
	case strings.HasPrefix(hostArch, "armv7"):
		return archconst.ARCH_ARM_V7
	case strings.HasPrefix(hostArch, "armv6"):
		return archconst.ARCH_ARM_V6
	case strings.HasPrefix(hostArch, "armv5"):
		return archconst.ARCH_ARM_V5
	case strings.HasPrefix(hostArch, "i686"),
		strings.HasPrefix(hostArch, "i386"),
		hostArch == "i86pc",
		hostArch == "x86pc",
		hostArch == "x86":
		return archconst.ARCH_X86
	case hostArch == "ppc64", hostArch == "powerpc64":
		return fallback(
			versionhelper.Arch(), // prefer build time value with micro arch info
			selectEndian(archconst.ARCH_PPC64_LE_V8, archconst.ARCH_PPC64_V8),
		)
	case hostArch == "ppc", hostArch == "powerpc":
		return fallback(
			versionhelper.Arch(), // prefer build time value with softfloat info
			selectEndian(archconst.ARCH_PPC_LE_SF, archconst.ARCH_PPC_SF),
		)
	case hostArch == "mips64":
		return selectEndian(archconst.ARCH_MIPS64_LE, archconst.ARCH_MIPS64)
	default:
		// uncertain, use build arch
		return versionhelper.Arch()
	}
}
