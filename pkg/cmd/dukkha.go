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
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

func NewRootCmd() *cobra.Command {
	// cli options
	var (
		logConfig = new(log.Config)

		workerCount  = int(1)
		failFast     = false
		forceColor   = false
		matrixFilter []string

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
		Use: "dukkha <tool-kind> <tool-name> <task-kind> <task-name>",
		Example: `dukkha buildah local build my-image
dukkha buildah in-docker build my-image`,

		SilenceErrors: true,
		SilenceUsage:  true,

		Args: cobra.ExactArgs(4),

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

			// setup global logger for debugging
			err := log.SetDefaultLogger(log.ConfigSet{*logConfig})
			if err != nil {
				return fmt.Errorf("failed to initialize logger: %w", err)
			}

			logger := log.Log.WithName("pre-run")

			// read all configration files
			err = readConfigRecursively(
				configPaths,
				!cmd.PersistentFlags().Changed("config"),
				config,
			)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			logger.V("initializing dukkha", log.Any("raw_config", config))
			_appCtx, err := initDukkha(appBaseCtx, config, failFast, forceColor, workerCount)
			if err != nil {
				return fmt.Errorf("failed to initialize dukkha: %w", err)
			}
			_appCtx.SetMatrixFilter(parseMatrixFilter(matrixFilter))

			appCtx = _appCtx

			logger.D("dukkha initialized", log.Any("init_config", config))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(appCtx, args)
		},
	}

	globalFlags := rootCmd.PersistentFlags()
	globalFlags.StringSliceVarP(
		&configPaths, "config", "c", []string{".dukkha.yaml", ".dukkha"},
		"path to your config files and directories, only files with .yaml extension are parsed",
	)

	globalFlags.IntVarP(&workerCount, "workers", "j", 1, "set parallel worker count")
	globalFlags.BoolVar(&failFast, "fail-fast", true, "cancel all task execution after one errored")
	globalFlags.BoolVar(&forceColor, "force-color", false, "force color output even when not given a tty")
	globalFlags.StringSliceVarP(&matrixFilter, "matrix", "m", nil, "set matrix filter, format: -m <name>=<value>")

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

	setupTaskCompletion(&appCtx, rootCmd)

	err := setupMatrixCompletion(&appCtx, rootCmd, "matrix")
	if err != nil {
		panic(fmt.Errorf("failed to setup matrix flag completion: %w", err))
	}

	return rootCmd
}

// initialize dukkha runtime, create a context for any task execution
func initDukkha(
	appBaseCtx context.Context,
	config *conf.Config,
	failFast bool,
	forceColor bool,
	workers int,
) (dukkha.Context, error) {
	logger := log.Log.WithName("init")

	// create global env per docs/environment-variables
	logger.V("creating global environment variables")
	globalEnv := createGlobalEnv(appBaseCtx)

	logger.D("resolving bootstrap config")
	err := config.Bootstrap.Resolve(&globalEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve bootstrap config: %w", err)
	}

	logger.D("ensuring cache directory exists",
		log.String("cache_dir", config.Bootstrap.CacheDir),
	)
	err = os.MkdirAll(config.Bootstrap.CacheDir, 0750)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("failed to ensure cache dir exists: %w", err)
	}

	globalEnv[constant.ENV_DUKKHA_CACHE_DIR] = config.Bootstrap.CacheDir

	// all global environment variables set
	// globalEnv MUST NOT be modified from now on

	appCtx := dukkha.NewConfigResolvingContext(
		appBaseCtx, globalEnv,
		config.Bootstrap.GetExecSpec,
		dukkha.GlobalInterfaceTypeHandler,
		failFast,
		term.IsTerminal(int(os.Stdout.Fd())) || forceColor,
		workers,
	)

	err = config.ResolveAfterBootstrap(appCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config: %w", err)
	}

	return appCtx, nil
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
