package cpuhelper

import (
	"sort"
	"strconv"
	"strings"
)

type X86 struct {
	Brand       string
	BrandDetail string

	Stepping    X86Stepping
	Model       X86Model
	Family      X86Family
	ProcessType X86ProcessType

	BrandIndex      uint8
	CacheLineSize   uint16
	MaxLogicalCPUID uint8
	InitialAPICID   uint8

	Features          X86Feature
	ExtendedFeatures1 X86FeatureExtendedBC
	ExtendedFeatures2 X86FeatureExtendedDA
	ExtraFeatures     X86FeatureExtra

	ThermalPowerFeatures             X86ThermalPowerFeature
	ThermalSensorInterruptThresholds int8

	CacheDescriptors []X86CacheDescriptor
}

func (cpu *X86) String() string {
	var sb strings.Builder

	sb.WriteString("brand: ")
	sb.WriteString(string(cpu.Brand))
	sb.WriteString("\n")

	sb.WriteString("vendor: ")
	sb.WriteString(cpu.Vendor().String())
	sb.WriteString("\n")

	sb.WriteString("name: ")
	sb.WriteString(cpu.BrandDetail)
	sb.WriteString("\n")

	sb.WriteString("microArch: ")
	sb.WriteString(strconv.FormatInt(int64(cpu.MicroArch()), 10))
	sb.WriteString("\n")

	sb.WriteString("features: [")
	allFeatures := cpu.AllFeatures()
	if len(allFeatures) == 0 {
		sb.WriteString("]\n")
	} else {
		sb.WriteString("\n  ")
		for i, feat := range allFeatures {
			if (i+1)%8 == 0 {
				sb.WriteString("\n  ")
			}

			sb.WriteString(feat)
			sb.WriteString(", ")
		}
		sb.WriteString("\n]\n")
	}

	sb.WriteString("thermalPowerFeatures:")
	allFeatures = cpu.ThermalPowerFeatures.Features()
	if len(allFeatures) == 0 {
		sb.WriteString(" []\n")
	} else {
		for _, feat := range allFeatures {
			sb.WriteString("\n- ")
			sb.WriteString(feat)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (cpu *X86) AllFeatures() (ret []string) {
	f0 := cpu.Features.Features()
	sort.Strings(f0)
	ret = append(ret, f0...)

	f1 := cpu.ExtendedFeatures1.Features()
	sort.Strings(f1)
	ret = append(ret, f1...)

	f2 := cpu.ExtendedFeatures2.Features()
	sort.Strings(f2)
	ret = append(ret, f2...)

	f3 := cpu.ExtraFeatures.Features()
	sort.Strings(f3)
	ret = append(ret, f3...)

	return
}

// MicroArch level (for x86-64)
// ref: https://en.wikipedia.org/wiki/X86-64#Microarchitecture_levels
func (cpu *X86) MicroArch() (ret int) {
	ok := cpu.Features.HasAll(X86Feature_CMOV) || cpu.ExtraFeatures.HasAll(X86ExtraFeature_CMOV)
	ok = ok && (cpu.Features.HasAll(X86Feature_CX8) || cpu.ExtraFeatures.HasAll(X86ExtraFeature_CX8))
	ok = ok && (cpu.Features.HasAll(X86Feature_FPU) || cpu.ExtraFeatures.HasAll(X86ExtraFeature_FPU))
	ok = ok && (cpu.Features.HasAll(X86Feature_FXSR) || cpu.ExtraFeatures.HasAll(X86ExtraFeature_FXSR))
	ok = ok && (cpu.Features.HasAll(X86Feature_MMX) || cpu.ExtraFeatures.HasAll(X86ExtraFeature_MMX))
	ok = ok && cpu.Features.HasAll(X86Feature_SEP, X86Feature_SSE, X86Feature_SSE2)

	if ok {
		ret = 1
	} else {
		return
	}

	if cpu.Features.HasAll(
		X86Feature_CX16, X86Feature_POPCNT, X86Feature_SSE3,
		X86Feature_SSE4_1, X86Feature_SSE4_2, X86Feature_SSSE3,
	) && cpu.ExtraFeatures.HasAll(
		X86ExtraFeature_LAHF_LM,
	) {
		ret = 2
	} else {
		return
	}

	if cpu.Features.HasAll(
		X86Feature_AVX, X86Feature_F16C, X86Feature_MOVBE,
		X86Feature_FMA, X86Feature_OSXSAVE,
	) && cpu.ExtraFeatures.HasAll(
		X86ExtraFeature_ABM,
	) && cpu.ExtendedFeatures1.HasAll(
		X86FeatureExtendedBC_AVX2, X86FeatureExtendedBC_BMI1, X86FeatureExtendedBC_BMI2,
	) {
		ret = 3
	} else {
		return
	}

	if cpu.ExtendedFeatures1.HasAll(
		X86FeatureExtendedBC_AVX512_F, X86FeatureExtendedBC_AVX512_BW, X86FeatureExtendedBC_AVX512_CD, X86FeatureExtendedBC_AVX512_DQ, X86FeatureExtendedBC_AVX512_VL,
	) {
		ret = 4
	} else {
		return
	}

	return
}

type X86Stepping uint8
type X86Model uint8
type X86Family uint16
type X86ProcessType uint8

// Vendor returns predefined Vendor value according to brand name
func (n X86) Vendor() Vendor {
	switch n.Brand {
	case "AMDisbetter!", "AuthenticAMD":
		return Vendor_AMD
	case "CentaurHauls", "VIA VIA VIA ":
		return Vendor_VIA
	case "CyrixInstead":
		return Vendor_Cyrix
	case "GenuineIntel":
		if strings.HasPrefix(n.BrandDetail, "VirtualApple ") {
			return Vendor_AppleRosetta2
		}

		return Vendor_Intel
	case "TransmetaCPU":
		return Vendor_Transmeta
	case "GenuineTMx86":
		return Vendor_Transmeta
	case "Geode by NSC":
		return Vendor_NSC
	case "NexGenDriven":
		return Vendor_NexGen
	case "HygonGenuine":
		return Vendor_Hygon
	case "RiseRiseRise":
		return Vendor_Rise
	case "SiS SiS SiS ":
		return Vendor_SiS
	case "Vortex86 SoC":
		return Vendor_DMP
	case "UMC UMC UMC ":
		return Vendor_UMC
	case "Genuine  RDC":
		return Vendor_RDC
	case "E2K MACHINE":
		return Vendor_MCST
	default:
		return HypervisorVendor(n.Brand)
	}
}

type X86Feature uint64

func (feat X86Feature) Features() (ret []string) {
	for i := 0; i < 64; i++ {
		if x := feat & (X86Feature(1) << i); x != 0 {
			str := x.String()
			if len(str) == 0 {
				continue
			}

			ret = append(ret, str)
		}
	}

	return
}

func (feat X86Feature) HasAll(features ...X86Feature) bool {
	for _, f := range features {
		if f&feat == 0 {
			return false
		}
	}

	return true
}

func (feat X86Feature) String() string {
	switch feat {
	case X86Feature_SSE3:
		return "SSE3"
	case X86Feature_PCLMULQDQ:
		return "PCLMULQDQ"
	case X86Feature_DTES64:
		return "DTES64"
	case X86Feature_MONITOR:
		return "MONITOR"
	case X86Feature_DSI_CPL:
		return "DSI_CPL"
	case X86Feature_VMX:
		return "VMX"
	case X86Feature_SMX:
		return "SMX"
	case X86Feature_EST:
		return "EST"
	case X86Feature_TM2:
		return "TM2"
	case X86Feature_SSSE3:
		return "SSSE3"
	case X86Feature_CNXT_ID:
		return "CNXT_ID"
	case X86Feature_SDBG:
		return "SDBG"
	case X86Feature_FMA:
		return "FMA"
	case X86Feature_CX16:
		return "CX16"
	case X86Feature_XTPR:
		return "XTPR"
	case X86Feature_PDCM:
		return "PDCM"
	case X86Feature_PCID:
		return "PCID"
	case X86Feature_DCA:
		return "DCA"
	case X86Feature_SSE4_1:
		return "SSE4.1"
	case X86Feature_SSE4_2:
		return "SSE4.2"
	case X86Feature_X2APIC:
		return "X2APIC"
	case X86Feature_MOVBE:
		return "MOVBE"
	case X86Feature_POPCNT:
		return "POPCNT"
	case X86Feature_TSC_DEADLINE:
		return "TSC_DEADLINE"
	case X86Feature_AES:
		return "AES"
	case X86Feature_XSAVE:
		return "XSAVE"
	case X86Feature_OSXSAVE:
		return "OSXSAVE"
	case X86Feature_AVX:
		return "AVX"
	case X86Feature_F16C:
		return "F16C"
	case X86Feature_RDRND:
		return "RDRND"
	case X86Feature_HYPERVISOR:
		return "HYPERVISOR"
	case X86Feature_FPU:
		return "FPU"
	case X86Feature_VME:
		return "VME"
	case X86Feature_DE:
		return "DE"
	case X86Feature_PSE:
		return "PSE"
	case X86Feature_TSC:
		return "TSC"
	case X86Feature_MSR:
		return "MSR"
	case X86Feature_PAE:
		return "PAE"
	case X86Feature_MCE:
		return "MCE"
	case X86Feature_CX8:
		return "CX8"
	case X86Feature_APIC:
		return "APIC"
	case X86Feature_SEP:
		return "SEP"
	case X86Feature_MTRR:
		return "MTRR"
	case X86Feature_PGE:
		return "PGE"
	case X86Feature_MCA:
		return "MCA"
	case X86Feature_CMOV:
		return "CMOV"
	case X86Feature_PAT:
		return "PAT"
	case X86Feature_PSE_36:
		return "PSE_36"
	case X86Feature_PSN:
		return "PSN"
	case X86Feature_CLFSH:
		return "CLFSH"
	case X86Feature_DS:
		return "DS"
	case X86Feature_ACPI:
		return "ACPI"
	case X86Feature_MMX:
		return "MMX"
	case X86Feature_FXSR:
		return "FXSR"
	case X86Feature_SSE:
		return "SSE"
	case X86Feature_SSE2:
		return "SSE2"
	case X86Feature_SS:
		return "SS"
	case X86Feature_HTT:
		return "HTT"
	case X86Feature_TM:
		return "TM"
	case X86Feature_IA64:
		return "IA64"
	case X86Feature_PBE:
		return "PBE"
	default:
		return ""
	}
}

const (
	// leaf1 ecx
	X86Feature_SSE3 X86Feature = 1 << iota
	X86Feature_PCLMULQDQ
	X86Feature_DTES64
	X86Feature_MONITOR
	X86Feature_DSI_CPL
	X86Feature_VMX
	X86Feature_SMX
	X86Feature_EST
	X86Feature_TM2
	X86Feature_SSSE3
	X86Feature_CNXT_ID
	X86Feature_SDBG
	X86Feature_FMA
	X86Feature_CX16
	X86Feature_XTPR
	X86Feature_PDCM
	_
	X86Feature_PCID
	X86Feature_DCA
	X86Feature_SSE4_1
	X86Feature_SSE4_2
	X86Feature_X2APIC
	X86Feature_MOVBE
	X86Feature_POPCNT
	X86Feature_TSC_DEADLINE
	X86Feature_AES
	X86Feature_XSAVE
	X86Feature_OSXSAVE
	X86Feature_AVX
	X86Feature_F16C
	X86Feature_RDRND
	X86Feature_HYPERVISOR

	// leaf1 edx
	X86Feature_FPU
	X86Feature_VME
	X86Feature_DE
	X86Feature_PSE
	X86Feature_TSC
	X86Feature_MSR
	X86Feature_PAE
	X86Feature_MCE
	X86Feature_CX8
	X86Feature_APIC
	_
	X86Feature_SEP // sysenter/sysexit
	X86Feature_MTRR
	X86Feature_PGE
	X86Feature_MCA
	X86Feature_CMOV
	X86Feature_PAT
	X86Feature_PSE_36
	X86Feature_PSN
	X86Feature_CLFSH
	_
	X86Feature_DS
	X86Feature_ACPI
	X86Feature_MMX
	X86Feature_FXSR
	X86Feature_SSE
	X86Feature_SSE2
	X86Feature_SS
	X86Feature_HTT
	X86Feature_TM
	X86Feature_IA64
	X86Feature_PBE
)

// X86FeatureExtendedBC (leaf7.0 ebx, ecx)
type X86FeatureExtendedBC uint64

func (feat X86FeatureExtendedBC) Features() (ret []string) {
	for i := 0; i < 64; i++ {
		if x := feat & (X86FeatureExtendedBC(1) << i); x != 0 {
			str := x.String()
			if len(str) == 0 {
				continue
			}

			ret = append(ret, str)
		}
	}

	return
}

func (feat X86FeatureExtendedBC) HasAll(features ...X86FeatureExtendedBC) bool {
	for _, f := range features {
		if f&feat == 0 {
			return false
		}
	}

	return true
}

func (feat X86FeatureExtendedBC) String() string {
	switch feat {
	case X86FeatureExtendedBC_FSGSBASE:
		return "FSGSBASE"
	case X86FeatureExtendedBC_IA32_TSC_ADJUST:
		return "IA32_TSC_ADJUST"
	case X86FeatureExtendedBC_BMI1:
		return "BMI1"
	case X86FeatureExtendedBC_HLE:
		return "HLE"
	case X86FeatureExtendedBC_AVX2:
		return "AVX2"
	case X86FeatureExtendedBC_SMEP:
		return "SMEP"
	case X86FeatureExtendedBC_BMI2:
		return "BMI2"
	case X86FeatureExtendedBC_ERMS:
		return "ERMS"
	case X86FeatureExtendedBC_INVPCID:
		return "INVPCID"
	case X86FeatureExtendedBC_RTM:
		return "RTM"
	case X86FeatureExtendedBC_PQM:
		return "PQM"
	case X86FeatureExtendedBC_DFPUCDS:
		return "DFPUCDS"
	case X86FeatureExtendedBC_MPX:
		return "MPX"
	case X86FeatureExtendedBC_PQE:
		return "PQE"
	case X86FeatureExtendedBC_AVX512_F:
		return "AVX512_F"
	case X86FeatureExtendedBC_AVX512_DQ:
		return "AVX512_DQ"
	case X86FeatureExtendedBC_RDSEED:
		return "RDSEED"
	case X86FeatureExtendedBC_ADX:
		return "ADX"
	case X86FeatureExtendedBC_SMAP:
		return "SMAP"
	case X86FeatureExtendedBC_AVX512_IFMA:
		return "AVX512_IFMA"
	case X86FeatureExtendedBC_PCOMMIT:
		return "PCOMMIT"
	case X86FeatureExtendedBC_CLFLUSHOPT:
		return "CLFLUSHOPT"
	case X86FeatureExtendedBC_CLWB:
		return "CLWB"
	case X86FeatureExtendedBC_INTEL_PT:
		return "INTEL_PT"
	case X86FeatureExtendedBC_AVX512_PF:
		return "AVX512_PF"
	case X86FeatureExtendedBC_AVX512_ER:
		return "AVX512_ER"
	case X86FeatureExtendedBC_AVX512_CD:
		return "AVX512_CD"
	case X86FeatureExtendedBC_SHA:
		return "SHA"
	case X86FeatureExtendedBC_AVX512_BW:
		return "AVX512_BW"
	case X86FeatureExtendedBC_AVX512_VL:
		return "AVX512_VL"
	case X86FeatureExtendedBC_PREFETCHWT1:
		return "PREFETCHWT1"
	case X86FeatureExtendedBC_AVX512_VBMI:
		return "AVX512_VBMI"
	case X86FeatureExtendedBC_UMIP:
		return "UMIP"
	case X86FeatureExtendedBC_PKU:
		return "PKU"
	case X86FeatureExtendedBC_OSPKE:
		return "OSPKE"
	case X86FeatureExtendedBC_WAITPKG:
		return "WAITPKG"
	case X86FeatureExtendedBC_AVX512_VBMI2:
		return "AVX512_VBMI2"
	case X86FeatureExtendedBC_CETSS:
		return "CETSS"
	case X86FeatureExtendedBC_GFNI:
		return "GFNI"
	case X86FeatureExtendedBC_VAES:
		return "VAES"
	case X86FeatureExtendedBC_VPCLMULQDQ:
		return "VPCLMULQDQ"
	case X86FeatureExtendedBC_AVX512_VNNI:
		return "AVX512_VNNI"
	case X86FeatureExtendedBC_AVX512_BITALG:
		return "AVX512_BITALG"
	case X86FeatureExtendedBC_TME_EN:
		return "TME_EN"
	case X86FeatureExtendedBC_AVX512_VPOPCNTDQ:
		return "AVX512_VPOPCNTDQ"
	case X86FeatureExtendedBC_INTEL_5_LEVEL_PAGING:
		return "INTEL_5_LEVEL_PAGING"
	case X86FeatureExtendedBC_MAWAU_0:
		return "MAWAU_0"
	case X86FeatureExtendedBC_MAWAU_1:
		return "MAWAU_1"
	case X86FeatureExtendedBC_MAWAU_2:
		return "MAWAU_2"
	case X86FeatureExtendedBC_MAWAU_3:
		return "MAWAU_3"
	case X86FeatureExtendedBC_MAWAU_4:
		return "MAWAU_4"
	case X86FeatureExtendedBC_RDPID:
		return "RDPID"
	case X86FeatureExtendedBC_KL:
		return "KL"
	case X86FeatureExtendedBC_CLDEMOTE:
		return "CLDEMOTE"
	case X86FeatureExtendedBC_MOVDIRI:
		return "MOVDIRI"
	case X86FeatureExtendedBC_MOVDIR64B:
		return "MOVDIR64B"
	case X86FeatureExtendedBC_ENQCMD:
		return "ENQCMD"
	case X86FeatureExtendedBC_SGX_LC:
		return "SGX_LC"
	case X86FeatureExtendedBC_PKS:
		return "PKS"
	default:
		return ""
	}
}

const (
	// leaf7.0 ebx
	X86FeatureExtendedBC_FSGSBASE X86FeatureExtendedBC = 1 << iota
	X86FeatureExtendedBC_IA32_TSC_ADJUST
	_
	X86FeatureExtendedBC_BMI1
	X86FeatureExtendedBC_HLE
	X86FeatureExtendedBC_AVX2
	_
	X86FeatureExtendedBC_SMEP
	X86FeatureExtendedBC_BMI2
	X86FeatureExtendedBC_ERMS
	X86FeatureExtendedBC_INVPCID
	X86FeatureExtendedBC_RTM
	X86FeatureExtendedBC_PQM
	X86FeatureExtendedBC_DFPUCDS
	X86FeatureExtendedBC_MPX
	X86FeatureExtendedBC_PQE
	X86FeatureExtendedBC_AVX512_F
	X86FeatureExtendedBC_AVX512_DQ
	X86FeatureExtendedBC_RDSEED
	X86FeatureExtendedBC_ADX
	X86FeatureExtendedBC_SMAP
	X86FeatureExtendedBC_AVX512_IFMA
	X86FeatureExtendedBC_PCOMMIT
	X86FeatureExtendedBC_CLFLUSHOPT
	X86FeatureExtendedBC_CLWB
	X86FeatureExtendedBC_INTEL_PT // Intel Processor Trace
	X86FeatureExtendedBC_AVX512_PF
	X86FeatureExtendedBC_AVX512_ER
	X86FeatureExtendedBC_AVX512_CD
	X86FeatureExtendedBC_SHA
	X86FeatureExtendedBC_AVX512_BW
	X86FeatureExtendedBC_AVX512_VL

	// leaf7.0 ecx
	X86FeatureExtendedBC_PREFETCHWT1
	X86FeatureExtendedBC_AVX512_VBMI
	X86FeatureExtendedBC_UMIP
	X86FeatureExtendedBC_PKU
	X86FeatureExtendedBC_OSPKE
	X86FeatureExtendedBC_WAITPKG
	X86FeatureExtendedBC_AVX512_VBMI2
	X86FeatureExtendedBC_CETSS
	X86FeatureExtendedBC_GFNI
	X86FeatureExtendedBC_VAES
	X86FeatureExtendedBC_VPCLMULQDQ
	X86FeatureExtendedBC_AVX512_VNNI
	X86FeatureExtendedBC_AVX512_BITALG
	X86FeatureExtendedBC_TME_EN
	X86FeatureExtendedBC_AVX512_VPOPCNTDQ
	_
	X86FeatureExtendedBC_INTEL_5_LEVEL_PAGING
	X86FeatureExtendedBC_MAWAU_0
	X86FeatureExtendedBC_MAWAU_1
	X86FeatureExtendedBC_MAWAU_2
	X86FeatureExtendedBC_MAWAU_3
	X86FeatureExtendedBC_MAWAU_4
	X86FeatureExtendedBC_RDPID
	X86FeatureExtendedBC_KL
	_
	X86FeatureExtendedBC_CLDEMOTE
	_
	X86FeatureExtendedBC_MOVDIRI
	X86FeatureExtendedBC_MOVDIR64B
	X86FeatureExtendedBC_ENQCMD
	X86FeatureExtendedBC_SGX_LC
	X86FeatureExtendedBC_PKS
)

// X86FeatureExtendedDA (leaf7.0 edx, leaf7.1 eax)
type X86FeatureExtendedDA uint64

func (feat X86FeatureExtendedDA) Features() (ret []string) {
	for i := 0; i < 64; i++ {
		if x := feat & (X86FeatureExtendedDA(1) << i); x != 0 {
			str := x.String()
			if len(str) == 0 {
				continue
			}

			ret = append(ret, str)
		}
	}

	return
}

func (feat X86FeatureExtendedDA) HasAll(features ...X86FeatureExtendedDA) bool {
	for _, f := range features {
		if f&feat == 0 {
			return false
		}
	}

	return true
}

func (feat X86FeatureExtendedDA) String() string {
	switch feat {
	case X86FeatureExtendedDA_AVX512_4VNNIW:
		return "AVX512_4VNNIW"
	case X86FeatureExtendedDA_AVX512_4FMAPS:
		return "AVX512_4FMAPS"
	case X86FeatureExtendedDA_FSRM:
		return "FSRM"
	case X86FeatureExtendedDA_AVX512_VP2INTERSECT:
		return "AVX512_VP2INTERSECT"
	case X86FeatureExtendedDA_SRBDS_CTRL:
		return "SRBDS_CTRL"
	case X86FeatureExtendedDA_MD_CLEAR:
		return "MD_CLEAR"
	case X86FeatureExtendedDA_RTM_ALWAYS_ABORT:
		return "RTM_ALWAYS_ABORT"
	case X86FeatureExtendedDA_TSX_FORCE_ABORT:
		return "TSX_FORCE_ABORT"
	case X86FeatureExtendedDA_SERIALIZE:
		return "SERIALIZE"
	case X86FeatureExtendedDA_HYBRID:
		return "HYBRID"
	case X86FeatureExtendedDA_TSXLDTRK:
		return "TSXLDTRK"
	case X86FeatureExtendedDA_PCONFIG:
		return "PCONFIG"
	case X86FeatureExtendedDA_LBR:
		return "LBR"
	case X86FeatureExtendedDA_CET_IBT:
		return "CET_IBT"
	case X86FeatureExtendedDA_AMX_BF16:
		return "AMX_BF16"
	case X86FeatureExtendedDA_AVX512_FP16:
		return "AVX512_FP16"
	case X86FeatureExtendedDA_AMX_TILE:
		return "AMX_TILE"
	case X86FeatureExtendedDA_AMX_INT8:
		return "AMX_INT8"
	case X86FeatureExtendedDA_SPEC_CTRL:
		return "SPEC_CTRL"
	case X86FeatureExtendedDA_STIBP:
		return "STIBP"
	case X86FeatureExtendedDA_L1D_FLUSH:
		return "L1D_FLUSH"
	case X86FeatureExtendedDA_IA32_ARCH_CAPABILITIES:
		return "IA32_ARCH_CAPABILITIES"
	case X86FeatureExtendedDA_IA32_CORE_CAPABILITIES:
		return "IA32_CORE_CAPABILITIES"
	case X86FeatureExtendedDA_SSBD:
		return "SSBD"
	case X86FeatureExtendedDA_AVX512_BF16:
		return "AVX512_BF16"
	default:
		return ""
	}
}

const (
	// leaf7.0 edx
	_ X86FeatureExtendedDA = 1 << iota
	_
	X86FeatureExtendedDA_AVX512_4VNNIW
	X86FeatureExtendedDA_AVX512_4FMAPS
	X86FeatureExtendedDA_FSRM
	_
	_
	_
	X86FeatureExtendedDA_AVX512_VP2INTERSECT
	X86FeatureExtendedDA_SRBDS_CTRL
	X86FeatureExtendedDA_MD_CLEAR
	X86FeatureExtendedDA_RTM_ALWAYS_ABORT
	_
	X86FeatureExtendedDA_TSX_FORCE_ABORT
	X86FeatureExtendedDA_SERIALIZE
	X86FeatureExtendedDA_HYBRID
	X86FeatureExtendedDA_TSXLDTRK
	_
	X86FeatureExtendedDA_PCONFIG
	X86FeatureExtendedDA_LBR
	X86FeatureExtendedDA_CET_IBT
	_
	X86FeatureExtendedDA_AMX_BF16
	X86FeatureExtendedDA_AVX512_FP16
	X86FeatureExtendedDA_AMX_TILE
	X86FeatureExtendedDA_AMX_INT8
	X86FeatureExtendedDA_SPEC_CTRL
	X86FeatureExtendedDA_STIBP
	X86FeatureExtendedDA_L1D_FLUSH
	X86FeatureExtendedDA_IA32_ARCH_CAPABILITIES
	X86FeatureExtendedDA_IA32_CORE_CAPABILITIES
	X86FeatureExtendedDA_SSBD

	// leaf7.1 eax
	_
	_
	_
	_
	_
	X86FeatureExtendedDA_AVX512_BF16
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
)

// X86FeatureExtra (leaf0x8000_0001.0 ecx, edx)
type X86FeatureExtra uint64

func (feat X86FeatureExtra) Features() (ret []string) {
	for i := 0; i < 64; i++ {
		if x := feat & (X86FeatureExtra(1) << i); x != 0 {
			str := x.String()
			if len(str) == 0 {
				continue
			}

			ret = append(ret, str)
		}
	}

	return
}

func (feat X86FeatureExtra) HasAll(features ...X86FeatureExtra) bool {
	for _, f := range features {
		if f&feat == 0 {
			return false
		}
	}

	return true
}

func (feat X86FeatureExtra) String() string {
	switch feat {
	case X86ExtraFeature_LAHF_LM:
		return "LAHF_LM"
	case X86ExtraFeature_CMP_LEGACY:
		return "CMP_LEGACY"
	case X86ExtraFeature_SVM:
		return "SVM"
	case X86ExtraFeature_EXTAPIC:
		return "EXTAPIC"
	case X86ExtraFeature_CR8_LEGACY:
		return "CR8_LEGACY"
	case X86ExtraFeature_ABM:
		return "ABM"
	case X86ExtraFeature_SSE4A:
		return "SSE4A"
	case X86ExtraFeature_MISALIGN_SSE:
		return "MISALIGN_SSE"
	case X86ExtraFeature_3D_NOW_PREFETCH:
		return "3D_NOW_PREFETCH"
	case X86ExtraFeature_OSVW:
		return "OSVW"
	case X86ExtraFeature_IBS:
		return "IBS"
	case X86ExtraFeature_XOP:
		return "XOP"
	case X86ExtraFeature_SKINIT:
		return "SKINIT"
	case X86ExtraFeature_WDT:
		return "WDT"
	case X86ExtraFeature_LWP:
		return "LWP"
	case X86ExtraFeature_FMA4:
		return "FMA4"
	case X86ExtraFeature_TCE:
		return "TCE"
	case X86ExtraFeature_NODEID_MSR:
		return "NODEID_MSR"
	case X86ExtraFeature_TBM:
		return "TBM"
	case X86ExtraFeature_TOPOEXT:
		return "TOPOEXT"
	case X86ExtraFeature_PERFCTR_CORE:
		return "PERFCTR_CORE"
	case X86ExtraFeature_PERFCTR_NB:
		return "PERFCTR_NB"
	case X86ExtraFeature_DBX:
		return "DBX"
	case X86ExtraFeature_PERFTSC:
		return "PERFTSC"
	case X86ExtraFeature_PCX_L2I:
		return "PCX_L2I"
	case X86ExtraFeature_FPU:
		return "FPU"
	case X86ExtraFeature_VME:
		return "VME"
	case X86ExtraFeature_DE:
		return "DE"
	case X86ExtraFeature_PSE:
		return "PSE"
	case X86ExtraFeature_TSC:
		return "TSC"
	case X86ExtraFeature_MSR:
		return "MSR"
	case X86ExtraFeature_PAE:
		return "PAE"
	case X86ExtraFeature_MCE:
		return "MCE"
	case X86ExtraFeature_CX8:
		return "CX8"
	case X86ExtraFeature_APIC:
		return "APIC"
	case X86ExtraFeature_SYSCALL:
		return "SYSCALL"
	case X86ExtraFeature_MTRR:
		return "MTRR"
	case X86ExtraFeature_PGE:
		return "PGE"
	case X86ExtraFeature_MCA:
		return "MCA"
	case X86ExtraFeature_CMOV:
		return "CMOV"
	case X86ExtraFeature_PAT:
		return "PAT"
	case X86ExtraFeature_PSE36:
		return "PSE36"
	case X86ExtraFeature_MP:
		return "MP"
	case X86ExtraFeature_NX:
		return "NX"
	case X86ExtraFeature_MMXEXT:
		return "MMXEXT"
	case X86ExtraFeature_MMX:
		return "MMX"
	case X86ExtraFeature_FXSR:
		return "FXSR"
	case X86ExtraFeature_FXSR_OPT:
		return "FXSR_OPT"
	case X86ExtraFeature_PDPE1GB:
		return "PDPE1GB"
	case X86ExtraFeature_RDTSCP:
		return "RDTSCP"
	case X86ExtraFeature_LM:
		return "LM"
	case X86ExtraFeature_3D_NOW_EXT:
		return "3D_NOW_EXT"
	case X86ExtraFeature_3D_NOW:
		return "3D_NOW"
	default:
		return ""
	}
}

const (
	// leaf0x8000_0001.0 ecx
	X86ExtraFeature_LAHF_LM X86FeatureExtra = 1 << iota
	X86ExtraFeature_CMP_LEGACY
	X86ExtraFeature_SVM
	X86ExtraFeature_EXTAPIC
	X86ExtraFeature_CR8_LEGACY
	X86ExtraFeature_ABM
	X86ExtraFeature_SSE4A
	X86ExtraFeature_MISALIGN_SSE
	X86ExtraFeature_3D_NOW_PREFETCH
	X86ExtraFeature_OSVW
	X86ExtraFeature_IBS
	X86ExtraFeature_XOP
	X86ExtraFeature_SKINIT
	X86ExtraFeature_WDT
	_
	X86ExtraFeature_LWP
	X86ExtraFeature_FMA4
	X86ExtraFeature_TCE
	_
	X86ExtraFeature_NODEID_MSR
	_
	X86ExtraFeature_TBM
	X86ExtraFeature_TOPOEXT
	X86ExtraFeature_PERFCTR_CORE
	X86ExtraFeature_PERFCTR_NB
	_
	X86ExtraFeature_DBX
	X86ExtraFeature_PERFTSC
	X86ExtraFeature_PCX_L2I
	_
	_
	_

	// leaf0x8000_0001.0 edx
	X86ExtraFeature_FPU
	X86ExtraFeature_VME
	X86ExtraFeature_DE
	X86ExtraFeature_PSE
	X86ExtraFeature_TSC
	X86ExtraFeature_MSR
	X86ExtraFeature_PAE
	X86ExtraFeature_MCE
	X86ExtraFeature_CX8
	X86ExtraFeature_APIC
	_
	X86ExtraFeature_SYSCALL // syscall/sysret
	X86ExtraFeature_MTRR
	X86ExtraFeature_PGE
	X86ExtraFeature_MCA
	X86ExtraFeature_CMOV
	X86ExtraFeature_PAT
	X86ExtraFeature_PSE36
	_
	X86ExtraFeature_MP
	X86ExtraFeature_NX
	_
	X86ExtraFeature_MMXEXT
	X86ExtraFeature_MMX
	X86ExtraFeature_FXSR
	X86ExtraFeature_FXSR_OPT
	X86ExtraFeature_PDPE1GB
	X86ExtraFeature_RDTSCP
	_
	X86ExtraFeature_LM
	X86ExtraFeature_3D_NOW_EXT
	X86ExtraFeature_3D_NOW
)

type X86ThermalPowerFeature uint32

func (feat X86ThermalPowerFeature) Features() (ret []string) {
	for i := 0; i < 32; i++ {
		if x := feat & (X86ThermalPowerFeature(1) << i); x != 0 {
			str := x.String()
			if len(str) == 0 {
				continue
			}

			ret = append(ret, str)
		}
	}

	return
}

func (feat X86ThermalPowerFeature) HasAll(features ...X86ThermalPowerFeature) bool {
	for _, f := range features {
		if f&feat == 0 {
			return false
		}
	}

	return true
}

func (feat X86ThermalPowerFeature) String() string {
	switch feat {
	case X86ThermalPowerFeature_DigitalThermalSensor:
		return "DigitalThermalSensor (DTS)"
	case X86ThermalPowerFeature_IntelTurboBoost:
		return "IntelTurboBoost (ITB)"
	case X86ThermalPowerFeature_AlwaysRunningAPICTimer:
		return "AlwaysRunningAPICTimer (ARAT)"
	case X86ThermalPowerFeature_PowerLimitNotification:
		return "PowerLimitNotification (PLN)"
	case X86ThermalPowerFeature_ExtendedClockModulationDuty:
		return "ExtendedClockModulationDuty (ECMD)"
	case X86ThermalPowerFeature_PackageThermalManagement:
		return "PackageThermalManagement (PTM)"
	case X86ThermalPowerFeature_HardwareCoordinationFeedback:
		return "HardwareCoordinationFeedback"
	case X86ThermalPowerFeature_ACNT2:
		return "ACNT2"
	case X86ThermalPowerFeature_PerformanceEnergyBias:
		return "PerformanceEnergyBias"
	default:
		return ""
	}
}

const (
	// eax (31:7 reserved)
	X86ThermalPowerFeature_DigitalThermalSensor X86ThermalPowerFeature = 1 << iota
	X86ThermalPowerFeature_IntelTurboBoost
	X86ThermalPowerFeature_AlwaysRunningAPICTimer
	_
	X86ThermalPowerFeature_PowerLimitNotification
	X86ThermalPowerFeature_ExtendedClockModulationDuty
	X86ThermalPowerFeature_PackageThermalManagement
	_
	_
	_
	_
	_
	_
	_
	_
	_

	// ecx (31:4 reserved)
	X86ThermalPowerFeature_HardwareCoordinationFeedback
	X86ThermalPowerFeature_ACNT2
	_
	X86ThermalPowerFeature_PerformanceEnergyBias
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
)

type X86CacheDescriptor struct {
	Level      int8         // Cache level
	CacheType  X86CacheType // Cache type
	CacheName  string       // Name
	CacheSize  int          // in KBytes (of page size for TLB)
	Ways       int16        // Associativity, 0 undefined, 0xFF fully associate
	LineSize   int16        // Cache line size in bytes
	Entries    int32        // number of entries for TLB
	Partioning uint16       // partitioning
}

type X86CacheType uint8

const (
	X86CacheType_NULL X86CacheType = iota
	X86CacheType_DATA_CACHE
	X86CacheType_INSTRUCTION_CACHE
	X86CacheType_UNIFIED_CACHE
	X86CacheType_TLB
	X86CacheType_DTLB
	X86CacheType_STLB
	X86CacheType_PREFETCH
)
