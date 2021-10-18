package debug

import (
	"fmt"
	"os"

	"arhat.dev/dukkha/pkg/cmd/utils"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewDebugTaskSpecCmd(ctx *dukkha.Context) *cobra.Command {
	var matrixFilter []string

	debugTaskSpecCmd := &cobra.Command{
		Use:   "spec",
		Short: "Print task field values at runtime",

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
			appCtx.SetMatrixFilter(utils.ParseMatrixFilter(matrixFilter))

			return debugTasks(*ctx, args,
				func(appCtx dukkha.Context, tool dukkha.Tool, task dukkha.Task) error {
					matrixSpecs, err := task.GetMatrixSpecs(appCtx)
					if err != nil {
						return fmt.Errorf("failed to get task matrix specs: %w", err)
					}

					execOpts := dukkha.CreateTaskExecOptions(0, len(matrixSpecs))
					tskCtx := appCtx.DeriveNew()
					tskCtx.SetTask(tool.Key(), task.Key())

					enc := yaml.NewEncoder(os.Stdout)
					enc.SetIndent(2)
					defer func() { _ = enc.Close() }()

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
							return fmt.Errorf("failed to create task matrix context: %w", err2)
						}

						err2 = task.DoAfterFieldsResolved(mCtx, -1, func() error {
							return enc.Encode(task)
						})
						if err2 != nil {
							return fmt.Errorf("failed to generate resolved yaml: %w", err2)
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
