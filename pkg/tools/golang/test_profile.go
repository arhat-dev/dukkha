package golang

import (
	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"
)

type testProfileSpec struct {
	rs.BaseField `yaml:"-"`

	// Directory to save all profile output files
	OutputDir string `yaml:"output_dir"`

	// Coverage profile
	Coverage testCoverageProfileSpec `yaml:"coverage"`

	// Goroutine Block profile
	Block testBlockProfileSpec `yaml:"block"`

	// CPU profile
	CPU testCPUProfileSpec `yaml:"cpu"`

	// Memory profile
	Memory testMemoryProfileSpec `yaml:"memory"`

	// Mutex profile
	Mutex testMutexProfileSpec `yaml:"mutex"`

	// Trace profile
	Trace testTraceProfileSpec `yaml:"trace"`
}

func (s testProfileSpec) generateArgs(cwdFS *fshelper.OSFS, compileTime bool) []string {
	var args []string

	prefix := getTestFlagPrefix(compileTime)
	if len(s.OutputDir) != 0 {
		outputDir, err := cwdFS.Abs(s.OutputDir)
		if err != nil {
			panic(err)
		}

		args = append(args, prefix+"outputdir", outputDir)
	}

	args = append(args, s.Coverage.generateArgs(compileTime)...)
	args = append(args, s.Block.generateArgs(compileTime)...)
	args = append(args, s.Memory.generateArgs(compileTime)...)
	args = append(args, s.CPU.generateArgs(compileTime)...)
	args = append(args, s.Mutex.generateArgs(compileTime)...)
	args = append(args, s.Trace.generateArgs(compileTime)...)

	return args
}
