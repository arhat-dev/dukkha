package debug

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/cmd/utils"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
)

func NewDebugTaskMatrixCmd(ctx *dukkha.Context) *cobra.Command {
	var matrixFilter []string

	debugTaskMatrixCmd := &cobra.Command{
		Use:   "matrix",
		Short: "Print task matrix at runtime",

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

			return debugTasks(appCtx, args,
				func(appCtx dukkha.Context, tool dukkha.Tool, task dukkha.Task) error {
					matrixSpecs, err := task.GetMatrixSpecs(appCtx)
					if err != nil {
						return fmt.Errorf("failed to get task matrix specs: %w", err)
					}

					for _, ms := range matrixSpecs {
						fmt.Fprintln(os.Stdout,
							"- { "+strings.Join(sliceutils.FormatStringMap(ms, ": ", false), ", ")+" }",
						)
					}

					return nil
				},
			)
		},
	}

	flags := debugTaskMatrixCmd.Flags()
	utils.RegisterMatrixFilterFlag(flags, &matrixFilter)

	err := utils.SetupTaskAndTaskMatrixCompletion(ctx, debugTaskMatrixCmd)
	if err != nil {
		panic(err)
	}

	debugTaskMatrixCmd.SetHelpCommand(&cobra.Command{
		SilenceUsage: true,
		Hidden:       true,
	})

	return debugTaskMatrixCmd
}
