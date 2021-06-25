package constant

func GetDockerArch(mArch string) string {
	return GetGolangArch(mArch)
}

func GetDockerArchVariant(mArch string) string {
	return map[string]string{
		ARCH_ARM_V5: "5",
		ARCH_ARM_V6: "6",
		ARCH_ARM_V7: "7",
	}[mArch]
}

func GetDockerHubArch(mArch string, other ...string) string {
	mKernel := KERNEL_LINUX
	if len(other) != 0 {
		mKernel = other[0]
	}

	arch := map[string]string{
		ARCH_AMD64: "amd64",
		ARCH_X86:   "i386",

		ARCH_ARM_V5: "arm32v5",
		ARCH_ARM_V6: "arm32v6",
		ARCH_ARM_V7: "arm32v7",

		ARCH_ARM64: "arm64v8",

		ARCH_MIPS64_LE:    "mips64le",
		ARCH_MIPS64_LE_HF: "mips64le",

		ARCH_PPC64_LE: "ppc64le",

		ARCH_S390X: "s390x",
	}[mArch]

	switch mKernel {
	case KERNEL_LINUX:
		return arch
	case KERNEL_WINDOWS:
		return "win" + arch
	default:
		return arch
	}
}
