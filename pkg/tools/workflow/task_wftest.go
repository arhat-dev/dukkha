package workflow

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindTest = "test"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindTest,
		func(toolName string) dukkha.Task {
			t := &TaskTest{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindTest, t)
			return t
		},
	)
}

type TaskTest struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`
}

func (w *TaskTest) GetExecSpecs(
	rc dukkha.TaskExecContext, options *dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}