package constant

func GetGolangArch(mArch string) string {
	return map[string]string{
		ARCH_X86:   "386",
		ARCH_AMD64: "amd64",
		ARCH_ARM64: "arm64",

		ARCH_ARM_V5: "arm",
		ARCH_ARM_V6: "arm",
		ARCH_ARM_V7: "arm",

		ARCH_MIPS:       "mips",
		ARCH_MIPS_HF:    "mips",
		ARCH_MIPS_LE:    "mipsle",
		ARCH_MIPS_LE_HF: "mipsle",

		ARCH_MIPS64:       "mips64",
		ARCH_MIPS64_HF:    "mips64",
		ARCH_MIPS64_LE:    "mips64le",
		ARCH_MIPS64_LE_HF: "mips64le",

		ARCH_PPC64:    "ppc64",
		ARCH_PPC64_LE: "ppc64le",

		ARCH_RISCV_64: "riscv64",
		ARCH_S390X:    "s390x",
		ARCH_IA64:     "ia64",
	}[mArch]
}
