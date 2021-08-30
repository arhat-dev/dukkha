package constant

func GetGolangArch(mArch string) (string, bool) {
	v, ok := map[string]string{
		ARCH_X86: "386",

		ARCH_AMD64: "amd64",

		ARCH_ARM_V5: "arm",
		ARCH_ARM_V6: "arm",
		ARCH_ARM_V7: "arm",

		ARCH_ARM64: "arm64",

		ARCH_MIPS:         "mips",
		ARCH_MIPS_SF:      "mips",
		ARCH_MIPS_LE:      "mipsle",
		ARCH_MIPS_LE_SF:   "mipsle",
		ARCH_MIPS64:       "mips64",
		ARCH_MIPS64_SF:    "mips64",
		ARCH_MIPS64_LE:    "mips64le",
		ARCH_MIPS64_LE_SF: "mips64le",

		ARCH_PPC:       "ppc",
		ARCH_PPC_SF:    "ppc",
		ARCH_PPC_LE:    "ppcle",
		ARCH_PPC_LE_SF: "ppcle",
		ARCH_PPC64:     "ppc64",
		ARCH_PPC64_LE:  "ppc64le",

		ARCH_RISCV_64: "riscv64",
		ARCH_S390X:    "s390x",
		ARCH_IA64:     "ia64",
	}[mArch]

	return v, ok
}
