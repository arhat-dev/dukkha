package constant

func GetArch(platform, mArch string) (string, bool) {
	aid := arch_id_of(mArch)
	pid := platform_id_of(platform)

	return archMapping[aid][pid], aid != _unknown_arch && pid != _unknown_platform
}

func GetDebianArch(mArch string) (string, bool) { return GetArch(Platform_Debian, mArch) }
func GetAlpineArch(mArch string) (string, bool) { return GetArch(Platform_Alpine, mArch) }
func GetGNUArch(mArch string) (string, bool)    { return GetArch(Platform_GNU, mArch) }

func GetGolangArch(mArch string) (string, bool) { return GetArch(Platform_Golang, mArch) }
func GetOciArch(mArch string) (string, bool)    { return GetArch(Platform_OCI, mArch) }
func GetDockerArch(mArch string) (string, bool) { return GetArch(Platform_Docker, mArch) }

func GetQemuArch(mArch string) (string, bool) { return GetArch(Platform_QEMU, mArch) }

func GetLLVMArch(mArch string) (string, bool) { return GetArch(Platform_LLVM, mArch) }
func GetZigArch(mArch string) (string, bool)  { return GetArch(Platform_Zig, mArch) }
func GetRustArch(mArch string) (string, bool) { return GetArch(Platform_Rust, mArch) }

