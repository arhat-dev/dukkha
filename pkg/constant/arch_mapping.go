package constant

import (
	"reflect"
	"strings"
)

func GetArch(platform, mArch string) (string, bool) {
	platform = strings.ToLower(platform)
	fieldName, ok := supportedPlatforms[platform]
	if !ok {
		return "", false
	}

	m, ok := archMapping[mArch]
	if ok {
		return reflect.ValueOf(m).FieldByName(fieldName).String(), true
	}

	return "", false
}

func GetDebianArch(mArch string) (string, bool) { return GetArch("debian", mArch) }
func GetAlpineArch(mArch string) (string, bool) { return GetArch("alpine", mArch) }
func GetGNUArch(mArch string) (string, bool)    { return GetArch("gnu", mArch) }

func GetGolangArch(mArch string) (string, bool) { return GetArch("golang", mArch) }
func GetOciArch(mArch string) (string, bool)    { return GetArch("oci", mArch) }
func GetDockerArch(mArch string) (string, bool) { return GetArch("docker", mArch) }

func GetQemuArch(mArch string) (string, bool) { return GetArch("qemu", mArch) }

func GetLLVMArch(mArch string) (string, bool) { return GetArch("llvm", mArch) }
func GetZigArch(mArch string) (string, bool)  { return GetArch("zig", mArch) }
func GetRustArch(mArch string) (string, bool) { return GetArch("rust", mArch) }

