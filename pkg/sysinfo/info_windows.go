// +build windows
//go:build windows

package sysinfo

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/version"
)

var (
	arch          string
	osName        string
	kernelVersion string
)

func Arch() string {
	return arch
}

func OSName() string {
	return osName
}

func OSVersion() string {
	// TODO: check os version using syscall
	return ""
}

func KernelVersion() string {
	return kernelVersion
}

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

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer func() {
		_ = k.Close()
	}()

	{
		// get product name as os name
		osName, _, _ = k.GetStringValue("ProductName")
	}

	{
		// build kernel version
		buildNumber, _, err := k.GetStringValue("CurrentBuildNumber")
		if err != nil {
			return
		}

		majorVersionNumber, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
		if err != nil {
			return
		}

		minorVersionNumber, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
		if err != nil {
			return
		}

		revision, _, err := k.GetIntegerValue("UBR")
		if err != nil {
			return
		}

		kernelVersion = fmt.Sprintf("%d.%d.%s.%d", majorVersionNumber, minorVersionNumber, buildNumber, revision)
	}
}
