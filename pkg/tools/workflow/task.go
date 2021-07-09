package workflow

import (
	"fmt"
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
	rc dukkha.RenderingContext,
	_useShell bool,
	_shellName string,
	_toolCmd []string,
) ([]dukkha.TaskExecSpec, error) {
	var ret []dukkha.TaskExecSpec
	for i, job := range w.Jobs {
		if len(job.Task) != 0 {
			// do task
			ref, err := dukkha.ParseTaskReference(job.Task, "")
			if err != nil {
				return nil, fmt.Errorf("invalid task reference at job#%d: %w", i, err)
			}

			_ = ref

			// TODO: deep copy current job, execute as a hook
			ret = append(ret, dukkha.TaskExecSpec{
				AlterExecFunc: func(
					replace map[string][]byte,
					stdin io.Reader,
					stdout, stderr io.Writer,
				) ([]dukkha.TaskExecSpec, error) {
					// TODO: implement
					return nil, nil
				},
			})
			continue
		}

		// run shell command
		ret = append(ret, dukkha.TaskExecSpec{
			Command: []string{},
		})
	}

	return ret, nil
}
