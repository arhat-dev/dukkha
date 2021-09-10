//go:build real
// +build real

package tools

import (
	"fmt"
	"io"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/templateutils"
	"mvdan.cc/sh/v3/syntax"
)

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
