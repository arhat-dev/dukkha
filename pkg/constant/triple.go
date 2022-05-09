package constant

import "arhat.dev/pkg/archconst"

// GetAlpineTripleName of matrix arch
// reference: https://more.musl.cc/10/x86_64-linux-musl/
func GetAlpineTripleName(mArch string) (string, bool) {
	// NOTE: some arch value in triple name is different from alpine arch
	// 		 so we cannot do the same conversion as GNU/LLVM triple does
	switch mArch {
	case archconst.ARCH_X86:
		return "i686-linux-musl", true
	case archconst.ARCH_X86_SF:
		return "i686-linux-musl", true

	case archconst.ARCH_AMD64:
		return "x86_64-linux-musl", true
	case archconst.ARCH_AMD64_V1:
		return "x86_64-linux-musl", true
	case archconst.ARCH_AMD64_V2:
		return "x86_64-linux-musl", true
	case archconst.ARCH_AMD64_V3:
		return "x86_64-linux-musl", true
	case archconst.ARCH_AMD64_V4:
		return "x86_64-linux-musl", true

	case archconst.ARCH_ARM:
		return "armv7l-linux-musleabihf", true
	case archconst.ARCH_ARM_V5:
		return "armv5l-linux-musleabi", true
	case archconst.ARCH_ARM_V6:
		return "armv6-linux-musleabihf", true
	case archconst.ARCH_ARM_V7:
		return "armv7l-linux-musleabihf", true

	case archconst.ARCH_ARM64:
		return "aarch64-linux-musl", true
	case archconst.ARCH_ARM64_V8:
		return "aarch64-linux-musl", true
	case archconst.ARCH_ARM64_V9:
		return "aarch64-linux-musl", true

	case archconst.ARCH_PPC:
		return "powerpc-linux-musl", true
	case archconst.ARCH_PPC_SF:
		return "powerpc-linux-muslsf", true
	case archconst.ARCH_PPC_LE:
		return "powerpcle-linux-musl", true
	case archconst.ARCH_PPC_LE_SF:
		return "powerpcle-linux-muslsf", true

	case archconst.ARCH_PPC64:
		return "powerpc64-linux-musl", true
	case archconst.ARCH_PPC64_LE:
		return "powerpc64le-linux-musl", true
	case archconst.ARCH_PPC64_V8:
		return "powerpc64-linux-musl", true
	case archconst.ARCH_PPC64_LE_V8:
		return "powerpc64le-linux-musl", true
	case archconst.ARCH_PPC64_V9:
		return "powerpc64-linux-musl", true
	case archconst.ARCH_PPC64_LE_V9:
		return "powerpc64le-linux-musl", true

	case archconst.ARCH_MIPS:
		return "mips-linux-musl", true
	case archconst.ARCH_MIPS_SF:
		return "mips-linux-muslsf", true
	case archconst.ARCH_MIPS_LE:
		return "mipsel-linux-musl", true
	case archconst.ARCH_MIPS_LE_SF:
		return "mipsel-linux-muslsf", true
	case archconst.ARCH_MIPS64:
		return "mips64-linux-musl", true
	case archconst.ARCH_MIPS64_SF:
		return "mips64-linux-musln32sf", true
	case archconst.ARCH_MIPS64_LE:
		return "mips64el-linux-musl", true
	case archconst.ARCH_MIPS64_LE_SF:
		return "mips64el-linux-musln32sf", true

	case archconst.ARCH_RISCV64:
		return "riscv64-linux-musl", true

	case archconst.ARCH_S390X:
		return "s390x-linux-musl", true

	case archconst.ARCH_IA64:
		return "", true
	default:
		return "", false
	}
}

