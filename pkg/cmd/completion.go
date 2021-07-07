package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/dukkha"

	_ "embed"
)

var (
	//go:embed completion_guide.txt
	completionGuide string
)

func setupTaskCompletion(appCtx **dukkha.Context, rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long:  completionGuide,

		SilenceUsage: true,
		Hidden:       true,

		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactValidArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				_ = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	rootCmd.AddCommand(cmd)
	rootCmd.ValidArgsFunction = func(
		cmd *cobra.Command, args []string, toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return handleTaskCompletion(**appCtx, args, toComplete)
	}

	rootCmd.SetHelpCommand(&cobra.Command{
		SilenceUsage: true,
		Hidden:       true,
	})
}

func setupMatrixCompletion(
	appCtx **dukkha.Context,
	rootCmd *cobra.Command,
	matrixFlagName string,
) error {
	return rootCmd.RegisterFlagCompletionFunc(
		matrixFlagName,
		func(
			cmd *cobra.Command, args []string, toComplete string,
		) ([]string, cobra.ShellCompDirective) {
			filter, _ := rootCmd.Flags().GetStringSlice(matrixFlagName)
			return handleMatrixFlagCompletion(
				**appCtx, filter, args, toComplete,
			)
		},
	)
}
