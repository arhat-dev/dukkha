package constant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// TODO: update this map when changing arch.go
	requiredArchMappingValues = map[string]string{
		ARCH_X86: "x86",

		ARCH_AMD64: "amd64",

		ARCH_AMD64_V1: "amd64v1",
		ARCH_AMD64_V2: "amd64v2",
		ARCH_AMD64_V3: "amd64v3",
		ARCH_AMD64_V4: "amd64v4",

		ARCH_ARM64: "arm64",

		ARCH_ARM_V5: "armv5",
		ARCH_ARM_V6: "armv6",
		ARCH_ARM_V7: "armv7",

		ARCH_MIPS:       "mips",
		ARCH_MIPS_SF:    "mipssf",
		ARCH_MIPS_LE:    "mipsle",
		ARCH_MIPS_LE_SF: "mipslesf",

		ARCH_MIPS64:       "mips64",
		ARCH_MIPS64_SF:    "mips64sf",
		ARCH_MIPS64_LE:    "mips64le",
		ARCH_MIPS64_LE_SF: "mips64lesf",

		ARCH_PPC:       "ppc",
		ARCH_PPC_SF:    "ppcsf",
		ARCH_PPC_LE:    "ppcle",
		ARCH_PPC_LE_SF: "ppclesf",

		ARCH_PPC64:       "ppc64",
		ARCH_PPC64_LE:    "ppc64le",
		ARCH_PPC64_V8:    "ppc64v8",
		ARCH_PPC64_V8_LE: "ppc64v8le",
		ARCH_PPC64_V9:    "ppc64v9",
		ARCH_PPC64_V9_LE: "ppc64v9le",

		ARCH_RISCV_64: "riscv64",

		ARCH_S390X: "s390x",

		ARCH_IA64: "ia64",
	}
)

func TestArchMapping(t *testing.T) {
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
			for mArch := range requiredArchMappingValues {
				_, ok := test.mappingFunc(mArch)

				assert.True(t, ok, mArch)
			}
		})
	}
}
