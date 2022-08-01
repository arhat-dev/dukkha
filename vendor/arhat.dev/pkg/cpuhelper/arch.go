package cpuhelper

import (
	"runtime"
	"strconv"
	"unsafe"

	"arhat.dev/pkg/archconst"
	"golang.org/x/sys/cpu"
)

func LittleEndian() bool {
	var x uint32 = 0x01020304
	return *(*byte)(unsafe.Pointer(&x)) == 0x04
}

// Bits returns number of bits in a int for current platform
func Bits() int {
	return strconv.IntSize
}

func selectEndian(little, big archconst.ArchValue) archconst.ArchValue {
	if LittleEndian() {
		return little
	}

	return big
}

func fallback[S ~string](v, def S) S {
	if len(v) == 0 {
		return def
	}
	return v
}

// ArchByCPUFeatures detects runtime cpu arch (value defined in package arhat.dev/pkg/archconst)
// by inspecting cpu features
//
// return value MAY be empty string if detection failed or unable to determine specific arch value
// supported archs are:
// - amd64v{1,2,3,4}
// - arm64
// - armv{6, 7} (linux only)
// - s390x
// - ppc64{, le}v{8, 9}
// - mips64{, le}
func ArchByCPUFeatures(data CPU) archconst.ArchValue {
	switch {
	case cpu.ARM64.HasFP, cpu.ARM64.HasASIMD: // arm64
		return archconst.ARCH_ARM64
	case cpu.ARM.HasNEON:
		// neon can only be found in armv7 and later
		// only optional in cortex-a5,a7 when implementing vfpv4-d16
		// 				 in cortex-a9 whem implementing vfpv3-d16
		return archconst.ARCH_ARM_V7
	case cpu.ARM.HasTHUMBEE:
		if cpu.ARM.HasVFPv3 || cpu.ARM.HasVFPv3D16 || cpu.ARM.HasVFPD32 || cpu.ARM.HasVFPv4 {
			return archconst.ARCH_ARM_V7
		}

		// thumbee requires thumb-2, and thumb-2 can only be found in armv6 and later
		return archconst.ARCH_ARM_V6
		// it's hard to tell the differences between from this point,
		// leave armv5,armv6 unhandled
		// case cpu.ARM.HasVFP:
		// 	// vfp can only be found in armv5 and later
		// 	return archconst.ARCH_ARM_V5
	case cpu.S390X.HasZARCH:
		return archconst.ARCH_S390X
	case cpu.PPC64.IsPOWER9:
		return selectEndian(archconst.ARCH_PPC64_LE_V9, archconst.ARCH_PPC64_V9)
	case cpu.PPC64.IsPOWER8:
		return selectEndian(archconst.ARCH_PPC64_LE_V8, archconst.ARCH_PPC64_V8)
	case cpu.MIPS64X.HasMSA:
		return selectEndian(archconst.ARCH_MIPS64_LE, archconst.ARCH_MIPS64)
	}

	if data == nil {
		data = Detect()
	}

	switch t := data.(type) {
	case X86:
		switch t.MicroArch() {
		case 4:
			if runtime.GOOS != "darwin" {
				return archconst.ARCH_AMD64_V3
			}

			return archconst.ARCH_AMD64_V4
		case 3:
			return archconst.ARCH_AMD64_V3
		case 2:
			return archconst.ARCH_AMD64_V2
		case 1:
			return archconst.ARCH_AMD64_V1
		default:
			if t.Features.HasAll(X86Feature_SSE) {
				return archconst.ARCH_AMD64_V1
			}

			return archconst.ARCH_X86
		}
	default:
		return ""
	}
}
