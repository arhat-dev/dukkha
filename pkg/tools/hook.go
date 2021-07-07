package tools

import (
	"fmt"
	"strings"

	"arhat.dev/pkg/log"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
)

type TaskHooks struct {
	field.BaseField

	// Before runs before the task execution start
	// if this hook failed, the whole task execution is canceled
	// and will run `After` hooks
	Before []Hook `yaml:"before"`

	// Matrix scope hooks

	// Before a specific matrix execution start
	BeforeMatrix []Hook `yaml:"before:matrix"`

	// AfterMatrixSuccess runs after a successful matrix execution
	AfterMatrixSuccess []Hook `yaml:"after:matrix:success"`

	// AfterMatrixFailure runs after a failed matrix execution
	AfterMatrixFailure []Hook `yaml:"after:matrix:failure"`

	// AfterMatrix runs after at any condition of the matrix execution
	// including success, failure
	AfterMatrix []Hook `yaml:"after:matrix"`

	// Task scope hooks again

	// AfterSuccess runs after a successful task execution
	// requires all matrix executions are successful
	AfterSuccess []Hook `yaml:"after:success"`

	// AfterFailure runs after a failed task execution
	// any failed matrix execution will cause this hook to run
	AfterFailure []Hook `yaml:"after:failure"`

	// After any condition of the task execution
	// including success, failure, canceled (hook `before` failure)
	After []Hook `yaml:"after"`
}

func (TaskHooks) GetFieldNameByStage(stage dukkha.TaskExecStage) string {
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

func (h *TaskHooks) Run(
	taskCtx dukkha.TaskExecContext,
	stage dukkha.TaskExecStage,
) error {
	logger := log.Log.WithName("TaskHooks").WithFields(
		log.String("stage", stage.String()),
	)

	logger.D("resolving hooks")
	err := h.ResolveFields(taskCtx, -1, h.GetFieldNameByStage(stage))
	if err != nil {
		return fmt.Errorf("failed to resolve hook spec: %w", err)
	}

	toRun, ok := map[dukkha.TaskExecStage][]Hook{
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
		return fmt.Errorf("unknown task exec stage: %d", stage)
	}

	hookCtx := taskCtx.DeriveNew()
	hookCtx.SetOutputPrefix(taskCtx.OutputPrefix() + " " + stage.String())

	for i := range toRun {
		err = toRun[i].Run(hookCtx.DeriveNew())
		if err != nil {
			return fmt.Errorf("action #%d failed: %w", i, err)
		}
	}

	return nil
}

type Hook struct {
	field.BaseField

	Task string `yaml:"task"`

	Other map[string]string `dukkha:"other"`
}

func (h *Hook) Run(hookCtx dukkha.Context) error {
	if len(h.Task) != 0 {
		ref, err := dukkha.ParseTaskReference(h.Task)
		if err != nil {
			return fmt.Errorf("invalid task reference %q: %w", h.Task, err)
		}

		if len(ref.MatrixFilter) != 0 {
			hookCtx.SetMatrixFilter(ref.MatrixFilter)
		}

		if !ref.HasToolName() && ref.ToolKind == hookCtx.CurrentTool().Kind() {
			toolName := hookCtx.CurrentTool().Name()
			return hookCtx.RunTask(
				ref.ToolKind, toolName,
				ref.TaskKind, ref.TaskName,
			)
		}

		return hookCtx.RunTask(
			ref.ToolKind, ref.ToolName,
			ref.TaskKind, ref.TaskName,
		)
	}

	switch {
	case len(h.Other) > 1:
		return fmt.Errorf("unexpected multiple entries in one hook spec")
	case len(h.Other) == 1:
	default:
		// no hook to run
		return nil
	}

	var (
		shell      string
		script     string
		isFilePath bool
	)

	for k, v := range h.Other {
		script = v

		switch {
		case strings.HasPrefix(k, "shell_file:"):
			shell = strings.SplitN(k, ":", 2)[1]
			isFilePath = true
		case k == "shell_file":
			shell = ""
			isFilePath = true
		case strings.HasPrefix(k, "shell:"):
			shell = strings.SplitN(k, ":", 2)[1]
			isFilePath = false
		case k == "shell":
			shell = ""
			isFilePath = false
		default:
			return fmt.Errorf("unknown action: %q", k)
		}
	}

	return hookCtx.RunShell(shell, script, isFilePath)
}
