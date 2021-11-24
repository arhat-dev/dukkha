package constant

// Arch values for matrix
// nolint:revive
const (
	ARCH_X86 = "x86"

	ARCH_AMD64 = "amd64"

	// ref: https://en.wikipedia.org/wiki/X86-64#Microarchitecture_levels
	ARCH_AMD64_V1 = "amd64v1" // alias of amd64
	ARCH_AMD64_V2 = "amd64v2" // 2009+
	ARCH_AMD64_V3 = "amd64v3" // 2015+
	ARCH_AMD64_V4 = "amd64v4" // avx512

	ARCH_ARM64 = "arm64"

	ARCH_ARM_V5 = "armv5"
	ARCH_ARM_V6 = "armv6"
	ARCH_ARM_V7 = "armv7"

	ARCH_MIPS       = "mips"
	ARCH_MIPS_SF    = "mipssf"
	ARCH_MIPS_LE    = "mipsle"
	ARCH_MIPS_LE_SF = "mipslesf"

	ARCH_MIPS64       = "mips64"
	ARCH_MIPS64_SF    = "mips64sf"
	ARCH_MIPS64_LE    = "mips64le"
	ARCH_MIPS64_LE_SF = "mips64lesf"

	ARCH_PPC       = "ppc"
	ARCH_PPC_SF    = "ppcsf"
	ARCH_PPC_LE    = "ppcle"
	ARCH_PPC_LE_SF = "ppclesf"

	ARCH_PPC64    = "ppc64"
	ARCH_PPC64_LE = "ppc64le"

	ARCH_PPC64_V8    = "ppc64v8"
	ARCH_PPC64_V8_LE = "ppc64v8le"
	ARCH_PPC64_V9    = "ppc64v9"
	ARCH_PPC64_V9_LE = "ppc64v9le"

	ARCH_RISCV_64 = "riscv64"

	ARCH_S390X = "s390x"

	ARCH_IA64 = "ia64"
)

type ArchMappingValues struct {
	Alpine string
	Debian string
	GNU    string

	Golang string
	Docker string
	OCI    string

	DockerHub string

	Qemu string

	// TODO
	LLVM string
	Zig  string
	Rust string
}

// mapping lower case values to ArchMappingValues Field Names
var supportedPlatforms = map[string]string{
	"alpine": "Alpine",
	"debian": "Debian",
	"gnu":    "GNU",

	"golang": "Golang",
	"docker": "Docker",
	"oci":    "OCI",

	"dockerhub": "DockerHub",

	"qemu": "Qemu",

	"llvm": "LLVM",
	"zig":  "Zig",
	"rust": "Rust",
}
