// +build windows
//go:build windows

package sysinfo

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/version"
)

func init() {
	defer func() {
		recover()
	}()

	kernel32dll := windows.NewLazySystemDLL("kernel32.dll")
	procGetNativeSystemInfo := kernel32dll.NewProc("GetNativeSystemInfo")

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

	procGetNativeSystemInfo.Call(uintptr(unsafe.Pointer(&info)))
	cpuArch := uint(info.wProcessorArchitecture)

	const (
		PROCESSOR_ARCHITECTURE_INTEL = 0
		PROCESSOR_ARCHITECTURE_ARM   = 5
		PROCESSOR_ARCHITECTURE_ARM64 = 12
		PROCESSOR_ARCHITECTURE_IA64  = 6
		PROCESSOR_ARCHITECTURE_AMD64 = 9
	)

	switch cpuArch {
	case PROCESSOR_ARCHITECTURE_INTEL:
		arch = constant.ARCH_X86
	case PROCESSOR_ARCHITECTURE_ARM:
		// usually armv7, can be armv6/armv5
		arch = version.Arch()
	case PROCESSOR_ARCHITECTURE_ARM64:
		arch = constant.ARCH_ARM64
	case PROCESSOR_ARCHITECTURE_IA64:
		arch = constant.ARCH_IA64
	case PROCESSOR_ARCHITECTURE_AMD64:
		arch = constant.ARCH_AMD64
	default:
		arch = version.Arch()
	}
}

var arch string

func Arch() string {
	return arch
}
