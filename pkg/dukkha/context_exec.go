package dukkha

import "github.com/fatih/color"

type ExecValues interface {
	SetOutputPrefix(s string)
	OutputPrefix() string

	SetTaskColors(prefixColor, outputColor *color.Color)
	PrefixColor() *color.Color
	OutputColor() *color.Color

	CurrentTool() Tool
}

func newContextExec() *contextExec {
	return &contextExec{}
}

var _ ExecValues = (*contextExec)(nil)

type contextExec struct {
	thisTool Tool

	toolName ToolName
	toolKind ToolKind

	taskKind TaskKind
	taskName TaskName

	outputPrefix string

	prefixColor *color.Color
	outputColor *color.Color
}

func (c *contextExec) deriveNew() *contextExec {
	return &contextExec{
		thisTool: nil,
		toolName: "",
		toolKind: "",

		taskKind: "",
		taskName: "",

		outputPrefix: c.outputPrefix,
		prefixColor:  c.prefixColor,
		outputColor:  c.outputColor,
	}
}

func (c *contextExec) setTask(k ToolKind, n ToolName, tK TaskKind, tN TaskName) {
	c.toolKind = k
	c.toolName = n
	c.taskKind = tK
	c.taskName = tN
}

func (c *contextExec) OutputPrefix() string {
	return c.outputPrefix
}

func (c *contextExec) SetOutputPrefix(s string) {
	c.outputPrefix = s
}

func (c *contextExec) SetTaskColors(prefixColor, outputColor *color.Color) {
	if c.prefixColor != nil || outputColor != nil {
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

func (c *contextExec) CurrentTool() Tool {
	return c.thisTool
}
