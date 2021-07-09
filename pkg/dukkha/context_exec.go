package dukkha

import "github.com/fatih/color"

type ExecValues interface {
	SetOutputPrefix(s string)
	OutputPrefix() string

	SetTaskColors(prefixColor, outputColor *color.Color)
	PrefixColor() *color.Color
	OutputColor() *color.Color

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

	prefixColor *color.Color
	outputColor *color.Color
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

func (c *contextExec) SetTaskColors(prefixColor, outputColor *color.Color) {
	if c.prefixColor != nil || c.outputColor != nil {
		return
	}

	c.prefixColor = prefixColor
	c.outputColor = outputColor
}

func (c *contextExec) PrefixColor() *color.Color {
	return c.prefixColor
}

func (c *contextExec) OutputColor() *color.Color {
	return c.outputColor
}

func (c *contextExec) CurrentTool() ToolKey {
	return ToolKey{Kind: c.toolKind, Name: c.toolName}
}

func (c *contextExec) CurrentTask() TaskKey {
	return TaskKey{Kind: c.taskKind, Name: c.taskName}
}
