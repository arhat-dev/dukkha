package tools

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

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
	options dukkha.TaskExecOptions,
) ([]dukkha.RunTaskOrRunShell, error) {
	// TODO: this func is only called by BaseTask with lock for now
	// 		 if we call it from other place, we need to add lock to it

	logger := log.Log.WithName("TaskHooks").WithFields(
		log.String("stage", stage.String()),
	)

	logger.D("resolving hooks for overview")
	// just to get a list of hook actions available
	err := h.ResolveFields(taskCtx, 1, h.GetFieldNameByStage(stage))
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

	var ret []dukkha.RunTaskOrRunShell
	for i := range toRun {
		ctx := hookCtx.DeriveNew()
		err = toRun[i].DoAfterFieldResolved(ctx, func(h *Hook) error {
			spec, err := h.GenSpecs(ctx, options, i)
			if err != nil {
				return err
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

type Hook struct {
	field.BaseField

	Name string `yaml:"name"`
	Task string `yaml:"task"`

	ContinueOnError bool `yaml:"continue_on_error"`

	// raw system cmd, not in shell
	Cmd []string `yaml:"cmd"`

	Other map[string]string `dukkha:"other"`

	mu sync.Mutex
}

func (h *Hook) DoAfterFieldResolved(mCtx dukkha.TaskExecContext, do func(h *Hook) error) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	err := h.ResolveFields(mCtx, -1, "")
	if err != nil {
		return fmt.Errorf("failed to resolve fields: %w", err)
	}

	return do(h)
}

func (h *Hook) GenSpecs(
	ctx dukkha.TaskExecContext, options dukkha.TaskExecOptions, index int,
) (dukkha.RunTaskOrRunShell, error) {
	hookID := "#" + strconv.FormatInt(int64(index), 10)
	if len(h.Name) != 0 {
		hookID = fmt.Sprintf("%s (%s)", h.Name, hookID)
	}

	if len(h.Task) != 0 {
		ref, err := dukkha.ParseTaskReference(h.Task, ctx.CurrentTool().Name)
		if err != nil {
			return nil, fmt.Errorf("%q: invalid task reference %q: %w", hookID, h.Task, err)
		}

		if len(ref.MatrixFilter) != 0 {
			ctx.SetMatrixFilter(ref.MatrixFilter)
		}

		tool, ok := ctx.GetTool(ref.ToolKey())
		if !ok {
			return nil, fmt.Errorf("%q: referenced tool %q not found", hookID, ref.ToolKey())
		}

		tsk, ok := tool.GetTask(ref.TaskKey())
		if !ok {
			return nil, fmt.Errorf("%q: referenced task %q not found", hookID, ref.TaskKey())
		}

		opts := options.Clone()
		// update tool specific values
		opts.ToolCmd = tool.GetCmd()
		opts.ShellName = tool.ShellName()
		opts.UseShell = tool.UseShell()
		opts.ContinueOnError = opts.ContinueOnError || h.ContinueOnError

		return &CompleteTaskExecSpecs{
			Context: ctx,
			Tool:    tool,
			Task:    tsk,
		}, nil
	}

	if len(h.Cmd) != 0 {
		return []dukkha.TaskExecSpec{
			{
				Env:         sliceutils.FormatStringMap(ctx.Env(), "="),
				Command:     sliceutils.NewStrings(h.Cmd),
				IgnoreError: h.ContinueOnError || options.ContinueOnError,
			},
		}, nil
	}

	switch {
	case len(h.Other) > 1:
		return nil, fmt.Errorf(
			"%q: unexpected multiple shell entries in one spec",
			hookID,
		)
	case len(h.Other) == 1:
		// ok
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
		case strings.HasPrefix(k, "shell:"):
			shell = strings.SplitN(k, ":", 2)[1]
			isFilePath = false
		case k == "shell":
			shell = ""
			isFilePath = false
		default:
			return nil, fmt.Errorf("%q: unknown action: %q", hookID, k)
		}
	}

	sh, ok := ctx.GetShell(shell)
	if !ok {
		return nil, fmt.Errorf("%q: shell %q not found", hookID, shell)
	}

	env, cmd, err := sh.GetExecSpec([]string{script}, isFilePath)
	if err != nil {
		return nil, fmt.Errorf("%q: failed to generate shell: %w", hookID, err)
	}
	ctx.AddEnv(env...)

	return []dukkha.TaskExecSpec{
		{
			Env:         sliceutils.FormatStringMap(ctx.Env(), "="),
			Command:     cmd,
			IgnoreError: h.ContinueOnError || options.ContinueOnError,
		},
	}, nil
}
