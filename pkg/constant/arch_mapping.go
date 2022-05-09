package constant

import (
	. "arhat.dev/pkg/archconst"
)

type PlatformValue string

const (
	PLATFORM_ALPINE PlatformValue = "Alpine"
	PLATFORM_DEBIAN PlatformValue = "Debian"
	PLATFORM_GNU    PlatformValue = "GNU"

	// golang family
	PLATFORM_GOLANG PlatformValue = "Golang"
	PLATFORM_DOCKER PlatformValue = "Docker"
	PLATFORM_OCI    PlatformValue = "OCI"

	PLATFORM_DOCKERHUB PlatformValue = "DockerHub"

	PLATFORM_QEMU PlatformValue = "Qemu"

	// llvm family
	PLATFORM_LLVM PlatformValue = "LLVM"
	PLATFORM_ZIG  PlatformValue = "Zig"
	PLATFORM_RUST PlatformValue = "Rust"
)

type platformID uint32

const (
	_unknown_platform platformID = iota

	platformID_Alpine
	platformID_Debian
	platformID_GNU

	platformID_Golang
	platformID_Docker
	platformID_OCI

	platformID_DockerHub

	platformID_Qemu

	platformID_LLVM
	platformID_Zig
	platformID_Rust

	platformID_COUNT
)

func platform_id_of(platform PlatformValue) platformID {
	switch platform {
	case PLATFORM_ALPINE:
		return platformID_Alpine
	case PLATFORM_DEBIAN:
		return platformID_Debian
	case PLATFORM_GNU:
		return platformID_GNU
	case PLATFORM_GOLANG:
		return platformID_Golang
	case PLATFORM_DOCKER:
		return platformID_Docker
	case PLATFORM_OCI:
		return platformID_OCI
	case PLATFORM_DOCKERHUB:
		return platformID_DockerHub
	case PLATFORM_QEMU:
		return platformID_Qemu
	case PLATFORM_LLVM:
		return platformID_LLVM
	case PLATFORM_ZIG:
		return platformID_Zig
	case PLATFORM_RUST:
		return platformID_Rust
	default:
		return _unknown_platform
	}
}

type archID uint32

func (id archID) String() string {
	switch id {
	case archID_X86:
		return ARCH_X86
	case archID_X86_SF:
		return ARCH_X86_SF

	case archID_AMD64:
		return ARCH_AMD64
	case archID_AMD64_V1:
		return ARCH_AMD64_V1
	case archID_AMD64_V2:
		return ARCH_AMD64_V2
	case archID_AMD64_V3:
		return ARCH_AMD64_V3
	case archID_AMD64_V4:
		return ARCH_AMD64_V4

	case archID_ARM:
		return ARCH_ARM
	case archID_ARM_V5:
		return ARCH_ARM_V5
	case archID_ARM_V6:
		return ARCH_ARM_V6
	case archID_ARM_V7:
		return ARCH_ARM_V7

	case archID_ARM64:
		return ARCH_ARM64
	case archID_ARM64_V8:
		return ARCH_ARM64_V8
	case archID_ARM64_V9:
		return ARCH_ARM64_V9

	case archID_MIPS:
		return ARCH_MIPS
	case archID_MIPS_SF:
		return ARCH_MIPS_SF

	case archID_MIPS_LE:
		return ARCH_MIPS_LE
	case archID_MIPS_LE_SF:
		return ARCH_MIPS_LE_SF

	case archID_MIPS64:
		return ARCH_MIPS64
	case archID_MIPS64_SF:
		return ARCH_MIPS64_SF

	case archID_MIPS64_LE:
		return ARCH_MIPS64_LE
	case archID_MIPS64_LE_SF:
		return ARCH_MIPS64_LE_SF

	case archID_PPC:
		return ARCH_PPC
	case archID_PPC_SF:
		return ARCH_PPC_SF

	case archID_PPC_LE:
		return ARCH_PPC_LE
	case archID_PPC_LE_SF:
		return ARCH_PPC_LE_SF

	case archID_PPC64:
		return ARCH_PPC64
	case archID_PPC64_V8:
		return ARCH_PPC64_V8
	case archID_PPC64_V9:
		return ARCH_PPC64_V9

	case archID_PPC64_LE:
		return ARCH_PPC64_LE
	case archID_PPC64_LE_V8:
		return ARCH_PPC64_LE_V8
	case archID_PPC64_LE_V9:
		return ARCH_PPC64_LE_V9

	case archID_RISCV64:
		return ARCH_RISCV64

	case archID_S390X:
		return ARCH_S390X

	case archID_IA64:
		return ARCH_IA64

	default:
		return "<unknown>"
	}
}

