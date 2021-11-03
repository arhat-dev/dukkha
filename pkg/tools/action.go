package tools

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"arhat.dev/rs"
	"mvdan.cc/sh/v3/syntax"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/templateutils"
)

// Action is a collection of all kinds of work can be done in a single step
// but only one kind of work is allowed in a single step
type Action struct {
	rs.BaseField `yaml:"-"`

	// Name of this action
	//
	// if set, can be used as value for `Next`
	//
	// Defaults to `#i` where i is the index of this action in slice (starting from 0)
	//
	// this field MUST NOT use any kind of rendering suffix, or will be set
	// to default
	Name string `yaml:"name"`

	// Env specific to this action
	Env dukkha.Env `yaml:"env"`

	// Idle does nothing but serves as a placeholder for preparation purpose
	// recommended usage of Idle action is to apply renderers like `template`
	// to do some task execution state related operation (e.g. set global
	// value with `dukkha.SetValue`)
	Idle interface{} `yaml:"idle,omitempty"`

	// Task reference of this action
	//
	// Task, Cmd, EmbeddedShell, ExternalShell are mutually exclusive
	Task string `yaml:"task,omitempty"`

	// EmbeddedShell using embedded shell
	//
	// Task, Cmd, EmbeddedShell, ExternalShell are mutually exclusive
	EmbeddedShell string `yaml:"shell,omitempty"`

	// EmbeddedShell script for this action
	//
	// Task, Cmd, EmbeddedShell, ExternalShell are mutually exclusive
	ExternalShell map[string]string `yaml:",omitempty" rs:"other"`

	// Cmd execution, not in any shell
	//
	// Task, Cmd, EmbeddedShell, ExternalShell are mutually exclusive
	Cmd []string `yaml:"cmd,omitempty"`

	// Chdir change working directory before executing command
	// this option only applies to Cmd, EmbeddedShell, ExternalShell action
	Chdir string `yaml:"chdir"`

	// ContuineOnError ignores error occurred in this action and continue
	// following actions in list (if any)
	ContinueOnError bool `yaml:"continue_on_error"`

	// Next action name
	// NOTE: this field is resolved after execution finished (right before leaving this action)
	//
	// Defaults to the next action in the same list
	Next *string `yaml:"next"`

	mu sync.Mutex
}

func (act *Action) DoAfterFieldResolved(
	mCtx dukkha.TaskExecContext, do func() error, tagNames ...string,
) error {
	act.mu.Lock()
	defer act.mu.Unlock()

	err := dukkha.ResolveEnv(act, mCtx, "Env", "env")
	if err != nil {
		return fmt.Errorf("failed to resolve action specific env: %w", err)
	}

	if len(tagNames) == 0 {
		tagNames = []string{
			"idle", "task", "shell", "cmd", "chdir",
			"continue_on_error", "ExternalShell",
		}
	}

	err = act.ResolveFields(mCtx, -1, tagNames...)
	if err != nil {
		return fmt.Errorf("failed to resolve fields: %w", err)
	}

	return do()
}

func (act *Action) GenSpecs(
	ctx dukkha.TaskExecContext, index int,
) (dukkha.RunTaskOrRunCmd, error) {
	actionID := "#" + strconv.FormatInt(int64(index), 10)
	if len(act.Name) != 0 {
		actionID = fmt.Sprintf("%s (%s)", act.Name, actionID)
	}

	switch {
	case act.Idle != nil:
		return []dukkha.TaskExecSpec{}, nil
	case len(act.Task) != 0:
		return act.genTaskActionSpecs(ctx, actionID)
	case len(act.Cmd) != 0:
		return act.genCmdActionSpecs(ctx, actionID)
	case len(act.EmbeddedShell) != 0:
		return act.genEmbeddedShellActionSpecs(ctx, actionID)
	default:
		return act.genExternalShellActionSpecs(ctx, actionID)
	}
}

func (act *Action) genTaskActionSpecs(
	ctx dukkha.TaskExecContext, hookID string,
) (dukkha.RunTaskOrRunCmd, error) {
	ref, err := dukkha.ParseTaskReference(act.Task, ctx.CurrentTool().Name)
	if err != nil {
		return nil, fmt.Errorf("%q: invalid task reference %q: %w", hookID, act.Task, err)
	}

	if ref.MatrixFilter != nil {
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

	return &TaskExecRequest{
		Context:     ctx,
		Tool:        tool,
		Task:        tsk,
		IgnoreError: act.ContinueOnError,
	}, nil
}

func (act *Action) genCmdActionSpecs(
	ctx dukkha.TaskExecContext, hookID string,
) (dukkha.RunTaskOrRunCmd, error) {
	_ = ctx
	_ = hookID
	return []dukkha.TaskExecSpec{
		{
			EnvOverride: act.Env.Clone(),
			Command:     sliceutils.NewStrings(act.Cmd),
			Chdir:       act.Chdir,
			IgnoreError: act.ContinueOnError,
		},
	}, nil
}

func (act *Action) genEmbeddedShellActionSpecs(
	ctx dukkha.TaskExecContext, hookID string,
) (dukkha.RunTaskOrRunCmd, error) {

	workingDir := act.Chdir
	script := act.EmbeddedShell

	ctx.AddEnv(true, act.Env...)

	return []dukkha.TaskExecSpec{{
		AlterExecFunc: func(
			replace dukkha.ReplaceEntries,
			stdin io.Reader,
			stdout, stderr io.Writer,
		) (dukkha.RunTaskOrRunCmd, error) {
			runner, err := templateutils.CreateEmbeddedShellRunner(
				workingDir, ctx, stdin, stdout, stderr,
			)
			if err != nil {
				return nil, fmt.Errorf("%q: failed to create embedded shell: %w", hookID, err)
			}

			parser := syntax.NewParser(syntax.Variant(syntax.LangBash))

			err = templateutils.RunScriptInEmbeddedShell(ctx, runner, parser, script)
			if err != nil {
				return nil, fmt.Errorf("%q: failed to run command in embedded shell: %w", hookID, err)
			}

			return nil, err
		},
		IgnoreError: act.ContinueOnError,
	}}, nil
}

func (act *Action) genExternalShellActionSpecs(
	ctx dukkha.TaskExecContext, hookID string,
) (dukkha.RunTaskOrRunCmd, error) {
	// check other shell
	_ = ctx
	switch {
	case len(act.ExternalShell) > 1:
		return nil, fmt.Errorf(
			"%q: unexpected multiple shell entries in one spec",
			hookID,
		)
	case len(act.ExternalShell) == 1:
		// ok
	default:
		// no hook to run
		return nil, nil
	}

	var (
		shell  string
		script string
	)

	for k, v := range act.ExternalShell {
		script = v

		switch {
		case strings.HasPrefix(k, "shell:"):
			shell = strings.SplitN(k, ":", 2)[1]
		default:
			return nil, fmt.Errorf("%q: unknown action: %q", hookID, k)
		}
	}

	return []dukkha.TaskExecSpec{
		{
			Command:     []string{script},
			EnvOverride: act.Env.Clone(),
			Chdir:       act.Chdir,
			UseShell:    true,
			ShellName:   shell,
			IgnoreError: act.ContinueOnError,
		},
	}, nil
}
