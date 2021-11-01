package debug

import (
	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/dukkha"
)

func NewDebugCmd(ctx *dukkha.Context) *cobra.Command {
	debugCmd := &cobra.Command{
		Use:           "debug",
		Short:         "Debug config and task definitions",
		SilenceErrors: true,
		SilenceUsage:  true,

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: true,
		},
	}

	return debugCmd
}
