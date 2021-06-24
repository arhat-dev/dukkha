package constant

func GetDebianArch(mArch string) string {
	return map[string]string{
		ARCH_X86:   "i386",
		ARCH_AMD64: "amd64",

		// ARCH_ARM_V5: "",
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

func GetDebianTripleName(mArch string) string {
	return map[string]string{
		ARCH_X86: "i686-linux-gnu",
		// ARCH_AMD64: "amd64",

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
		ARCH_MIPS64:       "mips64-linux-gnu",
		ARCH_MIPS64_HF:    "mips64-linux-gnu",
		ARCH_MIPS64_LE:    "mips64el-linux-gnu",
		ARCH_MIPS64_LE_HF: "mips64el-linux-gnu",

		ARCH_RISCV_64: "riscv64-linux-gnu",
		ARCH_S390X:    "s390x-linux-gnu",
	}[mArch]
}
