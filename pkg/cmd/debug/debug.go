package debug

import (
	"fmt"

	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/cmd/utils"
	"arhat.dev/dukkha/pkg/dukkha"
)

type cliOptions struct {
	ShowMatrix   bool
	MatrixFilter []string
}

func (opts *cliOptions) resolve() *reportOptions {
	return &reportOptions{
		cliOptions:   *opts,
		matrixFilter: utils.ParseMatrixFilter(opts.MatrixFilter),
	}
}

func NewDebugCmd(ctx *dukkha.Context) *cobra.Command {
	var (
		opts = &cliOptions{}
	)

	debugCmd := &cobra.Command{
		Use:           "debug",
		Short:         "Debug config and task definitions",
		Args:          cobra.RangeArgs(0, 4),
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: true,
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			appCtx := *ctx

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

			return printTasks(appCtx, toolKind, toolName, taskKind, taskName, opts.resolve())
		},
	}

	const (
		matrixFlagName = "matrix"
	)

	flags := debugCmd.Flags()
	flags.StringSliceVarP(&opts.MatrixFilter, matrixFlagName, "m", nil,
		"set matrix filter, format: `-m <name>=<value>` for matching, `-m <name>!=<value>` for ignoring",
	)

	flags.BoolVar(&opts.ShowMatrix, "show-matrix", true, "show matrix specs")

	debugCmd.ValidArgsFunction = func(
		cmd *cobra.Command, args []string, toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return utils.HandleCompletionTask(*ctx, args, toComplete)
	}

	err := debugCmd.RegisterFlagCompletionFunc(matrixFlagName,
		func(
			cmd *cobra.Command, args []string, toComplete string,
		) ([]string, cobra.ShellCompDirective) {
			filter, _ := debugCmd.Flags().GetStringSlice(matrixFlagName)
			return utils.HandleCompletionMatrix(*ctx, filter, args, toComplete)
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to register matrix filter autocompletion: %w", err))
	}

	return debugCmd
}

func printTasks(
	appCtx dukkha.Context,
	toolKind dukkha.ToolKind,
	toolName dukkha.ToolName,
	taskKind dukkha.TaskKind,
	taskName dukkha.TaskName,
	ropts *reportOptions,
) error {
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

	// gather tasks with same task kind
	type taskPartialKey struct {
		toolKind dukkha.ToolKind
		toolName dukkha.ToolName
		taskKind dukkha.TaskKind
	}

	var (
		taskReports = make(map[taskPartialKey][]*taskReport)
		errTasks    []dukkha.Task
	)

	for fk, tsk := range allTasks {
		pk := taskPartialKey{
			toolKind: fk.toolKind,
			toolName: fk.toolName,
			taskKind: fk.taskKind,
		}

		toolKey := dukkha.ToolKey{
			Kind: fk.toolKind,
			Name: fk.toolName,
		}
		tool, ok := appCtx.GetTool(toolKey)
		if !ok {
			return fmt.Errorf("unexpected tool %q not found", toolKey.String())
		}

		report, err := ropts.generateTaskReport(appCtx, tool, tsk)
		if err != nil {
			errTasks = append(errTasks, allTasks[fk])
			// TODO: print tasks failed to report
			_ = errTasks
		} else {
			taskReports[pk] = append(taskReports[pk], report)
		}
	}

	return nil
}
