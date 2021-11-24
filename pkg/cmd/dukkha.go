/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"arhat.dev/pkg/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/cmd/completion"
	"arhat.dev/dukkha/pkg/cmd/debug"
	"arhat.dev/dukkha/pkg/cmd/diff"
	"arhat.dev/dukkha/pkg/cmd/render"
	"arhat.dev/dukkha/pkg/cmd/run"
	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/dukkha"
)

// NewRootCmd creates the dukkha command with all sub commands added
func NewRootCmd() *cobra.Command {
	var (
		logConfig = new(log.Config)

		configPaths []string
		// merged config
		config = conf.NewConfig()

		appCtx dukkha.Context
	)

	appBaseCtx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		for range sigCh {
			cancel()
		}
	}()

	rootCmd := &cobra.Command{
		Use: "dukkha",

		SilenceErrors: true,
		SilenceUsage:  true,

		Args: cobra.ArbitraryArgs,

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableNoDescFlag:   true,
			DisableDescriptions: true,
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch {
			case strings.HasPrefix(cmd.Use, "version"),
				strings.HasPrefix(cmd.Use, "completion"):
				// they don't need to know config options at all
				return nil
			}

			// setup global logger for debugging
			err := log.SetDefaultLogger(log.ConfigSet{*logConfig})
			if err != nil {
				return fmt.Errorf("failed to initialize logger: %w", err)
			}

			logger := log.Log.WithName("pre-run")

			// read all configration files
			visitedPaths := make(map[string]struct{})
			err = conf.ReadConfigRecursively(
				afero.NewIOFS(afero.NewOsFs()),
				configPaths,
				!cmd.PersistentFlags().Changed("config"),
				&visitedPaths,
				config,
			)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			logger.V("initializing dukkha", log.Any("raw_config", config))

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working dir: %w", err)
			}

			_appCtx := dukkha.NewConfigResolvingContext(
				appBaseCtx, dukkha.GlobalInterfaceTypeHandler,
				createGlobalEnv(appBaseCtx, cwd),
			)

			_appCtx.AddListEnv(os.Environ()...)

			var needTasks bool
			switch {
			case strings.HasPrefix(cmd.Use, "render"),
				strings.HasPrefix(cmd.Use, "diff"):
				needTasks = false
			case strings.HasPrefix(cmd.Use, "debug"):
				needTasks = len(args) != 0
			default:
				// for sub commands: run
				needTasks = true
			}

			err = config.Resolve(_appCtx, needTasks)
			if err != nil {
				return fmt.Errorf("failed to resolve config: %w", err)
			}

			appCtx = _appCtx

			logger.D("dukkha initialized", log.Any("init_config", config))

			return nil
		},
	}

	globalFlags := rootCmd.PersistentFlags()
	globalFlags.StringSliceVarP(
		&configPaths, "config", "c", []string{".dukkha.yaml"},
		"path to your config files and directories, if a directory is provided"+
			"only files with .yaml extension in that directory are parsed",
	)

	// logging for debugging purpose
	globalFlags.StringVarP(
		&logConfig.Level, "log.level", "v",
		"info", "log level, one of [verbose, debug, info, error, silent]",
	)

	globalFlags.StringVar(
		&logConfig.Format, "log.format",
		"console", "log output format, one of [console, json]",
	)

	globalFlags.StringVar(
		&logConfig.File, "log.file",
		"stderr", "file path to write log output, including `stdout` and `stderr`",
	)

	debugCmd, debugCmdOpts := debug.NewDebugCmd(&appCtx)
	debugTaskCmd := debug.NewDebugTaskCmd(&appCtx, debugCmdOpts)
	debugTaskCmd.AddCommand(
		debug.NewDebugTaskListCmd(&appCtx, debugCmdOpts),
		debug.NewDebugTaskMatrixCmd(&appCtx, debugCmdOpts),
		debug.NewDebugTaskSpecCmd(&appCtx, debugCmdOpts),
	)

	debugCmd.AddCommand(
		debugTaskCmd,
	)

	rootCmd.AddCommand(
		// completion
		completion.NewCompletionCmd(),
		// dukkha render
		render.NewRenderCmd(&appCtx),
		// dukkha debug
		debugCmd,
		// dukkha run
		run.NewRunCmd(&appCtx),
		diff.NewDiffCmd(&appCtx),
	)

	return rootCmd
}
