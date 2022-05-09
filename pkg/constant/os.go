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
