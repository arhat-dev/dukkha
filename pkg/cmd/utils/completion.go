package utils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"arhat.dev/dukkha/pkg/dukkha"
)

const MatrixFilterFlagName = "matrix"

func RegisterMatrixFilterFlag(flags *pflag.FlagSet, matrixFilter *[]string) {
	flags.StringSliceVarP(matrixFilter, MatrixFilterFlagName, "m", nil,
		"set matrix filter, format: `-m <name>=<value>` for matching, `-m <name>!=<value>` for ignoring",
	)
}

func SetupTaskAndTaskMatrixCompletion(
	ctx *dukkha.Context,
	cmd *cobra.Command,
) error {
	cmd.ValidArgsFunction = func(
		cmd *cobra.Command, args []string, toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return handleTaskCompletion(*ctx, args, toComplete)
	}

	err := cmd.RegisterFlagCompletionFunc(MatrixFilterFlagName,
		func(
			cmd *cobra.Command, args []string, toComplete string,
		) ([]string, cobra.ShellCompDirective) {
			filter, _ := cmd.Flags().GetStringSlice(MatrixFilterFlagName)
			return handleTaskMatrixCompletion(*ctx, filter, args, toComplete)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to register matrix filter autocompletion: %w", err)
	}

	return nil
}
