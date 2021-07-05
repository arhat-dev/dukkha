package cmd

import (
	"context"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

func handleMatrixFlagCompletion(
	appCtx *context.Context,
	rf field.RenderingFunc,
	existingFilters []string,
	allTools map[tools.ToolKey]tools.Tool,
	toolSpecificTasks map[tools.ToolKey][]tools.Task,
	args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	type taskKey struct {
		taskKind string
		taskName string
	}

	var (
		targetTool tools.ToolKey
		targetTask taskKey
	)

	switch len(args) {
	case 3:
		targetTool.ToolKind, targetTool.ToolName = args[0], ""
		targetTask.taskKind, targetTask.taskName = args[1], args[2]
	case 4:
		targetTool.ToolKind, targetTool.ToolName = args[0], args[1]
		targetTask.taskKind, targetTask.taskName = args[2], args[3]
	default:
		return nil, cobra.ShellCompDirectiveError
	}

	tasks, ok := toolSpecificTasks[targetTool]
	if !ok {
		return nil, cobra.ShellCompDirectiveError
	}

	tool, ok := allTools[targetTool]
	if !ok {
		return nil, cobra.ShellCompDirectiveError
	}

	var task tools.Task
	for i, v := range tasks {
		if v.TaskKind() == targetTask.taskKind && v.TaskName() == targetTask.taskName {
			task = tasks[i]
			break
		}
	}

	if task == nil {
		return nil, cobra.ShellCompDirectiveError
	}

	ctx := field.WithRenderingValues(*appCtx, os.Environ(), tool.GetEnv())

	// DO NOT apply existing filter, new filters with same key
	// are merged together
	mSpecs, err := task.GetMatrixSpecs(ctx, rf, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	usedPairs := make(map[string]struct{})
	for _, v := range existingFilters {
		usedPairs[v] = struct{}{}
	}

	var values []string
	visited := make(map[string]struct{})
	for _, spec := range mSpecs {
		for k, v := range spec {
			val := k + "=" + v
			_, ok := usedPairs[val]
			if ok {
				continue
			}

			if _, ok := visited[val]; ok {
				continue
			}

			if !strings.HasPrefix(val, toComplete) {
				continue
			}

			values = append(values, val)
			visited[val] = struct{}{}
		}
	}
	sort.Strings(values)

	return values, cobra.ShellCompDirectiveNoFileComp
}

func parseMatrixFilter(arr []string) map[string][]string {
	mf := make(map[string][]string)
	for _, v := range arr {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 1 {
			continue
		}

		mf[parts[0]] = append(mf[parts[0]], parts[1])
	}

	return mf
}
