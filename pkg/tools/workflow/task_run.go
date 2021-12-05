package workflow

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindRun = "run"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindRun,
		func(toolName string) dukkha.Task {
			t := &TaskRun{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
			return t
		},
	)
}

type TaskRun struct {
	rs.BaseField `yaml:"-"`

	TaskName string `yaml:"name"`

	tools.BaseTask `yaml:",inline"`

	Jobs tools.Actions `yaml:"jobs"`
}

func (w *TaskRun) Kind() dukkha.TaskKind { return TaskKindRun }
func (w *TaskRun) Name() dukkha.TaskName { return dukkha.TaskName(w.TaskName) }

func (w *TaskRun) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: w.Kind(), Name: w.Name()}
}

func (w *TaskRun) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return tools.ResolveActions(rc, w, "Jobs", "jobs")
}
