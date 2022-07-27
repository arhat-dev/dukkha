package constant

import "arhat.dev/pkg/archconst"

// GetAlpineTripleName returns triple name for mArch on alpine, always assuming kernel=linux, libc=musl
// reference: https://more.musl.cc/10/x86_64-linux-musl/
// nolint:gocyclo
func GetAlpineTripleName(mArch string) (triple string, _ bool) {
	// NOTE: some arch value in triple name is different from alpine arch
	// 		 so we cannot do the same conversion as GNU/LLVM triple does
	aid := arch_id_of(mArch)
	if aid == _unknown_arch {
		return
	}

	switch aid {
	case archID_X86:
		return "i686-linux-musl", true
	case archID_X86_SF:
		return "i686-linux-musl", true

	case archID_AMD64:
		return "x86_64-linux-musl", true
	case archID_AMD64_V1:
		return "x86_64-linux-musl", true
	case archID_AMD64_V2:
		return "x86_64-linux-musl", true
	case archID_AMD64_V3:
		return "x86_64-linux-musl", true
	case archID_AMD64_V4:
		return "x86_64-linux-musl", true

	case archID_ARM:
		return "armv7l-linux-musleabihf", true
	case archID_ARM_V5:
		return "armv5l-linux-musleabi", true
	case archID_ARM_V6:
		return "armv6-linux-musleabihf", true
	case archID_ARM_V7:
		return "armv7l-linux-musleabihf", true

	case archID_ARM64:
		return "aarch64-linux-musl", true
	case archID_ARM64_V8:
		return "aarch64-linux-musl", true
	case archID_ARM64_V9:
		return "aarch64-linux-musl", true

	case archID_PPC:
		return "powerpc-linux-musl", true
	case archID_PPC_SF:
		return "powerpc-linux-muslsf", true
	case archID_PPC_LE:
		return "powerpcle-linux-musl", true
	case archID_PPC_LE_SF:
		return "powerpcle-linux-muslsf", true

	case archID_PPC64:
		return "powerpc64-linux-musl", true
	case archID_PPC64_LE:
		return "powerpc64le-linux-musl", true
	case archID_PPC64_V8:
		return "powerpc64-linux-musl", true
	case archID_PPC64_LE_V8:
		return "powerpc64le-linux-musl", true
	case archID_PPC64_V9:
		return "powerpc64-linux-musl", true
	case archID_PPC64_LE_V9:
		return "powerpc64le-linux-musl", true

	case archID_MIPS:
		return "mips-linux-musl", true
	case archID_MIPS_SF:
		return "mips-linux-muslsf", true
	case archID_MIPS_LE:
		return "mipsel-linux-musl", true
	case archID_MIPS_LE_SF:
		return "mipsel-linux-muslsf", true
	case archID_MIPS64:
		return "mips64-linux-musl", true
	case archID_MIPS64_SF:
		return "mips64-linux-musln32sf", true
	case archID_MIPS64_LE:
		return "mips64el-linux-musl", true
	case archID_MIPS64_LE_SF:
		return "mips64el-linux-musln32sf", true

	case archID_RISCV64:
		return "riscv64-linux-musl", true

	case archID_S390X:
		return "s390x-linux-musl", true

	case archID_IA64:
		return "", true
	default:
		return "", false
	}
}

