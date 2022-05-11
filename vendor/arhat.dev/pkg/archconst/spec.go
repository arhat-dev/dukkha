package archconst

import "arhat.dev/pkg/stringhelper"

// Spec of a cpu arch
type Spec struct {
	Name         ArchValue
	MicroArch    string
	LittleEndian bool
	SoftFloat    bool
}

// String is a wrapper of Format for Spec s
func (s Spec) String() (ret string) {
	ret, _ = Format[string, byte](s.Name, s.LittleEndian, s.SoftFloat, s.MicroArch)
	return
}

// Parse arch value into arch Spec
//
// if the provided arch value is unknown to this package, it returns the arch value as name, along with
// littleEndian = false, softfloat = false, microArch = ""
func Parse[B ~byte, T stringhelper.String[B]](arch T) (s Spec, ok bool) {
	ok = true
	switch {
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_AMD64):
		s.Name, s.MicroArch = ARCH_AMD64, "v1"
		goto AssumeLittleEndian
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_X86):
		s.Name, s.MicroArch = ARCH_X86, ""
		goto AssumeLittleEndian
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_ARM64):
		s.Name, s.MicroArch = ARCH_ARM64, "v8"
		goto AssumeLittleEndian
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_ARM):
		s.Name, s.MicroArch = ARCH_ARM, "v7"
		s.LittleEndian, s.SoftFloat = true, true
		arch = stringhelper.SliceStart[B](arch, len(s.Name))
		if stringhelper.HasPrefix[B, byte, T, string](arch, "be") {
			// armbe...
			s.LittleEndian = false
			arch = stringhelper.SliceStart[B](arch, 2)
		}

		if stringhelper.HasPrefix[B, byte, T, string](arch, "sf") {
			// not meaningful, as `arm` uses micro level to identify hard float support
			arch = stringhelper.SliceStart[B](arch, 2)
		}

		if len(arch) != 0 {
			s.MicroArch = stringhelper.Convert[string, B](arch)
		}

		switch s.MicroArch {
		case "v5", "v6":
		case "v7":
			s.SoftFloat = false
		}

		return
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_MIPS64):
		s.Name, s.MicroArch = ARCH_MIPS64, ""
		goto AssumeBigEndian
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_MIPS):
		s.Name, s.MicroArch = ARCH_MIPS, ""
		goto AssumeBigEndian
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_PPC64):
		s.Name, s.MicroArch = ARCH_PPC64, "v8"
		goto AssumeBigEndian
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_PPC):
		s.Name, s.MicroArch = ARCH_PPC, ""
		goto AssumeBigEndian
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_RISCV64):
		s.Name, s.MicroArch = ARCH_RISCV64, ""
		goto AssumeLittleEndian
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_S390X):
		s.Name, s.MicroArch = ARCH_S390X, ""
		goto AssumeBigEndian
	case stringhelper.HasPrefix[B, byte, T, string](arch, ARCH_IA64):
		s.Name, s.MicroArch = ARCH_RISCV64, ""
		goto AssumeLittleEndian
	default:
		s.Name = stringhelper.Convert[ArchValue, B](arch)
		ok = false
		return
	}

AssumeBigEndian:
	s.LittleEndian = false
	arch = stringhelper.SliceStart[B](arch, len(s.Name))
	if stringhelper.HasPrefix[B, byte, T, string](arch, "le") {
		s.LittleEndian = true
		arch = stringhelper.SliceStart[B](arch, 2)
	}

	goto AssumeHardFloat

AssumeLittleEndian:
	s.LittleEndian = true
	arch = stringhelper.SliceStart[B](arch, len(s.Name))
	if stringhelper.HasPrefix[B, byte, T, string](arch, "be") {
		s.LittleEndian = false
		arch = stringhelper.SliceStart[B](arch, 2)
	}

AssumeHardFloat:
	if stringhelper.HasPrefix[B, byte, T, string](arch, "sf") {
		s.SoftFloat = true
		arch = stringhelper.SliceStart[B](arch, 2)
	}

	if len(arch) != 0 {
		s.MicroArch = stringhelper.Convert[string, B](arch)
	}

	return
}

// Format arch value with name and variant info
func Format[R ~string, B ~byte, T stringhelper.String[B]](name T, littleEndian, softfloat bool, microArch string) (_ R, ok bool) {
	nameStr := stringhelper.Convert[R, B](name)

	switch nameStr {
	case ARCH_ARM: // default little-endian & micro arch indicates soft-float
		if !littleEndian {
			nameStr += "be"
		}

		if softfloat {
			switch microArch {
			case "v5":
			case "v6":
			case "v7":
				nameStr += "sf"
			default:
				// TODO: what ?
			}
		}

		return nameStr + R(microArch), true
	case ARCH_AMD64, ARCH_X86, ARCH_ARM64, ARCH_RISCV64: // default little-endian & hard-float
		if !littleEndian {
			nameStr += "be"
		}

		ok = true
	case ARCH_MIPS64, ARCH_MIPS, ARCH_PPC64, ARCH_PPC, ARCH_S390X: // default big endian & hard-float
		if littleEndian {
			nameStr += "le"
		}

		ok = true
	case ARCH_IA64: // selectable endianness & hard-float => assume little-endian
		ok = true
		if !littleEndian {
			nameStr += "be"
		}
	default: // unknown arch name
		return nameStr, false
	}

	if softfloat {
		nameStr += "sf"
	}

	return nameStr + R(microArch), ok
}
