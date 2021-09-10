package dukkha

import (
	"arhat.dev/dukkha/pkg/sliceutils"
)

type TermColor int

type TaskExecOptions interface {
	NextMatrixExecOptions(useShell bool, shellName string, toolCmd []string) TaskMatrixExecOptions
}

func CreateTaskExecOptions(id, totalMatrix int) TaskExecOptions {
	return &taskExecOpts{
		id:    id,
		seq:   -1,
		total: totalMatrix,
	}
}

type taskExecOpts struct {
	id    int
	seq   int
	total int
}

func (opts *taskExecOpts) NextMatrixExecOptions(
	useShell bool, shellName string, toolCmd []string,
) TaskMatrixExecOptions {
	opts.seq++

	ret := &taskMatrixExecOpts{
		id:    opts.id,
		seq:   opts.seq,
		total: opts.total,

		useShell:  useShell,
		shellName: shellName,
		toolCmd:   sliceutils.NewStrings(toolCmd),
	}

	return ret
}

type TaskMatrixExecOptions interface {
	ID() int
	Total() int

	UseShell() bool
	ShellName() string
	ToolCmd() []string

	Seq() int

	IsLast() bool
}

type taskMatrixExecOpts struct {
	id    int
	seq   int
	total int

	useShell  bool
	shellName string
	toolCmd   []string
}

func (opts *taskMatrixExecOpts) ID() int           { return opts.id }
func (opts *taskMatrixExecOpts) UseShell() bool    { return opts.useShell }
func (opts *taskMatrixExecOpts) ShellName() string { return opts.shellName }
func (opts *taskMatrixExecOpts) Seq() int          { return opts.seq }
func (opts *taskMatrixExecOpts) Total() int        { return opts.total }
func (opts *taskMatrixExecOpts) ToolCmd() []string { return opts.toolCmd }
func (opts *taskMatrixExecOpts) IsLast() bool      { return opts.seq == opts.total-1 }

type ExecValues interface {
	SetOutputPrefix(s string)
	OutputPrefix() string

	SetTaskColors(prefixColor, outputColor TermColor)
	PrefixColor() TermColor
	OutputColor() TermColor

	CurrentTool() ToolKey
	CurrentTask() TaskKey

	SetTask(k ToolKey, tK TaskKey)
}

func newContextExec() *contextExec {
	return &contextExec{}
}

var _ ExecValues = (*contextExec)(nil)

type contextExec struct {
	toolKind ToolKind
	toolName ToolName

	taskKind TaskKind
	taskName TaskName

	outputPrefix string

	prefixColor TermColor
	outputColor TermColor
}

func (c *contextExec) deriveNew() *contextExec {
	return &contextExec{
		toolKind: "",
		toolName: c.toolName,

		taskKind: "",
		taskName: "",

		outputPrefix: c.outputPrefix,
		prefixColor:  c.prefixColor,
		outputColor:  c.outputColor,
	}
}

func (c *contextExec) SetTask(k ToolKey, tK TaskKey) {
	c.toolKind = k.Kind
	c.toolName = k.Name
	c.taskKind = tK.Kind
	c.taskName = tK.Name
}

func (c *contextExec) OutputPrefix() string {
	return c.outputPrefix
}

func (c *contextExec) SetOutputPrefix(s string) {
	c.outputPrefix = s
}

func (c *contextExec) SetTaskColors(prefixColor, outputColor TermColor) {
	if c.prefixColor != 0 || c.outputColor != 0 {
		return
	}

	c.prefixColor = prefixColor
	c.outputColor = outputColor
}

func (c *contextExec) PrefixColor() TermColor {
	return c.prefixColor
}

func (c *contextExec) OutputColor() TermColor {
	return c.outputColor
}

func (c *contextExec) CurrentTool() ToolKey {
	return ToolKey{Kind: c.toolKind, Name: c.toolName}
}

func (c *contextExec) CurrentTask() TaskKey {
	return TaskKey{Kind: c.taskKind, Name: c.taskName}
}
