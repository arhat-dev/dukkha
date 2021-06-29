package constant

// https://github.com/tpoechtrager/osxcross

func GetAppleArch(mArch string) string {
	return map[string]string{
		ARCH_AMD64: "x86_64",

		// arm64 for m1 chip, arm64e for A12 - before m1
		ARCH_ARM64: "arm64",
	}[mArch]
}

func GetAppleTripleName(mArch, darwinVersion string) string {
	v := map[string]string{
		ARCH_AMD64: "x86_64-apple-darwin",

		ARCH_ARM64: "arm64-apple-darwin",
	}[mArch]

	return v + darwinVersion
}
