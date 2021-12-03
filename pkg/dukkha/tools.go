package dukkha

import "arhat.dev/pkg/fshelper"

type (
	ToolKind string
	ToolName string
)

type ToolKey struct {
	Kind ToolKind
	Name ToolName
}

func (k ToolKey) String() string { return string(k.Kind) + ":" + string(k.Name) }

// Tool implementation requirements
type Tool interface {
	Resolvable

	// Kind of the tool, e.g. golang, docker
	Kind() ToolKind

	// Name of the tool, e.g. my-tool
	Name() ToolName

	// Key
	Key() ToolKey

	// GetCmd get cli command to run this tool
	GetCmd() []string

	GetTask(TaskKey) (Task, bool)

	AllTasks() map[TaskKey]Task

	Init(cacheFS *fshelper.OSFS) error

	AddTasks(tasks []Task) error

	// Run task
	Run(taskCtx TaskExecContext, tsk TaskKey) error
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

func (c *contextTools) AddTool(k ToolKey, t Tool)  { c.tools[k] = t }
func (c *contextTools) AllTools() map[ToolKey]Tool { return c.tools }

func (c *contextTools) GetTool(k ToolKey) (Tool, bool) {
	t, ok := c.tools[k]
	return t, ok
}
