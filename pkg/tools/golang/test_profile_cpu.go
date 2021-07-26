package golang

import "arhat.dev/dukkha/pkg/field"

type testCPUProfileSpec struct {
	field.BaseField

	Enabled bool   `yaml:"enabled"`
	Output  string `yaml:"output"`
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
