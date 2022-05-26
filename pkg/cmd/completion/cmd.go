package completion

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"github.com/spf13/cobra"

	_ "embed"
)

var (
	//go:embed completion_guide.txt
	completionGuide string
)

func NewCompletionCmd(ctx *dukkha.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long:  completionGuide,

		SilenceUsage: true,
		Hidden:       true,

		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactValidArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			stdout := (*ctx).Stdout()
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(stdout)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(stdout)
			case "fish":
				_ = cmd.Root().GenFishCompletion(stdout, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletionWithDesc(stdout)
			default:
				return fmt.Errorf("unsupported shell %q", args[0])
			}

			return nil
		},
	}
}
