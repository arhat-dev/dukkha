//go:build amd64 || 386 || amd64p32

package cpuhelper

import (
	"unsafe"
)

// ret: eax, ebx, ecx, edx
func cpuid(eax, ecx uint32) [4]uint32

// ret: eax, edx
func xgetbv(arg0 uint32) [2]uint32

func detect() (ret X86) {
	vnd, maxFuncNum := leaf0(&ret)
	if maxFuncNum < 1 {
		return
	}

	leaves := [...]func(*X86, Vendor){
		leaf1,
		leaf2,
		leaf3,
		leaf4,
		nil, // 5
		leaf6,
		leaf7,
	}

	for i := 0; i < len(leaves); i++ {
		fn := leaves[i]
		if fn == nil {
			continue
		}

		if maxFuncNum < uint32(i+1) {
			break
		}

		fn(&ret, vnd)
	}

	maxFuncNum = leaf0x80000000() & 0x7FFF_FFFF

	extLeaves := [...]func(*X86, Vendor){
		leaf0x8000_0001,
		nil, // 2
		nil, // 3
		leaf0x8000_0004,
		leaf0x8000_0005,
		leaf0x8000_0006,
		// leaf0x8000_001f,
	}

	for i := 0; i < len(extLeaves); i++ {
		fn := extLeaves[i]
		if fn == nil {
			continue
		}

		if maxFuncNum < uint32(i+1) {
			break
		}

		fn(&ret, vnd)
	}

	return
}

// leaf0 (eax=0): Highest Function Parameter and Manufacturer ID
func leaf0(cpu *X86) (vnd Vendor, maxFuncNum uint32) {
	v := cpuid(0, 0)
	// brand stored in ebx, edx, ecx (notice the order)
	cpu.Brand = string(
		unsafe.Slice(
			(*byte)(unsafe.Pointer(&v[1])),
			unsafe.Sizeof(v[1]),
		),
	) + string(
		unsafe.Slice(
			(*byte)(unsafe.Pointer(&v[3])),
			unsafe.Sizeof(v[3]),
		),
	) + string(
		unsafe.Slice(
			(*byte)(unsafe.Pointer(&v[2])),
			unsafe.Sizeof(v[2]),
		),
	)

	return cpu.Vendor(), v[0]
}

// leaf1 (eax=1): Processor Info and Feature Bits
func leaf1(cpu *X86, vnd Vendor) {
	v := cpuid(1, 0)
	eax, ebx := v[0], v[1]
	cpu.Stepping = X86Stepping(eax & 0xF)
	cpu.Model = X86Model((eax >> 4) & 0xF)
	cpu.Family = X86Family((eax >> 8) & 0xF)
	cpu.ProcessType = X86ProcessType((eax >> 12) & 0x3)

	if cpu.Family == 0x6 || cpu.Family == 0xF {
		cpu.Model |= X86Model((eax>>16)&0xF) << 4
	}

	if cpu.Family == 0xF {
		cpu.Family += X86Family((eax >> 20) & 0xFF)
	}

	cpu.BrandIndex = uint8(ebx & 0xFF)
	cpu.CacheLineSize = uint16((ebx>>8)&0xFF) << 3
	cpu.MaxLogicalCPUID = uint8((ebx >> 16) & 0xFF)
	cpu.InitialAPICID = uint8((ebx >> 24))

	cpu.Features = *(*X86Feature)(unsafe.Pointer(&v[2]))
}

// leaf2 (eax=2): Cache and TLB Descriptor information
func leaf2(cpu *X86, vnd Vendor) {
	if vnd != Vendor_Intel {
		return
	}

	v := cpuid(2, 0)
	sz := int(unsafe.Sizeof(uint32(0))) * len(v)
	data := unsafe.Slice(
		(*byte)(unsafe.Pointer(&v[0])),
		sz,
	)

	for i := 0; i < sz; i++ {
		if (i%4 == 0) && (data[i+3]&(1<<7) != 0) {
			i += 4
			continue
		}

		if data[i] == 0xFF {
			// use leaf4 for cache info
			cpu.CacheDescriptors = cpu.CacheDescriptors[0:0]
			break
		}

		cpu.CacheDescriptors = append(cpu.CacheDescriptors, leaf2GetDescriptor(int16(data[i])))
	}
}

