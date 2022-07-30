package workflow

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindRun = "run"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindRun, tools.NewTask[TaskRun, *TaskRun])
}

type WorkflowRun struct {
	Jobs tools.Actions `yaml:"jobs"`

	parent tools.BaseTaskType
}

type TaskRun struct {
	tools.BaseTask[WorkflowRun, *WorkflowRun]
}

func (w *WorkflowRun) ToolKind() dukkha.ToolKind       { return ToolKind }
func (w *WorkflowRun) Kind() dukkha.TaskKind           { return TaskKindRun }
func (w *WorkflowRun) LinkParent(p tools.BaseTaskType) { w.parent = p }

func (w *WorkflowRun) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return tools.ResolveActions(rc, w.parent, &w.Jobs, "jobs")
}
