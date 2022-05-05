package constant

import "arhat.dev/pkg/archconst"

// GetAlpineTripleName of matrix arch
// reference: https://more.musl.cc/10/x86_64-linux-musl/
func GetAlpineTripleName(mArch string) (string, bool) {
	// NOTE: some arch value in triple name is different from alpine arch
	// 		 so we cannot do the same conversion as GNU/LLVM triple does
	v, ok := map[string]string{
		archconst.ARCH_X86:      "i686-linux-musl",
		archconst.ARCH_AMD64:    "x86_64-linux-musl",
		archconst.ARCH_AMD64_V1: "x86_64-linux-musl",
		archconst.ARCH_AMD64_V2: "x86_64-linux-musl",
		archconst.ARCH_AMD64_V3: "x86_64-linux-musl",
		archconst.ARCH_AMD64_V4: "x86_64-linux-musl",

		archconst.ARCH_ARM_V5: "armv5l-linux-musleabi",
		archconst.ARCH_ARM_V6: "armv6-linux-musleabihf",
		archconst.ARCH_ARM_V7: "armv7l-linux-musleabihf",
		archconst.ARCH_ARM64:  "aarch64-linux-musl",

		archconst.ARCH_PPC:       "powerpc-linux-musl",
		archconst.ARCH_PPC_SF:    "powerpc-linux-muslsf",
		archconst.ARCH_PPC_LE:    "powerpcle-linux-musl",
		archconst.ARCH_PPC_LE_SF: "powerpcle-linux-muslsf",

		archconst.ARCH_PPC64:       "powerpc64-linux-musl",
		archconst.ARCH_PPC64_LE:    "powerpc64le-linux-musl",
		archconst.ARCH_PPC64_V8:    "powerpc64-linux-musl",
		archconst.ARCH_PPC64_LE_V8: "powerpc64le-linux-musl",
		archconst.ARCH_PPC64_V9:    "powerpc64-linux-musl",
		archconst.ARCH_PPC64_LE_V9: "powerpc64le-linux-musl",

		archconst.ARCH_MIPS:         "mips-linux-musl",
		archconst.ARCH_MIPS_SF:      "mips-linux-muslsf",
		archconst.ARCH_MIPS_LE:      "mipsel-linux-musl",
		archconst.ARCH_MIPS_LE_SF:   "mipsel-linux-muslsf",
		archconst.ARCH_MIPS64:       "mips64-linux-musl",
		archconst.ARCH_MIPS64_SF:    "mips64-linux-musln32sf",
		archconst.ARCH_MIPS64_LE:    "mips64el-linux-musl",
		archconst.ARCH_MIPS64_LE_SF: "mips64el-linux-musln32sf",

		archconst.ARCH_RISCV_64: "riscv64-linux-musl",
		archconst.ARCH_S390X:    "s390x-linux-musl",

		archconst.ARCH_IA64: "",
	}[mArch]

	return v, ok
}