func leaf2GetDescriptor(d int16) X86CacheDescriptor {
	switch d {
	case 0x01:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB", 4, 4, -1, 32, 0}
	case 0x02:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB", 4 * 1024, 0xFF, -1, 2, 0}
	case 0x03:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB", 4, 4, -1, 64, 0}
	case 0x04:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB", 4 * 1024, 4, -1, 8, 0}
	case 0x05:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB1", 4 * 1024, 4, -1, 32, 0}
	case 0x06:
		return X86CacheDescriptor{1, X86CacheType_INSTRUCTION_CACHE, "1st-level instruction cache", 8, 4, 32, -1, 0}
	case 0x08:
		return X86CacheDescriptor{1, X86CacheType_INSTRUCTION_CACHE, "1st-level instruction cache", 16, 4, 32, -1, 0}
	case 0x09:
		return X86CacheDescriptor{1, X86CacheType_INSTRUCTION_CACHE, "1st-level instruction cache", 32, 4, 64, -1, 0}
	case 0x0A:
		return X86CacheDescriptor{1, X86CacheType_DATA_CACHE, "1st-level data cache", 8, 2, 32, -1, 0}
	case 0x0B:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB", 4 * 1024, 4, -1, 4, 0}
	case 0x0C:
		return X86CacheDescriptor{1, X86CacheType_DATA_CACHE, "1st-level data cache", 16, 4, 32, -1, 0}
	case 0x0D:
		return X86CacheDescriptor{1, X86CacheType_DATA_CACHE, "1st-level data cache", 16, 4, 64, -1, 0}
	case 0x0E:
		return X86CacheDescriptor{1, X86CacheType_DATA_CACHE, "1st-level data cache", 24, 6, 64, -1, 0}
	case 0x1D:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 128, 2, 64, -1, 0}
	case 0x21:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 256, 8, 64, -1, 0}
	case 0x22:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 512, 4, 64, -1, 2}
	case 0x23:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 1 * 1024, 8, 64, -1, 2}
	case 0x24:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 1 * 1024, 16, 64, -1, 0}
	case 0x25:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 2 * 1024, 8, 64, -1, 2}
	case 0x29:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "2nd-level cache", 4 * 1024, 8, 64, -1, 2}
	case 0x2C:
		return X86CacheDescriptor{1, X86CacheType_DATA_CACHE, "1st-level cache", 32, 8, 64, -1, 0}
	case 0x30:
		return X86CacheDescriptor{1, X86CacheType_INSTRUCTION_CACHE, "1st-level instruction cache", 32, 8, 64, -1, 0}
	case 0x40:
		return X86CacheDescriptor{
			-1, X86CacheType_DATA_CACHE,
			"No 2nd-level cache or, if processor contains a valid 2nd-level cache, no 3rd-level cache",
			-1, -1, -1, -1, 0,
		}
	case 0x41:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 128, 4, 32, -1, 0}
	case 0x42:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 256, 4, 32, -1, 0}
	case 0x43:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 512, 4, 32, -1, 0}
	case 0x44:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 1 * 1024, 4, 32, -1, 0}
	case 0x45:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 2 * 1024, 4, 32, -1, 0}
	case 0x46:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 4 * 1024, 4, 64, -1, 0}
	case 0x47:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 8 * 1024, 8, 64, -1, 0}
	case 0x48:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 3 * 1024, 12, 64, -1, 0}
	case 0x49:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 4 * 1024, 16, 64, -1, 0}
	// (Intel Xeon processor MP, Family 0FH, Model 06H)
	case (0x49 | (1 << 8)):
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 4 * 1024, 16, 64, -1, 0}
	case 0x4A:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 6 * 1024, 12, 64, -1, 0}
	case 0x4B:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 8 * 1024, 16, 64, -1, 0}
	case 0x4C:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 12 * 1024, 12, 64, -1, 0}
	case 0x4D:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 16 * 1024, 16, 64, -1, 0}
	case 0x4E:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "3nd-level cache", 6 * 1024, 24, 64, -1, 0}
	case 0x4F:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB", 4, -1, -1, 32, 0}
	case 0x50:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB: 4 KByte and 2-MByte or 4-MByte pages", 4, -1, -1, 64, 0}
	case 0x51:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB: 4 KByte and 2-MByte or 4-MByte pages", 4, -1, -1, 128, 0}
	case 0x52:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB: 4 KByte and 2-MByte or 4-MByte pages", 4, -1, -1, 256, 0}
	case 0x55:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB: 2-MByte or 4-MByte pages", 2 * 1024, 0xFF, -1, 7, 0}
	case 0x56:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB0", 4 * 1024, 4, -1, 16, 0}
	case 0x57:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB0", 4, 4, -1, 16, 0}
	case 0x59:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB0", 4, 0xFF, -1, 16, 0}
	case 0x5A:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB0 2-MByte or 4 MByte pages", 2 * 1024, 4, -1, 32, 0}
	case 0x5B:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB 4 KByte and 4 MByte pages", 4, -1, -1, 64, 0}
	case 0x5C:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB 4 KByte and 4 MByte pages", 4, -1, -1, 128, 0}
	case 0x5D:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB 4 KByte and 4 MByte pages", 4, -1, -1, 256, 0}
	case 0x60:
		return X86CacheDescriptor{1, X86CacheType_DATA_CACHE, "1st-level data cache", 16, 8, 64, -1, 0}
	case 0x61:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB", 4, 0xFF, -1, 48, 0}
	case 0x63:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB", 1 * 1024 * 1024, 4, -1, 4, 0}
	case 0x66:
		return X86CacheDescriptor{1, X86CacheType_DATA_CACHE, "1st-level data cache", 8, 4, 64, -1, 0}
	case 0x67:
		return X86CacheDescriptor{1, X86CacheType_DATA_CACHE, "1st-level data cache", 16, 4, 64, -1, 0}
	case 0x68:
		return X86CacheDescriptor{1, X86CacheType_DATA_CACHE, "1st-level data cache", 32, 4, 64, -1, 0}
	case 0x70:
		return X86CacheDescriptor{1, X86CacheType_INSTRUCTION_CACHE, "Trace cache (size in K of uop)", 12, 8, -1, -1, 0}
	case 0x71:
		return X86CacheDescriptor{1, X86CacheType_INSTRUCTION_CACHE, "Trace cache (size in K of uop)", 16, 8, -1, -1, 0}
	case 0x72:
		return X86CacheDescriptor{1, X86CacheType_INSTRUCTION_CACHE, "Trace cache (size in K of uop)", 32, 8, -1, -1, 0}
	case 0x76:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB: 2M/4M pages", 2 * 1024, 0xFF, -1, 8, 0}
	case 0x78:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 1 * 1024, 4, 64, -1, 0}
	case 0x79:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 128, 8, 64, -1, 2}
	case 0x7A:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 256, 8, 64, -1, 2}
	case 0x7B:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 512, 8, 64, -1, 2}
	case 0x7C:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 1 * 1024, 8, 64, -1, 2}
	case 0x7D:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 2 * 1024, 8, 64, -1, 0}
	case 0x7F:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 512, 2, 64, -1, 0}
	case 0x80:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 512, 8, 64, -1, 0}
	case 0x82:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 256, 8, 32, -1, 0}
	case 0x83:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 512, 8, 32, -1, 0}
	case 0x84:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 1 * 1024, 8, 32, -1, 0}
	case 0x85:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 2 * 1024, 8, 32, -1, 0}
	case 0x86:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 512, 4, 32, -1, 0}
	case 0x87:
		return X86CacheDescriptor{2, X86CacheType_DATA_CACHE, "2nd-level cache", 1 * 1024, 8, 64, -1, 0}
	case 0xA0:
		return X86CacheDescriptor{-1, X86CacheType_DTLB, "DTLB", 4, 0xFF, -1, 32, 0}
	case 0xB0:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB", 4, 4, -1, 128, 0}
	case 0xB1:
		return X86CacheDescriptor{
			-1, X86CacheType_TLB,
			"Instruction TLB 2M pages 4 way 8 entries or 4M pages 4-way, 4 entries",
			2 * 1024, 4, -1, 8, 0,
		}
	case 0xB2:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB", 4, 4, -1, 64, 0}
	case 0xB3:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB", 4, 4, -1, 128, 0}
	case 0xB4:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB1", 4, 4, -1, 256, 0}
	case 0xB5:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB", 4, 8, -1, 64, 0}
	case 0xB6:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Instruction TLB", 4, 8, -1, 128, 0}
	case 0xBA:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB1", 4, 4, -1, 64, 0}
	case 0xC0:
		return X86CacheDescriptor{-1, X86CacheType_TLB, "Data TLB: 4 KByte and 4 MByte pages", 4, 4, -1, 8, 0}
	case 0xC1:
		return X86CacheDescriptor{-1, X86CacheType_STLB, "Shared 2nd-Level TLB: 4Kbyte and 2Mbyte pages", 4, 8, -1, 1024, 0}
	case 0xC2:
		return X86CacheDescriptor{-1, X86CacheType_DTLB, "DTLB 4KByte/2 MByte pages", 4, 4, -1, 16, 0}
	case 0xC3:
		return X86CacheDescriptor{
			-1, X86CacheType_STLB,
			"Shared 2nd-Level TLB: 4 KByte /2 MByte pages, 6-way associative, 1536 entries. Also 1GBbyte pages, 4-way,16 entries.",
			4, 6, -1, 1536, 0,
		}
	case 0xCA:
		return X86CacheDescriptor{-1, X86CacheType_STLB, "Shared 2nd-Level TLB", 4, 4, -1, 512, 0}
	case 0xD0:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 512, 4, 64, -1, 0}
	case 0xD1:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 1 * 1024, 4, 64, -1, 0}
	case 0xD2:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 2 * 1024, 4, 64, -1, 0}
	case 0xD6:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 1 * 1024, 8, 64, -1, 0}
	case 0xD7:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 2 * 1024, 8, 64, -1, 0}
	case 0xD8:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 4 * 1024, 8, 64, -1, 0}
	case 0xDC:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 1 * 1536, 12, 64, -1, 0}
	case 0xDD:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 3 * 1024, 12, 64, -1, 0}
	case 0xDE:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 6 * 1024, 12, 64, -1, 0}
	case 0xE2:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 2 * 1024, 16, 64, -1, 0}
	case 0xE3:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 4 * 1024, 16, 64, -1, 0}
	case 0xE4:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 8 * 1024, 16, 64, -1, 0}
	case 0xEA:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 12 * 1024, 24, 64, -1, 0}
	case 0xEB:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 18 * 1024, 24, 64, -1, 0}
	case 0xEC:
		return X86CacheDescriptor{3, X86CacheType_DATA_CACHE, "3nd-level cache", 24 * 1024, 24, 64, -1, 0}
	case 0xF0:
		return X86CacheDescriptor{-1, X86CacheType_PREFETCH, "", 64, -1, -1, -1, 0}
	case 0xF1:
		return X86CacheDescriptor{-1, X86CacheType_PREFETCH, "", 128, -1, -1, -1, 0}
	case 0xFF:
		return X86CacheDescriptor{
			-1, X86CacheType_NULL,
			"CPUID leaf 2 does not report cache descriptor information, use CPUID leaf 4 to query cache parameters",
			-1, -1, -1, -1, 0,
		}
	default:
		return X86CacheDescriptor{}
	}
}