const (
	_unknown_arch archID = iota

	archID_X86
	archID_X86_SF

	archID_AMD64
	archID_AMD64_V1
	archID_AMD64_V2
	archID_AMD64_V3
	archID_AMD64_V4

	archID_ARM
	archID_ARM_V5
	archID_ARM_V6
	archID_ARM_V7

	archID_ARM64
	archID_ARM64_V8
	archID_ARM64_V9

	archID_MIPS
	archID_MIPS_SF

	archID_MIPS_LE
	archID_MIPS_LE_SF

	archID_MIPS64
	archID_MIPS64_SF

	archID_MIPS64_LE
	archID_MIPS64_LE_SF

	archID_PPC
	archID_PPC_SF

	archID_PPC_LE
	archID_PPC_LE_SF

	archID_PPC64
	archID_PPC64_V8
	archID_PPC64_V9

	archID_PPC64_LE
	archID_PPC64_LE_V8
	archID_PPC64_LE_V9

	archID_RISCV64

	archID_S390X

	archID_IA64

	archID_COUNT
)

func arch_id_of(arch string) archID {
	switch arch {
	case ARCH_X86:
		return archID_X86
	case ARCH_X86_SF:
		return archID_X86_SF
	case ARCH_AMD64:
		return archID_AMD64
	case ARCH_AMD64_V1:
		return archID_AMD64_V1
	case ARCH_AMD64_V2:
		return archID_AMD64_V2
	case ARCH_AMD64_V3:
		return archID_AMD64_V3
	case ARCH_AMD64_V4:
		return archID_AMD64_V4
	case ARCH_ARM:
		return archID_ARM
	case ARCH_ARM_V5:
		return archID_ARM_V5
	case ARCH_ARM_V6:
		return archID_ARM_V6
	case ARCH_ARM_V7:
		return archID_ARM_V7
	case ARCH_ARM64:
		return archID_ARM64
	case ARCH_ARM64_V8:
		return archID_ARM64_V8
	case ARCH_ARM64_V9:
		return archID_ARM64_V9
	case ARCH_MIPS:
		return archID_MIPS
	case ARCH_MIPS_SF:
		return archID_MIPS_SF
	case ARCH_MIPS_LE:
		return archID_MIPS_LE
	case ARCH_MIPS_LE_SF:
		return archID_MIPS_LE_SF
	case ARCH_MIPS64:
		return archID_MIPS64
	case ARCH_MIPS64_SF:
		return archID_MIPS64_SF
	case ARCH_MIPS64_LE:
		return archID_MIPS64_LE
	case ARCH_MIPS64_LE_SF:
		return archID_MIPS64_LE_SF
	case ARCH_PPC:
		return archID_PPC
	case ARCH_PPC_SF:
		return archID_PPC_SF
	case ARCH_PPC_LE:
		return archID_PPC_LE
	case ARCH_PPC_LE_SF:
		return archID_PPC_LE_SF
	case ARCH_PPC64:
		return archID_PPC64
	case ARCH_PPC64_V8:
		return archID_PPC64_V8
	case ARCH_PPC64_V9:
		return archID_PPC64_V9
	case ARCH_PPC64_LE:
		return archID_PPC64_LE
	case ARCH_PPC64_LE_V8:
		return archID_PPC64_LE_V8
	case ARCH_PPC64_LE_V9:
		return archID_PPC64_LE_V9
	case ARCH_RISCV64:
		return archID_RISCV64
	case ARCH_S390X:
		return archID_S390X
	case ARCH_IA64:
		return archID_IA64
	default:
		return _unknown_arch
	}
}

func GetArch(platform PlatformValue, mArch string) (string, bool) {
	aid := arch_id_of(mArch)
	pid := platform_id_of(platform)

	return archMapping[aid][pid], aid != _unknown_arch && pid != _unknown_platform
}