// nolint:gocyclo
func GetDebianTripleName(mArch string, targetKernel, targetLibc string) (triple string, _ bool) {
	// TODO: adjust triple name according to target kernel
	_ = targetKernel

	aid := arch_id_of(mArch)
	if aid == _unknown_arch {
		return
	}

	// NOTE: some arch value in triple name is different from debian arch
	// 		 so we cannot do the same conversion as GNU/LLVM triple does
	switch targetLibc {
	case LIBC_MUSL:
		// https://packages.debian.org/buster/musl-dev
		// check list of files
		switch aid {
		case archID_X86:
			return "i386-linux-musl", true
		case archID_X86_SF:
			return "i386-linux-musl", true

		case archID_AMD64:
			return "x86_64-linux-musl", true
		case archID_AMD64_V1:
			return "x86_64-linux-musl", true
		case archID_AMD64_V2:
			return "x86_64-linux-musl", true
		case archID_AMD64_V3:
			return "x86_64-linux-musl", true
		case archID_AMD64_V4:
			return "x86_64-linux-musl", true

		case archID_ARM:
			return "arm-linux-musleabihf", true
		case archID_ARM_V5:
			return "arm-linux-musleabi", true
		case archID_ARM_V6:
			return "arm-linux-musleabi", true
		case archID_ARM_V7:
			return "arm-linux-musleabihf", true

		case archID_ARM64:
			return "aarch64-linux-musl", true
		case archID_ARM64_V8:
			return "aarch64-linux-musl", true
		case archID_ARM64_V9:
			return "aarch64-linux-musl", true

		case archID_MIPS:
			return "mips-linux-musl", true
		case archID_MIPS_SF:
			return "mips-linux-musl", true
		case archID_MIPS_LE:
			return "mipsel-linux-musl", true
		case archID_MIPS_LE_SF:
			return "mipsel-linux-musl", true
		case archID_MIPS64:
			return "mips64-linux-musl", true
		case archID_MIPS64_SF:
			return "mips64-linux-musl", true
		case archID_MIPS64_LE:
			return "mips64el-linux-musl", true
		case archID_MIPS64_LE_SF:
			return "mips64el-linux-musl", true

		case archID_S390X:
			return "s390x-linux-musl", true

		// http://ftp.ports.debian.org/debian-ports/pool-riscv64/main/m/musl/
		// download one musl-dev package
		// list package contents with following commands
		//
		// $ ar -x musl-dev_1.2.2-3_riscv64.deb
		// $ tar -tvf data.tar.xz
		case archID_RISCV64:
			return "riscv64-linux-musl", true

		case archID_IA64:
			return "", true

		case archID_PPC:
			return "", true
		case archID_PPC_SF:
			return "", true
		case archID_PPC_LE:
			return "", true
		case archID_PPC_LE_SF:
			return "", true

		case archID_PPC64:
			return "", true
		case archID_PPC64_LE:
			return "", true
		case archID_PPC64_V8:
			return "", true
		case archID_PPC64_LE_V8:
			return "", true
		case archID_PPC64_V9:
			return "", true
		case archID_PPC64_LE_V9:
			return "", true
		default:
			return "", false
		}
	case LIBC_MSVC:
		switch aid {
		// https://packages.debian.org/buster/mingw-w64-i686-dev
		// check list of files
		case archID_X86:
			return "i686-w64-mingw32", true
		case archID_X86_SF:
			return "i686-w64-mingw32", true
		// https://packages.debian.org/buster/mingw-w64-x86-64-dev
		// check list of files
		case archID_AMD64:
			return "x86_64-w64-mingw32", true
		case archID_AMD64_V1:
			return "x86_64-w64-mingw32", true
		case archID_AMD64_V2:
			return "x86_64-w64-mingw32", true
		case archID_AMD64_V3:
			return "x86_64-w64-mingw32", true
		case archID_AMD64_V4:
			return "x86_64-w64-mingw32", true

		case archID_IA64:
			return "", true

		case archID_ARM:
			return "", true
		case archID_ARM_V5:
			return "", true
		case archID_ARM_V6:
			return "", true
		case archID_ARM_V7:
			return "", true

		case archID_ARM64:
			return "", true
		case archID_ARM64_V8:
			return "", true
		case archID_ARM64_V9:
			return "", true

		case archID_PPC:
			return "", true
		case archID_PPC_SF:
			return "", true
		case archID_PPC_LE:
			return "", true
		case archID_PPC_LE_SF:
			return "", true

		case archID_PPC64:
			return "", true
		case archID_PPC64_LE:
			return "", true
		case archID_PPC64_V8:
			return "", true
		case archID_PPC64_LE_V8:
			return "", true
		case archID_PPC64_V9:
			return "", true
		case archID_PPC64_LE_V9:
			return "", true

		case archID_MIPS:
			return "", true
		case archID_MIPS_SF:
			return "", true
		case archID_MIPS_LE:
			return "", true
		case archID_MIPS_LE_SF:
			return "", true
		case archID_MIPS64:
			return "", true
		case archID_MIPS64_SF:
			return "", true
		case archID_MIPS64_LE:
			return "", true
		case archID_MIPS64_LE_SF:
			return "", true

		case archID_RISCV64:
			return "", true
		case archID_S390X:
			return "", true
		default:
			return "", false
		}

	case LIBC_GNU:
		fallthrough
	default:
		switch aid {
		case archID_X86:
			return "i686-linux-gnu", true
		case archID_X86_SF:
			return "i686-linux-gnu", true

		case archID_AMD64:
			return "x86_64-linux-gnu", true
		case archID_AMD64_V1:
			return "x86_64-linux-gnu", true
		case archID_AMD64_V2:
			return "x86_64-linux-gnu", true
		case archID_AMD64_V3:
			return "x86_64-linux-gnu", true
		case archID_AMD64_V4:
			return "x86_64-linux-gnu", true

		case archID_ARM:
			return "arm-linux-gnueabihf", true
		case archID_ARM_V5:
			return "arm-linux-gnueabi", true
		case archID_ARM_V6:
			return "arm-linux-gnueabi", true
		case archID_ARM_V7:
			return "arm-linux-gnueabihf", true

		case archID_ARM64:
			return "aarch64-linux-gnu", true
		case archID_ARM64_V8:
			return "aarch64-linux-gnu", true
		case archID_ARM64_V9:
			return "aarch64-linux-gnu", true

		case archID_PPC64:
			return "powerpc64-linux-gnu", true
		case archID_PPC64_LE:
			return "powerpc64le-linux-gnu", true
		case archID_PPC64_V8:
			return "powerpc64-linux-gnu", true
		case archID_PPC64_LE_V8:
			return "powerpc64le-linux-gnu", true
		case archID_PPC64_V9:
			return "powerpc64-linux-gnu", true
		case archID_PPC64_LE_V9:
			return "powerpc64le-linux-gnu", true

		case archID_MIPS:
			return "mips-linux-gnu", true
		case archID_MIPS_SF:
			return "mips-linux-gnu", true
		case archID_MIPS_LE:
			return "mipsel-linux-gnu", true
		case archID_MIPS_LE_SF:
			return "mipsel-linux-gnu", true
		case archID_MIPS64:
			return "mips64-linux-gnuabi64", true
		case archID_MIPS64_SF:
			return "mips64-linux-gnuabi64", true
		case archID_MIPS64_LE:
			return "mips64el-linux-gnuabi64", true
		case archID_MIPS64_LE_SF:
			return "mips64el-linux-gnuabi64", true

		case archID_RISCV64:
			return "riscv64-linux-gnu", true
		case archID_S390X:
			return "s390x-linux-gnu", true

		case archID_IA64:
			return "", true

		case archID_PPC:
			return "", true
		case archID_PPC_SF:
			return "", true
		case archID_PPC_LE:
			return "", true
		case archID_PPC_LE_SF:
			return "", true
		default:
			return "", false
		}
	}
}

