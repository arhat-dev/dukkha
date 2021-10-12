package golang

import "arhat.dev/rs"

type testTraceProfileSpec struct {
	rs.BaseField `yaml:"-"`

	// Write test execution trace
	Enabled bool `yaml:"enabled"`

	// Output filename of trace profile
	//
	// go test -trace
	//
	// defaults to trace.out if not set and `enabled` is true
	Output string `yaml:"output"`
}

func (s testTraceProfileSpec) generateArgs(compileTime bool) []string {
	if !s.Enabled {
		return nil
	}

	var args []string
	prefix := getTestFlagPrefix(compileTime)
	if len(s.Output) != 0 {
		args = append(args, prefix+"trace", s.Output)
	} else {
		args = append(args, prefix+"trace", "trace.out")
	}

	return args
}
