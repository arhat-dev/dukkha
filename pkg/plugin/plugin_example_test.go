package plugin

import (
	"arhat.dev/dukkha/pkg/dukkha"
)

// NewRenderer_{renderer-default-name}
// nolint:revive
func NewRenderer_foo(name string) dukkha.Renderer {
	return nil
}

// NewTool_{tool-kind}
// nolint:revive
func NewTool_foo_tool() dukkha.Tool {
	return nil
}

// NewTask_{task-kind}
// e.g. You have defined a task `foo-task`, its factory func should be NewTask_foo_task
// nolint:revive
func NewTask_foo_tool_foo_task(name string) dukkha.Task {
	return nil
}
