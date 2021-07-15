package constant

func GetDebianArch(mArch string) string {
	return map[string]string{
		ARCH_X86:   "i386",
		ARCH_AMD64: "amd64",

		ARCH_ARM_V5: "armel",
		ARCH_ARM_V6: "armel",
		ARCH_ARM_V7: "armhf",
		ARCH_ARM64:  "arm64",

		ARCH_PPC64:    "ppc64",
		ARCH_PPC64_LE: "ppc64el",

		ARCH_MIPS:         "mips",
		ARCH_MIPS_HF:      "mips",
		ARCH_MIPS_LE:      "mipsel",
		ARCH_MIPS_LE_HF:   "mipsel",
		ARCH_MIPS64:       "mips64",
		ARCH_MIPS64_HF:    "mips64",
		ARCH_MIPS64_LE:    "mips64el",
		ARCH_MIPS64_LE_HF: "mips64el",

		ARCH_RISCV_64: "riscv64",
		ARCH_S390X:    "s390x",
	}[mArch]
}

func GetDebianTripleName(mArch, targetKernel, targetLibc string) string {
	// TODO: adjust triple name according to target kernel
	_ = targetKernel

	switch targetLibc {
	case "musl":
		// https://packages.debian.org/buster/musl-dev
		// check list of files
		return map[string]string{
			ARCH_X86:   "i386-linux-musl",
			ARCH_AMD64: "x86_64-linux-musl",

			ARCH_ARM_V5: "arm-linux-musleabi",
			ARCH_ARM_V6: "arm-linux-musleabi",
			ARCH_ARM_V7: "arm-linux-musleabihf",
			ARCH_ARM64:  "aarch64-linux-musl",

			ARCH_MIPS:         "mips-linux-musl",
			ARCH_MIPS_HF:      "mips-linux-musl",
			ARCH_MIPS_LE:      "mipsel-linux-musl",
			ARCH_MIPS_LE_HF:   "mipsel-linux-musl",
			ARCH_MIPS64:       "mips64-linux-musl",
			ARCH_MIPS64_HF:    "mips64-linux-musl",
			ARCH_MIPS64_LE:    "mips64el-linux-musl",
			ARCH_MIPS64_LE_HF: "mips64el-linux-musl",

			ARCH_S390X: "s390x-linux-musl",

			// http://ftp.ports.debian.org/debian-ports/pool-riscv64/main/m/musl/
			// download one musl-dev package
			// list package contents with following commands
			//
			// $ ar -x musl-dev_1.2.2-3_riscv64.deb
			// $ tar -tvf data.tar.xz
			ARCH_RISCV_64: "riscv64-linux-musl",
		}[mArch]
	case "msvc":
		return map[string]string{
			// https://packages.debian.org/buster/mingw-w64-i686-dev
			// check list of files
			ARCH_X86: "i686-w64-mingw32",
			// https://packages.debian.org/buster/mingw-w64-x86-64-dev
			// check list of files
			ARCH_AMD64: "x86_64-w64-mingw32",
		}[mArch]
	default:
	}

	return map[string]string{
		ARCH_X86:   "i686-linux-gnu",
		ARCH_AMD64: "x86_64-linux-gnu",

		ARCH_ARM_V5: "arm-linux-gnueabi",
		ARCH_ARM_V6: "arm-linux-gnueabi",
		ARCH_ARM_V7: "arm-linux-gnueabihf",
		ARCH_ARM64:  "aarch64-linux-gnu",

		ARCH_PPC64:    "powerpc64-linux-gnu",
		ARCH_PPC64_LE: "powerpc64le-linux-gnu",

		ARCH_MIPS:         "mips-linux-gnu",
		ARCH_MIPS_HF:      "mips-linux-gnu",
		ARCH_MIPS_LE:      "mipsel-linux-gnu",
		ARCH_MIPS_LE_HF:   "mipsel-linux-gnu",
		ARCH_MIPS64:       "mips64-linux-gnuabi64",
		ARCH_MIPS64_HF:    "mips64-linux-gnuabi64",
		ARCH_MIPS64_LE:    "mips64el-linux-gnuabi64",
		ARCH_MIPS64_LE_HF: "mips64el-linux-gnuabi64",

		ARCH_RISCV_64: "riscv64-linux-gnu",
		ARCH_S390X:    "s390x-linux-gnu",
	}[mArch]
}
