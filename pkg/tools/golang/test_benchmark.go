package golang

import (
	"strconv"
	"time"

	"arhat.dev/rs"
)

type testBenchmarkSpec struct {
	rs.BaseField

	// Run benchmarks during test execution
	Enabled bool `yaml:"enabled"`

	// Duration of each benchmark run
	Duration time.Duration `yaml:"duration"`

	// Count of benchmark run
	Count int `yaml:"count"`

	// Run only regexp matched benchmarks
	//
	// go test -bench
	//
	// defaults to `.` (all)
	Match string `yaml:"match"`

	// Memory benchmark settings
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
	rs.BaseField

	// Enbaled by default if benchmark is enabled
	Enabled *bool `yaml:"enabled"`

	// Allocs  *bool `yaml:""`
}

func (s benchmarkMemorySpec) generateArgs(compileTime bool) []string {
	if s.Enabled != nil && !*s.Enabled {
		return nil
	}

	var args []string

	prefix := getTestFlagPrefix(compileTime)

	args = append(args, prefix+"benchmem")

	return args
}
