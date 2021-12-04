package templateutils

import "arhat.dev/dukkha/pkg/dukkha"

func createHostNS(rc dukkha.RenderingContext) *hostNS {
	return &hostNS{rc: rc}
}

type hostNS struct {
	rc dukkha.RenderingContext
}

// Arch get HOST_ARCH
func (ns *hostNS) Arch() string {
	return ns.rc.HostArch()
}

// Kernel get HOST_KERNEL
func (ns *hostNS) Kernel() string {
	return ns.rc.HostKernel()
}

// KernelVersion get HOST_KERNEL_VERSION
func (ns *hostNS) KernelVersion() string {
	return ns.rc.HostKernelVersion()
}

// OS get HOST_OS
func (ns *hostNS) OS() string {
	return ns.rc.HostOS()
}

// OSVersion get HOST_OS_VERSION
func (ns *hostNS) OSVersion() string {
	return ns.rc.HostOSVersion()
}
