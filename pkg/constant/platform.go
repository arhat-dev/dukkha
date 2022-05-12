package constant

const (
	// linux/gnu family
	Platform_Alpine = "alpine"
	Platform_Debian = "debian"
	Platform_Ubuntu = "ubuntu"
	Platform_GNU    = "gnu"

	// windows family
	Platform_WindowsNT    = "windows_nt"    // native windows (libc: MSVC)
	Platform_WindowsMINGW = "windows_mingw" // posix windows (libc: mingw)

	// darwin family
	Platform_MacOS   = "macos"
	Platform_iOS     = "ios"
	Platform_WatchOS = "watchos"

	// golang family
	Platform_Golang = "golang"
	Platform_Docker = "docker"
	Platform_OCI    = "oci"

	Platform_DockerHub = "dockerhub"

	// llvm family
	Platform_LLVM = "llvm"
	Platform_Zig  = "zig"
	Platform_Rust = "rust"

	Platform_QEMU = "qemu"
)

type platformID uint32

func (pid platformID) String() string {
	switch pid {
	case platformID_Alpine:
		return Platform_Alpine
	case platformID_Debian:
		return Platform_Debian
	case platformID_Ubuntu:
		return Platform_Ubuntu
	case platformID_GNU:
		return Platform_GNU

	case platformID_WindowsNT:
		return Platform_WindowsNT
	case platformID_WindowsMINGW:
		return Platform_WindowsMINGW

	case platformID_MacOS:
		return Platform_MacOS
	case platformID_iOS:
		return Platform_iOS
	case platformID_WatchOS:
		return Platform_WatchOS

	case platformID_Golang:
		return Platform_Golang
	case platformID_Docker:
		return Platform_Docker
	case platformID_OCI:
		return Platform_OCI

	case platformID_DockerHub:
		return Platform_DockerHub

	case platformID_LLVM:
		return Platform_LLVM
	case platformID_Zig:
		return Platform_Zig
	case platformID_Rust:
		return Platform_Rust

	case platformID_QEMU:
		return Platform_QEMU

	default:
		return "<unknown>"
	}
}

const (
	_unknown_platform platformID = iota

	platformID_Alpine
	platformID_Debian
	platformID_Ubuntu
	platformID_GNU

	platformID_WindowsNT
	platformID_WindowsMINGW

	platformID_MacOS
	platformID_iOS
	platformID_WatchOS

	platformID_Golang
	platformID_Docker
	platformID_OCI

	platformID_DockerHub

	platformID_LLVM
	platformID_Zig
	platformID_Rust

	platformID_QEMU

	platformID_COUNT
)

func platform_id_of(platform string) platformID {
	switch platform {
	case Platform_Alpine:
		return platformID_Alpine
	case Platform_Debian:
		return platformID_Debian
	case Platform_Ubuntu:
		return platformID_Ubuntu
	case Platform_GNU:
		return platformID_GNU

	case Platform_WindowsNT:
		return platformID_WindowsNT
	case Platform_WindowsMINGW:
		return platformID_WindowsMINGW

	case Platform_MacOS:
		return platformID_MacOS
	case Platform_iOS:
		return platformID_iOS
	case Platform_WatchOS:
		return platformID_WatchOS

	case Platform_Golang:
		return platformID_Golang
	case Platform_Docker:
		return platformID_Docker
	case Platform_OCI:
		return platformID_OCI

	case Platform_DockerHub:
		return platformID_DockerHub

	case Platform_LLVM:
		return platformID_LLVM
	case Platform_Zig:
		return platformID_Zig
	case Platform_Rust:
		return platformID_Rust

	case Platform_QEMU:
		return platformID_QEMU
	default:
		return _unknown_platform
	}
}
