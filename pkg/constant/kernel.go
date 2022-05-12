package constant

// Kernel names
// currently they are defined the same as GOOS values
//
// nolint:revive
const (
	KERNEL_Windows    = "windows"
	KERNEL_Linux      = "linux"
	KERNEL_Darwin     = "darwin"
	KERNEL_FreeBSD    = "freebsd"
	KERNEL_NetBSD     = "netbsd"
	KERNEL_OpenBSD    = "openbsd"
	KERNEL_Solaris    = "solaris"
	KERNEL_Illumos    = "illumos"
	KERNEL_JavaScript = "js"
	KERNEL_Aix        = "aix"
	KERNEL_Android    = "android"
	KERNEL_iOS        = "ios"
	KERNEL_Plan9      = "plan9"
)

type kernelID uint32

func (kid kernelID) String() string {
	switch kid {
	case kernelID_Windows:
		return KERNEL_Windows
	case kernelID_Linux:
		return KERNEL_Linux
	case kernelID_Darwin:
		return KERNEL_Darwin
	case kernelID_FreeBSD:
		return KERNEL_FreeBSD
	case kernelID_NetBSD:
		return KERNEL_NetBSD
	case kernelID_OpenBSD:
		return KERNEL_OpenBSD
	case kernelID_Solaris:
		return KERNEL_Solaris
	case kernelID_Illumos:
		return KERNEL_Illumos
	case kernelID_JavaScript:
		return KERNEL_JavaScript
	case kernelID_Aix:
		return KERNEL_Aix
	case kernelID_Android:
		return KERNEL_Android
	case kernelID_iOS:
		return KERNEL_iOS
	case kernelID_Plan9:
		return KERNEL_Plan9
	default:
		return "<unknown>"
	}
}

const (
	_unknown_kernel kernelID = iota

	kernelID_Windows
	kernelID_Linux
	kernelID_Darwin
	kernelID_FreeBSD
	kernelID_NetBSD
	kernelID_OpenBSD
	kernelID_Solaris
	kernelID_Illumos
	kernelID_JavaScript
	kernelID_Aix
	kernelID_Android
	kernelID_iOS
	kernelID_Plan9

	kernelID_COUNT
)

func kernel_id_of(kernel string) kernelID {
	switch kernel {
	case KERNEL_Windows:
		return kernelID_Windows
	case KERNEL_Linux:
		return kernelID_Linux
	case KERNEL_Darwin:
		return kernelID_Darwin
	case KERNEL_FreeBSD:
		return kernelID_FreeBSD
	case KERNEL_NetBSD:
		return kernelID_NetBSD
	case KERNEL_OpenBSD:
		return kernelID_OpenBSD
	case KERNEL_Solaris:
		return kernelID_Solaris
	case KERNEL_Illumos:
		return kernelID_Illumos
	case KERNEL_JavaScript:
		return kernelID_JavaScript
	case KERNEL_Aix:
		return kernelID_Aix
	case KERNEL_Android:
		return kernelID_Android
	case KERNEL_iOS:
		return kernelID_iOS
	case KERNEL_Plan9:
		return kernelID_Plan9
	default:
		return _unknown_kernel
	}
}
