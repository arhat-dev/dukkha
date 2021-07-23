package constant

func GetQemuArch(mArch string) string {
	return map[string]string{
		ARCH_X86:   "i386",
		ARCH_AMD64: "x86_64",

		ARCH_ARM64: "aarch64",

		ARCH_ARM_V5: "arm",
		ARCH_ARM_V6: "arm",
		ARCH_ARM_V7: "arm",

		ARCH_MIPS:       "mips",
		ARCH_MIPS_SF:    "mips",
		ARCH_MIPS_LE:    "mipsel",
		ARCH_MIPS_LE_SF: "mipsel",

		ARCH_MIPS64:       "mips64",
		ARCH_MIPS64_SF:    "mips64",
		ARCH_MIPS64_LE:    "mips64el",
		ARCH_MIPS64_LE_SF: "mips64el",

		ARCH_PPC:      "ppc",
		ARCH_PPC64:    "ppc64",
		ARCH_PPC64_LE: "ppc64le",

		ARCH_S390X:    "s390x",
		ARCH_RISCV_64: "riscv64",
	}[mArch]
}
