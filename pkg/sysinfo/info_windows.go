//go:build windows
// +build windows

package sysinfo

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func OSName() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer func() { _ = k.Close() }()

	// get product name as os name
	osName, _, _ := k.GetStringValue("ProductName")
	return osName
}

func OSVersion() string {
	// TODO: check os version using syscall
	return ""
}

func KernelVersion() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer func() { _ = k.Close() }()

	// build kernel version
	buildNumber, _, err := k.GetStringValue("CurrentBuildNumber")
	if err != nil {
		return ""
	}

	majorVersionNumber, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return ""
	}

	minorVersionNumber, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	if err != nil {
		return ""
	}

	revision, _, err := k.GetIntegerValue("UBR")
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%d.%d.%s.%d", majorVersionNumber, minorVersionNumber, buildNumber, revision)
}