func GetDebianTripleName(mArch string, targetKernel, targetLibc string) (v string, ok bool) {
	// TODO: adjust triple name according to target kernel
	_ = targetKernel

	// NOTE: some arch value in triple name is different from debian arch
	// 		 so we cannot do the same conversion as GNU/LLVM triple does
	switch targetLibc {
	case LIBC_MUSL:
		// https://packages.debian.org/buster/musl-dev
		// check list of files
		switch mArch {
		case archconst.ARCH_X86:
			return "i386-linux-musl", true
		case archconst.ARCH_X86_SF:
			return "i386-linux-musl", true

		case archconst.ARCH_AMD64:
			return "x86_64-linux-musl", true
		case archconst.ARCH_AMD64_V1:
			return "x86_64-linux-musl", true
		case archconst.ARCH_AMD64_V2:
			return "x86_64-linux-musl", true
		case archconst.ARCH_AMD64_V3:
			return "x86_64-linux-musl", true
		case archconst.ARCH_AMD64_V4:
			return "x86_64-linux-musl", true

		case archconst.ARCH_ARM:
			return "arm-linux-musleabihf", true
		case archconst.ARCH_ARM_V5:
			return "arm-linux-musleabi", true
		case archconst.ARCH_ARM_V6:
			return "arm-linux-musleabi", true
		case archconst.ARCH_ARM_V7:
			return "arm-linux-musleabihf", true

		case archconst.ARCH_ARM64:
			return "aarch64-linux-musl", true
		case archconst.ARCH_ARM64_V8:
			return "aarch64-linux-musl", true
		case archconst.ARCH_ARM64_V9:
			return "aarch64-linux-musl", true

		case archconst.ARCH_MIPS:
			return "mips-linux-musl", true
		case archconst.ARCH_MIPS_SF:
			return "mips-linux-musl", true
		case archconst.ARCH_MIPS_LE:
			return "mipsel-linux-musl", true
		case archconst.ARCH_MIPS_LE_SF:
			return "mipsel-linux-musl", true
		case archconst.ARCH_MIPS64:
			return "mips64-linux-musl", true
		case archconst.ARCH_MIPS64_SF:
			return "mips64-linux-musl", true
		case archconst.ARCH_MIPS64_LE:
			return "mips64el-linux-musl", true
		case archconst.ARCH_MIPS64_LE_SF:
			return "mips64el-linux-musl", true

		case archconst.ARCH_S390X:
			return "s390x-linux-musl", true

		// http://ftp.ports.debian.org/debian-ports/pool-riscv64/main/m/musl/
		// download one musl-dev package
		// list package contents with following commands
		//
		// $ ar -x musl-dev_1.2.2-3_riscv64.deb
		// $ tar -tvf data.tar.xz
		case archconst.ARCH_RISCV64:
			return "riscv64-linux-musl", true

		case archconst.ARCH_IA64:
			return "", true

		case archconst.ARCH_PPC:
			return "", true
		case archconst.ARCH_PPC_SF:
			return "", true
		case archconst.ARCH_PPC_LE:
			return "", true
		case archconst.ARCH_PPC_LE_SF:
			return "", true

		case archconst.ARCH_PPC64:
			return "", true
		case archconst.ARCH_PPC64_LE:
			return "", true
		case archconst.ARCH_PPC64_V8:
			return "", true
		case archconst.ARCH_PPC64_LE_V8:
			return "", true
		case archconst.ARCH_PPC64_V9:
			return "", true
		case archconst.ARCH_PPC64_LE_V9:
			return "", true
		default:
			return "", false
		}
	case LIBC_MSVC:
		switch mArch {
		// https://packages.debian.org/buster/mingw-w64-i686-dev
		// check list of files
		case archconst.ARCH_X86:
			return "i686-w64-mingw32", true
		case archconst.ARCH_X86_SF:
			return "i686-w64-mingw32", true
		// https://packages.debian.org/buster/mingw-w64-x86-64-dev
		// check list of files
		case archconst.ARCH_AMD64:
			return "x86_64-w64-mingw32", true
		case archconst.ARCH_AMD64_V1:
			return "x86_64-w64-mingw32", true
		case archconst.ARCH_AMD64_V2:
			return "x86_64-w64-mingw32", true
		case archconst.ARCH_AMD64_V3:
			return "x86_64-w64-mingw32", true
		case archconst.ARCH_AMD64_V4:
			return "x86_64-w64-mingw32", true

		case archconst.ARCH_IA64:
			return "", true

		case archconst.ARCH_ARM:
			return "", true
		case archconst.ARCH_ARM_V5:
			return "", true
		case archconst.ARCH_ARM_V6:
			return "", true
		case archconst.ARCH_ARM_V7:
			return "", true

		case archconst.ARCH_ARM64:
			return "", true
		case archconst.ARCH_ARM64_V8:
			return "", true
		case archconst.ARCH_ARM64_V9:
			return "", true

		case archconst.ARCH_PPC:
			return "", true
		case archconst.ARCH_PPC_SF:
			return "", true
		case archconst.ARCH_PPC_LE:
			return "", true
		case archconst.ARCH_PPC_LE_SF:
			return "", true

		case archconst.ARCH_PPC64:
			return "", true
		case archconst.ARCH_PPC64_LE:
			return "", true
		case archconst.ARCH_PPC64_V8:
			return "", true
		case archconst.ARCH_PPC64_LE_V8:
			return "", true
		case archconst.ARCH_PPC64_V9:
			return "", true
		case archconst.ARCH_PPC64_LE_V9:
			return "", true

		case archconst.ARCH_MIPS:
			return "", true
		case archconst.ARCH_MIPS_SF:
			return "", true
		case archconst.ARCH_MIPS_LE:
			return "", true
		case archconst.ARCH_MIPS_LE_SF:
			return "", true
		case archconst.ARCH_MIPS64:
			return "", true
		case archconst.ARCH_MIPS64_SF:
			return "", true
		case archconst.ARCH_MIPS64_LE:
			return "", true
		case archconst.ARCH_MIPS64_LE_SF:
			return "", true

		case archconst.ARCH_RISCV64:
			return "", true
		case archconst.ARCH_S390X:
			return "", true
		default:
			return "", false
		}

	case LIBC_GNU:
		fallthrough
	default:
		switch mArch {
		case archconst.ARCH_X86:
			return "i686-linux-gnu", true
		case archconst.ARCH_X86_SF:
			return "i686-linux-gnu", true

		case archconst.ARCH_AMD64:
			return "x86_64-linux-gnu", true
		case archconst.ARCH_AMD64_V1:
			return "x86_64-linux-gnu", true
		case archconst.ARCH_AMD64_V2:
			return "x86_64-linux-gnu", true
		case archconst.ARCH_AMD64_V3:
			return "x86_64-linux-gnu", true
		case archconst.ARCH_AMD64_V4:
			return "x86_64-linux-gnu", true

		case archconst.ARCH_ARM:
			return "arm-linux-gnueabihf", true
		case archconst.ARCH_ARM_V5:
			return "arm-linux-gnueabi", true
		case archconst.ARCH_ARM_V6:
			return "arm-linux-gnueabi", true
		case archconst.ARCH_ARM_V7:
			return "arm-linux-gnueabihf", true

		case archconst.ARCH_ARM64:
			return "aarch64-linux-gnu", true
		case archconst.ARCH_ARM64_V8:
			return "aarch64-linux-gnu", true
		case archconst.ARCH_ARM64_V9:
			return "aarch64-linux-gnu", true

		case archconst.ARCH_PPC64:
			return "powerpc64-linux-gnu", true
		case archconst.ARCH_PPC64_LE:
			return "powerpc64le-linux-gnu", true
		case archconst.ARCH_PPC64_V8:
			return "powerpc64-linux-gnu", true
		case archconst.ARCH_PPC64_LE_V8:
			return "powerpc64le-linux-gnu", true
		case archconst.ARCH_PPC64_V9:
			return "powerpc64-linux-gnu", true
		case archconst.ARCH_PPC64_LE_V9:
			return "powerpc64le-linux-gnu", true

		case archconst.ARCH_MIPS:
			return "mips-linux-gnu", true
		case archconst.ARCH_MIPS_SF:
			return "mips-linux-gnu", true
		case archconst.ARCH_MIPS_LE:
			return "mipsel-linux-gnu", true
		case archconst.ARCH_MIPS_LE_SF:
			return "mipsel-linux-gnu", true
		case archconst.ARCH_MIPS64:
			return "mips64-linux-gnuabi64", true
		case archconst.ARCH_MIPS64_SF:
			return "mips64-linux-gnuabi64", true
		case archconst.ARCH_MIPS64_LE:
			return "mips64el-linux-gnuabi64", true
		case archconst.ARCH_MIPS64_LE_SF:
			return "mips64el-linux-gnuabi64", true

		case archconst.ARCH_RISCV64:
			return "riscv64-linux-gnu", true
		case archconst.ARCH_S390X:
			return "s390x-linux-gnu", true

		case archconst.ARCH_IA64:
			return "", true

		case archconst.ARCH_PPC:
			return "", true
		case archconst.ARCH_PPC_SF:
			return "", true
		case archconst.ARCH_PPC_LE:
			return "", true
		case archconst.ARCH_PPC_LE_SF:
			return "", true
		default:
			return "", false
		}
	}
}