// Ref:
// - https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/data/cputable
// - https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/data/ostable
func GetGNUTripleName(mArch, targetKernel, targetLibc string) (triple string, _ bool) {
	_ = targetKernel

	aid := arch_id_of(mArch)
	if aid == _unknown_arch {
		return
	}

	arch := archMapping[aid][platformID_GNU]
	abi := targetLibc

	switch targetLibc {
	case LIBC_MUSL:
		switch aid {
		case archID_ARM_V5, archID_ARM_V6:
			abi = "musleabi"
		case archID_ARM_V7, archID_ARM:
			abi = "musleabihf"
		default:
			abi = "musl"
		}
	case LIBC_MSVC:
		// TODO
		_ = abi
	case LIBC_GNU:
		fallthrough
	default:
		switch aid {
		case archID_ARM_V5, archID_ARM_V6:
			abi = "gnueabi"
		case archID_ARM_V7, archID_ARM:
			abi = "gnueabihf"
		case archID_MIPS64, archID_MIPS64_SF,
			archID_MIPS64_LE, archID_MIPS64_LE_SF:
			abi = "gnuabi64"
		default:
			abi = "gnu"
		}
	}

	return arch + "-linux-" + abi, true
}

// ref: https://llvm.org/doxygen/Triple_8h_source.html
func GetLLVMTripleName(mArch, targetKernel, targetLibc string) (triple string, _ bool) {
	aid := arch_id_of(mArch)
	if aid == _unknown_arch {
		return
	}

	arch := archMapping[aid][platformID_LLVM]
	sys := targetKernel
	switch targetKernel {
	case KERNEL_Windows:
		sys = "windows"
	case KERNEL_Linux:
		sys = "linux"
	case KERNEL_Darwin:
		sys = "darwin"
	case KERNEL_FreeBSD:
		sys = "freebsd"
	case KERNEL_NetBSD:
		sys = "darwin"
	case KERNEL_OpenBSD:
		sys = "openbsd"
	case KERNEL_Solaris:
		sys = "solaris"
	case KERNEL_Illumos:
		sys = "illumos"
	case KERNEL_JavaScript:
		sys = "js"
	case KERNEL_Aix:
		sys = "aix"
	case KERNEL_Android:
		sys = "android"
	case KERNEL_iOS:
		sys = "ios"
	case KERNEL_Plan9:
		sys = "plan9"
	default:
		// sys = targetKernel
	}

	abi := targetLibc
	switch targetLibc {
	case LIBC_MUSL:
		switch aid {
		case archID_ARM_V5, archID_ARM_V6:
			abi = "musleabi"
		case archID_ARM_V7, archID_ARM:
			abi = "musleabihf"
		default:
			abi = "musl"
		}
	case LIBC_MSVC:
		switch aid {
		case _unknown_arch: // TODO: add special cases
			_ = abi
		default:
			abi = "msvc"
		}
	case LIBC_GNU:
		fallthrough
	default:
		switch aid {
		case archID_ARM_V5, archID_ARM_V6:
			abi = "gnueabi"
		case archID_ARM_V7, archID_ARM:
			abi = "gnueabihf"
		case archID_MIPS64, archID_MIPS64_SF,
			archID_MIPS64_LE, archID_MIPS64_LE_SF:
			// TODO: is it gnu or gnuabi64?
			abi = "gnuabi64"
		default:
			abi = "gnu"
		}
	}

	// <arch>-<vendor>-<sys>-<abi>
	return arch + "-unknown-" + sys + "-" + abi, true
}

