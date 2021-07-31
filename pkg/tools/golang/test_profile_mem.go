package golang

import (
	"strconv"

	"arhat.dev/dukkha/pkg/field"
)

type testMemoryProfileSpec struct {
	field.BaseField

	// Profile memory during test execution
	Enabled bool `yaml:"enabled"`

	// Rate of memory profile
	//
	// go test -memprofilerate
	Rate int `yaml:"rate"`

	// Output filename of memory profile
	//
	// go test -memprofile
	//
	// defaults to mem.out if not set and `enabled` is true
	Output string `yaml:"output"`
}

func (s testMemoryProfileSpec) generateArgs(compileTime bool) []string {
	if !s.Enabled {
		return nil
	}

	var args []string
	prefix := getTestFlagPrefix(compileTime)
	if len(s.Output) != 0 {
		args = append(args, prefix+"memprofile", s.Output)
	} else {
		args = append(args, prefix+"memprofile", "mem.out")
	}

	if s.Rate != 0 {
		args = append(args, prefix+"memprofilerate", strconv.FormatInt(int64(s.Rate), 10))
	}

	return args
}