// Ref:
// - https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/data/cputable
// - https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/data/ostable
func GetGNUTripleName(mArch, targetKernel, targetLibc string) (string, bool) {
	_ = targetKernel

	arch, ok := GetArch(Platform_GNU, mArch)
	if !ok {
		return "", false
	}

	abi := targetLibc
	switch targetLibc {
	case LIBC_MUSL:
		switch mArch {
		case archconst.ARCH_ARM_V5, archconst.ARCH_ARM_V6:
			abi = "musleabi"
		case archconst.ARCH_ARM_V7, archconst.ARCH_ARM:
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
		switch mArch {
		case archconst.ARCH_ARM_V5, archconst.ARCH_ARM_V6:
			abi = "gnueabi"
		case archconst.ARCH_ARM_V7, archconst.ARCH_ARM:
			abi = "gnueabihf"
		case archconst.ARCH_MIPS64, archconst.ARCH_MIPS64_SF,
			archconst.ARCH_MIPS64_LE, archconst.ARCH_MIPS64_LE_SF:
			abi = "gnuabi64"
		default:
			abi = "gnu"
		}
	}

	return arch + "-linux-" + abi, ok
}

// ref: https://llvm.org/doxygen/Triple_8h_source.html
func GetLLVMTripleName(mArch, targetKernel, targetLibc string) (string, bool) {
	arch, ok := GetArch(Platform_LLVM, mArch)
	if !ok {
		return "", false
	}

	sys := targetKernel
	switch targetKernel {
	case KERNEL_WINDOWS:
		sys = "windows"
	case KERNEL_LINUX:
		sys = "linux"
	case KERNEL_DARWIN:
		sys = "darwin"
	case KERNEL_FREEBSD:
		sys = "freebsd"
	case KERNEL_NETBSD:
		sys = "darwin"
	case KERNEL_OPENBSD:
		sys = "openbsd"
	case KERNEL_SOLARIS:
		sys = "solaris"
	case KERNEL_ILLUMOS:
		sys = "illumos"
	case KERNEL_JAVASCRIPT:
		sys = "js"
	case KERNEL_AIX:
		sys = "aix"
	case KERNEL_ANDROID:
		sys = "android"
	case KERNEL_IOS:
		sys = "ios"
	case KERNEL_PLAN9:
		sys = "plan9"
	default:
		// sys = targetKernel
	}

	abi := targetLibc
	switch targetLibc {
	case LIBC_MUSL:
		switch mArch {
		case archconst.ARCH_ARM_V5, archconst.ARCH_ARM_V6:
			abi = "musleabi"
		case archconst.ARCH_ARM_V7, archconst.ARCH_ARM:
			abi = "musleabihf"
		default:
			abi = "musl"
		}
	case LIBC_MSVC:
		switch mArch {
		case "":
			// TODO: add special cases
			_ = abi
		default:
			abi = "msvc"
		}
	case LIBC_GNU:
		fallthrough
	default:
		switch mArch {
		case archconst.ARCH_ARM_V5, archconst.ARCH_ARM_V6:
			abi = "gnueabi"
		case archconst.ARCH_ARM_V7, archconst.ARCH_ARM:
			abi = "gnueabihf"
		case archconst.ARCH_MIPS64, archconst.ARCH_MIPS64_SF,
			archconst.ARCH_MIPS64_LE, archconst.ARCH_MIPS64_LE_SF:
			// TODO: is it gnu or gnuabi64?
			abi = "gnuabi64"
		default:
			abi = "gnu"
		}
	}

	// <arch>-<vendor>-<sys>-<abi>
	return arch + "-unknown-" + sys + "-" + abi, ok
}
