package golang

import (
	"strings"

	"arhat.dev/dukkha/pkg/field"
)

type testCoverageProfileSpec struct {
	field.BaseField

	Enabled  bool     `yaml:"enabled"`
	Output   string   `yaml:"output"`
	Mode     string   `yaml:"mode"`
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
