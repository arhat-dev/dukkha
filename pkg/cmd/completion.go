package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/tools"

	_ "embed"
)

var (
	//go:embed completion_guide.txt
	completionGuide string
)

func setupCompletion(
	rootCmd *cobra.Command,
	allTools *map[tools.ToolKey]tools.Tool,
	toolSpecificTasks *map[tools.ToolKey][]tools.Task,
) {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long:  completionGuide,

		DisableFlagsInUseLine: true,

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
		return handleTaskCompletion(args, toComplete, allTools, toolSpecificTasks)
	}

	err := rootCmd.RegisterFlagCompletionFunc(
		"matrix",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// TODO
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to register flag completion for --matrix: %w", err))
	}
}