// GetZigTripleName returns a valid value to `--target` of `zig cc`
//
// targetKernel defaults to linux kernel
//
// targetLibc defaults to musl
func GetZigTripleName(mArch, targetKernel, targetLibc string) (triple string, _ bool) {
	aid := arch_id_of(mArch)
	if aid == _unknown_arch {
		return
	}

	spec, ok := archconst.Parse[byte](mArch)
	if !ok {
		return
	}

	var abi string
	switch targetLibc {
	case LIBC_MUSL, "":
		switch spec.Name {
		default:
			abi = "musl"
		case archconst.ARCH_ARM:
			if spec.SoftFloat {
				abi = "musleabi"
			} else {
				abi = "musleabihf"
			}
		}
	case LIBC_GNU:
		switch spec.Name {
		default:
			abi = "gnu"
		case archconst.ARCH_ARM, archconst.ARCH_MIPS, archconst.ARCH_PPC:
			// arm-linux-gnueabi
			// arm-linux-gnueabihf
			// mipsel-linux-gnueabi
			// mipsel-linux-gnueabihf
			// mips-linux-gnueabi
			// mips-linux-gnueabihf
			// powerpc-linux-gnueabi
			// powerpc-linux-gnueabihf

			if spec.SoftFloat {
				abi = "gnueabi"
			} else {
				abi = "gnueabihf"
			}
		case archconst.ARCH_MIPS64:
			// mips64el-linux-gnuabi64
			// mips64el-linux-gnuabin32
			// mips64-linux-gnuabi64
			// mips64-linux-gnuabin32

			// TODO: is it correct?
			if spec.SoftFloat {
				abi = "gnuabin32"
			} else {
				abi = "gnuabi64"
			}
		}
	case LIBC_MSVC:
		// currently no msvc support in zig
		fallthrough
	default:
		return
	}

	var os string
	switch targetKernel {
	case KERNEL_Linux, "":
		os = "linux"
	case KERNEL_Darwin:
		os = "macos"
	case KERNEL_Windows:
		os = "windows"
	default:
		return
	}

	arch := archMapping[aid][platformID_Zig]
	switch aid {
	case archID_ARM_V5, archID_ARM_V6, archID_ARM_V7:
		if spec.LittleEndian {
			arch = "arm"
		} else {
			arch = "armeb"
		}
	case archID_PPC_LE, archID_PPC_LE_SF:
		arch = "powerpc"
	}

	return arch + "-" + os + "-" + abi, true
}