func GetDebianTripleName(mArch, targetKernel, targetLibc string) (string, bool) {
	// TODO: adjust triple name according to target kernel
	_ = targetKernel

	// NOTE: some arch value in triple name is different from debian arch
	// 		 so we cannot do the same conversion as GNU/LLVM triple does
	switch targetLibc {
	case LIBC_MUSL:
		// https://packages.debian.org/buster/musl-dev
		// check list of files
		v, ok := map[string]string{
			archconst.ARCH_X86:      "i386-linux-musl",
			archconst.ARCH_AMD64:    "x86_64-linux-musl",
			archconst.ARCH_AMD64_V1: "x86_64-linux-musl",
			archconst.ARCH_AMD64_V2: "x86_64-linux-musl",
			archconst.ARCH_AMD64_V3: "x86_64-linux-musl",
			archconst.ARCH_AMD64_V4: "x86_64-linux-musl",

			archconst.ARCH_ARM_V5: "arm-linux-musleabi",
			archconst.ARCH_ARM_V6: "arm-linux-musleabi",
			archconst.ARCH_ARM_V7: "arm-linux-musleabihf",
			archconst.ARCH_ARM64:  "aarch64-linux-musl",

			archconst.ARCH_MIPS:         "mips-linux-musl",
			archconst.ARCH_MIPS_SF:      "mips-linux-musl",
			archconst.ARCH_MIPS_LE:      "mipsel-linux-musl",
			archconst.ARCH_MIPS_LE_SF:   "mipsel-linux-musl",
			archconst.ARCH_MIPS64:       "mips64-linux-musl",
			archconst.ARCH_MIPS64_SF:    "mips64-linux-musl",
			archconst.ARCH_MIPS64_LE:    "mips64el-linux-musl",
			archconst.ARCH_MIPS64_LE_SF: "mips64el-linux-musl",

			archconst.ARCH_S390X: "s390x-linux-musl",

			// http://ftp.ports.debian.org/debian-ports/pool-riscv64/main/m/musl/
			// download one musl-dev package
			// list package contents with following commands
			//
			// $ ar -x musl-dev_1.2.2-3_riscv64.deb
			// $ tar -tvf data.tar.xz
			archconst.ARCH_RISCV_64: "riscv64-linux-musl",

			archconst.ARCH_IA64: "",

			archconst.ARCH_PPC:       "",
			archconst.ARCH_PPC_SF:    "",
			archconst.ARCH_PPC_LE:    "",
			archconst.ARCH_PPC_LE_SF: "",

			archconst.ARCH_PPC64:       "",
			archconst.ARCH_PPC64_LE:    "",
			archconst.ARCH_PPC64_V8:    "",
			archconst.ARCH_PPC64_LE_V8: "",
			archconst.ARCH_PPC64_V9:    "",
			archconst.ARCH_PPC64_LE_V9: "",
		}[mArch]

		return v, ok
	case LIBC_MSVC:
		v, ok := map[string]string{
			// https://packages.debian.org/buster/mingw-w64-i686-dev
			// check list of files
			archconst.ARCH_X86: "i686-w64-mingw32",
			// https://packages.debian.org/buster/mingw-w64-x86-64-dev
			// check list of files
			archconst.ARCH_AMD64:    "x86_64-w64-mingw32",
			archconst.ARCH_AMD64_V1: "x86_64-w64-mingw32",
			archconst.ARCH_AMD64_V2: "x86_64-w64-mingw32",
			archconst.ARCH_AMD64_V3: "x86_64-w64-mingw32",
			archconst.ARCH_AMD64_V4: "x86_64-w64-mingw32",

			archconst.ARCH_IA64: "",

			archconst.ARCH_ARM_V5: "",
			archconst.ARCH_ARM_V6: "",
			archconst.ARCH_ARM_V7: "",
			archconst.ARCH_ARM64:  "",

			archconst.ARCH_PPC:       "",
			archconst.ARCH_PPC_SF:    "",
			archconst.ARCH_PPC_LE:    "",
			archconst.ARCH_PPC_LE_SF: "",

			archconst.ARCH_PPC64:       "",
			archconst.ARCH_PPC64_LE:    "",
			archconst.ARCH_PPC64_V8:    "",
			archconst.ARCH_PPC64_LE_V8: "",
			archconst.ARCH_PPC64_V9:    "",
			archconst.ARCH_PPC64_LE_V9: "",

			archconst.ARCH_MIPS:         "",
			archconst.ARCH_MIPS_SF:      "",
			archconst.ARCH_MIPS_LE:      "",
			archconst.ARCH_MIPS_LE_SF:   "",
			archconst.ARCH_MIPS64:       "",
			archconst.ARCH_MIPS64_SF:    "",
			archconst.ARCH_MIPS64_LE:    "",
			archconst.ARCH_MIPS64_LE_SF: "",

			archconst.ARCH_RISCV_64: "",
			archconst.ARCH_S390X:    "",
		}[mArch]

		return v, ok
	case LIBC_GNU:
		fallthrough
	default:
		v, ok := map[string]string{
			archconst.ARCH_X86:      "i686-linux-gnu",
			archconst.ARCH_AMD64:    "x86_64-linux-gnu",
			archconst.ARCH_AMD64_V1: "x86_64-linux-gnu",
			archconst.ARCH_AMD64_V2: "x86_64-linux-gnu",
			archconst.ARCH_AMD64_V3: "x86_64-linux-gnu",
			archconst.ARCH_AMD64_V4: "x86_64-linux-gnu",

			archconst.ARCH_ARM_V5: "arm-linux-gnueabi",
			archconst.ARCH_ARM_V6: "arm-linux-gnueabi",
			archconst.ARCH_ARM_V7: "arm-linux-gnueabihf",
			archconst.ARCH_ARM64:  "aarch64-linux-gnu",

			archconst.ARCH_PPC64:       "powerpc64-linux-gnu",
			archconst.ARCH_PPC64_LE:    "powerpc64le-linux-gnu",
			archconst.ARCH_PPC64_V8:    "powerpc64-linux-gnu",
			archconst.ARCH_PPC64_LE_V8: "powerpc64le-linux-gnu",
			archconst.ARCH_PPC64_V9:    "powerpc64-linux-gnu",
			archconst.ARCH_PPC64_LE_V9: "powerpc64le-linux-gnu",

			archconst.ARCH_MIPS:         "mips-linux-gnu",
			archconst.ARCH_MIPS_SF:      "mips-linux-gnu",
			archconst.ARCH_MIPS_LE:      "mipsel-linux-gnu",
			archconst.ARCH_MIPS_LE_SF:   "mipsel-linux-gnu",
			archconst.ARCH_MIPS64:       "mips64-linux-gnuabi64",
			archconst.ARCH_MIPS64_SF:    "mips64-linux-gnuabi64",
			archconst.ARCH_MIPS64_LE:    "mips64el-linux-gnuabi64",
			archconst.ARCH_MIPS64_LE_SF: "mips64el-linux-gnuabi64",

			archconst.ARCH_RISCV_64: "riscv64-linux-gnu",
			archconst.ARCH_S390X:    "s390x-linux-gnu",

			archconst.ARCH_IA64: "",

			archconst.ARCH_PPC:       "",
			archconst.ARCH_PPC_SF:    "",
			archconst.ARCH_PPC_LE:    "",
			archconst.ARCH_PPC_LE_SF: "",
		}[mArch]

		return v, ok
	}
}

// Ref:
// - https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/data/cputable
// - https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/data/ostable
func GetGNUTripleName(mArch, targetKernel, targetLibc string) (string, bool) {
	_ = targetKernel

	arch, ok := GetArch("gnu", mArch)
	if !ok {
		return "", false
	}

	abi := targetLibc
	switch targetLibc {
	case LIBC_MUSL:
		switch mArch {
		case archconst.ARCH_ARM_V5, archconst.ARCH_ARM_V6:
			abi = "musleabi"
		case archconst.ARCH_ARM_V7:
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
		case archconst.ARCH_ARM_V7:
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
	arch, ok := GetArch("llvm", mArch)
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
		case archconst.ARCH_ARM_V7:
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
		case archconst.ARCH_ARM_V7:
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
