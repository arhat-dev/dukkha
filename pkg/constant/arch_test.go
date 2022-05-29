package constant

import (
	"fmt"
	"testing"

	. "arhat.dev/pkg/archconst"
	"github.com/stretchr/testify/assert"
)

func TestSimpleArch(t *testing.T) {
	for _, test := range []struct {
		arch     string
		expected string
	}{
		{ARCH_AMD64, ARCH_AMD64},
		{ARCH_AMD64_V1, ARCH_AMD64},
		{ARCH_AMD64_V2, ARCH_AMD64},
		{ARCH_AMD64_V3, ARCH_AMD64},
		{ARCH_AMD64_V4, ARCH_AMD64},

		{ARCH_ARM64, ARCH_ARM64},
		{ARCH_ARM64_V8, ARCH_ARM64},
		{ARCH_ARM64_V9, ARCH_ARM64},

		{ARCH_ARM, ARCH_ARM_V7},
		{ARCH_ARM_V5, ARCH_ARM_V5},
		{ARCH_ARM_V6, ARCH_ARM_V6},
		{ARCH_ARM_V7, ARCH_ARM_V7},

		{ARCH_PPC64, ARCH_PPC64},
		{ARCH_PPC64_V8, ARCH_PPC64},
		{ARCH_PPC64_V9, ARCH_PPC64},

		{ARCH_PPC64_LE, ARCH_PPC64_LE},
		{ARCH_PPC64_LE_V8, ARCH_PPC64_LE},
		{ARCH_PPC64_LE_V9, ARCH_PPC64_LE},
	} {
		t.Run(test.arch, func(t *testing.T) {
			assert.EqualValues(t, test.expected, SimpleArch(test.arch))
		})
	}
}

func TestCrossPlatform(t *testing.T) {
	for _, test := range []struct {
		hostKernel, hostArch, targetKernel, targetArch string

		ret bool
	}{
		// same kernel, same arch name => true
		{KERNEL_Linux, ARCH_AMD64, KERNEL_Linux, ARCH_AMD64_V2, false},
		{KERNEL_Linux, ARCH_AMD64_V1, KERNEL_Linux, ARCH_AMD64_V2, false},
		{KERNEL_Linux, ARCH_AMD64_V2, KERNEL_Linux, ARCH_AMD64_V2, false},
		{KERNEL_Linux, ARCH_AMD64_V3, KERNEL_Linux, ARCH_AMD64_V2, false},
		{KERNEL_Linux, ARCH_MIPS64, KERNEL_Linux, ARCH_MIPS64_SF, false},
		{KERNEL_Linux, ARCH_PPC64_LE, KERNEL_Linux, ARCH_PPC64_LE_V9, false},

		// different kernel => true
		{KERNEL_Linux, ARCH_AMD64_V2, KERNEL_FreeBSD, ARCH_AMD64_V2, true},

		// different arch => true
		{KERNEL_Linux, ARCH_AMD64_V2, KERNEL_FreeBSD, ARCH_ARM64, true},

		// different endian => true
		{KERNEL_Linux, ARCH_MIPS64, KERNEL_FreeBSD, ARCH_MIPS64_LE, true},
	} {
		name := fmt.Sprintf("%s/%s-%s/%s", test.hostKernel, test.hostArch, test.targetKernel, test.targetArch)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.ret, CrossPlatform(test.targetKernel, test.targetArch, test.hostKernel, test.hostArch))
		})
	}
}

func TestArchMapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		mappingFunc func(mArch string) (string, bool)
	}{
		{
			name:        "AlpineArch",
			mappingFunc: GetAlpineArch,
		},
		{
			name:        "AlpineTripleName",
			mappingFunc: GetAlpineTripleName,
		},

		{
			name:        "AppleArch",
			mappingFunc: GetAppleArch,
		},

		{
			name:        "DebianArch",
			mappingFunc: GetDebianArch,
		},
		{
			name: "DebianTripleName GNU",
			mappingFunc: func(mArch string) (string, bool) {
				return GetDebianTripleName(mArch, "", LIBC_GNU)
			},
		},
		{
			name: "DebianTripleName MUSL",
			mappingFunc: func(mArch string) (string, bool) {
				return GetDebianTripleName(mArch, "", LIBC_MUSL)
			},
		},
		{
			name: "DebianTripleName MSVC",
			mappingFunc: func(mArch string) (string, bool) {
				return GetDebianTripleName(mArch, "", LIBC_MSVC)
			},
		},

		{
			name:        "GNUArch",
			mappingFunc: GetGNUArch,
		},
		{
			name: "GNUTripleName GNU",
			mappingFunc: func(mArch string) (string, bool) {
				return GetGNUTripleName(mArch, "", LIBC_GNU)
			},
		},
		{
			name: "GNUTripleName MUSL",
			mappingFunc: func(mArch string) (string, bool) {
				return GetGNUTripleName(mArch, "", LIBC_MUSL)
			},
		},
		{
			name: "GNUTripleName MSVC",
			mappingFunc: func(mArch string) (string, bool) {
				return GetGNUTripleName(mArch, "", LIBC_MSVC)
			},
		},
		{
			name: "GetLLVMTripleName GNU",
			mappingFunc: func(mArch string) (string, bool) {
				return GetLLVMTripleName(mArch, "", LIBC_GNU)
			},
		},
		{
			name: "GetLLVMTripleName MUSL",
			mappingFunc: func(mArch string) (string, bool) {
				return GetLLVMTripleName(mArch, "", LIBC_MUSL)
			},
		},
		{
			name: "GetLLVMTripleName MSVC",
			mappingFunc: func(mArch string) (string, bool) {
				return GetLLVMTripleName(mArch, "", LIBC_MSVC)
			},
		},

		{
			name:        "GolangArch",
			mappingFunc: GetGolangArch,
		},

		{
			name:        "DockerArch",
			mappingFunc: GetDockerArch,
		},
		{
			name:        "DockerArchVariant",
			mappingFunc: GetDockerArchVariant,
		},
		{
			name: "DockerHubArch linux",
			mappingFunc: func(mArch string) (string, bool) {
				return GetDockerHubArch(mArch, "linux")
			},
		},
		{
			name: "DockerHubArch windows",
			mappingFunc: func(mArch string) (string, bool) {
				return GetDockerHubArch(mArch, "windows")
			},
		},
		{
			name:        "OciArch",
			mappingFunc: GetOciArch,
		},
		{
			name:        "OciArchVariant",
			mappingFunc: GetOciArchVariant,
		},
		{
			name:        "QemuArch",
			mappingFunc: GetQemuArch,
		},
		{
			name:        "LLVMArch",
			mappingFunc: GetLLVMArch,
		},
		{
			name:        "ZigArch",
			mappingFunc: GetZigArch,
		},
		{
			name:        "RustArch",
			mappingFunc: GetRustArch,
		},
		{
			name: "DockerHubArch",
			mappingFunc: func(mArch string) (string, bool) {
				return GetDockerHubArch(mArch, "linux")
			},
		},
		{
			name: "AppleArch",
			mappingFunc: func(mArch string) (string, bool) {
				return GetAppleTripleName(mArch, "")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for aid := _unknown_arch + 1; aid < archID_COUNT; aid++ {
				_, ok := test.mappingFunc(aid.String())
				assert.True(t, ok, aid.String())
			}
		})
	}
}
