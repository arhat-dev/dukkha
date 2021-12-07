package constant

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
