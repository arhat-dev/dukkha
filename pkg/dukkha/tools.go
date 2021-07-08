package dukkha

import (
	"arhat.dev/dukkha/pkg/field"
)

type (
	ToolKind string
	ToolName string
)

// nolint:revive
type Tool interface {
	field.Field

	// Kind of the tool, e.g. golang, docker
	Kind() ToolKind

	Name() ToolName

	Init(kind ToolKind, cachdDir string) error

	ResolveTasks(tasks []Task) error

	Run(taskCtx TaskExecContext) error
}

type ToolManager interface {
	AddTool(kind ToolKind, name ToolName, impl Tool)
}

type ToolUser interface {
	AllTools() map[ToolKey]Tool
	GetTool(kind ToolKind, name ToolName) (Tool, bool)
}

type ToolKey struct {
	Kind ToolKind
	Name ToolName
}

func newContextTools() *contextTools {
	return &contextTools{
		tools: make(map[ToolKey]Tool),
	}
}

type contextTools struct {
	tools map[ToolKey]Tool
}

func (c *contextTools) AddTool(
	kind ToolKind,
	name ToolName,
	impl Tool,
) {
	c.tools[ToolKey{Kind: kind, Name: name}] = impl
}

func (c *contextTools) AllTools() map[ToolKey]Tool {
	return c.tools
}

func (c *contextTools) GetTool(kind ToolKind, name ToolName) (Tool, bool) {
	t, ok := c.tools[ToolKey{Kind: kind, Name: name}]
	return t, ok
}
