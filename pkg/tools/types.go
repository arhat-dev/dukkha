package tools

import (
	"fmt"

	"arhat.dev/dukkha/pkg/field"
)

type Interface interface {
	Name() string
}

func CreateToolKey(toolName, toolID string) (*ToolKey, error) {
	if len(toolName) == 0 {
		return nil, fmt.Errorf("missing tool name")
	}

	return &ToolKey{
		toolName: toolName,
		toolID:   toolID,
	}, nil
}

type ToolKey struct {
	toolName string
	toolID   string
}

func (k ToolKey) String() string {
	return joinReplaceEmpty(
		":",
		[]string{"<undefined-tool-name>", ""},
		k.toolName, k.toolID,
	)
}

type ToolConfig interface {
	field.Interface

	Name() string
	ID() string
}

type TaskTypeKey struct {
	tool ToolKey

	taskType string
}

func (k TaskTypeKey) String() string {
	return joinReplaceEmpty(
		":",
		[]string{"", "<undefined-task-type>"},
		k.tool.String(), k.taskType,
	)
}

func CreateTaskTypeKey(toolName, toolID, taskType string) (*TaskTypeKey, error) {
	if len(taskType) == 0 {
		return nil, fmt.Errorf("missing task type")
	}

	toolKey, err := CreateToolKey(toolName, toolID)
	if toolKey == nil {
		return nil, err
	}

	return &TaskTypeKey{
		tool:     *toolKey,
		taskType: taskType,
	}, nil
}

type TaskConfig interface {
	field.Interface

	Name() string
}

func NewTaskBase(name string) *TaskBase {
	return &TaskBase{
		name: name,
	}
}

type TaskBase struct {
	name string
}

func (t *TaskBase) Name() string {
	return t.name
}
