package tools

import (
	"fmt"

	"arhat.dev/pkg/log"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
)

type TaskHooks struct {
	rs.BaseField `yaml:"-"`

	// Before runs before the task execution start
	// if this hook failed, the whole task execution is canceled
	// and will run `After` hooks
	//
	// This hook MUST NOT have any reference to matrix information
	Before []Action `yaml:"before,omitempty"`

	// Matrix scope hooks

	// Before a specific matrix execution start
	//
	// This hook May have reference to matrix information
	BeforeMatrix []Action `yaml:"before:matrix,omitempty"`

	// AfterMatrixSuccess runs after a successful matrix execution
	//
	// This hook May have reference to matrix information
	AfterMatrixSuccess []Action `yaml:"after:matrix:success,omitempty"`

	// AfterMatrixFailure runs after a failed matrix execution
	//
	// This hook May have reference to matrix information
	AfterMatrixFailure []Action `yaml:"after:matrix:failure,omitempty"`

	// AfterMatrix runs after at any condition of the matrix execution
	// including success, failure
	//
	// This hook May have reference to matrix information
	AfterMatrix []Action `yaml:"after:matrix,omitempty"`

	// Task scope hooks again

	// AfterSuccess runs after a successful task execution
	// requires all matrix executions are successful
	//
	// This hook MUST NOT have any reference to matrix information
	AfterSuccess []Action `yaml:"after:success,omitempty"`

	// AfterFailure runs after a failed task execution
	// any failed matrix execution will cause this hook to run
	//
	// This hook MUST NOT have any reference to matrix information
	AfterFailure []Action `yaml:"after:failure,omitempty"`

	// After any condition of the task execution
	// including success, failure, canceled (hook `before` failure)
	//
	// This hook MUST NOT have any reference to matrix information
	After []Action `yaml:"after,omitempty"`
}

func (*TaskHooks) GetFieldNameByStage(stage dukkha.TaskExecStage) string {
	return map[dukkha.TaskExecStage]string{
		dukkha.StageBefore: "Before",

		dukkha.StageBeforeMatrix:       "BeforeMatrix",
		dukkha.StageAfterMatrixSuccess: "AfterMatrixSuccess",
		dukkha.StageAfterMatrixFailure: "AfterMatrixFailure",
		dukkha.StageAfterMatrix:        "AfterMatrix",

		dukkha.StageAfterSuccess: "AfterSuccess",
		dukkha.StageAfterFailure: "AfterFailure",
		dukkha.StageAfter:        "After",
	}[stage]
}

func (h *TaskHooks) GenSpecs(
	taskCtx dukkha.TaskExecContext,
	stage dukkha.TaskExecStage,
) ([]dukkha.RunTaskOrRunCmd, error) {
	// TODO: this func is only called by BaseTask with lock for now
	// 		 if we call it from other places, we need to use lock here

	logger := log.Log.WithName("TaskHooks").WithFields(
		log.String("stage", stage.String()),
	)

	logger.D("resolving hooks for overview")
	// just to get a list of hook actions available
	err := h.ResolveFields(taskCtx, 1, h.GetFieldNameByStage(stage))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve hook spec: %w", err)
	}

	toRun, ok := map[dukkha.TaskExecStage][]Action{
		dukkha.StageBefore: h.Before,

		dukkha.StageBeforeMatrix:       h.BeforeMatrix,
		dukkha.StageAfterMatrixSuccess: h.AfterMatrixSuccess,
		dukkha.StageAfterMatrixFailure: h.AfterMatrixFailure,
		dukkha.StageAfterMatrix:        h.AfterMatrix,

		dukkha.StageAfterSuccess: h.AfterSuccess,
		dukkha.StageAfterFailure: h.AfterFailure,
		dukkha.StageAfter:        h.After,
	}[stage]
	if !ok {
		return nil, fmt.Errorf("unknown task exec stage: %d", stage)
	}

	hookCtx := taskCtx.DeriveNew()
	prefix := taskCtx.OutputPrefix() + stage.String() + ": "
	hookCtx.SetOutputPrefix(prefix)

	var ret []dukkha.RunTaskOrRunCmd
	for i := range toRun {
		ctx := hookCtx.DeriveNew()
		err = toRun[i].DoAfterFieldResolved(ctx, func(h *Action) error {
			spec, err2 := h.GenSpecs(ctx, i)
			if err2 != nil {
				return err2
			}

			ret = append(ret, spec)
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf(
				"failed to generate action #%d exec specs: %w",
				i, err,
			)
		}
	}

	return ret, nil
}
