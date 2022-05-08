package constant

import "arhat.dev/pkg/archconst"

// https://github.com/tpoechtrager/osxcross

// TODO: determine apple arch
func GetAppleArch(mArch string) (string, bool) {
	v, ok := map[string]string{
		archconst.ARCH_AMD64:    "x86_64",
		archconst.ARCH_AMD64_V1: "x86_64",
		archconst.ARCH_AMD64_V2: "x86_64",
		archconst.ARCH_AMD64_V3: "x86_64",
		archconst.ARCH_AMD64_V4: "x86_64",

		// arm64 for m1 chip, arm64e for A12 - before m1
		archconst.ARCH_ARM64: "arm64",

		archconst.ARCH_X86:    "",
		archconst.ARCH_X86_SF: "",

		archconst.ARCH_ARM_V5: "",
		archconst.ARCH_ARM_V6: "",
		archconst.ARCH_ARM_V7: "",

		archconst.ARCH_PPC:       "",
		archconst.ARCH_PPC_SF:    "",
		archconst.ARCH_PPC_LE:    "",
		archconst.ARCH_PPC_LE_SF: "",

		archconst.ARCH_PPC64:       "",
		archconst.ARCH_PPC64_LE:    "",
		archconst.ARCH_PPC64_V8:    "",
		archconst.ARCH_PPC64_LE_V8: "",
		archconst.ARCH_PPC64_V9:    "",
		archconst.ARCH_PPC64_LE_V9: "",

		archconst.ARCH_MIPS:         "",
		archconst.ARCH_MIPS_SF:      "",
		archconst.ARCH_MIPS_LE:      "",
		archconst.ARCH_MIPS_LE_SF:   "",
		archconst.ARCH_MIPS64:       "",
		archconst.ARCH_MIPS64_SF:    "",
		archconst.ARCH_MIPS64_LE:    "",
		archconst.ARCH_MIPS64_LE_SF: "",

		archconst.ARCH_RISCV_64: "",
		archconst.ARCH_S390X:    "",

		archconst.ARCH_IA64: "",
	}[mArch]

	return v, ok
}

func GetAppleTripleName(mArch, darwinVersion string) (string, bool) {
	v, ok := map[string]string{
		archconst.ARCH_AMD64:    "x86_64-apple-darwin",
		archconst.ARCH_AMD64_V1: "x86_64-apple-darwin",
		archconst.ARCH_AMD64_V2: "x86_64-apple-darwin",
		archconst.ARCH_AMD64_V3: "x86_64-apple-darwin",
		archconst.ARCH_AMD64_V4: "x86_64-apple-darwin",

		archconst.ARCH_ARM64: "arm64-apple-darwin",

		archconst.ARCH_X86:    "",
		archconst.ARCH_X86_SF: "",

		archconst.ARCH_ARM_V5: "",
		archconst.ARCH_ARM_V6: "",
		archconst.ARCH_ARM_V7: "",

		archconst.ARCH_PPC:       "",
		archconst.ARCH_PPC_SF:    "",
		archconst.ARCH_PPC_LE:    "",
		archconst.ARCH_PPC_LE_SF: "",

		archconst.ARCH_PPC64:       "",
		archconst.ARCH_PPC64_LE:    "",
		archconst.ARCH_PPC64_V8:    "",
		archconst.ARCH_PPC64_LE_V8: "",
		archconst.ARCH_PPC64_V9:    "",
		archconst.ARCH_PPC64_LE_V9: "",

		archconst.ARCH_MIPS:         "",
		archconst.ARCH_MIPS_SF:      "",
		archconst.ARCH_MIPS_LE:      "",
		archconst.ARCH_MIPS_LE_SF:   "",
		archconst.ARCH_MIPS64:       "",
		archconst.ARCH_MIPS64_SF:    "",
		archconst.ARCH_MIPS64_LE:    "",
		archconst.ARCH_MIPS64_LE_SF: "",

		archconst.ARCH_RISCV_64: "",
		archconst.ARCH_S390X:    "",

		archconst.ARCH_IA64: "",
	}[mArch]

	return v + darwinVersion, ok
}
