package constant

func GetDockerArchVariant(mArch string) (string, bool) {
	v, ok := map[string]string{
		ARCH_ARM_V5: "v5",
		ARCH_ARM_V6: "v6",
		ARCH_ARM_V7: "v7",
		ARCH_ARM64:  "v8",

		ARCH_X86:   "",
		ARCH_AMD64: "",

		ARCH_PPC:       "",
		ARCH_PPC_SF:    "softfloat",
		ARCH_PPC_LE:    "",
		ARCH_PPC_LE_SF: "softfloat",

		ARCH_PPC64:    "",
		ARCH_PPC64_LE: "",

		ARCH_MIPS:         "",
		ARCH_MIPS_SF:      "softfloat",
		ARCH_MIPS_LE:      "",
		ARCH_MIPS_LE_SF:   "softfloat",
		ARCH_MIPS64:       "",
		ARCH_MIPS64_SF:    "softfloat",
		ARCH_MIPS64_LE:    "",
		ARCH_MIPS64_LE_SF: "softfloat",

		ARCH_RISCV_64: "",
		ARCH_S390X:    "",

		ARCH_IA64: "",
	}[mArch]

	return v, ok
}

func GetOciArchVariant(mArch string) (string, bool) {
	v, ok := map[string]string{
		ARCH_ARM_V5: "v5",
		ARCH_ARM_V6: "v6",
		ARCH_ARM_V7: "v7",

		ARCH_ARM64: "v8",

		ARCH_X86:   "",
		ARCH_AMD64: "",

		ARCH_PPC:       "",
		ARCH_PPC_SF:    "softfloat",
		ARCH_PPC_LE:    "",
		ARCH_PPC_LE_SF: "softfloat",

		ARCH_PPC64:    "",
		ARCH_PPC64_LE: "",

		ARCH_MIPS:         "",
		ARCH_MIPS_SF:      "softfloat",
		ARCH_MIPS_LE:      "",
		ARCH_MIPS_LE_SF:   "softfloat",
		ARCH_MIPS64:       "",
		ARCH_MIPS64_SF:    "softfloat",
		ARCH_MIPS64_LE:    "",
		ARCH_MIPS64_LE_SF: "softfloat",

		ARCH_RISCV_64: "",
		ARCH_S390X:    "",

		ARCH_IA64: "",
	}[mArch]

	return v, ok
}