package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"arhat.dev/pkg/textquery"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/cmd/utils"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/tools"
)

func NewDebugTaskSpecCmd(ctx *dukkha.Context, opts *Options) *cobra.Command {
	var (
		matrixFilter []string
	)

	debugTaskSpecCmd := &cobra.Command{
		Use:   "spec",
		Short: "Show task spec",
		Long: "Tasks may contain step by step jobs (e.g. actions in hooks), " +
			"if one such step is using rendering suffix, its value will not be printed",

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
			appCtx = appCtx.DeriveNew()
			appCtx.SetMatrixFilter(matrix.ParseMatrixFilter(matrixFilter))

			stdout, stderr := appCtx.Stdout(), appCtx.Stderr()

			query, err := opts.getQuery()
			if err != nil {
				return err
			}

			return forEachTask(appCtx, args,
				func(appCtx dukkha.Context, tool dukkha.Tool, task dukkha.Task, _, _ int) error {
					matrixSpecs, err := task.GetMatrixSpecs(appCtx)
					if err != nil {
						return fmt.Errorf("create task matrix specs: %w", err)
					}

					execOpts := dukkha.CreateTaskExecOptions(0, len(matrixSpecs))
					tskCtx := appCtx.DeriveNew()
					tskCtx.SetTask(tool.Key(), task.Key())

					for _, ms := range matrixSpecs {
						mCtx, mExecOpts, err2 := tools.CreateTaskMatrixContext(
							&tools.TaskExecRequest{
								Context: tskCtx,
								Tool:    tool,
								Task:    task,
							},
							ms, execOpts,
						)
						_ = mExecOpts
						if err2 != nil {
							return fmt.Errorf("creating task matrix context: %w", err2)
						}

						err2 = task.DoAfterFieldsResolved(mCtx, -1, true, func() error {
							var buf bytes.Buffer
							enc := yaml.NewEncoder(&buf)
							enc.SetIndent(2)
							defer func() { _ = enc.Close() }()

							dec := yaml.NewDecoder(&buf)

							err = enc.Encode(task)
							if err != nil {
								return err
							}

							var data any
							err = dec.Decode(&data)
							if err != nil {
								return err
							}

							if query != nil {
								var ret []any
								ret, err = textquery.RunQuery(query, data, nil)
								if err != nil {
									return err
								}

								switch len(ret) {
								case 0:
									data = nil
								case 1:
									data = ret[0]
								default:
									data = ret
								}
							}

							jenc := json.NewEncoder(&buf)
							jenc.SetIndent("", "  ")
							err = jenc.Encode(data)
							if err != nil {
								return err
							}

							err = opts.writeHeader(stdout, stderr, TaskHeaderLineData{
								ToolKind: tool.Kind(),
								ToolName: tool.Name(),
								TaskKind: task.Kind(),
								TaskName: task.Name(),
								Matrix:   ms,
							}.json())
							if err != nil {
								return err
							}

							_, err = io.Copy(stdout, &buf)
							return err
						})
						if err2 != nil {
							return fmt.Errorf("generate task spec: %w", err2)
						}
					}

					return nil
				},
			)
		},
	}

	flags := debugTaskSpecCmd.Flags()
	utils.RegisterMatrixFilterFlag(flags, &matrixFilter)
	err := utils.SetupTaskAndTaskMatrixCompletion(ctx, debugTaskSpecCmd)
	if err != nil {
		panic(err)
	}

	debugTaskSpecCmd.SetHelpCommand(&cobra.Command{
		SilenceUsage: true,
		Hidden:       true,
	})

	return debugTaskSpecCmd
}
