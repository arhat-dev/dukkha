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
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindRun, t)
			return t
		},
	)
}

type TaskRun struct {
	rs.BaseField `yaml:"-"`

	tools.BaseTask `yaml:",inline"`

	Jobs tools.Actions `yaml:"jobs"`
}

func (w *TaskRun) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return tools.ResolveActions(rc, w, "Jobs", "jobs", options)
}
