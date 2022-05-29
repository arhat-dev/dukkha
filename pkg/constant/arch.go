package constant

import (
	. "arhat.dev/pkg/archconst"
)

// SimpleArch generalizes cpu micro arch for practical use
//
// - amd64v{1, 2, 3, 4} => amd64
// - arm64v{8, 9} => arm64
// - ppc64{, le}v{8, 9} => ppc64{, le}
//
// Exceptions:
// - armv{5, 6, 7} are kept as is
// - arm => armv7
func SimpleArch(arch string) string {
	s, ok := Parse[byte](arch)
	if !ok {
		// unknown
		return arch
	}

	switch s.Name {
	case ARCH_ARM:
		if len(s.MicroArch) == 0 {
			s.MicroArch = "v7"
		}
	default:
		s.MicroArch = ""
	}

	return s.String()
}

// HardFloadArch returns hardfloat version of arch
func HardFloadArch(arch string) string {
	spec, ok := Parse[byte](arch)
	if !ok {
		// unknown, assume hardfloat arch
		return arch
	}

	spec.SoftFloat = false
	return spec.String()
}

// SoftFloadArch returns hardfloat version of arch
func SoftFloadArch(arch string) string {
	spec, ok := Parse[byte](arch)
	if !ok {
		// unknown, assume softfloat arch
		return arch
	}

	spec.SoftFloat = true
	return spec.String()
}

func CrossPlatform(
	targetKernel, targetArch,
	hostKernel, hostArch string,
) bool {
	var (
		host, target Spec

		ok bool
	)

	if hostKernel != targetKernel {
		return true
	}

	target, ok = Parse[byte](targetArch)
	if !ok {
		// is cross platform if not a exact match (for unknown arch)
		return targetArch != hostArch
	}

	host, ok = Parse[byte](hostArch)
	if !ok {
		// is cross platform if not a exact match (for unknown arch)
		return targetArch != hostArch
	}

	// check cpu compatibility
	// TODO: check micro arch compatibility?

	if host.Name != target.Name || host.LittleEndian != target.LittleEndian {
		return true
	}

	if host.SoftFloat {
		return !target.SoftFloat
	}

	// here we assume hardfloat host supports softfloat target
	return false
}

type archID uint32

// nolint:gocyclo
func (id archID) String() string {
	switch id {
	case archID_X86:
		return ARCH_X86
	case archID_X86_SF:
		return ARCH_X86_SF

	case archID_AMD64:
		return ARCH_AMD64
	case archID_AMD64_V1:
		return ARCH_AMD64_V1
	case archID_AMD64_V2:
		return ARCH_AMD64_V2
	case archID_AMD64_V3:
		return ARCH_AMD64_V3
	case archID_AMD64_V4:
		return ARCH_AMD64_V4

	case archID_ARM:
		return ARCH_ARM
	case archID_ARM_V5:
		return ARCH_ARM_V5
	case archID_ARM_V6:
		return ARCH_ARM_V6
	case archID_ARM_V7:
		return ARCH_ARM_V7

	case archID_ARM64:
		return ARCH_ARM64
	case archID_ARM64_V8:
		return ARCH_ARM64_V8
	case archID_ARM64_V9:
		return ARCH_ARM64_V9

	case archID_MIPS:
		return ARCH_MIPS
	case archID_MIPS_SF:
		return ARCH_MIPS_SF

	case archID_MIPS_LE:
		return ARCH_MIPS_LE
	case archID_MIPS_LE_SF:
		return ARCH_MIPS_LE_SF

	case archID_MIPS64:
		return ARCH_MIPS64
	case archID_MIPS64_SF:
		return ARCH_MIPS64_SF

	case archID_MIPS64_LE:
		return ARCH_MIPS64_LE
	case archID_MIPS64_LE_SF:
		return ARCH_MIPS64_LE_SF

	case archID_PPC:
		return ARCH_PPC
	case archID_PPC_SF:
		return ARCH_PPC_SF

	case archID_PPC_LE:
		return ARCH_PPC_LE
	case archID_PPC_LE_SF:
		return ARCH_PPC_LE_SF

	case archID_PPC64:
		return ARCH_PPC64
	case archID_PPC64_V8:
		return ARCH_PPC64_V8
	case archID_PPC64_V9:
		return ARCH_PPC64_V9

	case archID_PPC64_LE:
		return ARCH_PPC64_LE
	case archID_PPC64_LE_V8:
		return ARCH_PPC64_LE_V8
	case archID_PPC64_LE_V9:
		return ARCH_PPC64_LE_V9

	case archID_RISCV64:
		return ARCH_RISCV64

	case archID_S390X:
		return ARCH_S390X

	case archID_IA64:
		return ARCH_IA64

	default:
		return "<unknown>"
	}
}

