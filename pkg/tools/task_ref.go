package tools

import (
	"fmt"
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
)

type TaskReference struct {
	rs.BaseField

	// Ref is a task reference in following format:
	//
	//	<tool-kind>{:<tool-name>}:<task-kind>(<task-name>)
	//
	// when tool-name is not set, current tool-name is used
	//
	// e.g. when running `dukkha run workflow local run foo`, the tool-name is `local`.
	// 		say we reference `golang:build(bar)` in its job list, we will
	//		actually run `golang:local:build(bar)` in this case
	//
	Ref string `yaml:"ref"`

	// MatrixFilter for filtering matrix to be used
	//	- nil filter (no value or nil) means run all matrix tasks of the referenced task
	//	- empty filter (`{}`) means use current matrix entry (if exists)
	//	- non-empty cases are like task matrix
	MatrixFilter *matrix.Spec `yaml:"matrix_filter"`
}

func (tr *TaskReference) genTaskExecReq(
	ctx dukkha.TaskExecContext,
	hookID string,
	continueOnError bool,
) (*TaskExecRequest, error) {
	var (
		toolKind dukkha.ToolKind
		toolName dukkha.ToolName
		taskKind dukkha.TaskKind
		taskName dukkha.TaskName
	)

	tk, name, ok := strings.Cut(strings.TrimSpace(tr.Ref), "(")
	if !ok {
		return nil, fmt.Errorf("invalid task reference %q: missing task call `(<name>)`", tr.Ref)
	}

	taskName = dukkha.TaskName(strings.TrimSuffix(name, ")"))

	// <tool-kind>{:<tool-name>}:<task-kind>
	parts := strings.Split(tk, ":")
	toolKind = dukkha.ToolKind(parts[0])

	switch len(parts) {
	case 2:
		// no tool name set, use the default tool name
		// no matter what kind the tool is
		//
		// current task
		// 		buildah:in-docker:build 	# tool name is `in-docker`
		// has task reference in hook:
		// 		buildah:login(foo)    	# same kind
		// 		golang:build(bar)		# different kind
		// will actually be treated as
		// 		buildah:in-docker:login(foo)	# same kind
		//		golang:in-docker:build(bar)		# different kind

		toolName = ctx.CurrentTool().Name
		taskKind = dukkha.TaskKind(parts[1])
	case 3:
		toolName = dukkha.ToolName(parts[1])
		taskKind = dukkha.TaskKind(parts[2])
	default:
		return nil, fmt.Errorf(
			"invalid task reference %q: expecting <tool-kind>{:<tool-name>}:<task-kind>", tr.Ref,
		)
	}

	if tr.MatrixFilter == nil {
		// not set or nil, reset filter
		ctx.SetMatrixFilter(matrix.Filter{})
	} else if !tr.MatrixFilter.IsEmpty() {
		ctx.SetMatrixFilter(tr.MatrixFilter.AsFilter())
	} // else { /* empty filter, keep current filter */ }

	toolKey, taskKey := dukkha.ToolKey{
		Kind: toolKind,
		Name: toolName,
	}, dukkha.TaskKey{
		Kind: taskKind,
		Name: taskName,
	}

	tool, ok := ctx.GetTool(toolKey)
	if !ok {
		return nil, fmt.Errorf("%q: referenced tool %q not found", hookID, toolKey)
	}

	tsk, ok := tool.GetTask(taskKey)
	if !ok {
		return nil, fmt.Errorf("%q: referenced task %q not found", hookID, taskKey)
	}

	return &TaskExecRequest{
		Context:     ctx,
		Tool:        tool,
		Task:        tsk,
		IgnoreError: continueOnError,
	}, nil
}
