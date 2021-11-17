package dukkha

type TaskManager interface {
	AddToolSpecificTasks(kind ToolKind, name ToolName, tasks []Task)
}

type TaskUser interface {
	GetToolSpecificTasks(k ToolKey) ([]Task, bool)
	AllToolSpecificTasks() map[ToolKey][]Task
}

func newContextTasks() *contextTasks {
	return &contextTasks{
		toolSpecificTasks: make(map[ToolKey][]Task),
	}
}

type contextTasks struct {
	toolSpecificTasks map[ToolKey][]Task
}

func (c *contextTasks) AddToolSpecificTasks(k ToolKind, n ToolName, tasks []Task) {
	toolKey := ToolKey{Kind: k, Name: n}

	c.toolSpecificTasks[toolKey] = append(
		c.toolSpecificTasks[toolKey], tasks...,
	)
}

func (c *contextTasks) GetToolSpecificTasks(k ToolKey) ([]Task, bool) {
	tasks, ok := c.toolSpecificTasks[k]
	return tasks, ok
}

func (c *contextTasks) AllToolSpecificTasks() map[ToolKey][]Task {
	return c.toolSpecificTasks
}
