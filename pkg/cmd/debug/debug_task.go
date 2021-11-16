package debug

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"arhat.dev/pkg/sorthelper"
	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
)

func NewDebugTaskCmd(ctx *dukkha.Context, options *Options) *cobra.Command {
	debugTaskCmd := &cobra.Command{
		Use:   "task",
		Short: "Show task related configuration in json",

		SilenceErrors: true,
		SilenceUsage:  true,

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: true,
		},
	}

	return debugTaskCmd
}

type TaskHeaderLineData struct {
	dukkha.ToolKind
	dukkha.ToolName
	dukkha.TaskKind
	dukkha.TaskName

	Matrix map[string]string
}

func (s TaskHeaderLineData) json() string {
	kind := string(s.ToolKind) + ":" + string(s.TaskKind)

	var parts []string
	parts = append(parts, `"kind": "`+kind+`"`)
	parts = append(parts, `"tool_name": "`+string(s.ToolName)+`"`)
	if len(s.TaskName) != 0 {
		parts = append(parts, `"name": "`+string(s.TaskName)+`"`)
	}

	if len(s.Matrix) != 0 {
		parts = append(parts, `"matrix": { "`+
			strings.Join(sliceutils.FormatStringMap(s.Matrix, `": "`, false), `", "`)+`" }`,
		)
	}

	return `{ ` + strings.Join(parts, `, `) + ` }`
}

type singleTaskDebugActionFunc func(
	appCtx dukkha.Context,
	tool dukkha.Tool,
	task dukkha.Task,
	idx int,
	count int,
) error

func forEachTask(
	appCtx dukkha.Context,
	args []string,
	debugSingleTask singleTaskDebugActionFunc,
) error {
	var (
		toolKind dukkha.ToolKind
		toolName dukkha.ToolName

		taskKind dukkha.TaskKind
		taskName dukkha.TaskName
	)
	switch len(args) {
	case 0:
		// all
		// print non task related info
		// TODO: implement
		return nil
	case 4:
		// <tool-kind> <tool-name> <task-kind> <task-name>
		taskName = dukkha.TaskName(args[3])
		fallthrough
	case 3:
		// <tool-kind> <tool-name> <task-kind>
		taskKind = dukkha.TaskKind(args[2])
		fallthrough
	case 2:
		// <tool-kind> <tool-name>
		// print tasks accessible by this tool
		toolName = dukkha.ToolName(args[1])
		fallthrough
	case 1:
		// <tool-kind>
		// print tool related tasks
		toolKind = dukkha.ToolKind(args[0])
	}

	if len(toolKind) == 0 {
		return fmt.Errorf("invalid no tool kind provided")
	}

	var tools []dukkha.Tool
	if len(toolName) == 0 {
		// no tool name, get all tools with this kind
		for k, v := range appCtx.AllTools() {
			if toolKind != k.Kind {
				continue
			}

			tools = append(tools, v)
		}
	} else {
		key := dukkha.ToolKey{
			Kind: toolKind,
			Name: toolName,
		}

		tool, ok := appCtx.GetTool(key)
		if !ok {
			return fmt.Errorf("tool %q not found", key.String())
		}

		tools = append(tools, tool)
	}

	type taskFullKey struct {
		toolKind dukkha.ToolKind
		toolName dukkha.ToolName
		taskKind dukkha.TaskKind
		taskName dukkha.TaskName
	}

	// ensure tasks are unique
	allTasks := make(map[taskFullKey]dukkha.Task)
	for _, tool := range tools {
		for _, tv := range tool.AllTasks() {
			// filter out unmatched tasks
			switch {
			case len(taskKind) != 0 && taskKind != tv.Kind(),
				len(taskName) != 0 && taskName != tv.Name():
				continue
			default:
				allTasks[taskFullKey{
					toolKind: tool.Kind(),
					toolName: tool.Name(),
					taskKind: tv.Kind(),
					taskName: tv.Name(),
				}] = tv
			}
		}
	}

	// get tasks sorted by
	// tool kind -> tool name -> task kind -> task name
	var (
		taskKeys = make([]taskFullKey, len(allTasks))
		tasks    = make([]dukkha.Task, len(allTasks))
	)

	i := 0
	for fk := range allTasks {
		taskKeys[i] = fk
		tasks[i] = allTasks[fk]
		i++
	}

	sortStub := sorthelper.NewCustomSortable(
		func(i, j int) {
			taskKeys[i], taskKeys[j] = taskKeys[j], taskKeys[i]
			tasks[i], tasks[j] = tasks[j], tasks[i]
		},
		func(i, j int) bool {
			a := taskKeys[i]
			b := taskKeys[j]

			// compare tool kind
			switch {
			case a.toolKind < b.toolKind:
				return true
			case a.toolKind == b.toolKind:
			default:
				return false
			}

			// compare tool name
			switch {
			case a.toolName < b.toolName:
				return true
			case a.toolName == b.toolName:
			default:
				return false
			}

			// compare task kind
			switch {
			case a.taskKind < b.taskKind:
				return true
			case a.taskKind == b.taskKind:
			default:
				return false
			}

			// compare task name
			switch {
			case a.taskName < b.taskName:
				return true
			case a.taskName == b.taskName:
			default:
				return false
			}

			return false
		},
		func() int { return len(taskKeys) },
	)

	sort.Sort(sortStub)

	for i, tsk := range tasks {
		toolKey := dukkha.ToolKey{
			Kind: taskKeys[i].toolKind,
			Name: taskKeys[i].toolName,
		}
		tool, ok := appCtx.GetTool(toolKey)
		if !ok {
			return fmt.Errorf("unexpected tool %q not found", toolKey.String())
		}

		err := debugSingleTask(appCtx, tool, tsk, i, len(tasks))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	return nil
}
