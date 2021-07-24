package constant

func GetAlpineArch(mArch string) (string, bool) {
	v, ok := map[string]string{
		ARCH_X86:   "x86",
		ARCH_AMD64: "x86_64",

		ARCH_ARM_V5: "armv5l",
		ARCH_ARM_V6: "armhf",
		ARCH_ARM_V7: "armv7",
		ARCH_ARM64:  "aarch64",

		ARCH_PPC:       "ppc",
		ARCH_PPC_SF:    "ppc",
		ARCH_PPC_LE:    "ppcle",
		ARCH_PPC_LE_SF: "ppcle",

		ARCH_PPC64:    "ppc64",
		ARCH_PPC64_LE: "ppc64le",

		ARCH_MIPS:         "mips",
		ARCH_MIPS_SF:      "mips",
		ARCH_MIPS_LE:      "mipsel",
		ARCH_MIPS_LE_SF:   "mipsel",
		ARCH_MIPS64:       "mips64",
		ARCH_MIPS64_SF:    "mips64",
		ARCH_MIPS64_LE:    "mips64el",
		ARCH_MIPS64_LE_SF: "mips64el",

		ARCH_RISCV_64: "riscv64",
		ARCH_S390X:    "s390x",

		ARCH_IA64: "",
	}[mArch]

	return v, ok
}

// GetAlpineTripleName of matrix arch
// reference: https://more.musl.cc/10/x86_64-linux-musl/
func GetAlpineTripleName(mArch string) (string, bool) {
	v, ok := map[string]string{
		ARCH_X86:   "i686-linux-musl",
		ARCH_AMD64: "x86_64-linux-musl",

		ARCH_ARM_V5: "armv5l-linux-musleabi",
		ARCH_ARM_V6: "armv6-linux-musleabi",
		ARCH_ARM_V7: "armv7l-linux-musleabihf",
		ARCH_ARM64:  "aarch64-linux-musl",

		ARCH_PPC:       "powerpc-linux-musl",
		ARCH_PPC_SF:    "powerpc-linux-muslsf",
		ARCH_PPC_LE:    "powerpcle-linux-musl",
		ARCH_PPC_LE_SF: "powerpcle-linux-muslsf",

		ARCH_PPC64:    "powerpc64-linux-musl",
		ARCH_PPC64_LE: "powerpc64le-linux-musl",

		ARCH_MIPS:         "mips-linux-musl",
		ARCH_MIPS_SF:      "mips-linux-muslsf",
		ARCH_MIPS_LE:      "mipsel-linux-musl",
		ARCH_MIPS_LE_SF:   "mipsel-linux-muslsf",
		ARCH_MIPS64:       "mips64-linux-musl",
		ARCH_MIPS64_SF:    "mips64-linux-musln32sf",
		ARCH_MIPS64_LE:    "mips64el-linux-musl",
		ARCH_MIPS64_LE_SF: "mips64el-linux-musln32sf",

		ARCH_RISCV_64: "riscv64-linux-musl",
		ARCH_S390X:    "s390x-linux-musl",

		ARCH_IA64: "",
	}[mArch]

	return v, ok
}
