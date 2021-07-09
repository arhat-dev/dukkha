package tools

import (
	"fmt"
	"strings"

	"arhat.dev/pkg/log"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
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

func (h *TaskHooks) GenSpecs(
	taskCtx dukkha.TaskExecContext,
	stage dukkha.TaskExecStage,
) ([][]dukkha.TaskExecSpec, error) {
	logger := log.Log.WithName("TaskHooks").WithFields(
		log.String("stage", stage.String()),
	)

	logger.D("resolving hooks")
	err := h.ResolveFields(taskCtx, -1, h.GetFieldNameByStage(stage))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve hook spec: %w", err)
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
		return nil, fmt.Errorf("unknown task exec stage: %d", stage)
	}

	hookCtx := taskCtx.DeriveNew()
	prefix := taskCtx.OutputPrefix() + stage.String() + ": "
	hookCtx.SetOutputPrefix(prefix)

	var ret [][]dukkha.TaskExecSpec
	for i := range toRun {
		specs, err := toRun[i].GenSpecs(hookCtx.DeriveNew())
		if err != nil {
			return nil, fmt.Errorf(
				"failed to generate action #%d exec specs: %w",
				i, err,
			)
		}
		ret = append(ret, specs)
	}

	return ret, nil
}

type Hook struct {
	field.BaseField

	Task string `yaml:"task"`

	Other map[string]string `dukkha:"other"`
}

func (h *Hook) GenSpecs(ctx dukkha.Context) ([]dukkha.TaskExecSpec, error) {
	if len(h.Task) != 0 {
		ref, err := dukkha.ParseTaskReference(h.Task, ctx.CurrentTool().Name)
		if err != nil {
			return nil, fmt.Errorf("invalid task reference %q: %w", h.Task, err)
		}

		if len(ref.MatrixFilter) != 0 {
			ctx.SetMatrixFilter(ref.MatrixFilter)
		}

		tool, ok := ctx.GetTool(ref.ToolKey())
		if !ok {
			return nil, fmt.Errorf("referenced tool %q not found", ref.ToolKey())
		}

		tsk, ok := tool.GetTask(ref.TaskKey())
		if !ok {
			return nil, fmt.Errorf("referenced task %q not found", ref.TaskKey())
		}

		specs, err := tsk.GetExecSpecs(ctx, tool.UseShell(), tool.ShellName(), tool.GetCmd())
		if err != nil {
			return nil, fmt.Errorf("failed to generate task exec specs: %w", err)
		}

		return specs, nil
	}

	switch {
	case len(h.Other) > 1:
		return nil, fmt.Errorf("unexpected multiple entries in one hook spec")
	case len(h.Other) == 1:
	default:
		// no hook to run
		return nil, nil
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
			return nil, fmt.Errorf("unknown action: %q", k)
		}
	}

	sh, ok := ctx.GetShell(shell)
	if !ok {
		return nil, fmt.Errorf("shell %q not found", shell)
	}

	env, cmd, err := sh.GetExecSpec([]string{script}, isFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate shell ")
	}
	ctx.AddEnv(env...)

	return []dukkha.TaskExecSpec{
		{
			Env:     sliceutils.FormatStringMap(ctx.Env(), "="),
			Command: cmd,
		},
	}, nil
}
