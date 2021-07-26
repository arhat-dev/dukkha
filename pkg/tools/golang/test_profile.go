package golang

import (
	"path/filepath"

	"arhat.dev/dukkha/pkg/field"
)

type testProfileSpec struct {
	field.BaseField

	OutputDir string `yaml:"output_dir"`

	Coverage testCoverageProfileSpec `yaml:"coverage"`
	Block    testBlockProfileSpec    `yaml:"block"`
	CPU      testCPUProfileSpec      `yaml:"cpu"`
	Memory   testMemoryProfileSpec   `yaml:"memory"`
	Mutex    testMutexProfileSpec    `yaml:"mutex"`
	Trace    testTraceProfileSpec    `yaml:"trace"`
}

func (s testProfileSpec) generateArgs(dukkhaWorkDir string, compileTime bool) []string {
	var args []string

	prefix := getTestFlagPrefix(compileTime)
	if len(s.OutputDir) != 0 {
		if filepath.IsAbs(s.OutputDir) {
			args = append(args, prefix+"outputdir", s.OutputDir)
		} else {
			args = append(args, prefix+"outputdir", filepath.Join(dukkhaWorkDir, s.OutputDir))
		}
	}

	args = append(args, s.Coverage.generateArgs(compileTime)...)
	args = append(args, s.Block.generateArgs(compileTime)...)
	args = append(args, s.Memory.generateArgs(compileTime)...)
	args = append(args, s.CPU.generateArgs(compileTime)...)
	args = append(args, s.Mutex.generateArgs(compileTime)...)
	args = append(args, s.Trace.generateArgs(compileTime)...)

	return args
}
