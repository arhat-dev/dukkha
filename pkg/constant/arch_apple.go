package constant

import "arhat.dev/pkg/archconst"

// https://github.com/tpoechtrager/osxcross

// TODO: determine apple arch
// nolint:gocyclo
func GetAppleArch(mArch string) (string, bool) {
	switch mArch {
	case archconst.ARCH_AMD64:
		return "x86_64", true
	case archconst.ARCH_AMD64_V1:
		return "x86_64", true
	case archconst.ARCH_AMD64_V2:
		return "x86_64", true
	case archconst.ARCH_AMD64_V3:
		return "x86_64", true
	case archconst.ARCH_AMD64_V4:
		return "x86_64", true

	// arm64 for m1 chip, arm64e for A12 - before m1
	case archconst.ARCH_ARM64:
		return "arm64", true
	case archconst.ARCH_ARM64_V8:
		return "arm64", true
	case archconst.ARCH_ARM64_V9:
		return "arm64", true

	case archconst.ARCH_X86:
		return "", true
	case archconst.ARCH_X86_SF:
		return "", true

	case archconst.ARCH_ARM:
		return "", true
	case archconst.ARCH_ARM_V5:
		return "", true
	case archconst.ARCH_ARM_V6:
		return "", true
	case archconst.ARCH_ARM_V7:
		return "", true

	case archconst.ARCH_PPC:
		return "", true
	case archconst.ARCH_PPC_SF:
		return "", true
	case archconst.ARCH_PPC_LE:
		return "", true
	case archconst.ARCH_PPC_LE_SF:
		return "", true

	case archconst.ARCH_PPC64:
		return "", true
	case archconst.ARCH_PPC64_LE:
		return "", true
	case archconst.ARCH_PPC64_V8:
		return "", true
	case archconst.ARCH_PPC64_LE_V8:
		return "", true
	case archconst.ARCH_PPC64_V9:
		return "", true
	case archconst.ARCH_PPC64_LE_V9:
		return "", true

	case archconst.ARCH_MIPS:
		return "", true
	case archconst.ARCH_MIPS_SF:
		return "", true
	case archconst.ARCH_MIPS_LE:
		return "", true
	case archconst.ARCH_MIPS_LE_SF:
		return "", true
	case archconst.ARCH_MIPS64:
		return "", true
	case archconst.ARCH_MIPS64_SF:
		return "", true
	case archconst.ARCH_MIPS64_LE:
		return "", true
	case archconst.ARCH_MIPS64_LE_SF:
		return "", true

	case archconst.ARCH_RISCV64:
		return "", true
	case archconst.ARCH_S390X:
		return "", true

	case archconst.ARCH_IA64:
		return "", true
	default:
		return "", false
	}
}

// nolint:gocyclo
func GetAppleTripleName(mArch, darwinVersion string) (string, bool) {
	switch mArch {
	case archconst.ARCH_AMD64:
		return "x86_64-apple-darwin" + darwinVersion, true
	case archconst.ARCH_AMD64_V1:
		return "x86_64-apple-darwin" + darwinVersion, true
	case archconst.ARCH_AMD64_V2:
		return "x86_64-apple-darwin" + darwinVersion, true
	case archconst.ARCH_AMD64_V3:
		return "x86_64-apple-darwin" + darwinVersion, true
	case archconst.ARCH_AMD64_V4:
		return "x86_64-apple-darwin" + darwinVersion, true

	case archconst.ARCH_ARM64:
		return "arm64-apple-darwin" + darwinVersion, true
	case archconst.ARCH_ARM64_V8:
		return "arm64-apple-darwin" + darwinVersion, true
	case archconst.ARCH_ARM64_V9:
		return "arm64-apple-darwin" + darwinVersion, true

	case archconst.ARCH_X86:
		return "", true
	case archconst.ARCH_X86_SF:
		return "", true

	case archconst.ARCH_ARM:
		return "", true
	case archconst.ARCH_ARM_V5:
		return "", true
	case archconst.ARCH_ARM_V6:
		return "", true
	case archconst.ARCH_ARM_V7:
		return "", true

	case archconst.ARCH_PPC:
		return "", true
	case archconst.ARCH_PPC_SF:
		return "", true
	case archconst.ARCH_PPC_LE:
		return "", true
	case archconst.ARCH_PPC_LE_SF:
		return "", true

	case archconst.ARCH_PPC64:
		return "", true
	case archconst.ARCH_PPC64_LE:
		return "", true
	case archconst.ARCH_PPC64_V8:
		return "", true
	case archconst.ARCH_PPC64_LE_V8:
		return "", true
	case archconst.ARCH_PPC64_V9:
		return "", true
	case archconst.ARCH_PPC64_LE_V9:
		return "", true

	case archconst.ARCH_MIPS:
		return "", true
	case archconst.ARCH_MIPS_SF:
		return "", true
	case archconst.ARCH_MIPS_LE:
		return "", true
	case archconst.ARCH_MIPS_LE_SF:
		return "", true
	case archconst.ARCH_MIPS64:
		return "", true
	case archconst.ARCH_MIPS64_SF:
		return "", true
	case archconst.ARCH_MIPS64_LE:
		return "", true
	case archconst.ARCH_MIPS64_LE_SF:
		return "", true

	case archconst.ARCH_RISCV64:
		return "", true
	case archconst.ARCH_S390X:
		return "", true

	case archconst.ARCH_IA64:
		return "", true
	default:
		return "", false
	}
}
