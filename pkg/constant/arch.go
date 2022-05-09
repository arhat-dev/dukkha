package constant

import "arhat.dev/pkg/archconst"

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

// SimpleArch generalizes cpu micro arch for practical use
//
// - amd64v{1, 2, 3, 4} => amd64
// - arm64v{8, 9} => arm64
// - ppc64{, le}v{8, 9} => ppc64{, le}
//
// Exceptions:
// - armv{5, 6, 7} are kept as is
// - arm => armv7
func SimpleArch[T ~string](arch T) T {
	s, ok := archconst.Split(arch)
	if !ok {
		// unknown
		return arch
	}

	switch s.Name {
	case archconst.ARCH_AMD64, archconst.ARCH_ARM64, archconst.ARCH_PPC64:
		s.MicroArch = ""
	case archconst.ARCH_ARM:
		if len(s.MicroArch) == 0 {
			s.MicroArch = "v7"
		}
	}

	return T(s.String())
}
