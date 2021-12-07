package cpuhelper

import (
	"syscall"
	"unsafe"

	"arhat.dev/pkg/archconst"
	"arhat.dev/pkg/versionhelper"
)

func Arch() string {
	hostArch := ArchByCPUFeatures()
	if len(hostArch) != 0 {
		return hostArch
	}

	kernel32dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return versionhelper.Arch()
	}

	procGetNativeSystemInfo, err := kernel32dll.FindProc("GetNativeSystemInfo")
	if err != nil {
		return versionhelper.Arch()
	}

	type systemInfo struct {
		wProcessorArchitecture      uint16
		wReserved                   uint16
		dwPageSize                  uint32
		lpMinimumApplicationAddress uintptr
		lpMaximumApplicationAddress uintptr
		dwActiveProcessorMask       uintptr
		dwNumberOfProcessors        uint32
		dwProcessorType             uint32
		dwAllocationGranularity     uint32
		wProcessorLevel             uint16
		wProcessorRevision          uint16
	}

	var info systemInfo

	const (
		PROCESSOR_ARCHITECTURE_INTEL = 0
		PROCESSOR_ARCHITECTURE_ARM   = 5
		PROCESSOR_ARCHITECTURE_ARM64 = 12
		PROCESSOR_ARCHITECTURE_IA64  = 6
		PROCESSOR_ARCHITECTURE_AMD64 = 9

		_unused_unknown = 128
	)

	info.wProcessorArchitecture = _unused_unknown

	// ref: https://docs.microsoft.com/en-us/windows/win32/api/sysinfoapi/nf-sysinfoapi-getnativesysteminfo
	r1, _, err := procGetNativeSystemInfo.Call(uintptr(unsafe.Pointer(&info)))
	if r1 != 0 {
		_ = err
		return versionhelper.Arch()
	}

	cpuArch := uint(info.wProcessorArchitecture)

	switch cpuArch {
	case PROCESSOR_ARCHITECTURE_INTEL:
		// zero, be default value, prefer value in versionhelper
		return archconst.ARCH_X86
	case PROCESSOR_ARCHITECTURE_ARM:
		// usually armv7, can be armv6/armv5
		return versionhelper.Arch()
	case PROCESSOR_ARCHITECTURE_ARM64:
		return archconst.ARCH_ARM64
	case PROCESSOR_ARCHITECTURE_IA64:
		return archconst.ARCH_IA64
	case PROCESSOR_ARCHITECTURE_AMD64:
		return archconst.ARCH_AMD64_V1
	default:
		return versionhelper.Arch()
	}
}
