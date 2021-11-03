package tools

import (
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
	Before Actions `yaml:"before,omitempty"`

	// Matrix scope hooks

	// Before a specific matrix execution start
	//
	// This hook May have reference to matrix information
	BeforeMatrix Actions `yaml:"before:matrix,omitempty"`

	// AfterMatrixSuccess runs after a successful matrix execution
	//
	// This hook May have reference to matrix information
	AfterMatrixSuccess Actions `yaml:"after:matrix:success,omitempty"`

	// AfterMatrixFailure runs after a failed matrix execution
	//
	// This hook May have reference to matrix information
	AfterMatrixFailure Actions `yaml:"after:matrix:failure,omitempty"`

	// AfterMatrix runs after at any condition of the matrix execution
	// including success, failure
	//
	// This hook May have reference to matrix information
	AfterMatrix Actions `yaml:"after:matrix,omitempty"`

	// Task scope hooks again

	// AfterSuccess runs after a successful task execution
	// requires all matrix executions are successful
	//
	// This hook MUST NOT have any reference to matrix information
	AfterSuccess Actions `yaml:"after:success,omitempty"`

	// AfterFailure runs after a failed task execution
	// any failed matrix execution will cause this hook to run
	//
	// This hook MUST NOT have any reference to matrix information
	AfterFailure Actions `yaml:"after:failure,omitempty"`

	// After any condition of the task execution
	// including success, failure, canceled (hook `before` failure)
	//
	// This hook MUST NOT have any reference to matrix information
	After Actions `yaml:"after,omitempty"`
}

func (*TaskHooks) getTagNameByStage(stage dukkha.TaskExecStage) [2]string {
	return map[dukkha.TaskExecStage][2]string{
		dukkha.StageBefore: {"Before", "before"},

		dukkha.StageBeforeMatrix:       {"BeforeMatrix", "before:matrix"},
		dukkha.StageAfterMatrixSuccess: {"AfterMatrixSuccess", "after:matrix:success"},
		dukkha.StageAfterMatrixFailure: {"AfterMatrixFailure", "after:matrix:failure"},
		dukkha.StageAfterMatrix:        {"AfterMatrix", "after:matrix"},

		dukkha.StageAfterSuccess: {"AfterSuccess", "after:success"},
		dukkha.StageAfterFailure: {"AfterFailure", "after:failure"},
		dukkha.StageAfter:        {"After", "after"},
	}[stage]
}

func (h *TaskHooks) GenSpecs(
	taskCtx dukkha.TaskExecContext,
	stage dukkha.TaskExecStage,
) ([]dukkha.TaskExecSpec, error) {
	// TODO: this func is only called by BaseTask with lock for now
	// 		 if we call it from other places, we need to use lock in
	// 		 DoAfterFieldsResolved
	ftNames := h.getTagNameByStage(stage)
	return ResolveActions(taskCtx.DeriveNew(), h, ftNames[0], ftNames[1], nil)
}

func (h *TaskHooks) DoAfterFieldsResolved(
	ctx dukkha.RenderingContext, depth int, do func() error, names ...string,
) error {
	err := h.ResolveFields(ctx, depth, names...)
	if err != nil {
		return err
	}

	return do()
}
