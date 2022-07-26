package dukkha

import (
	"fmt"
	"io"
)

type TaskExecStage uint8

const (
	StageBefore TaskExecStage = iota + 1

	StageBeforeMatrix
	StageAfterMatrixSuccess
	StageAfterMatrixFailure
	StageAfterMatrix

	StageAfterSuccess
	StageAfterFailure
	StageAfter
)

func (s TaskExecStage) String() string {
	switch s {
	case StageBefore:
		return "before"
	case StageBeforeMatrix:
		return "before:matrix"
	case StageAfterMatrixSuccess:
		return "after:matrix:success"
	case StageAfterMatrixFailure:
		return "after:matrix:failure"
	case StageAfterMatrix:
		return "after:matrix"
	case StageAfterSuccess:
		return "after:success"
	case StageAfterFailure:
		return "after:failure"
	case StageAfter:
		return "after"
	default:
		panic(fmt.Errorf("unknown task exec stage %d", s))
	}
}

// one of `tools.TaskExecRequest`, `[]dukkha.TaskExecSpec`
type RunTaskOrRunCmd interface{}

type ReplaceEntries map[string]ReplaceEntry

type ReplaceEntry struct {
	Data []byte
	Err  error
}

// TaskExecSpec is the specification
type TaskExecSpec struct {
	// StdoutAsReplace to replace same string in following TaskExecSpecs
	// use output to stdout of this exec
	StdoutAsReplace string

	// ShowStdout when StdoutAsReplace is set
	ShowStdout bool

	FixStdoutValueForReplace func(data []byte) []byte

	// StderrAsReplace to replace same string in following TaskExecSpecs
	// use output to stderr of this exec
	StderrAsReplace string

	// ShowStderr when StderrAsReplace is set
	ShowStderr bool

	FixStderrValueForReplace func(data []byte) []byte

	Chdir string

	// EnvSuggest to provide environment variables when not set by user
	EnvSuggest NameValueList
	// EnvOverride to override existing environment variables
	EnvOverride NameValueList

	Command []string

	AlterExecFunc func(
		replace ReplaceEntries,
		stdin io.Reader, stdout, stderr io.Writer,
	) (RunTaskOrRunCmd, error)

	Stdin io.Reader

	// IgnoreError to ignore error generated after running this spec
	// this option applies to all sub specs (as returned in AlterExecFunc)
	IgnoreError bool

	// UseShell if true, write command to local script cache
	// and execute with the target shell (as referenced by ShellName)
	UseShell bool

	// ShellName to reference a shell to execute command
	// when `UseShell` is true
	//
	// the availability of the shells denpends on `shells` in dukkha config
	// a special shell name is `embedded`, which will use the built-in shell
	// for command evaluation
	ShellName string
}
