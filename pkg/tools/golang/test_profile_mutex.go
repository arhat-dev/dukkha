package golang

import (
	"strconv"

	"arhat.dev/dukkha/pkg/field"
)

type testMutexProfileSpec struct {
	field.BaseField

	Enabled  bool   `yaml:"enabled"`
	Fraction int    `yaml:"fraction"`
	Output   string `yaml:"output"`
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
