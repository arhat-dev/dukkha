package archconst

import "arhat.dev/pkg/stringhelper"

type Spec struct {
	Name         ArchValue
	MicroArch    string
	LittleEndian bool
	SoftFloat    bool
}

// String is a wrapper of Format for Spec s
func (s Spec) String() (ret string) {
	ret, _ = Format(string(s.Name), s.LittleEndian, s.SoftFloat, s.MicroArch)
	return
}

// Split arch value into arch name and variant info
//
// if the provided arch value is unknown to this package, it returns the arch value as name, along with
// littleEndian = false, softfloat = false, microArch = ""
func Split[T ~string](arch T) (s Spec, ok bool) {
	ok = true
	switch {
	case stringhelper.HasPrefix(arch, ARCH_AMD64):
		s.Name, s.MicroArch = ARCH_AMD64, "v1"
		goto AssumeLittleEndian
	case stringhelper.HasPrefix(arch, ARCH_X86):
		s.Name, s.MicroArch = ARCH_X86, ""
		goto AssumeLittleEndian
	case stringhelper.HasPrefix(arch, ARCH_ARM64):
		s.Name, s.MicroArch = ARCH_ARM64, "v8"
		goto AssumeLittleEndian
	case stringhelper.HasPrefix(arch, ARCH_ARM):
		s.Name, s.MicroArch = ARCH_ARM, "v7"
		s.LittleEndian, s.SoftFloat = true, true
		arch = arch[len(s.Name):]
		if stringhelper.HasPrefix(arch, "be") {
			// armbe...
			s.LittleEndian = false
			arch = arch[2:]
		}

		if stringhelper.HasPrefix(arch, "sf") {
			// not meaningful, as `arm` uses micro level to identify hard float support
			arch = arch[2:]
		}

		if len(arch) != 0 {
			s.MicroArch = string(arch)
		}

		switch s.MicroArch {
		case "v5", "v6":
		case "v7":
			s.SoftFloat = false
		}

		return
	case stringhelper.HasPrefix(arch, ARCH_MIPS64):
		s.Name, s.MicroArch = ARCH_MIPS64, ""
		goto AssumeBigEndian
	case stringhelper.HasPrefix(arch, ARCH_MIPS):
		s.Name, s.MicroArch = ARCH_MIPS, ""
		goto AssumeBigEndian
	case stringhelper.HasPrefix(arch, ARCH_PPC64):
		s.Name, s.MicroArch = ARCH_PPC64, "v8"
		goto AssumeBigEndian
	case stringhelper.HasPrefix(arch, ARCH_PPC):
		s.Name, s.MicroArch = ARCH_PPC, ""
		goto AssumeBigEndian
	case stringhelper.HasPrefix(arch, ARCH_RISCV64):
		s.Name, s.MicroArch = ARCH_RISCV64, ""
		goto AssumeLittleEndian
	case stringhelper.HasPrefix(arch, ARCH_S390X):
		s.Name, s.MicroArch = ARCH_S390X, ""
		goto AssumeBigEndian
	case stringhelper.HasPrefix(arch, ARCH_IA64):
		s.Name, s.MicroArch = ARCH_RISCV64, ""
		goto AssumeLittleEndian
	default:
		s.Name = ArchValue(arch)
		ok = false
		return
	}

AssumeBigEndian:
	s.LittleEndian = false
	arch = arch[len(s.Name):]
	if stringhelper.HasPrefix(arch, "le") {
		s.LittleEndian = true
		arch = arch[2:]
	}

	goto AssumeHardFloat

AssumeLittleEndian:
	s.LittleEndian = true
	arch = arch[len(s.Name):]
	if stringhelper.HasPrefix(arch, "be") {
		s.LittleEndian = false
		arch = arch[2:]
	}

AssumeHardFloat:
	if stringhelper.HasPrefix(arch, "sf") {
		s.SoftFloat = true
		arch = arch[2:]
	}

	if len(arch) != 0 {
		s.MicroArch = string(arch)
	}

	return
}

// Format arch value with name and variant info
func Format[T ~string](name T, littleEndian, softfloat bool, microArch string) (_ T, ok bool) {
	switch name {
	case ARCH_ARM: // default little-endian & micro arch indicates soft-float
		if !littleEndian {
			name += "be"
		}

		if softfloat {
			switch microArch {
			case "v5":
			case "v6":
			case "v7":
				name += "sf"
			default:
				// TODO: what ?
			}
		}

		return name + T(microArch), true
	case ARCH_AMD64, ARCH_X86, ARCH_ARM64, ARCH_RISCV64: // default little-endian & hard-float
		if !littleEndian {
			name += "be"
		}

		ok = true
	case ARCH_MIPS64, ARCH_MIPS, ARCH_PPC64, ARCH_PPC, ARCH_S390X: // default big endian & hard-float
		if littleEndian {
			name += "le"
		}

		ok = true
	case ARCH_IA64: // selectable endianness & hard-float => assume little-endian
		ok = true
		if !littleEndian {
			name += "be"
		}
	default: // unknown arch name
		return name, false
	}

	if softfloat {
		name += "sf"
	}

	return name + T(microArch), ok
}