// leaf3 (eax=3): Processor Serial Number
func leaf3(cpu *X86, vnd Vendor) {
	if vnd != Vendor_Intel {
		return
	}

	// TODO: introduced on Intel Pentium III, but due to privacy concerns, this feature is no longer implemented on later models (the PSN feature bit is always cleared)
}

// leaf4 (eax=4): Intel thread/core and cache topology
func leaf4(cpu *X86, vnd Vendor) {
	if vnd != Vendor_Intel {
		return
	}

	var v [4]uint32
	for cacheID := uint32(0); true; cacheID++ {
		v = cpuid(4, cacheID)
		eax, ebx, ecx := v[0], v[1], v[2]

		typ := X86CacheType(eax & 0xF)
		if typ == X86CacheType_NULL {
			break
		}

		cpu.CacheDescriptors = append(cpu.CacheDescriptors, X86CacheDescriptor{
			CacheType: typ,
			CacheName: "",
			Level:     int8(eax>>5) & 0x7,
			// CacheSize: ,
			Ways:       int16((ebx>>22)&0x3FF + 1),
			LineSize:   int16((ebx & 0xFFF) + 1),
			Entries:    int32(ecx + 1),
			Partioning: uint16((ebx>>12)&0x3FF + 1),
		})
	}
}

// leaf6 (eax=6): Thermal and power management
func leaf6(cpu *X86, vnd Vendor) {
	v := cpuid(6, 0)

	cpu.ThermalPowerFeatures = X86ThermalPowerFeature(v[0]&0xFFFF) | X86ThermalPowerFeature(v[2]&0xFFFF)<<16
	cpu.ThermalSensorInterruptThresholds = int8(v[1] & 0xF)
}

