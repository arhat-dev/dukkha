package golang

import (
	"strconv"

	"arhat.dev/rs"
)

type testBlockProfileSpec struct {
	rs.BaseField `yaml:"-"`

	// Profile goroutine blocking during test execution
	Enabled bool `yaml:"enabled"`

	// Rate of block profile
	//
	// go test -blockprofilerate
	//
	// no default
	Rate int `yaml:"rate"`

	// Output filename of block profile
	//
	// go test -blockprofile
	//
	// defaults to `block.out` if not set and `enabled` is true
	Output string `yaml:"output"`
}

func (s testBlockProfileSpec) generateArgs(compileTime bool) []string {
	if !s.Enabled {
		return nil
	}

	var args []string

	prefix := getTestFlagPrefix(compileTime)
	if len(s.Output) != 0 {
		args = append(args, prefix+"blockprofile", s.Output)
	} else {
		args = append(args, prefix+"blockprofile", "block.out")
	}

	if s.Rate != 0 {
		args = append(args, prefix+"blockprofilerate", strconv.FormatInt(int64(s.Rate), 10))
	}

	return args
}
