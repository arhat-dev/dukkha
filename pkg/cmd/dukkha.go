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
	"strings"

	"arhat.dev/pkg/log"
	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/tools"
)

func NewRootCmd() *cobra.Command {
	// cli options
	// nolint:revive
	var (
		cli_logConfig    = new(log.Config)
		cli_workerCount  int
		cli_failFast     bool
		cli_matrixFilter []string
	)

	// config file
	var (
		configPaths []string
		config      = conf.NewConfig()
	)

	// runtime data
	var (
		_appCtx = context.Background()
		appCtx  = &_appCtx

		renderingMgr = renderer.NewManager()

		allTools  = make(map[tools.ToolKey]tools.Tool)
		allShells = make(map[tools.ToolKey]*tools.BaseTool)

		toolSpecificTasks = make(map[tools.ToolKey][]tools.Task)
	)

	rootCmd := &cobra.Command{
		Use: "dukkha <tool-kind> {tool-name} <task-kind> <task-name>",
		Example: `dukkha docker build my-image
dukkha docker non-default-tool build my-image`,

		SilenceErrors: true,
		SilenceUsage:  true,

		Args: cobra.RangeArgs(3, 4),

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableNoDescFlag:   true,
			DisableDescriptions: true,
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch {
			case strings.HasPrefix(cmd.Use, "version"),
				strings.HasPrefix(cmd.Use, "completion"):
				return nil
			}

			// set global env per docs/environment-variables
			err := populateGlobalEnv(*appCtx)
			if err != nil {
				return fmt.Errorf("failed to set essential global env: %w", err)
			}

			// setup global logger for debugging
			err = log.SetDefaultLogger(log.ConfigSet{*cli_logConfig})
			if err != nil {
				return err
			}

			// read all configration files
			err = readConfigRecursively(
				configPaths,
				cmd.PersistentFlags().Changed("config"),
				config,
			)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			err = resolveConfig(*appCtx, renderingMgr, config, &allShells, &allTools, &toolSpecificTasks)
			if err != nil {
				return fmt.Errorf("failed to resolve config: %w", err)
			}

			mf := parseMatrixFilter(cli_matrixFilter)
			ctx := constant.WithWorkerCount(*appCtx, cli_workerCount)
			if len(mf) != 0 {
				ctx = constant.WithMatrixFilter(ctx, mf)
			}

			ctx = constant.WithFailFast(ctx, cli_failFast)
			appCtx = &ctx

			log.Log.D("application configured", log.Any("config", config))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(*appCtx, args, allShells, allTools)
		},
	}

	globalFlags := rootCmd.PersistentFlags()

	globalFlags.StringSliceVarP(
		&configPaths,
		"config", "c",
		[]string{".dukkha.yaml", ".dukkha"},
		"path to your config files and directories, only files with .yaml extension are parsed",
	)

	globalFlags.IntVarP(
		&cli_workerCount,
		"workers", "j",
		1,
		"set parallel worker count",
	)

	globalFlags.StringSliceVarP(
		&cli_matrixFilter,
		"matrix", "m",
		nil,
		"set matrix filter, format: <name>=<value>",
	)

	globalFlags.BoolVar(
		&cli_failFast,
		"fail-fast",
		true,
		"cancel other task execution when one errored",
	)

	// logging for debug purpose
	globalFlags.StringVarP(
		&cli_logConfig.Level,
		"log.level", "v",
		"info",
		"log level, one of [verbose, debug, info, error, silent]",
	)

	globalFlags.StringVar(
		&cli_logConfig.Format,
		"log.format",
		"console",
		"log output format, one of [console, json]",
	)

	globalFlags.StringVar(
		&cli_logConfig.File,
		"log.file",
		"stderr",
		"file path to write log output, including `stdout` and `stderr`",
	)

	setupCompletion(appCtx, rootCmd, renderingMgr.Render, &allTools, &toolSpecificTasks)

	return rootCmd
}

func run(
	appCtx context.Context,
	args []string,
	allShells map[tools.ToolKey]*tools.BaseTool,
	allTools map[tools.ToolKey]tools.Tool,
) error {
	type taskKey struct {
		taskKind string
		taskName string
	}

	var (
		targetTool tools.ToolKey
		targetTask taskKey
	)
	switch n := len(args); n {
	case 3:
		targetTool.ToolKind, targetTool.ToolName = args[0], ""
		targetTask.taskKind, targetTask.taskName = args[1], args[2]
	case 4:
		targetTool.ToolKind, targetTool.ToolName = args[0], args[1]
		targetTask.taskKind, targetTask.taskName = args[2], args[3]
	default:
		return fmt.Errorf("expecting 3 or 4 args, got %d", n)
	}

	tool, ok := allTools[targetTool]
	if !ok {
		return fmt.Errorf("tool %q with name %q not found", targetTool.ToolKind, targetTool.ToolName)
	}

	return tool.Run(appCtx, allTools, allShells, targetTask.taskKind, targetTask.taskName)
}
