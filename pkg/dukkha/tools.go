package dukkha

import (
	"arhat.dev/dukkha/pkg/field"
)

type (
	ToolKind string
	ToolName string

	ToolKey struct {
		Kind ToolKind
		Name ToolName
	}
)

func (k *ToolKey) String() string {
	return string(k.Kind) + ":" + string(k.Name)
}

// nolint:revive
type Tool interface {
	field.Field

	// Kind of the tool, e.g. golang, docker
	Kind() ToolKind

	Name() ToolName

	Key() ToolKey

	GetCmd() []string

	GetEnv() []string

	UseShell() bool

	ShellName() string

	GetTask(TaskKey) (Task, bool)

	Init(kind ToolKind, cachdDir string) error

	ResolveTasks(tasks []Task) error

	Run(taskCtx TaskExecContext) error
}

type ToolManager interface {
	AddTool(k ToolKey, t Tool)
}

type ToolUser interface {
	AllTools() map[ToolKey]Tool
	GetTool(k ToolKey) (Tool, bool)
}

func newContextTools() *contextTools {
	return &contextTools{
		tools: make(map[ToolKey]Tool),
	}
}

type contextTools struct {
	tools map[ToolKey]Tool
}

func (c *contextTools) AddTool(k ToolKey, t Tool) {
	c.tools[k] = t
}

func (c *contextTools) AllTools() map[ToolKey]Tool {
	return c.tools
}

func (c *contextTools) GetTool(k ToolKey) (Tool, bool) {
	t, ok := c.tools[k]
	return t, ok
}
