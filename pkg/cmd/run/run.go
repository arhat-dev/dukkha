package run

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"arhat.dev/dukkha/pkg/cmd/utils"
	"arhat.dev/dukkha/pkg/dukkha"
)

func NewRunCmd(ctx *dukkha.Context) *cobra.Command {
	var (
		workerCount  = int(1)
		failFast     = false
		forceColor   = false
		matrixFilter []string

		translateANSIStream = false
		retainANSIStyle     = false
	)

	runCmd := &cobra.Command{
		Use:   "run <tool-kind> <tool-name> <task-kind> <task-name>",
		Short: "Run your task",
		Example: `dukkha run buildah local build my-image
dukkha run golang in-docker build my-executable`,

		SilenceErrors: true,
		SilenceUsage:  true,

		Args: cobra.ExactArgs(4),

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableNoDescFlag:   true,
			DisableDescriptions: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			appCtx := *ctx

			stdoutIsPty := term.IsTerminal(int(os.Stdout.Fd()))

			translateANSIFlag := cmd.Flag("translate-ansi-stream")
			var actualTranslateANSIStream bool
			// only translate ansi stream when
			// 	- (automatic) flag --translate-ansi-stream not set and stdout is not a pty
			actualTranslateANSIStream = (translateANSIFlag == nil || !translateANSIFlag.Changed) && !stdoutIsPty
			// 	- (manual) flag --translate-ansi-stream set to true
			actualTranslateANSIStream = actualTranslateANSIStream || translateANSIStream

			actualRetainANSIStyle := actualTranslateANSIStream && retainANSIStyle

			appCtx.SetRuntimeOptions(dukkha.RuntimeOptions{
				FailFast:            failFast,
				ColorOutput:         stdoutIsPty || forceColor,
				TranslateANSIStream: actualTranslateANSIStream,
				RetainANSIStyle:     actualRetainANSIStyle,
				Workers:             workerCount,
			})

			appCtx.SetMatrixFilter(utils.ParseMatrixFilter(matrixFilter))

			return run(appCtx, args)
		},
	}

	flags := runCmd.Flags()

	const (
		matrixFlagName = "matrix"
	)

	flags.IntVarP(&workerCount, "workers", "j", 1, "set parallel worker count")
	flags.BoolVar(&failFast, "fail-fast", true, "cancel all task execution after one errored")
	flags.BoolVar(&forceColor, "force-color", false, "force color output even when not given a tty")
	flags.StringSliceVarP(&matrixFilter, matrixFlagName, "m", nil,
		"set matrix filter, format: `-m <name>=<value>` for matching, `-m <name>!=<value>` for ignoring",
	)
	flags.BoolVar(&translateANSIStream, "translate-ansi-stream", false,
		"when set to true, will translate ansi stream to plain text before write to stdout/stderr, "+
			"when set to false, do nothing to the ansi stream, "+
			"when not set, will behavior as set to true if stdout/stderr is not a pty environment",
	)
	flags.BoolVar(&retainANSIStyle, "retain-ansi-style", retainANSIStyle,
		"when set to true, will retain ansi style when write to stdout/stderr, only effective "+
			"when ansi stream is going to be translated",
	)

	runCmd.ValidArgsFunction = func(
		cmd *cobra.Command, args []string, toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return utils.HandleCompletionTask(*ctx, args, toComplete)
	}

	runCmd.SetHelpCommand(&cobra.Command{
		SilenceUsage: true,
		Hidden:       true,
	})

	err := runCmd.RegisterFlagCompletionFunc(
		matrixFlagName,
		func(
			cmd *cobra.Command, args []string, toComplete string,
		) ([]string, cobra.ShellCompDirective) {
			filter, _ := runCmd.Flags().GetStringSlice(matrixFlagName)
			return utils.HandleCompletionMatrix(
				*ctx, filter, args, toComplete,
			)
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to setup matrix flag completion: %w", err))
	}

	return runCmd
}

func run(appCtx dukkha.Context, args []string) error {
	// defensive check, arg count should be guarded by cobra
	if len(args) != 4 {
		return fmt.Errorf("expecting 4 args, got %d", len(args))
	}

	return appCtx.RunTask(
		dukkha.ToolKey{
			Kind: dukkha.ToolKind(args[0]),
			Name: dukkha.ToolName(args[1]),
		},
		dukkha.TaskKey{
			Kind: dukkha.TaskKind(args[2]),
			Name: dukkha.TaskName(args[3]),
		},
	)
}
