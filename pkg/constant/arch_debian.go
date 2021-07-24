package constant

func GetDebianArch(mArch string) (string, bool) {
	v, ok := map[string]string{
		ARCH_X86:   "i386",
		ARCH_AMD64: "amd64",

		ARCH_ARM_V5: "armel",
		ARCH_ARM_V6: "armel",
		ARCH_ARM_V7: "armhf",
		ARCH_ARM64:  "arm64",

		ARCH_PPC64:    "ppc64",
		ARCH_PPC64_LE: "ppc64el",

		ARCH_MIPS:         "mips",
		ARCH_MIPS_SF:      "mips",
		ARCH_MIPS_LE:      "mipsel",
		ARCH_MIPS_LE_SF:   "mipsel",
		ARCH_MIPS64:       "mips64",
		ARCH_MIPS64_SF:    "mips64",
		ARCH_MIPS64_LE:    "mips64el",
		ARCH_MIPS64_LE_SF: "mips64el",

		ARCH_RISCV_64: "riscv64",
		ARCH_S390X:    "s390x",

		ARCH_IA64: "",

		ARCH_PPC:       "",
		ARCH_PPC_SF:    "",
		ARCH_PPC_LE:    "",
		ARCH_PPC_LE_SF: "",
	}[mArch]

	return v, ok
}

func GetDebianTripleName(mArch, targetKernel, targetLibc string) (string, bool) {
	// TODO: adjust triple name according to target kernel
	_ = targetKernel

	switch targetLibc {
	case LIBC_MUSL:
		// https://packages.debian.org/buster/musl-dev
		// check list of files
		v, ok := map[string]string{
			ARCH_X86:   "i386-linux-musl",
			ARCH_AMD64: "x86_64-linux-musl",

			ARCH_ARM_V5: "arm-linux-musleabi",
			ARCH_ARM_V6: "arm-linux-musleabi",
			ARCH_ARM_V7: "arm-linux-musleabihf",
			ARCH_ARM64:  "aarch64-linux-musl",

			ARCH_MIPS:         "mips-linux-musl",
			ARCH_MIPS_SF:      "mips-linux-musl",
			ARCH_MIPS_LE:      "mipsel-linux-musl",
			ARCH_MIPS_LE_SF:   "mipsel-linux-musl",
			ARCH_MIPS64:       "mips64-linux-musl",
			ARCH_MIPS64_SF:    "mips64-linux-musl",
			ARCH_MIPS64_LE:    "mips64el-linux-musl",
			ARCH_MIPS64_LE_SF: "mips64el-linux-musl",

			ARCH_S390X: "s390x-linux-musl",

			// http://ftp.ports.debian.org/debian-ports/pool-riscv64/main/m/musl/
			// download one musl-dev package
			// list package contents with following commands
			//
			// $ ar -x musl-dev_1.2.2-3_riscv64.deb
			// $ tar -tvf data.tar.xz
			ARCH_RISCV_64: "riscv64-linux-musl",

			ARCH_IA64: "",

			ARCH_PPC:       "",
			ARCH_PPC_SF:    "",
			ARCH_PPC_LE:    "",
			ARCH_PPC_LE_SF: "",

			ARCH_PPC64:    "",
			ARCH_PPC64_LE: "",
		}[mArch]

		return v, ok
	case LIBC_MSVC:
		v, ok := map[string]string{
			// https://packages.debian.org/buster/mingw-w64-i686-dev
			// check list of files
			ARCH_X86: "i686-w64-mingw32",
			// https://packages.debian.org/buster/mingw-w64-x86-64-dev
			// check list of files
			ARCH_AMD64: "x86_64-w64-mingw32",

			ARCH_IA64: "",

			ARCH_ARM_V5: "",
			ARCH_ARM_V6: "",
			ARCH_ARM_V7: "",
			ARCH_ARM64:  "",

			ARCH_PPC:       "",
			ARCH_PPC_SF:    "",
			ARCH_PPC_LE:    "",
			ARCH_PPC_LE_SF: "",

			ARCH_PPC64:    "",
			ARCH_PPC64_LE: "",

			ARCH_MIPS:         "",
			ARCH_MIPS_SF:      "",
			ARCH_MIPS_LE:      "",
			ARCH_MIPS_LE_SF:   "",
			ARCH_MIPS64:       "",
			ARCH_MIPS64_SF:    "",
			ARCH_MIPS64_LE:    "",
			ARCH_MIPS64_LE_SF: "",

			ARCH_RISCV_64: "",
			ARCH_S390X:    "",
		}[mArch]

		return v, ok
	case LIBC_GLIBC:
		fallthrough
	default:
		v, ok := map[string]string{
			ARCH_X86:   "i686-linux-gnu",
			ARCH_AMD64: "x86_64-linux-gnu",

			ARCH_ARM_V5: "arm-linux-gnueabi",
			ARCH_ARM_V6: "arm-linux-gnueabi",
			ARCH_ARM_V7: "arm-linux-gnueabihf",
			ARCH_ARM64:  "aarch64-linux-gnu",

			ARCH_PPC64:    "powerpc64-linux-gnu",
			ARCH_PPC64_LE: "powerpc64le-linux-gnu",

			ARCH_MIPS:         "mips-linux-gnu",
			ARCH_MIPS_SF:      "mips-linux-gnu",
			ARCH_MIPS_LE:      "mipsel-linux-gnu",
			ARCH_MIPS_LE_SF:   "mipsel-linux-gnu",
			ARCH_MIPS64:       "mips64-linux-gnuabi64",
			ARCH_MIPS64_SF:    "mips64-linux-gnuabi64",
			ARCH_MIPS64_LE:    "mips64el-linux-gnuabi64",
			ARCH_MIPS64_LE_SF: "mips64el-linux-gnuabi64",

			ARCH_RISCV_64: "riscv64-linux-gnu",
			ARCH_S390X:    "s390x-linux-gnu",

			ARCH_IA64: "",

			ARCH_PPC:       "",
			ARCH_PPC_SF:    "",
			ARCH_PPC_LE:    "",
			ARCH_PPC_LE_SF: "",
		}[mArch]

		return v, ok
	}
}