func GetDockerHubArch(mArch, mKernel string) (string, bool) {
	arch, ok := GetArch("dockerhub", mArch)

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

// Ref:
// for GNU values: https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/data/cputable
var archMapping = map[string]ArchMappingValues{
	ARCH_X86: {
		Alpine: "x86",
		Debian: "i386",
		GNU:    "i686",

		Golang: "386",
		Docker: "386",
		OCI:    "386",

		Qemu: "i386",

		DockerHub: "i386",

		LLVM: "x86",
	},
	ARCH_AMD64: {
		Alpine: "x86_64",
		Debian: "amd64",
		GNU:    "x86_64",

		Golang: "amd64",
		Docker: "amd64",
		OCI:    "amd64",

		Qemu: "x86_64",

		DockerHub: "amd64",

		LLVM: "x86_64",
	},

	ARCH_ARM_V5: {
		Alpine: "armv5l",
		Debian: "armel",
		GNU:    "arm",

		Golang: "arm",
		Docker: "arm",
		OCI:    "arm",

		Qemu: "arm",

		DockerHub: "arm32v5",

		LLVM: "armv5",
	},
	ARCH_ARM_V6: {
		Alpine: "armhf",
		Debian: "armel",
		GNU:    "arm",

		Golang: "arm",
		Docker: "arm",
		OCI:    "arm",

		Qemu: "arm",

		DockerHub: "arm32v6",

		LLVM: "armv6",
	},
	ARCH_ARM_V7: {
		Alpine: "armv7",
		Debian: "armhf",
		GNU:    "arm",

		Golang: "arm",
		Docker: "arm",
		OCI:    "arm",

		Qemu: "arm",

		DockerHub: "arm32v7",

		LLVM: "armv7",
	},
	ARCH_ARM64: {
		Alpine: "aarch64",
		Debian: "arm64",
		GNU:    "aarch64",

		Golang: "arm64",
		Docker: "arm64",
		OCI:    "arm64",

		Qemu: "aarch64",

		DockerHub: "arm64v8",

		LLVM: "aarch64",
	},

	ARCH_PPC: {
		Alpine: "",
		Debian: "powerpc",
		GNU:    "powerpc",

		Golang: "ppc",
		Docker: "ppc",
		OCI:    "ppc",

		Qemu: "ppc",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_PPC_SF: {
		Alpine: "",
		Debian: "powerpc",
		GNU:    "powerpc",

		Golang: "ppc",
		Docker: "ppc",
		OCI:    "ppc",

		Qemu: "ppc",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_PPC_LE: {
		Alpine: "",
		Debian: "powerpcel",
		GNU:    "powerpcle",

		Golang: "",
		Docker: "",
		OCI:    "",

		Qemu: "",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_PPC_LE_SF: {
		Alpine: "",
		Debian: "powerpcel",
		GNU:    "powerpcle",

		Golang: "",
		Docker: "",
		OCI:    "",

		Qemu: "",

		DockerHub: "",

		LLVM: "",
	},

	ARCH_PPC64: {
		Alpine: "ppc64",
		Debian: "ppc64",
		GNU:    "powerpc64",

		Golang: "ppc64",
		Docker: "ppc64",
		OCI:    "ppc64",

		Qemu: "ppc64",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_PPC64_LE: {
		Alpine: "ppc64le",
		Debian: "ppc64el",
		GNU:    "powerpc64le",

		Golang: "ppc64le",
		Docker: "ppc64le",
		OCI:    "ppc64le",

		Qemu: "ppc64le",

		DockerHub: "ppc64le",

		LLVM: "ppc64le",
	},

	ARCH_MIPS: {
		Alpine: "mips",
		Debian: "mips",
		GNU:    "mips",

		Golang: "mips",
		Docker: "mips",
		OCI:    "mips",

		Qemu: "mips",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_MIPS_SF: {
		Alpine: "mips",
		Debian: "mips",
		GNU:    "mips",

		Golang: "mips",
		Docker: "mips",
		OCI:    "mips",

		Qemu: "mips",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_MIPS_LE: {
		Alpine: "mipsel",
		Debian: "mipsel",
		GNU:    "mipsel",

		Golang: "mipsle",
		Docker: "mipsle",
		OCI:    "mipsle",

		Qemu: "mipsel",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_MIPS_LE_SF: {
		Alpine: "mipsel",
		Debian: "mipsel",
		GNU:    "mipsel",

		Golang: "mipsle",
		Docker: "mipsle",
		OCI:    "mipsle",

		Qemu: "mipsel",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_MIPS64: {
		Alpine: "mips64",
		Debian: "mips64",
		GNU:    "mips64",

		Golang: "mips64",
		Docker: "mips64",
		OCI:    "mips64",

		Qemu: "mips64",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_MIPS64_SF: {
		Alpine: "mips64",
		Debian: "mips64",
		GNU:    "mips64",

		Golang: "mips64",
		Docker: "mips64",
		OCI:    "mips64",

		Qemu: "mips64",

		DockerHub: "",

		LLVM: "",
	},
	ARCH_MIPS64_LE: {
		Alpine: "mips64el",
		Debian: "mips64el",
		GNU:    "mips64el",

		Golang: "mips64le",
		Docker: "mips64le",
		OCI:    "mips64le",

		Qemu: "mips64el",

		DockerHub: "mips64le",

		LLVM: "mips64el",
	},
	ARCH_MIPS64_LE_SF: {
		Alpine: "mips64el",
		Debian: "mips64el",
		GNU:    "mips64el",

		Golang: "mips64le",
		Docker: "mips64le",
		OCI:    "mips64le",

		Qemu: "mips64el",

		DockerHub: "mips64le",

		LLVM: "mips64el",
	},

	ARCH_RISCV_64: {
		Alpine: "riscv64",
		Debian: "riscv64",
		GNU:    "riscv64",

		Golang: "riscv64",
		Docker: "riscv64",
		OCI:    "riscv64",

		Qemu:      "riscv64",
		DockerHub: "riscv64",

		LLVM: "",
	},
	ARCH_S390X: {
		Alpine: "s390x",
		Debian: "s390x",
		GNU:    "s390x",

		Golang: "s390x",
		Docker: "s390x",
		OCI:    "s390x",

		Qemu: "s390x",

		DockerHub: "s390x",

		LLVM: "systemz",
	},

	ARCH_IA64: {
		Alpine: "",
		Debian: "ia64",
		GNU:    "ia64",

		Golang: "",
		Docker: "",
		OCI:    "",

		Qemu: "",

		DockerHub: "",

		LLVM: "",
	},
}
