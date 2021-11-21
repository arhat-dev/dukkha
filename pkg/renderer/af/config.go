package af

import (
	"arhat.dev/rs"
)

// InputSpec is the alternative yaml input schema to af renderer
type inputSpec struct {
	rs.BaseField

	// Archive is the local file path to the archive
	Archive string `yaml:"archive"`

	// Path is the in archive path of the target file to extract
	Path string `yaml:"path"`

	// Flatten generate a map of files in archive
	Flatten bool `yaml:"flatten"`

	// Password for password protected archive files
	Password string `yaml:"password"`
}

func (s *inputSpec) ScopeUniqueID() string {
	// return append([]byte(s.Archive), "|...|"+s.Path...))
	return s.Archive
}
