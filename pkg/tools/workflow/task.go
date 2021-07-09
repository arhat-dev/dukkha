package workflow

import (
	"io"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindRun = "run"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindRun,
		func(toolName string) dukkha.Task {
			t := &TaskRun{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindRun)
			return t
		},
	)
}

type TaskRun struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Jobs []tools.Hook `yaml:"jobs"`
}

func (w *TaskRun) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var ret []dukkha.TaskExecSpec
	for i, job := range w.Jobs {
		specs, err := job.GenSpecs(rc, dukkha.TaskExecOptions{}, i)
		if err != nil {
			return nil, err
		}

		ret = append(ret, dukkha.TaskExecSpec{
			AlterExecFunc: func(
				replace map[string][]byte,
				stdin io.Reader, stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunShell, error) {
				return specs, nil
			},
		})
	}

	return ret, nil
}