const (
	_unknown_arch archID = iota

	archID_X86
	archID_X86_SF

	archID_AMD64
	archID_AMD64_V1
	archID_AMD64_V2
	archID_AMD64_V3
	archID_AMD64_V4

	archID_ARM
	archID_ARM_V5
	archID_ARM_V6
	archID_ARM_V7

	archID_ARM64
	archID_ARM64_V8
	archID_ARM64_V9

	archID_MIPS
	archID_MIPS_SF

	archID_MIPS_LE
	archID_MIPS_LE_SF

	archID_MIPS64
	archID_MIPS64_SF

	archID_MIPS64_LE
	archID_MIPS64_LE_SF

	archID_PPC
	archID_PPC_SF

	archID_PPC_LE
	archID_PPC_LE_SF

	archID_PPC64
	archID_PPC64_V8
	archID_PPC64_V9

	archID_PPC64_LE
	archID_PPC64_LE_V8
	archID_PPC64_LE_V9

	archID_RISCV64

	archID_S390X

	archID_IA64

	archID_COUNT
)

// nolint:gocyclo
func arch_id_of(arch string) archID {
	switch arch {
	case ARCH_X86:
		return archID_X86
	case ARCH_X86_SF:
		return archID_X86_SF
	case ARCH_AMD64:
		return archID_AMD64
	case ARCH_AMD64_V1:
		return archID_AMD64_V1
	case ARCH_AMD64_V2:
		return archID_AMD64_V2
	case ARCH_AMD64_V3:
		return archID_AMD64_V3
	case ARCH_AMD64_V4:
		return archID_AMD64_V4
	case ARCH_ARM:
		return archID_ARM
	case ARCH_ARM_V5:
		return archID_ARM_V5
	case ARCH_ARM_V6:
		return archID_ARM_V6
	case ARCH_ARM_V7:
		return archID_ARM_V7
	case ARCH_ARM64:
		return archID_ARM64
	case ARCH_ARM64_V8:
		return archID_ARM64_V8
	case ARCH_ARM64_V9:
		return archID_ARM64_V9
	case ARCH_MIPS:
		return archID_MIPS
	case ARCH_MIPS_SF:
		return archID_MIPS_SF
	case ARCH_MIPS_LE:
		return archID_MIPS_LE
	case ARCH_MIPS_LE_SF:
		return archID_MIPS_LE_SF
	case ARCH_MIPS64:
		return archID_MIPS64
	case ARCH_MIPS64_SF:
		return archID_MIPS64_SF
	case ARCH_MIPS64_LE:
		return archID_MIPS64_LE
	case ARCH_MIPS64_LE_SF:
		return archID_MIPS64_LE_SF
	case ARCH_PPC:
		return archID_PPC
	case ARCH_PPC_SF:
		return archID_PPC_SF
	case ARCH_PPC_LE:
		return archID_PPC_LE
	case ARCH_PPC_LE_SF:
		return archID_PPC_LE_SF
	case ARCH_PPC64:
		return archID_PPC64
	case ARCH_PPC64_V8:
		return archID_PPC64_V8
	case ARCH_PPC64_V9:
		return archID_PPC64_V9
	case ARCH_PPC64_LE:
		return archID_PPC64_LE
	case ARCH_PPC64_LE_V8:
		return archID_PPC64_LE_V8
	case ARCH_PPC64_LE_V9:
		return archID_PPC64_LE_V9
	case ARCH_RISCV64:
		return archID_RISCV64
	case ARCH_S390X:
		return archID_S390X
	case ARCH_IA64:
		return archID_IA64
	default:
		return _unknown_arch
	}
}