func GetDebianArch(mArch string) (string, bool) { return GetArch(PLATFORM_DEBIAN, mArch) }
func GetAlpineArch(mArch string) (string, bool) { return GetArch(PLATFORM_ALPINE, mArch) }
func GetGNUArch(mArch string) (string, bool)    { return GetArch(PLATFORM_GNU, mArch) }

func GetGolangArch(mArch string) (string, bool) { return GetArch(PLATFORM_GOLANG, mArch) }
func GetOciArch(mArch string) (string, bool)    { return GetArch(PLATFORM_OCI, mArch) }
func GetDockerArch(mArch string) (string, bool) { return GetArch(PLATFORM_DOCKER, mArch) }

func GetQemuArch(mArch string) (string, bool) { return GetArch(PLATFORM_QEMU, mArch) }

func GetLLVMArch(mArch string) (string, bool) { return GetArch(PLATFORM_LLVM, mArch) }
func GetZigArch(mArch string) (string, bool)  { return GetArch(PLATFORM_ZIG, mArch) }
func GetRustArch(mArch string) (string, bool) { return GetArch(PLATFORM_RUST, mArch) }

func GetDockerHubArch(mArch, mKernel string) (string, bool) {
	arch, ok := GetArch(PLATFORM_DOCKERHUB, mArch)

	switch mKernel {
	case KERNEL_WINDOWS:
		if !ok {
			return "", false
		}
		return "win" + arch, ok
	case KERNEL_LINUX:
		fallthrough
	default:
		return arch, ok
	}
}

type platformArchMapping [platformID_COUNT]string

