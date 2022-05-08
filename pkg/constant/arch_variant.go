package constant

import "arhat.dev/pkg/archconst"

func GetDockerArchVariant(mArch string) (string, bool) {
	v, ok := map[string]string{
		archconst.ARCH_ARM_V5: "v5",
		archconst.ARCH_ARM_V6: "v6",
		archconst.ARCH_ARM_V7: "v7",
		archconst.ARCH_ARM64:  "v8",

		archconst.ARCH_X86:    "",
		archconst.ARCH_X86_SF: "",

		archconst.ARCH_AMD64:    "",
		archconst.ARCH_AMD64_V1: "v1",
		archconst.ARCH_AMD64_V2: "v2",
		archconst.ARCH_AMD64_V3: "v3",
		archconst.ARCH_AMD64_V4: "v4",

		archconst.ARCH_PPC:       "",
		archconst.ARCH_PPC_SF:    "softfloat",
		archconst.ARCH_PPC_LE:    "",
		archconst.ARCH_PPC_LE_SF: "softfloat",

		archconst.ARCH_PPC64:       "",
		archconst.ARCH_PPC64_LE:    "",
		archconst.ARCH_PPC64_V8:    "power8",
		archconst.ARCH_PPC64_LE_V8: "power8",
		archconst.ARCH_PPC64_V9:    "power9",
		archconst.ARCH_PPC64_LE_V9: "power9",

		archconst.ARCH_MIPS:         "",
		archconst.ARCH_MIPS_SF:      "softfloat",
		archconst.ARCH_MIPS_LE:      "",
		archconst.ARCH_MIPS_LE_SF:   "softfloat",
		archconst.ARCH_MIPS64:       "",
		archconst.ARCH_MIPS64_SF:    "softfloat",
		archconst.ARCH_MIPS64_LE:    "",
		archconst.ARCH_MIPS64_LE_SF: "softfloat",

		archconst.ARCH_RISCV_64: "",
		archconst.ARCH_S390X:    "",

		archconst.ARCH_IA64: "",
	}[mArch]

	return v, ok
}

func GetOciArchVariant(mArch string) (string, bool) {
	v, ok := map[string]string{
		archconst.ARCH_ARM_V5: "v5",
		archconst.ARCH_ARM_V6: "v6",
		archconst.ARCH_ARM_V7: "v7",

		archconst.ARCH_ARM64: "v8",

		archconst.ARCH_X86:    "",
		archconst.ARCH_X86_SF: "",

		archconst.ARCH_AMD64:    "",
		archconst.ARCH_AMD64_V1: "v1",
		archconst.ARCH_AMD64_V2: "v2",
		archconst.ARCH_AMD64_V3: "v3",
		archconst.ARCH_AMD64_V4: "v4",

		archconst.ARCH_PPC:       "",
		archconst.ARCH_PPC_SF:    "softfloat",
		archconst.ARCH_PPC_LE:    "",
		archconst.ARCH_PPC_LE_SF: "softfloat",

		archconst.ARCH_PPC64:       "",
		archconst.ARCH_PPC64_LE:    "",
		archconst.ARCH_PPC64_V8:    "power8",
		archconst.ARCH_PPC64_LE_V8: "power8",
		archconst.ARCH_PPC64_V9:    "power9",
		archconst.ARCH_PPC64_LE_V9: "power9",

		archconst.ARCH_MIPS:         "",
		archconst.ARCH_MIPS_SF:      "softfloat",
		archconst.ARCH_MIPS_LE:      "",
		archconst.ARCH_MIPS_LE_SF:   "softfloat",
		archconst.ARCH_MIPS64:       "",
		archconst.ARCH_MIPS64_SF:    "softfloat",
		archconst.ARCH_MIPS64_LE:    "",
		archconst.ARCH_MIPS64_LE_SF: "softfloat",

		archconst.ARCH_RISCV_64: "",
		archconst.ARCH_S390X:    "",

		archconst.ARCH_IA64: "",
	}[mArch]

	return v, ok
}
