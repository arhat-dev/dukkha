package tools

import (
	"fmt"
	"os"
	"strings"

	"arhat.dev/pkg/exechelper"
	"github.com/fatih/color"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/output"
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

type TaskExecStage uint8

const (
	StageBefore TaskExecStage = iota + 1

	StageBeforeMatrix
	StageAfterMatrixSuccess
	StageAfterMatrixFailure
	StageAfterMatrix

	StageAfterSuccess
	StageAfterFailure
	StageAfter
)

func (s TaskExecStage) String() string {
	return map[TaskExecStage]string{
		StageBefore: "before",

		StageBeforeMatrix:       "before:matrix",
		StageAfterMatrixSuccess: "after:matrix:success",
		StageAfterMatrixFailure: "after:matrix:failure",
		StageAfterMatrix:        "after:matrix",

		StageAfterSuccess: "after:success",
		StageAfterFailure: "after:failure",
		StageAfter:        "after",
	}[s]
}

func (h *TaskHooks) Run(
	ctx *field.RenderingContext,
	rf field.RenderingFunc,
	stage TaskExecStage,
	prefix string,
	prefixColor, outputColor *color.Color,
	thisTool Tool,
	allTools map[ToolKey]Tool,
	allShells map[ToolKey]*BaseTool,
) error {
	// TODO: resolve specific hook only
	err := h.ResolveFields(ctx, rf, 1, true)
	if err != nil {
		return fmt.Errorf("failed to resolve hooks: %w", err)
	}

	toRun, ok := map[TaskExecStage][]Hook{
		StageBefore: h.Before,

		StageBeforeMatrix:       h.BeforeMatrix,
		StageAfterMatrixSuccess: h.AfterMatrixSuccess,
		StageAfterMatrixFailure: h.AfterMatrixFailure,
		StageAfterMatrix:        h.AfterMatrix,

		StageAfterSuccess: h.AfterSuccess,
		StageAfterFailure: h.AfterFailure,
		StageAfter:        h.After,
	}[stage]
	if !ok {
		return fmt.Errorf("unknown task exec stage: %d", stage)
	}

	for i := range toRun {
		err := toRun[i].ResolveFields(ctx, rf, -1, false)
		if err != nil {
			return fmt.Errorf("failed to resolve fields: %w", err)
		}

		err = toRun[i].Run(ctx, prefix, prefixColor, outputColor, thisTool, allTools, allShells)
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

func (h *Hook) Run(
	ctx *field.RenderingContext,
	prefix string,
	prefixColor, outputColor *color.Color,
	thisTool Tool,
	allTools map[ToolKey]Tool,
	allShells map[ToolKey]*BaseTool,
) error {
	if len(h.Task) != 0 {
		parts := strings.Split(h.Task, ":")

		var (
			taskKind string
			taskName string
		)

		key := ToolKey{
			ToolKind: parts[0],
			ToolName: "",
		}

		switch len(parts) {
		case 3:
			taskKind = parts[1]
			taskName = parts[2]

			if key.ToolKind == thisTool.ToolKind() {
				// same kind, but no tool name provided, use same tool to handle it
				return thisTool.Run(ctx.Context(), allTools, allShells, taskKind, taskName)
			}
		case 4:
			key.ToolName = parts[1]
			taskKind = parts[2]
			taskName = parts[3]
		default:
			return fmt.Errorf("invalid task reference: %q", h.Task)
		}

		// has tool name or using a different tool kind, find target tool to handle it
		tool, ok := allTools[key]
		if !ok {
			return fmt.Errorf("tool %q not found", key.ToolKind+":"+key.ToolName)
		}

		return tool.Run(ctx.Context(), allTools, allShells, taskKind, taskName)
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
		shellKey   *ToolKey
		script     string
		isFilePath bool
	)

	for k, v := range h.Other {
		script = v

		switch {
		case strings.HasPrefix(k, "shell_file:"):
			shellKey = &ToolKey{ToolKind: "shell", ToolName: strings.SplitN(k, ":", 2)[1]}
			isFilePath = true
		case k == "shell_file":
			shellKey = &ToolKey{ToolKind: "shell", ToolName: ""}
			isFilePath = true
		case strings.HasPrefix(k, "shell:"):
			shellKey = &ToolKey{ToolKind: "shell", ToolName: strings.SplitN(k, ":", 2)[1]}
			isFilePath = false
		case k == "shell":
			shellKey = &ToolKey{ToolKind: "shell", ToolName: ""}
			isFilePath = false
		default:
			return fmt.Errorf("unknown action: %q", k)
		}
	}

	if shellKey == nil {
		return nil
	}

	sh, ok := allShells[*shellKey]
	if !ok {
		return fmt.Errorf("shell %q not found", shellKey.ToolName)
	}

	scriptCtx := ctx.Clone()
	env, cmd, err := sh.GetExecSpec([]string{script}, isFilePath)
	if err != nil {
		return err
	}

	scriptCtx.AddEnv(env...)

	p, err := exechelper.Do(exechelper.Spec{
		Context: scriptCtx.Context(),
		Command: cmd,
		Env:     scriptCtx.Values().Env,

		Stdin:  os.Stdin,
		Stderr: output.PrefixWriter(prefix, prefixColor, outputColor, os.Stderr),
		Stdout: output.PrefixWriter(prefix, prefixColor, outputColor, os.Stdout),
	})
	if err != nil {
		return fmt.Errorf("failed to run script: %w", err)
	}

	code, err := p.Wait()
	if err != nil {
		return fmt.Errorf("command exited with code %d: %w", code, err)
	}

	return nil
}