// Ref:
// for GNU values: https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/data/cputable
var archMapping = [archID_COUNT]platformArchMapping{
	archID_X86: {
		platformID_Alpine: "x86",
		platformID_Debian: "i386",
		platformID_GNU:    "i686",

		platformID_Golang: "386",
		platformID_Docker: "386",
		platformID_OCI:    "386",

		platformID_Qemu: "i386",

		platformID_DockerHub: "i386",

		platformID_LLVM: "x86",
	},
	archID_X86_SF: {
		platformID_Alpine: "x86",
		platformID_Debian: "i386",
		platformID_GNU:    "i686",

		platformID_Golang: "386",
		platformID_Docker: "386",
		platformID_OCI:    "386",

		platformID_Qemu: "i386",

		platformID_DockerHub: "i386",

		platformID_LLVM: "x86",
	},
	archID_AMD64: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_Qemu: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_AMD64_V1: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_Qemu: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_AMD64_V2: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_Qemu: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_AMD64_V3: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_Qemu: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_AMD64_V4: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_Qemu: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_ARM: {
		platformID_Alpine: "armv7",
		platformID_Debian: "armhf",
		platformID_GNU:    "arm",

		platformID_Golang: "arm",
		platformID_Docker: "arm",
		platformID_OCI:    "arm",

		platformID_Qemu: "arm",

		platformID_DockerHub: "arm32v7",

		platformID_LLVM: "armv7",
	},
	archID_ARM_V5: {
		platformID_Alpine: "armv5l",
		platformID_Debian: "armel",
		platformID_GNU:    "arm",

		platformID_Golang: "arm",
		platformID_Docker: "arm",
		platformID_OCI:    "arm",

		platformID_Qemu: "arm",

		platformID_DockerHub: "arm32v5",

		platformID_LLVM: "armv5",
	},
	archID_ARM_V6: {
		platformID_Alpine: "armhf",
		platformID_Debian: "armel",
		platformID_GNU:    "arm",

		platformID_Golang: "arm",
		platformID_Docker: "arm",
		platformID_OCI:    "arm",

		platformID_Qemu: "arm",

		platformID_DockerHub: "arm32v6",

		platformID_LLVM: "armv6",
	},
	archID_ARM_V7: {
		platformID_Alpine: "armv7",
		platformID_Debian: "armhf",
		platformID_GNU:    "arm",

		platformID_Golang: "arm",
		platformID_Docker: "arm",
		platformID_OCI:    "arm",

		platformID_Qemu: "arm",

		platformID_DockerHub: "arm32v7",

		platformID_LLVM: "armv7",
	},
	archID_ARM64: {
		platformID_Alpine: "aarch64",
		platformID_Debian: "arm64",
		platformID_GNU:    "aarch64",

		platformID_Golang: "arm64",
		platformID_Docker: "arm64",
		platformID_OCI:    "arm64",

		platformID_Qemu: "aarch64",

		platformID_DockerHub: "arm64v8",

		platformID_LLVM: "aarch64",
	},
	archID_ARM64_V8: {
		platformID_Alpine: "aarch64",
		platformID_Debian: "arm64",
		platformID_GNU:    "aarch64",

		platformID_Golang: "arm64",
		platformID_Docker: "arm64",
		platformID_OCI:    "arm64",

		platformID_Qemu: "aarch64",

		platformID_DockerHub: "arm64v8",

		platformID_LLVM: "aarch64",
	},
	// TODO: revise once standardized
	archID_ARM64_V9: {
		platformID_Alpine: "aarch64",
		platformID_Debian: "arm64",
		platformID_GNU:    "aarch64",

		platformID_Golang: "arm64",
		platformID_Docker: "arm64",
		platformID_OCI:    "arm64",

		platformID_Qemu: "aarch64",

		platformID_DockerHub: "arm64v8",

		platformID_LLVM: "aarch64",
	},
	archID_PPC: {
		platformID_Alpine: "",
		platformID_Debian: "powerpc",
		platformID_GNU:    "powerpc",

		platformID_Golang: "ppc",
		platformID_Docker: "ppc",
		platformID_OCI:    "ppc",

		platformID_Qemu: "ppc",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC_SF: {
		platformID_Alpine: "",
		platformID_Debian: "powerpc",
		platformID_GNU:    "powerpc",

		platformID_Golang: "ppc",
		platformID_Docker: "ppc",
		platformID_OCI:    "ppc",

		platformID_Qemu: "ppc",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC_LE: {
		platformID_Alpine: "",
		platformID_Debian: "powerpcel",
		platformID_GNU:    "powerpcle",

		platformID_Golang: "",
		platformID_Docker: "",
		platformID_OCI:    "",

		platformID_Qemu: "",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC_LE_SF: {
		platformID_Alpine: "",
		platformID_Debian: "powerpcel",
		platformID_GNU:    "powerpcle",

		platformID_Golang: "",
		platformID_Docker: "",
		platformID_OCI:    "",

		platformID_Qemu: "",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC64: {
		platformID_Alpine: "ppc64",
		platformID_Debian: "ppc64",
		platformID_GNU:    "powerpc64",

		platformID_Golang: "ppc64",
		platformID_Docker: "ppc64",
		platformID_OCI:    "ppc64",

		platformID_Qemu: "ppc64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC64_V8: {
		platformID_Alpine: "ppc64",
		platformID_Debian: "ppc64",
		platformID_GNU:    "powerpc64",

		platformID_Golang: "ppc64",
		platformID_Docker: "ppc64",
		platformID_OCI:    "ppc64",

		platformID_Qemu: "ppc64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC64_V9: {
		platformID_Alpine: "ppc64",
		platformID_Debian: "ppc64",
		platformID_GNU:    "powerpc64",

		platformID_Golang: "ppc64",
		platformID_Docker: "ppc64",
		platformID_OCI:    "ppc64",

		platformID_Qemu: "ppc64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC64_LE: {
		platformID_Alpine: "ppc64le",
		platformID_Debian: "ppc64el",
		platformID_GNU:    "powerpc64le",

		platformID_Golang: "ppc64le",
		platformID_Docker: "ppc64le",
		platformID_OCI:    "ppc64le",

		platformID_Qemu: "ppc64le",

		platformID_DockerHub: "ppc64le",

		platformID_LLVM: "ppc64le",
	},
	archID_PPC64_LE_V8: {
		platformID_Alpine: "ppc64le",
		platformID_Debian: "ppc64el",
		platformID_GNU:    "powerpc64le",

		platformID_Golang: "ppc64le",
		platformID_Docker: "ppc64le",
		platformID_OCI:    "ppc64le",

		platformID_Qemu: "ppc64le",

		platformID_DockerHub: "ppc64le",

		platformID_LLVM: "ppc64le",
	},
	archID_PPC64_LE_V9: {
		platformID_Alpine: "ppc64le",
		platformID_Debian: "ppc64el",
		platformID_GNU:    "powerpc64le",

		platformID_Golang: "ppc64le",
		platformID_Docker: "ppc64le",
		platformID_OCI:    "ppc64le",

		platformID_Qemu: "ppc64le",

		platformID_DockerHub: "ppc64le",

		platformID_LLVM: "ppc64le",
	},
	archID_MIPS: {
		platformID_Alpine: "mips",
		platformID_Debian: "mips",
		platformID_GNU:    "mips",

		platformID_Golang: "mips",
		platformID_Docker: "mips",
		platformID_OCI:    "mips",

		platformID_Qemu: "mips",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS_SF: {
		platformID_Alpine: "mips",
		platformID_Debian: "mips",
		platformID_GNU:    "mips",

		platformID_Golang: "mips",
		platformID_Docker: "mips",
		platformID_OCI:    "mips",

		platformID_Qemu: "mips",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS_LE: {
		platformID_Alpine: "mipsel",
		platformID_Debian: "mipsel",
		platformID_GNU:    "mipsel",

		platformID_Golang: "mipsle",
		platformID_Docker: "mipsle",
		platformID_OCI:    "mipsle",

		platformID_Qemu: "mipsel",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS_LE_SF: {
		platformID_Alpine: "mipsel",
		platformID_Debian: "mipsel",
		platformID_GNU:    "mipsel",

		platformID_Golang: "mipsle",
		platformID_Docker: "mipsle",
		platformID_OCI:    "mipsle",

		platformID_Qemu: "mipsel",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS64: {
		platformID_Alpine: "mips64",
		platformID_Debian: "mips64",
		platformID_GNU:    "mips64",

		platformID_Golang: "mips64",
		platformID_Docker: "mips64",
		platformID_OCI:    "mips64",

		platformID_Qemu: "mips64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS64_SF: {
		platformID_Alpine: "mips64",
		platformID_Debian: "mips64",
		platformID_GNU:    "mips64",

		platformID_Golang: "mips64",
		platformID_Docker: "mips64",
		platformID_OCI:    "mips64",

		platformID_Qemu: "mips64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS64_LE: {
		platformID_Alpine: "mips64el",
		platformID_Debian: "mips64el",
		platformID_GNU:    "mips64el",

		platformID_Golang: "mips64le",
		platformID_Docker: "mips64le",
		platformID_OCI:    "mips64le",

		platformID_Qemu: "mips64el",

		platformID_DockerHub: "mips64le",

		platformID_LLVM: "mips64el",
	},
	archID_MIPS64_LE_SF: {
		platformID_Alpine: "mips64el",
		platformID_Debian: "mips64el",
		platformID_GNU:    "mips64el",

		platformID_Golang: "mips64le",
		platformID_Docker: "mips64le",
		platformID_OCI:    "mips64le",

		platformID_Qemu: "mips64el",

		platformID_DockerHub: "mips64le",

		platformID_LLVM: "mips64el",
	},
	archID_RISCV64: {
		platformID_Alpine: "riscv64",
		platformID_Debian: "riscv64",
		platformID_GNU:    "riscv64",

		platformID_Golang: "riscv64",
		platformID_Docker: "riscv64",
		platformID_OCI:    "riscv64",

		platformID_Qemu:      "riscv64",
		platformID_DockerHub: "riscv64",

		platformID_LLVM: "",
	},
	archID_S390X: {
		platformID_Alpine: "s390x",
		platformID_Debian: "s390x",
		platformID_GNU:    "s390x",

		platformID_Golang: "s390x",
		platformID_Docker: "s390x",
		platformID_OCI:    "s390x",

		platformID_Qemu: "s390x",

		platformID_DockerHub: "s390x",

		platformID_LLVM: "systemz",
	},
	archID_IA64: {
		platformID_Alpine: "",
		platformID_Debian: "ia64",
		platformID_GNU:    "ia64",

		platformID_Golang: "",
		platformID_Docker: "",
		platformID_OCI:    "",

		platformID_Qemu: "",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
}
