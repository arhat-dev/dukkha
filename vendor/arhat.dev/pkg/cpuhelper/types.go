package cpuhelper

// Detect cpu details
func Detect() CPU {
	return detect()
}

type CPU interface{}

type CPUFeatures interface{}

type Vendor uint32

// nolint:gocyclo
func (vnd Vendor) String() string {
	switch vnd {
	case Vendor_AMD:
		return "AMD"
	case Vendor_VIA:
		return "VIA"
	case Vendor_Intel:
		return "Intel"
	case Vendor_Transmeta:
		return "Transmeta"
	case Vendor_NSC:
		return "NSC"
	case Vendor_NexGen:
		return "NexGen"
	case Vendor_Cyrix:
		return "Cyrix"
	case Vendor_Rise:
		return "Rise"
	case Vendor_SiS:
		return "SiS"
	case Vendor_DMP:
		return "DMP"
	case Vendor_UMC:
		return "UMC"
	case Vendor_Hygon:
		return "Hygon"
	case Vendor_RDC:
		return "RDC"
	case Vendor_MCST:
		return "MCST"

	// arm vendor
	case Vendor_Ampere:
		return "Ampere"
	case Vendor_ARM:
		return "ARM"
	case Vendor_Broadcom:
		return "Broadcom"
	case Vendor_Cavium:
		return "Cavium"
	case Vendor_DEC:
		return "DEC"
	case Vendor_Fujitsu:
		return "Fujitsu"
	case Vendor_Infineon:
		return "Infineon"
	case Vendor_Motorola:
		return "Motorola"
	case Vendor_NVIDIA:
		return "NVIDIA"
	case Vendor_AMCC:
		return "AMCC"
	case Vendor_Qualcomm:
		return "Qualcomm"
	case Vendor_Marvell:
		return "Marvell"
	case Vendor_Apple:
		return "Apple"

	// hypervisor
	case Vendor_KVM:
		return "KVM"
	case Vendor_MSVM:
		return "MSVM"
	case Vendor_VMware:
		return "VMware"
	case Vendor_XenHVM:
		return "XenHVM"
	case Vendor_Bhyve:
		return "Bhyve"
	case Vendor_QEMU:
		return "QEMU"
	case Vendor_Parallels:
		return "Parallels"
	case Vendor_QNXVM:
		return "QNXVM"
	case Vendor_ACRN:
		return "ACRN"
	case Vendor_AppleRosetta2:
		return "AppleRosetta2"

	default:
		return ""
	}
}

// nolint:revive
const (
	Vendor_Unknown Vendor = iota

	// x86 vendor
	Vendor_AMD
	Vendor_VIA
	Vendor_Intel
	Vendor_Transmeta
	Vendor_NSC // National Semiconductor
	Vendor_NexGen
	Vendor_Cyrix
	Vendor_Rise
	Vendor_SiS
	Vendor_DMP
	Vendor_UMC
	Vendor_Hygon
	Vendor_RDC
	Vendor_MCST

	// arm vendor
	Vendor_Ampere
	Vendor_ARM
	Vendor_Broadcom
	Vendor_Cavium
	Vendor_DEC
	Vendor_Fujitsu
	Vendor_Infineon
	Vendor_Motorola
	Vendor_NVIDIA
	Vendor_AMCC
	Vendor_Qualcomm
	Vendor_Marvell
	Vendor_Apple

	// hypervisor
	Vendor_KVM
	Vendor_MSVM
	Vendor_VMware
	Vendor_XenHVM
	Vendor_Bhyve
	Vendor_QEMU
	Vendor_Parallels
	Vendor_QNXVM
	Vendor_ACRN
	Vendor_AppleRosetta2
)

func HypervisorVendor(brand string) Vendor {
	switch brand {
	case "bhyve bhyve ":
		return Vendor_Bhyve
	case "KVMKVMKVMKVM":
		return Vendor_KVM
	case "TCGTCGTCGTCG":
		return Vendor_QEMU
	case "Microsoft Hv":
		return Vendor_MSVM
	case " lrpepyh  vr":
		return Vendor_Parallels
	case "VMwareVMware":
		return Vendor_VMware
	case "XenVMMXenVMM":
		return Vendor_XenHVM
	case "ACRNACRNACRN":
		return Vendor_ACRN
	case " QNXQVMBSQG ":
		return Vendor_QNXVM
	case "GenuineIntel":
		// TODO: this is not useful
		return Vendor_AppleRosetta2
	default:
		return Vendor_Unknown
	}
}
