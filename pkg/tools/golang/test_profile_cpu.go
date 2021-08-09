package golang

import "arhat.dev/rs"

type testCPUProfileSpec struct {
	rs.BaseField

	// Profile cpu during test execution
	Enabled bool `yaml:"enabled"`

	// Output filename of cpu profile
	//
	// go test -cpuprofile
	//
	// defaults to `cpu.out` if not set and `enabled` is true
	Output string `yaml:"output"`
}

func (s testCPUProfileSpec) generateArgs(compileTime bool) []string {
	if !s.Enabled {
		return nil
	}

	var args []string
	prefix := getTestFlagPrefix(compileTime)
	if len(s.Output) != 0 {
		args = append(args, prefix+"cpuprofile", s.Output)
	} else {
		args = append(args, prefix+"cpuprofile", "cpu.out")
	}

	return args
}
