package constant

// https://github.com/tpoechtrager/osxcross

// TODO: determine apple arch
func GetAppleArch(mArch string) (string, bool) {
	v, ok := map[string]string{
		ARCH_AMD64: "x86_64",

		// arm64 for m1 chip, arm64e for A12 - before m1
		ARCH_ARM64: "arm64",

		ARCH_X86: "",

		ARCH_ARM_V5: "",
		ARCH_ARM_V6: "",
		ARCH_ARM_V7: "",

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

		ARCH_IA64: "",
	}[mArch]

	return v, ok
}

func GetAppleTripleName(mArch, darwinVersion string) (string, bool) {
	v, ok := map[string]string{
		ARCH_AMD64: "x86_64-apple-darwin",

		ARCH_ARM64: "arm64-apple-darwin",

		ARCH_X86: "",

		ARCH_ARM_V5: "",
		ARCH_ARM_V6: "",
		ARCH_ARM_V7: "",

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

		ARCH_IA64: "",
	}[mArch]

	return v + darwinVersion, ok
}
