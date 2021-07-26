package workflow

import (
	"io"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
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
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Jobs []tools.Action `yaml:"jobs"`
}

func (w *TaskRun) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return w.next(rc, options, 0)
}

func (w *TaskRun) next(
	mCtx dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	index int,
) ([]dukkha.TaskExecSpec, error) {
	var (
		thisAction dukkha.RunTaskOrRunCmd
		hasJob     = false
	)

	var err error
	// depth = 1 to get job list only
	err = w.DoAfterFieldsResolved(mCtx, 1, func() error {
		if index >= len(w.Jobs) {
			return nil
		}

		hasJob = true

		// resolve single job (Hook)
		return w.Jobs[index].DoAfterFieldResolved(mCtx, func(h *tools.Action) error {
			thisAction, err = h.GenSpecs(mCtx, index)
			return err
		})
	}, "Jobs")
	if err != nil || !hasJob {
		return nil, err
	}

	return []dukkha.TaskExecSpec{
		{
			Env: sliceutils.NewStrings(w.Env),
			AlterExecFunc: func(
				replace map[string][]byte,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				return thisAction, nil
			},
		},
		{
			Env: sliceutils.NewStrings(w.Env),
			AlterExecFunc: func(
				replace map[string][]byte,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				return w.next(mCtx, options, index+1)
			},
		},
	}, nil
}
