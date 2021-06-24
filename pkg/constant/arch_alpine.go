package constant

func GetAlpineArch(mArch string) string {
	return map[string]string{
		ARCH_X86:   "x86",
		ARCH_AMD64: "x86_64",

		// ARCH_ARM_V5: "", // alpine has no armv5 support
		ARCH_ARM_V6: "armhf",
		ARCH_ARM_V7: "armv7",
		ARCH_ARM64:  "aarch64",

		ARCH_PPC64:    "ppc64",
		ARCH_PPC64_LE: "ppc64le",

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

func GetAlpineTripleName(mArch string) string {
	return map[string]string{
		ARCH_X86: "i686-linux-musl",
		// ARCH_AMD64: "",

		// ARCH_ARM_V5: "",
		ARCH_ARM_V6: "armel-linux-musleabi",
		ARCH_ARM_V7: "armv7l-linux-musleabihf",
		ARCH_ARM64:  "aarch64-linux-musl",

		ARCH_PPC64:    "powerpc64-linux-musl",
		ARCH_PPC64_LE: "powerpc64le-linux-musl",

		ARCH_MIPS:         "mips-linux-musl",
		ARCH_MIPS_HF:      "mips-linux-musl",
		ARCH_MIPS_LE:      "mipsel-linux-musl",
		ARCH_MIPS_LE_HF:   "mipsel-linux-musl",
		ARCH_MIPS64:       "mips64-linux-musl",
		ARCH_MIPS64_HF:    "mips64-linux-musl",
		ARCH_MIPS64_LE:    "mips64el-linux-musl",
		ARCH_MIPS64_LE_HF: "mips64el-linux-musl",

		ARCH_RISCV_64: "riscv64-linux-musl",
		ARCH_S390X:    "s390x-linux-musl",
	}[mArch]
}
