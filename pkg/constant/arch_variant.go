package constant

import "arhat.dev/pkg/archconst"

// nolint:gocyclo
func GetDockerArchVariant(mArch string) (string, bool) {
	switch mArch {
	case archconst.ARCH_ARM:
		return "", true
	case archconst.ARCH_ARM_V5:
		return "v5", true
	case archconst.ARCH_ARM_V6:
		return "v6", true
	case archconst.ARCH_ARM_V7:
		return "v7", true

	case archconst.ARCH_ARM64:
		return "", true
	case archconst.ARCH_ARM64_V8:
		return "v8", true
	case archconst.ARCH_ARM64_V9:
		return "v9", true

	case archconst.ARCH_X86:
		return "", true
	case archconst.ARCH_X86_SF:
		return "", true

	case archconst.ARCH_AMD64:
		// TODO: set v1?
		return "", true
	case archconst.ARCH_AMD64_V1:
		return "v1", true
	case archconst.ARCH_AMD64_V2:
		return "v2", true
	case archconst.ARCH_AMD64_V3:
		return "v3", true
	case archconst.ARCH_AMD64_V4:
		return "v4", true

	case archconst.ARCH_PPC:
		return "", true
	case archconst.ARCH_PPC_SF:
		return "softfloat", true
	case archconst.ARCH_PPC_LE:
		return "", true
	case archconst.ARCH_PPC_LE_SF:
		return "softfloat", true

	case archconst.ARCH_PPC64:
		return "", true
	case archconst.ARCH_PPC64_LE:
		return "", true
	case archconst.ARCH_PPC64_V8:
		return "power8", true
	case archconst.ARCH_PPC64_LE_V8:
		return "power8", true
	case archconst.ARCH_PPC64_V9:
		return "power9", true
	case archconst.ARCH_PPC64_LE_V9:
		return "power9", true

	case archconst.ARCH_MIPS:
		return "", true
	case archconst.ARCH_MIPS_SF:
		return "softfloat", true
	case archconst.ARCH_MIPS_LE:
		return "", true
	case archconst.ARCH_MIPS_LE_SF:
		return "softfloat", true

	case archconst.ARCH_MIPS64:
		return "", true
	case archconst.ARCH_MIPS64_SF:
		return "softfloat", true
	case archconst.ARCH_MIPS64_LE:
		return "", true
	case archconst.ARCH_MIPS64_LE_SF:
		return "softfloat", true

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

// currently it's the same as docker arch variant
// nolint:gocyclo
func GetOciArchVariant(mArch string) (string, bool) {
	switch mArch {
	case archconst.ARCH_ARM:
		return "", true
	case archconst.ARCH_ARM_V5:
		return "v5", true
	case archconst.ARCH_ARM_V6:
		return "v6", true
	case archconst.ARCH_ARM_V7:
		return "v7", true

	case archconst.ARCH_ARM64:
		return "", true
	case archconst.ARCH_ARM64_V8:
		return "v8", true
	case archconst.ARCH_ARM64_V9:
		return "v9", true

	case archconst.ARCH_X86:
		return "", true
	case archconst.ARCH_X86_SF:
		return "", true

	case archconst.ARCH_AMD64:
		// TODO: set v1?
		return "", true
	case archconst.ARCH_AMD64_V1:
		return "v1", true
	case archconst.ARCH_AMD64_V2:
		return "v2", true
	case archconst.ARCH_AMD64_V3:
		return "v3", true
	case archconst.ARCH_AMD64_V4:
		return "v4", true

	case archconst.ARCH_PPC:
		return "", true
	case archconst.ARCH_PPC_SF:
		return "softfloat", true
	case archconst.ARCH_PPC_LE:
		return "", true
	case archconst.ARCH_PPC_LE_SF:
		return "softfloat", true

	case archconst.ARCH_PPC64:
		return "", true
	case archconst.ARCH_PPC64_LE:
		return "", true
	case archconst.ARCH_PPC64_V8:
		return "power8", true
	case archconst.ARCH_PPC64_LE_V8:
		return "power8", true
	case archconst.ARCH_PPC64_V9:
		return "power9", true
	case archconst.ARCH_PPC64_LE_V9:
		return "power9", true

	case archconst.ARCH_MIPS:
		return "", true
	case archconst.ARCH_MIPS_SF:
		return "softfloat", true
	case archconst.ARCH_MIPS_LE:
		return "", true
	case archconst.ARCH_MIPS_LE_SF:
		return "softfloat", true

	case archconst.ARCH_MIPS64:
		return "", true
	case archconst.ARCH_MIPS64_SF:
		return "softfloat", true
	case archconst.ARCH_MIPS64_LE:
		return "", true
	case archconst.ARCH_MIPS64_LE_SF:
		return "softfloat", true

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
