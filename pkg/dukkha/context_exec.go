package dukkha

import (
	"github.com/muesli/termenv"
)

// RuntimeOptions for task execution
type RuntimeOptions struct {
	FailFast            bool
	ColorOutput         bool
	TranslateANSIStream bool
	RetainANSIStyle     bool
	Workers             int
}

type TaskExecOptions interface {
	NextMatrixExecOptions() TaskMatrixExecOptions
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

func (opts *taskExecOpts) NextMatrixExecOptions() TaskMatrixExecOptions {
	opts.seq++

	ret := &taskMatrixExecOpts{
		id:    opts.id,
		seq:   opts.seq,
		total: opts.total,
	}

	return ret
}

type TaskMatrixExecOptions interface {
	ID() int

	Total() int

	Seq() int

	IsLast() bool
}

type taskMatrixExecOpts struct {
	id    int
	seq   int
	total int
}

func (opts *taskMatrixExecOpts) ID() int      { return opts.id }
func (opts *taskMatrixExecOpts) Seq() int     { return opts.seq }
func (opts *taskMatrixExecOpts) Total() int   { return opts.total }
func (opts *taskMatrixExecOpts) IsLast() bool { return opts.seq == opts.total-1 }

type ExecValues interface {
	SetOutputPrefix(s string)
	OutputPrefix() string

	SetTaskColors(prefixColor, outputColor termenv.Color)
	PrefixColor() termenv.Color
	OutputColor() termenv.Color

	CurrentTool() ToolKey
	CurrentTask() TaskKey

	SetTask(k ToolKey, tK TaskKey)

	TranslateANSIStream() bool
	RetainANSIStyle() bool
	ColorOutput() bool
	FailFast() bool
	ClaimWorkers(n int) int

	SetState(s TaskExecState)
	State() TaskExecState
}

type TaskExecState int

const (
	TaskExecPending TaskExecState = iota
	TaskExecNotStarted
	TaskExecWorking
	TaskExecSucceeded
	TaskExecFailed
	TaskExecCanceled
)

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

	prefixColor termenv.Color
	outputColor termenv.Color

	state TaskExecState

	runtimeOpts RuntimeOptions
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

		runtimeOpts: c.runtimeOpts,
	}
}

func (c *contextExec) SetTask(k ToolKey, tK TaskKey) {
	c.toolKind = k.Kind
	c.toolName = k.Name
	c.taskKind = tK.Kind
	c.taskName = tK.Name
}

func (c *contextExec) OutputPrefix() string     { return c.outputPrefix }
func (c *contextExec) SetOutputPrefix(s string) { c.outputPrefix = s }

func (c *contextExec) SetTaskColors(prefixColor, outputColor termenv.Color) {
	if c.prefixColor != nil || c.outputColor != nil {
		return
	}

	c.prefixColor = prefixColor
	c.outputColor = outputColor
}

func (c *contextExec) PrefixColor() termenv.Color { return c.prefixColor }
func (c *contextExec) OutputColor() termenv.Color { return c.outputColor }

func (c *contextExec) CurrentTool() ToolKey { return ToolKey{Kind: c.toolKind, Name: c.toolName} }
func (c *contextExec) CurrentTask() TaskKey { return TaskKey{Kind: c.taskKind, Name: c.taskName} }

func (c *contextExec) SetRuntimeOptions(opts RuntimeOptions) { c.runtimeOpts = opts }

func (c *contextExec) FailFast() bool            { return c.runtimeOpts.FailFast }
func (c *contextExec) ColorOutput() bool         { return c.runtimeOpts.ColorOutput }
func (c *contextExec) TranslateANSIStream() bool { return c.runtimeOpts.TranslateANSIStream }
func (c *contextExec) RetainANSIStyle() bool     { return c.runtimeOpts.RetainANSIStyle }

func (c *contextExec) ClaimWorkers(n int) int {
	if c.runtimeOpts.Workers > n {
		return n
	}

	// TODO: limit workers
	return c.runtimeOpts.Workers
}

func (c *contextExec) SetState(s TaskExecState) { c.state = s }
func (c *contextExec) State() TaskExecState     { return c.state }
