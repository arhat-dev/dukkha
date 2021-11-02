package dukkha

import "arhat.dev/rs"

type (
	ToolKind string
	ToolName string

	ToolKey struct {
		Kind ToolKind
		Name ToolName
	}
)

func (k ToolKey) String() string {
	return string(k.Kind) + ":" + string(k.Name)
}

// nolint:revive
type Tool interface {
	rs.Field

	// Kind of the tool, e.g. golang, docker
	Kind() ToolKind

	Name() ToolName

	Key() ToolKey

	GetCmd() []string

	GetEnv() Env

	UseShell() bool

	ShellName() string

	GetTask(TaskKey) (Task, bool)

	AllTasks() map[TaskKey]Task

	Init(kind ToolKind, cachdDir string) error

	ResolveTasks(tasks []Task) error

	Run(taskCtx TaskExecContext) error

	DoAfterFieldsResolved(
		mCtx TaskExecContext,
		depth int,
		do func() error,
		tagNames ...string,
	) error
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