// leaf7 (eax=7, ecx=0,1): Extended Features
func leaf7(cpu *X86, vnd Vendor) {
	v := cpuid(7, 0)

	cpu.ExtendedFeatures1 = *(*X86FeatureExtendedBC)(unsafe.Pointer(&v[1]))

	d := v[3]

	v = cpuid(7, 1)
	cpu.ExtendedFeatures2 = X86FeatureExtendedDA(d) | X86FeatureExtendedDA(v[0])<<32
}

func leaf0x80000000() (maxExtFuncNums uint32) {
	v := cpuid(0x80000000, 0)
	return v[0]
}

// EAX=80000001h: Extended Processor Info and Feature Bits
func leaf0x8000_0001(cpu *X86, vnd Vendor) {
	v := cpuid(0x80000001, 0)

	cpu.ExtraFeatures = *(*X86FeatureExtra)(unsafe.Pointer(&v[2]))
}

// Processor Brand String
func leaf0x8000_0004(cpu *X86, vnd Vendor) {
	v0 := cpuid(0x80000002, 0)
	v1 := cpuid(0x80000003, 0)
	v2 := cpuid(0x80000004, 0)

	cpu.BrandDetail = string(unsafe.Slice(
		(*byte)(unsafe.Pointer(&v0[0])),
		unsafe.Sizeof(v0[0])*uintptr(len(v0)),
	))
	cpu.BrandDetail += string(unsafe.Slice(
		(*byte)(unsafe.Pointer(&v1[0])),
		unsafe.Sizeof(v1[0])*uintptr(len(v1)),
	))
	cpu.BrandDetail += string(unsafe.Slice(
		(*byte)(unsafe.Pointer(&v2[0])),
		unsafe.Sizeof(v2[0])*uintptr(len(v2)),
	))
}

// AMD L1 Cache and TLB Information
func leaf0x8000_0005(cpu *X86, vnd Vendor) {
	if vnd != Vendor_AMD {
		return
	}

	// TODO
}

func leaf0x8000_0006(cpu *X86, vnd Vendor) {
	switch vnd {
	case Vendor_AMD:
		// TODO
	case Vendor_Intel:
		// TODO
	default:
		return
	}
}

// AMD Encrypted Memory Capabilities
func leaf0x8000_001f(cpu *X86, vnd Vendor) {
	if vnd != Vendor_AMD {
		return
	}

	// TODO
}
