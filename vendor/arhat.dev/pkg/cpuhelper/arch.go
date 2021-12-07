package cpuhelper

import (
	"strconv"
	"unsafe"

	"arhat.dev/pkg/archconst"
	"golang.org/x/sys/cpu"
)

func LittleEndian() bool {
	var x uint32 = 0x01020304
	return *(*byte)(unsafe.Pointer(&x)) == 0x04
}

func Bits() int {
	return strconv.IntSize
}

func selectEndian(little, big string) string {
	if LittleEndian() {
		return little
	}

	return big
}

func fallback(v, def string) string {
	if len(v) == 0 {
		return def
	}
	return v
}

func ArchByCPUFeatures() string {
	if !cpu.Initialized {
		return ""
	}

	switch {
	case cpu.X86.HasSSE2: // amd64
		// Microarchitecture levels
		// Ref: https://en.wikipedia.org/wiki/X86-64
		switch {
		case cpu.X86.HasAVX512F,
			cpu.X86.HasAVX512VL,
			cpu.X86.HasAVX512DQ,
			cpu.X86.HasAVX512CD,
			cpu.X86.HasAVX512BW:
			return archconst.ARCH_AMD64_V4
		case cpu.X86.HasAVX2,
			cpu.X86.HasAVX,
			cpu.X86.HasBMI1,
			cpu.X86.HasBMI2,
			cpu.X86.HasFMA,
			cpu.X86.HasOSXSAVE:
			return archconst.ARCH_AMD64_V3
		case cpu.X86.HasCX16,
			cpu.X86.HasPOPCNT,
			cpu.X86.HasSSE3,
			cpu.X86.HasSSE41,
			cpu.X86.HasSSE42,
			cpu.X86.HasSSSE3:
			return archconst.ARCH_AMD64_V2
		}

		return archconst.ARCH_AMD64_V1
	case cpu.ARM64.HasFP, cpu.ARM64.HasASIMD: // arm64
		return archconst.ARCH_ARM64
	case cpu.MIPS64X.HasMSA:
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
	case cpu.S390X.HasZARCH:
		return archconst.ARCH_S390X
	// case cpu.ARM.HasVFP:
	// 	// vfp can only be found in armv5 and later
	// 	return archconst.ARCH_ARM_V5
	case cpu.PPC64.IsPOWER8:
		return selectEndian(archconst.ARCH_PPC64_V8_LE, archconst.ARCH_PPC64_V8)
	case cpu.PPC64.IsPOWER9:
		return selectEndian(archconst.ARCH_PPC64_V9_LE, archconst.ARCH_PPC64_V9)
	case cpu.MIPS64X.HasMSA:
		return selectEndian(archconst.ARCH_MIPS64_LE, archconst.ARCH_MIPS64)
	}

	return ""
}
