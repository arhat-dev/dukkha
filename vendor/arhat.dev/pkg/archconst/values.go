package archconst

// Arch values
// format: [arch name]{endian variant}{`sf` (soft float mark)}{micro arch}
//
// arch name is required, usually a valid GOARCH value (exceptions: `x86`)
// when {endian variant} not set, default endian is assumed
// when `sf` absent, hard float support is assumed
// if there is no micro arch specified, the lowest micro arch is assumed
//
// example `mips`:
// - `mipsle` -> [arch name] = mips, {endian variant} = little-endian (le) as original mips uses big-endian
// - `mipslesf` -> soft float version of `mipsle`
//
// example `amd64`
// - `amd64v3` -> [arch name] = amd64, {micro level} = v3
//
// example `ppc64`: `ppc64`
// - `ppc64le` -> little-endian version of ppc64
// - `ppc64lev9` -> `ppc64le` with `v9` micro arch
//
// nolint:revive
const (
	/*
		X86 (i386)
	*/
	ARCH_X86 = "x86"

	ARCH_X86_SF = "x86sf"

	/*
		AMD64 (x86_64)
	*/
	ARCH_AMD64 = "amd64"

	ARCH_AMD64_V1 = "amd64v1"
	ARCH_AMD64_V2 = "amd64v2"
	ARCH_AMD64_V3 = "amd64v3"
	ARCH_AMD64_V4 = "amd64v4"

	/*
		ARM
	*/

	ARCH_ARM_V5 = "armv5"
	ARCH_ARM_V6 = "armv6"
	ARCH_ARM_V7 = "armv7"

	/*
		ARM64 (aarch64)
	*/
	ARCH_ARM64 = "arm64"

	ARCH_ARM64_V8 = "arm64v8"
	ARCH_ARM64_V9 = "arm64v9"

	/*
		MIPS
	*/

	ARCH_MIPS    = "mips"
	ARCH_MIPS_SF = "mipssf"

	ARCH_MIPS_LE    = "mipsle"
	ARCH_MIPS_LE_SF = "mipslesf"

	/*
		MIPS64
	*/

	ARCH_MIPS64    = "mips64"
	ARCH_MIPS64_SF = "mips64sf"

	ARCH_MIPS64_LE    = "mips64le"
	ARCH_MIPS64_LE_SF = "mips64lesf"

	/*
		PowerPC
	*/

	ARCH_PPC       = "ppc"
	ARCH_PPC_SF    = "ppcsf"
	ARCH_PPC_LE    = "ppcle"
	ARCH_PPC_LE_SF = "ppclesf"

	/*
		PowerPC 64
	*/

	ARCH_PPC64    = "ppc64"
	ARCH_PPC64_V8 = "ppc64v8"
	ARCH_PPC64_V9 = "ppc64v9"

	ARCH_PPC64_LE    = "ppc64le"
	ARCH_PPC64_LE_V8 = "ppc64lev8"
	ARCH_PPC64_LE_V9 = "ppc64lev9"

	/*
		RISCV64
	*/

	ARCH_RISCV_64 = "riscv64"

	/*
		S390X (64bit S390)
	*/

	ARCH_S390X = "s390x"

	/*
		IA64
	*/

	ARCH_IA64 = "ia64"
)
