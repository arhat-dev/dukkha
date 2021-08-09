package golang

import (
	"strings"

	"arhat.dev/rs"
)

type testCoverageProfileSpec struct {
	rs.BaseField

	Enabled bool `yaml:"enabled"`

	// Output file of the coverage
	//
	// go test -coverprofile
	//
	// defaults to cover.out if not set and `enabled` is true
	Output string `yaml:"output"`

	// Mode of coverage
	//
	// go test -covermode
	//
	// defaults to `atomic` if not set and `enabled` is true
	Mode string `yaml:"mode"`

	// Packages to coverage
	//
	// go test -coverpkg
	//
	// no default (use golang default behavior)
	Packages []string `yaml:"packages"`
}

func (s testCoverageProfileSpec) generateArgs(compileTime bool) []string {
	if !s.Enabled {
		return nil
	}

	var args []string

	if compileTime {
		args = append(args, "-cover")

		if len(s.Mode) != 0 {
			args = append(args, "-covermode", s.Mode)
		} else {
			args = append(args, "-covermode", "atomic")
		}

		if len(s.Packages) != 0 {
			args = append(args, "-coverpkg", strings.Join(s.Packages, ","))
		}
	}

	prefix := getTestFlagPrefix(compileTime)
	if len(s.Output) != 0 {
		args = append(args, prefix+"coverprofile", s.Output)
	} else {
		args = append(args, prefix+"coverprofile", "cover.out")
	}

	return args
}
