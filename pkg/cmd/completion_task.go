package cmd

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/tools"
)

func handleTaskCompletion(
	args []string,
	toComplete string,
	allTools *map[tools.ToolKey]tools.Tool,
	toolSpecificTasks *map[tools.ToolKey][]tools.Task,
) ([]string, cobra.ShellCompDirective) {
	alreadyFallthrough := false
	var alreadyComplete bool
	var ret []string
	switch len(args) {
	case 0:
		ret, alreadyComplete = tryFindToolKinds(*allTools, toComplete)
		if !alreadyComplete {
			break
		}

		args = append(args, toComplete)
		alreadyFallthrough = true
		fallthrough
	case 1:
		toolKind := args[0]
		// case 1: trying to use default tool, expecting task kind
		if len(toComplete) != 0 {
			ret, alreadyComplete = tryFindDefaultToolTaskKinds(*toolSpecificTasks, toolKind, toComplete)
			if len(ret) != 0 {
				goto endCase1
			}

			ret, alreadyComplete = tryFindToolNames(*allTools, toolKind, toComplete)
			goto endCase1
		} else {
			ret, alreadyComplete = tryFindToolNames(*allTools, toolKind, toComplete)
			goto endCase1
		}

	endCase1:
		if alreadyFallthrough {
			break
		}

		if !alreadyComplete {
			break
		}

		args = append(args, toComplete)
		alreadyFallthrough = true
		fallthrough
	case 2:
		toolKind := args[0]

		// case 1: arg1 is tool name, expecting task kind
		ret, alreadyComplete = tryFindTaskKindsWithToolName(
			*toolSpecificTasks, toolKind, args[1], toComplete,
		)
		if len(ret) != 0 {
			goto endCase2
		}

		// case 2: arg1 is task kind
		ret, alreadyComplete = tryFindTaskNamesWithDefaultTool(
			*toolSpecificTasks, toolKind, args[1], toComplete,
		)
		if len(ret) != 0 {
			goto endCase2
		}
	endCase2:
		if alreadyFallthrough {
			break
		}

		if !alreadyComplete {
			break
		}

		ret = []string{}
		args = append(args, toComplete)
		fallthrough
	case 3:
		targetToolKind := args[0]

		// case 1: already a valid invocation
		targetToolName := ""
		targetTaskKind := args[1]
		targetTaskName := args[2]
		key := tools.ToolKey{ToolKind: targetToolKind, ToolName: targetToolName}

		defaultToolTasks, ok := (*toolSpecificTasks)[key]
		if ok {
			for _, v := range defaultToolTasks {
				if v.TaskKind() == targetTaskKind && v.TaskName() == targetTaskName {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
			}
		}

		// case 2: missing task name
		targetToolName = args[1]
		targetTaskKind = args[2]

		key = tools.ToolKey{ToolKind: targetToolKind, ToolName: targetToolName}
		toolTasks, ok := (*toolSpecificTasks)[key]
		if !ok {
			// no such tasks
			return nil, cobra.ShellCompDirectiveNoSpace
		}

		hasLongerCanditates := false
		visited := make(map[string]struct{})
		for _, v := range toolTasks {
			if v.TaskKind() != targetTaskKind {
				continue
			}

			if _, ok := visited[v.TaskName()]; ok {
				continue
			}

			if !strings.HasPrefix(v.TaskName(), toComplete) {
				continue
			}

			if !hasLongerCanditates {
				hasLongerCanditates = len(toComplete) < len(v.TaskName())
			}

			ret = append(ret, v.TaskName())
			visited[v.TaskName()] = struct{}{}
		}

		if _, ok := visited[toComplete]; !ok || hasLongerCanditates {
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
	allTools map[tools.ToolKey]tools.Tool,
	toComplete string,
) (ret []string, alreadyComplete bool) {
	hasLongerCanditates := false
	visited := make(map[string]struct{})
	for k := range allTools {
		if _, ok := visited[k.ToolKind]; ok {
			continue
		}

		if !strings.HasPrefix(k.ToolKind, toComplete) {
			continue
		}

		if !hasLongerCanditates {
			hasLongerCanditates = len(toComplete) < len(k.ToolKind)
		}

		ret = append(ret, k.ToolKind)
		visited[k.ToolKind] = struct{}{}
	}

	if _, ok := visited[toComplete]; !ok {
		return ret, false
	}

	return nil, true
}

func tryFindDefaultToolTaskKinds(
	toolSpecificTasks map[tools.ToolKey][]tools.Task,
	toolKind string,
	toComplete string,
) (ret []string, alreadyComplete bool) {
	key := tools.ToolKey{ToolKind: toolKind, ToolName: ""}
	toolTasks, ok := toolSpecificTasks[key]
	if !ok {
		return nil, false
	}

	hasLongerCanditates := false
	visited := make(map[string]struct{})
	for _, v := range toolTasks {
		if _, ok := visited[v.TaskKind()]; ok {
			continue
		}

		if !strings.HasPrefix(v.TaskKind(), toComplete) {
			continue
		}

		if !hasLongerCanditates {
			hasLongerCanditates = len(toComplete) < len(v.TaskKind())
		}

		ret = append(ret, v.TaskKind())
		visited[v.TaskKind()] = struct{}{}
	}

	if _, ok := visited[toComplete]; !ok || hasLongerCanditates {
		return ret, false
	}

	return nil, true
}

func tryFindToolNames(
	allTools map[tools.ToolKey]tools.Tool,
	toolKind string,
	toComplete string,
) (ret []string, alreadyComplete bool) {
	hasLongerCanditates := false
	visited := make(map[string]struct{})
	for k := range allTools {
		if len(k.ToolName) == 0 {
			continue
		}

		if k.ToolKind != toolKind {
			continue
		}

		if _, ok := visited[k.ToolName]; ok {
			continue
		}

		if !strings.HasPrefix(k.ToolName, toComplete) {
			continue
		}

		if !hasLongerCanditates {
			hasLongerCanditates = len(toComplete) < len(k.ToolName)
		}

		ret = append(ret, k.ToolName)
		visited[k.ToolName] = struct{}{}
	}

	if _, ok := visited[toComplete]; !ok {
		return ret, false
	}

	return nil, true
}

func tryFindTaskNamesWithDefaultTool(
	toolSpecificTasks map[tools.ToolKey][]tools.Task,
	toolKind string,
	taskKind string,
	toComplete string,
) (ret []string, alreadyComplete bool) {
	key := tools.ToolKey{ToolKind: toolKind, ToolName: ""}
	tasks, ok := toolSpecificTasks[key]
	if !ok {
		return nil, false
	}

	hasLongerCanditates := false
	visited := make(map[string]struct{})
	for _, v := range tasks {
		if taskKind != v.TaskKind() {
			continue
		}

		if _, ok := visited[v.TaskName()]; ok {
			continue
		}

		if !strings.HasPrefix(v.TaskName(), toComplete) {
			continue
		}

		if !hasLongerCanditates {
			hasLongerCanditates = len(toComplete) < len(v.TaskName())
		}

		ret = append(ret, v.TaskName())
		visited[v.TaskName()] = struct{}{}
	}

	if _, ok := visited[toComplete]; !ok || hasLongerCanditates {
		return ret, false
	}

	return nil, true
}

func tryFindTaskKindsWithToolName(
	toolSpecificTasks map[tools.ToolKey][]tools.Task,
	toolKind string,
	toolName string,
	toComplete string,
) (ret []string, alreadyComplete bool) {
	key := tools.ToolKey{ToolKind: toolKind, ToolName: toolName}
	tasks, ok := toolSpecificTasks[key]
	if !ok {
		return nil, false
	}

	hasLongerCanditates := false
	visited := make(map[string]struct{})
	for _, v := range tasks {
		if _, ok := visited[v.TaskKind()]; ok {
			continue
		}

		if !strings.HasPrefix(v.TaskKind(), toComplete) {
			continue
		}

		if !hasLongerCanditates {
			hasLongerCanditates = len(toComplete) < len(v.TaskKind())
		}

		ret = append(ret, v.TaskKind())
		visited[v.TaskKind()] = struct{}{}
	}

	if _, ok := visited[toComplete]; !ok || hasLongerCanditates {
		return ret, false
	}

	return nil, true
}
