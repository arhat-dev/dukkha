package golang

import (
	"strconv"

	"arhat.dev/dukkha/pkg/field"
)

type testMutexProfileSpec struct {
	field.BaseField

	// Profile mutex during test execution
	Enabled bool `yaml:"enabled"`

	// Fraction number
	//
	// go test -mutexprofilefraction
	Fraction int `yaml:"fraction"`

	// Output filename of mutex profile
	//
	// go test -mutexprofile
	//
	// defaults to `mutex.out` if not set and `enabled` is true
	Output string `yaml:"output"`
}

func (s testMutexProfileSpec) generateArgs(compileTime bool) []string {
	if !s.Enabled {
		return nil
	}

	var args []string
	prefix := getTestFlagPrefix(compileTime)
	if len(s.Output) != 0 {
		args = append(args, prefix+"mutexprofile", s.Output)
	} else {
		args = append(args, prefix+"mutexprofile", "mutex.out")
	}

	if s.Fraction != 0 {
		args = append(args, prefix+"mutexprofilefraction", strconv.FormatInt(int64(s.Fraction), 10))
	}

	return args
}
