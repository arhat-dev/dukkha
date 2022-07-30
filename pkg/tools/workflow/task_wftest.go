package workflow

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindTest = "test"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindTest,
		tools.NewTask[TaskTest, *TaskTest],
	)
}

type TaskTest struct {
	tools.BaseTask[WorkflowTest, *WorkflowTest]
}

// nolint:revive
type WorkflowTest struct{}

func (w *WorkflowTest) ToolKind() dukkha.ToolKind       { return ToolKind }
func (w *WorkflowTest) Kind() dukkha.TaskKind           { return TaskKindTest }
func (w *WorkflowTest) LinkParent(p tools.BaseTaskType) {}

func (w *WorkflowTest) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}
