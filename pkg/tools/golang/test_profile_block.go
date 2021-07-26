package golang

import (
	"strconv"

	"arhat.dev/dukkha/pkg/field"
)

type testBlockProfileSpec struct {
	field.BaseField

	Enabled bool   `yaml:"enabled"`
	Rate    int    `yaml:"rate"`
	Output  string `yaml:"output"`
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
