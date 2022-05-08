package constant

import (
	"testing"

	"arhat.dev/pkg/archconst"
	"github.com/stretchr/testify/assert"
)

var (
	// TODO: update this map when changing arch.go
	requiredArchMappingValues = map[string]string{
		archconst.ARCH_X86:    "x86",
		archconst.ARCH_X86_SF: "x86sf",

		archconst.ARCH_AMD64: "amd64",

		archconst.ARCH_AMD64_V1: "amd64v1",
		archconst.ARCH_AMD64_V2: "amd64v2",
		archconst.ARCH_AMD64_V3: "amd64v3",
		archconst.ARCH_AMD64_V4: "amd64v4",

		archconst.ARCH_ARM64: "arm64",

		archconst.ARCH_ARM_V5: "armv5",
		archconst.ARCH_ARM_V6: "armv6",
		archconst.ARCH_ARM_V7: "armv7",

		archconst.ARCH_MIPS:       "mips",
		archconst.ARCH_MIPS_SF:    "mipssf",
		archconst.ARCH_MIPS_LE:    "mipsle",
		archconst.ARCH_MIPS_LE_SF: "mipslesf",

		archconst.ARCH_MIPS64:       "mips64",
		archconst.ARCH_MIPS64_SF:    "mips64sf",
		archconst.ARCH_MIPS64_LE:    "mips64le",
		archconst.ARCH_MIPS64_LE_SF: "mips64lesf",

		archconst.ARCH_PPC:       "ppc",
		archconst.ARCH_PPC_SF:    "ppcsf",
		archconst.ARCH_PPC_LE:    "ppcle",
		archconst.ARCH_PPC_LE_SF: "ppclesf",

		archconst.ARCH_PPC64:       "ppc64",
		archconst.ARCH_PPC64_LE:    "ppc64le",
		archconst.ARCH_PPC64_V8:    "ppc64v8",
		archconst.ARCH_PPC64_LE_V8: "ppc64lev8",
		archconst.ARCH_PPC64_V9:    "ppc64v9",
		archconst.ARCH_PPC64_LE_V9: "ppc64lev9",

		archconst.ARCH_RISCV_64: "riscv64",

		archconst.ARCH_S390X: "s390x",

		archconst.ARCH_IA64: "ia64",
	}
)

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
			for mArch := range requiredArchMappingValues {
				_, ok := test.mappingFunc(mArch)

				assert.True(t, ok, mArch)
			}
		})
	}
}
