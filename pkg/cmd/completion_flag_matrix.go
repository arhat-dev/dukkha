package cmd

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/dukkha"
)

func handleMatrixFlagCompletion(
	appCtx dukkha.Context,
	existingFilters []string,
	args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) != 4 {
		return nil, cobra.ShellCompDirectiveError
	}

	toolKind, toolName := dukkha.ToolKind(args[0]), dukkha.ToolName(args[1])
	taskKind, taskName := dukkha.TaskKind(args[2]), dukkha.TaskName(args[3])

	tasks, ok := appCtx.GetToolSpecificTasks(toolKind, toolName)
	if !ok {
		return nil, cobra.ShellCompDirectiveError
	}

	_, ok = appCtx.GetTool(toolKind, toolName)
	if !ok {
		return nil, cobra.ShellCompDirectiveError
	}

	var task dukkha.Task
	for i, v := range tasks {
		if v.Kind() == taskKind && v.Name() == taskName {
			task = tasks[i]
			break
		}
	}

	if task == nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// DO NOT apply existing filter, new filters with same key
	// are merged together
	mSpecs, err := task.GetMatrixSpecs(appCtx)
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
