package constant

func GetOciArch(mArch string) string {
	return GetGolangArch(mArch)
}

func GetOciArchVariant(mArch string) string {
	return map[string]string{
		ARCH_ARM_V5: "v5",
		ARCH_ARM_V6: "v6",
		ARCH_ARM_V7: "v7",

		ARCH_ARM64: "v8",
	}[mArch]
}
