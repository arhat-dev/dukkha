package golang

import (
	"strconv"
	"time"

	"arhat.dev/dukkha/pkg/field"
)

type testBenchmarkSpec struct {
	field.BaseField

	Enabled  bool          `yaml:"enabled"`
	Duration time.Duration `yaml:"duration"`
	Count    int           `yaml:"count"`
	Match    string        `yaml:"match"`

	Memory benchmarkMemorySpec `yaml:"memory"`
}

func (s testBenchmarkSpec) generateArgs(compileTime bool) []string {
	if !s.Enabled {
		return nil
	}

	var args []string

	prefix := getTestFlagPrefix(compileTime)
	if len(s.Match) != 0 {
		args = append(args, prefix+"bench", s.Match)
	} else {
		args = append(args, prefix+"bench", ".")
	}

	if s.Duration != 0 {
		args = append(args, prefix+"benchtime", s.Duration.String())
	}

	if s.Count != 0 {
		args = append(args, prefix+"benchtime", strconv.FormatInt(int64(s.Count), 10)+"x")
	}

	args = append(args, s.Memory.generateArgs(compileTime)...)

	return args
}

type benchmarkMemorySpec struct {
	field.BaseField

	Enabled bool `yaml:"enabled"`
	// Allocs  *bool `yaml:""`
}

func (s benchmarkMemorySpec) generateArgs(compileTime bool) []string {
	if !s.Enabled {
		return nil
	}

	var args []string

	prefix := getTestFlagPrefix(compileTime)

	args = append(args, prefix+"benchmem")

	return args
}
