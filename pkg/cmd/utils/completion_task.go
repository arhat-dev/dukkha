package utils

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/dukkha"
)

func handleTaskCompletion(
	appCtx dukkha.Context,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	var (
		ret []string
	)

	switch len(args) {
	case 0:
		ret = tryFindToolKinds(
			appCtx.AllTools(), toComplete,
		)
	case 1:
		toolKind := args[0]
		// case 1: trying to use default tool, expecting task kind
		ret = tryFindToolNames(
			appCtx.AllTools(),
			dukkha.ToolKind(toolKind),
			toComplete,
		)
	case 2:
		toolKind := args[0]

		// arg1 is tool name, expecting task kind
		ret = tryFindTaskKindsWithToolName(
			appCtx.AllToolSpecificTasks(),
			dukkha.ToolKind(toolKind),
			dukkha.ToolName(args[1]),
			toComplete,
		)
	case 3:
		// missing task name
		targetTaskKind := dukkha.TaskKind(args[2])

		key := dukkha.ToolKey{
			Kind: dukkha.ToolKind(args[0]),
			Name: dukkha.ToolName(args[1]),
		}
		toolTasks, ok := appCtx.GetToolSpecificTasks(key)
		if !ok {
			// no such tasks
			return nil, cobra.ShellCompDirectiveNoSpace
		}

		hasLongerCanditates := false
		visited := make(map[dukkha.TaskName]struct{})
		for _, v := range toolTasks {
			if v.Kind() != targetTaskKind {
				continue
			}

			if _, ok := visited[v.Name()]; ok {
				continue
			}

			if !strings.HasPrefix(string(v.Name()), toComplete) {
				continue
			}

			if !hasLongerCanditates {
				hasLongerCanditates = len(toComplete) < len(v.Name())
			}

			ret = append(ret, string(v.Name()))
			visited[v.Name()] = struct{}{}
		}

		if _, ok := visited[dukkha.TaskName(toComplete)]; !ok || hasLongerCanditates {
			break
		}
	default:
		return []string{"-m"}, cobra.ShellCompDirectiveNoFileComp
	}

	if len(ret) == 0 {
		return nil, cobra.ShellCompDirectiveNoSpace
	}

	sort.Strings(ret)

	return ret, cobra.ShellCompDirectiveNoFileComp
}

func tryFindToolKinds(
	allTools map[dukkha.ToolKey]dukkha.Tool,
	toComplete string,
) (ret []string) {
	hasLongerCanditates := false
	visited := make(map[dukkha.ToolKind]struct{})
	for k := range allTools {
		if _, ok := visited[k.Kind]; ok {
			continue
		}

		if !strings.HasPrefix(string(k.Kind), toComplete) {
			continue
		}

		if !hasLongerCanditates {
			hasLongerCanditates = len(toComplete) < len(k.Kind)
		}

		ret = append(ret, string(k.Kind))
		visited[k.Kind] = struct{}{}
	}

	if _, ok := visited[dukkha.ToolKind(toComplete)]; !ok {
		return ret
	}

	return []string{toComplete}
}

func tryFindToolNames(
	allTools map[dukkha.ToolKey]dukkha.Tool,
	toolKind dukkha.ToolKind,
	toComplete string,
) (ret []string) {
	hasLongerCanditates := false
	visited := make(map[dukkha.ToolName]struct{})
	for k := range allTools {
		if len(k.Name) == 0 {
			continue
		}

		if k.Kind != toolKind {
			continue
		}

		if _, ok := visited[k.Name]; ok {
			continue
		}

		if !strings.HasPrefix(string(k.Name), toComplete) {
			continue
		}

		if !hasLongerCanditates {
			hasLongerCanditates = len(toComplete) < len(k.Name)
		}

		ret = append(ret, string(k.Name))
		visited[k.Name] = struct{}{}
	}

	if _, ok := visited[dukkha.ToolName(toComplete)]; !ok {
		return ret
	}

	return []string{toComplete}
}

func tryFindTaskKindsWithToolName(
	toolSpecificTasks map[dukkha.ToolKey][]dukkha.Task,
	toolKind dukkha.ToolKind,
	toolName dukkha.ToolName,
	toComplete string,
) (ret []string) {
	key := dukkha.ToolKey{Kind: toolKind, Name: toolName}
	tasks, ok := toolSpecificTasks[key]
	if !ok {
		return nil
	}

	hasLongerCanditates := false
	visited := make(map[dukkha.TaskKind]struct{})
	for _, v := range tasks {
		if _, ok := visited[v.Kind()]; ok {
			continue
		}

		if !strings.HasPrefix(string(v.Kind()), toComplete) {
			continue
		}

		if !hasLongerCanditates {
			hasLongerCanditates = len(toComplete) < len(v.Kind())
		}

		ret = append(ret, string(v.Kind()))
		visited[v.Kind()] = struct{}{}
	}

	if _, ok := visited[dukkha.TaskKind(toComplete)]; !ok || hasLongerCanditates {
		return ret
	}

	return []string{toComplete}
}