func GetDockerHubArch(mArch, mKernel string) (string, bool) {
	arch, ok := GetArch(Platform_DockerHub, mArch)

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
		platformID_Ubuntu: "i386",
		platformID_GNU:    "i686",

		platformID_Golang: "386",
		platformID_Docker: "386",
		platformID_OCI:    "386",

		platformID_QEMU: "i386",

		platformID_DockerHub: "i386",

		platformID_LLVM: "x86",
	},
	archID_X86_SF: {
		platformID_Alpine: "x86",
		platformID_Debian: "i386",
		platformID_Ubuntu: "i386",
		platformID_GNU:    "i686",

		platformID_Golang: "386",
		platformID_Docker: "386",
		platformID_OCI:    "386",

		platformID_QEMU: "i386",

		platformID_DockerHub: "i386",

		platformID_LLVM: "x86",
	},
	archID_AMD64: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_Ubuntu: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_QEMU: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_AMD64_V1: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_Ubuntu: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_QEMU: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_AMD64_V2: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_Ubuntu: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_QEMU: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_AMD64_V3: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_Ubuntu: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_QEMU: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_AMD64_V4: {
		platformID_Alpine: "x86_64",
		platformID_Debian: "amd64",
		platformID_Ubuntu: "amd64",
		platformID_GNU:    "x86_64",

		platformID_Golang: "amd64",
		platformID_Docker: "amd64",
		platformID_OCI:    "amd64",

		platformID_QEMU: "x86_64",

		platformID_DockerHub: "amd64",

		platformID_LLVM: "x86_64",
	},
	archID_ARM: {
		platformID_Alpine: "armv7",
		platformID_Debian: "armhf",
		platformID_Ubuntu: "armhf",
		platformID_GNU:    "arm",

		platformID_Golang: "arm",
		platformID_Docker: "arm",
		platformID_OCI:    "arm",

		platformID_QEMU: "arm",

		platformID_DockerHub: "arm32v7",

		platformID_LLVM: "armv7",
	},
	archID_ARM_V5: {
		platformID_Alpine: "armv5l",
		platformID_Debian: "armel",
		platformID_Ubuntu: "armel",
		platformID_GNU:    "arm",

		platformID_Golang: "arm",
		platformID_Docker: "arm",
		platformID_OCI:    "arm",

		platformID_QEMU: "arm",

		platformID_DockerHub: "arm32v5",

		platformID_LLVM: "armv5",
	},
	archID_ARM_V6: {
		platformID_Alpine: "armhf",
		platformID_Debian: "armel",
		platformID_Ubuntu: "armel",
		platformID_GNU:    "arm",

		platformID_Golang: "arm",
		platformID_Docker: "arm",
		platformID_OCI:    "arm",

		platformID_QEMU: "arm",

		platformID_DockerHub: "arm32v6",

		platformID_LLVM: "armv6",
	},
	archID_ARM_V7: {
		platformID_Alpine: "armv7",
		platformID_Debian: "armhf",
		platformID_Ubuntu: "armhf",
		platformID_GNU:    "arm",

		platformID_Golang: "arm",
		platformID_Docker: "arm",
		platformID_OCI:    "arm",

		platformID_QEMU: "arm",

		platformID_DockerHub: "arm32v7",

		platformID_LLVM: "armv7",
	},
	archID_ARM64: {
		platformID_Alpine: "aarch64",
		platformID_Debian: "arm64",
		platformID_Ubuntu: "arm64",
		platformID_GNU:    "aarch64",

		platformID_Golang: "arm64",
		platformID_Docker: "arm64",
		platformID_OCI:    "arm64",

		platformID_QEMU: "aarch64",

		platformID_DockerHub: "arm64v8",

		platformID_LLVM: "aarch64",
	},
	archID_ARM64_V8: {
		platformID_Alpine: "aarch64",
		platformID_Debian: "arm64",
		platformID_Ubuntu: "arm64",
		platformID_GNU:    "aarch64",

		platformID_Golang: "arm64",
		platformID_Docker: "arm64",
		platformID_OCI:    "arm64",

		platformID_QEMU: "aarch64",

		platformID_DockerHub: "arm64v8",

		platformID_LLVM: "aarch64",
	},
	// TODO: revise once standardized
	archID_ARM64_V9: {
		platformID_Alpine: "aarch64",
		platformID_Debian: "arm64",
		platformID_Ubuntu: "arm64",
		platformID_GNU:    "aarch64",

		platformID_Golang: "arm64",
		platformID_Docker: "arm64",
		platformID_OCI:    "arm64",

		platformID_QEMU: "aarch64",

		platformID_DockerHub: "arm64v8",

		platformID_LLVM: "aarch64",
	},
	archID_PPC: {
		platformID_Alpine: "",
		platformID_Debian: "powerpc",
		platformID_Ubuntu: "powerpc",
		platformID_GNU:    "powerpc",

		platformID_Golang: "ppc",
		platformID_Docker: "ppc",
		platformID_OCI:    "ppc",

		platformID_QEMU: "ppc",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC_SF: {
		platformID_Alpine: "",
		platformID_Debian: "powerpc",
		platformID_Ubuntu: "powerpc",
		platformID_GNU:    "powerpc",

		platformID_Golang: "ppc",
		platformID_Docker: "ppc",
		platformID_OCI:    "ppc",

		platformID_QEMU: "ppc",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC_LE: {
		platformID_Alpine: "",
		platformID_Debian: "powerpcel",
		platformID_Ubuntu: "powerpcel",
		platformID_GNU:    "powerpcle",

		platformID_Golang: "",
		platformID_Docker: "",
		platformID_OCI:    "",

		platformID_QEMU: "",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC_LE_SF: {
		platformID_Alpine: "",
		platformID_Debian: "powerpcel",
		platformID_Ubuntu: "powerpcel",
		platformID_GNU:    "powerpcle",

		platformID_Golang: "",
		platformID_Docker: "",
		platformID_OCI:    "",

		platformID_QEMU: "",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC64: {
		platformID_Alpine: "ppc64",
		platformID_Debian: "ppc64",
		platformID_Ubuntu: "ppc64",
		platformID_GNU:    "powerpc64",

		platformID_Golang: "ppc64",
		platformID_Docker: "ppc64",
		platformID_OCI:    "ppc64",

		platformID_QEMU: "ppc64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC64_V8: {
		platformID_Alpine: "ppc64",
		platformID_Debian: "ppc64",
		platformID_Ubuntu: "ppc64",
		platformID_GNU:    "powerpc64",

		platformID_Golang: "ppc64",
		platformID_Docker: "ppc64",
		platformID_OCI:    "ppc64",

		platformID_QEMU: "ppc64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC64_V9: {
		platformID_Alpine: "ppc64",
		platformID_Debian: "ppc64",
		platformID_Ubuntu: "ppc64",
		platformID_GNU:    "powerpc64",

		platformID_Golang: "ppc64",
		platformID_Docker: "ppc64",
		platformID_OCI:    "ppc64",

		platformID_QEMU: "ppc64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_PPC64_LE: {
		platformID_Alpine: "ppc64le",
		platformID_Debian: "ppc64el",
		platformID_Ubuntu: "ppc64el",
		platformID_GNU:    "powerpc64le",

		platformID_Golang: "ppc64le",
		platformID_Docker: "ppc64le",
		platformID_OCI:    "ppc64le",

		platformID_QEMU: "ppc64le",

		platformID_DockerHub: "ppc64le",

		platformID_LLVM: "ppc64le",
	},
	archID_PPC64_LE_V8: {
		platformID_Alpine: "ppc64le",
		platformID_Debian: "ppc64el",
		platformID_Ubuntu: "ppc64el",
		platformID_GNU:    "powerpc64le",

		platformID_Golang: "ppc64le",
		platformID_Docker: "ppc64le",
		platformID_OCI:    "ppc64le",

		platformID_QEMU: "ppc64le",

		platformID_DockerHub: "ppc64le",

		platformID_LLVM: "ppc64le",
	},
	archID_PPC64_LE_V9: {
		platformID_Alpine: "ppc64le",
		platformID_Debian: "ppc64el",
		platformID_Ubuntu: "ppc64el",
		platformID_GNU:    "powerpc64le",

		platformID_Golang: "ppc64le",
		platformID_Docker: "ppc64le",
		platformID_OCI:    "ppc64le",

		platformID_QEMU: "ppc64le",

		platformID_DockerHub: "ppc64le",

		platformID_LLVM: "ppc64le",
	},
	archID_MIPS: {
		platformID_Alpine: "mips",
		platformID_Debian: "mips",
		platformID_Ubuntu: "mips",
		platformID_GNU:    "mips",

		platformID_Golang: "mips",
		platformID_Docker: "mips",
		platformID_OCI:    "mips",

		platformID_QEMU: "mips",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS_SF: {
		platformID_Alpine: "mips",
		platformID_Debian: "mips",
		platformID_Ubuntu: "mips",
		platformID_GNU:    "mips",

		platformID_Golang: "mips",
		platformID_Docker: "mips",
		platformID_OCI:    "mips",

		platformID_QEMU: "mips",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS_LE: {
		platformID_Alpine: "mipsel",
		platformID_Debian: "mipsel",
		platformID_Ubuntu: "mipsel",
		platformID_GNU:    "mipsel",

		platformID_Golang: "mipsle",
		platformID_Docker: "mipsle",
		platformID_OCI:    "mipsle",

		platformID_QEMU: "mipsel",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS_LE_SF: {
		platformID_Alpine: "mipsel",
		platformID_Debian: "mipsel",
		platformID_Ubuntu: "mipsel",
		platformID_GNU:    "mipsel",

		platformID_Golang: "mipsle",
		platformID_Docker: "mipsle",
		platformID_OCI:    "mipsle",

		platformID_QEMU: "mipsel",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS64: {
		platformID_Alpine: "mips64",
		platformID_Debian: "mips64",
		platformID_Ubuntu: "mips64",
		platformID_GNU:    "mips64",

		platformID_Golang: "mips64",
		platformID_Docker: "mips64",
		platformID_OCI:    "mips64",

		platformID_QEMU: "mips64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS64_SF: {
		platformID_Alpine: "mips64",
		platformID_Debian: "mips64",
		platformID_Ubuntu: "mips64",
		platformID_GNU:    "mips64",

		platformID_Golang: "mips64",
		platformID_Docker: "mips64",
		platformID_OCI:    "mips64",

		platformID_QEMU: "mips64",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
	archID_MIPS64_LE: {
		platformID_Alpine: "mips64el",
		platformID_Debian: "mips64el",
		platformID_Ubuntu: "mips64el",
		platformID_GNU:    "mips64el",

		platformID_Golang: "mips64le",
		platformID_Docker: "mips64le",
		platformID_OCI:    "mips64le",

		platformID_QEMU: "mips64el",

		platformID_DockerHub: "mips64le",

		platformID_LLVM: "mips64el",
	},
	archID_MIPS64_LE_SF: {
		platformID_Alpine: "mips64el",
		platformID_Debian: "mips64el",
		platformID_Ubuntu: "mips64el",
		platformID_GNU:    "mips64el",

		platformID_Golang: "mips64le",
		platformID_Docker: "mips64le",
		platformID_OCI:    "mips64le",

		platformID_QEMU: "mips64el",

		platformID_DockerHub: "mips64le",

		platformID_LLVM: "mips64el",
	},
	archID_RISCV64: {
		platformID_Alpine: "riscv64",
		platformID_Debian: "riscv64",
		platformID_Ubuntu: "riscv64",
		platformID_GNU:    "riscv64",

		platformID_Golang: "riscv64",
		platformID_Docker: "riscv64",
		platformID_OCI:    "riscv64",

		platformID_QEMU:      "riscv64",
		platformID_DockerHub: "riscv64",

		platformID_LLVM: "",
	},
	archID_S390X: {
		platformID_Alpine: "s390x",
		platformID_Debian: "s390x",
		platformID_Ubuntu: "s390x",
		platformID_GNU:    "s390x",

		platformID_Golang: "s390x",
		platformID_Docker: "s390x",
		platformID_OCI:    "s390x",

		platformID_QEMU: "s390x",

		platformID_DockerHub: "s390x",

		platformID_LLVM: "systemz",
	},
	archID_IA64: {
		platformID_Alpine: "",
		platformID_Debian: "ia64",
		platformID_Ubuntu: "ia64",
		platformID_GNU:    "ia64",

		platformID_Golang: "",
		platformID_Docker: "",
		platformID_OCI:    "",

		platformID_QEMU: "",

		platformID_DockerHub: "",

		platformID_LLVM: "",
	},
}
