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

	"arhat.dev/dukkha/pkg/cmd/render"
	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer/echo"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/http"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/template"
)

func NewRootCmd() *cobra.Command {
	// cli options
	var (
		logConfig = new(log.Config)

		workerCount  = int(1)
		failFast     = false
		forceColor   = false
		matrixFilter []string

		translateANSIStream = false
		retainANSIStyle     = false

		configPaths []string
		// merged config
		config = conf.NewConfig()

		pluginConfigPaths []string
		plugins           = conf.NewPluginConfig()

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

			stdoutIsPty := term.IsTerminal(int(os.Stdout.Fd()))

			translateANSIFlag := cmd.Flag("translate-ansi-stream")
			var actualTranslateANSIStream bool
			// only translate ansi stream when
			// 	- (automatic) flag --translate-ansi-stream not set and stdout is not a pty
			actualTranslateANSIStream = (translateANSIFlag == nil || !translateANSIFlag.Changed) && !stdoutIsPty
			// 	- (manual) flag --translate-ansi-stream set to true
			actualTranslateANSIStream = actualTranslateANSIStream || translateANSIStream

			actualRetainANSIStyle := actualTranslateANSIStream && retainANSIStyle

			// prepare context
			_appCtx := dukkha.NewConfigResolvingContext(
				appBaseCtx, dukkha.ContextOptions{
					InterfaceTypeHandler: dukkha.GlobalInterfaceTypeHandler,
					FailFast:             failFast,
					ColorOutput:          stdoutIsPty || forceColor,
					TranslateANSIStream:  actualTranslateANSIStream,
					RetainANSIStyle:      actualRetainANSIStyle,
					Workers:              workerCount,
					GlobalEnv:            createGlobalEnv(appBaseCtx),
				},
			)

			{
				_appCtx.AddListEnv(os.Environ()...)

				logger.V("creating essential renderers")

				// TODO: let user decide what renderers to use
				// 		 resolve renderers first?
				_appCtx.AddRenderer(echo.DefaultName, echo.NewDefault(""))
				_appCtx.AddRenderer(env.DefaultName, env.NewDefault(""))
				_appCtx.AddRenderer(shell.DefaultName, shell.NewDefault(""))
				_appCtx.AddRenderer(template.DefaultName, template.NewDefault(""))
				_appCtx.AddRenderer(file.DefaultName, file.NewDefault(""))
				_appCtx.AddRenderer(http.DefaultName, http.NewDefault(""))

				essentialRenderers := _appCtx.AllRenderers()
				logger.D("initializing essential renderers",
					log.Int("count", len(essentialRenderers)),
				)

				for name, r := range essentialRenderers {
					// using default config, no need to resolve fields

					err := r.Init(_appCtx)
					if err != nil {
						return fmt.Errorf("failed to initialize essential renderer %q: %w", name, err)
					}
				}
			}

			// read all plugins config files
			visitedPaths := make(map[string]struct{})
			err = readConfigRecursively(
				pluginConfigPaths,
				!cmd.PersistentFlags().Changed("plugins-config"),
				&visitedPaths,
				plugins,
				readAndMergePluginConfigFile,
			)
			if err != nil {
				return fmt.Errorf("failed to load plugins config: %w", err)
			}

			err = plugins.ResolveAndRegisterPlugins(_appCtx)
			if err != nil {
				return fmt.Errorf("failed to resolve and register plugins: %w", err)
			}

			// read all configration files
			visitedPaths = make(map[string]struct{})
			err = readConfigRecursively(
				configPaths,
				!cmd.PersistentFlags().Changed("config"),
				&visitedPaths,
				config,
				readAndMergeConfigFile,
			)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			logger.V("initializing dukkha", log.Any("raw_config", config))

			var resolveTasks bool
			switch {
			case strings.HasPrefix(cmd.Use, "render"):
				resolveTasks = false
			default:
				resolveTasks = true
			}

			err = config.Resolve(_appCtx, resolveTasks)
			if err != nil {
				return fmt.Errorf("failed to resolve config: %w", err)
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
		&pluginConfigPaths, "plugins-config", "p", []string{".dukkha-plugins.yaml"},
		"path to your plugin config files and directories, behave like --config, but for plugins",
	)

	globalFlags.StringSliceVarP(
		&configPaths, "config", "c", []string{".dukkha.yaml", ".dukkha"},
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

	flags := rootCmd.Flags()

	flags.IntVarP(&workerCount, "workers", "j", 1, "set parallel worker count")
	flags.BoolVar(&failFast, "fail-fast", true, "cancel all task execution after one errored")
	flags.BoolVar(&forceColor, "force-color", false, "force color output even when not given a tty")
	flags.StringSliceVarP(&matrixFilter, "matrix", "m", nil, "set matrix filter, format: -m <name>=<value>")
	flags.BoolVar(&translateANSIStream, "translate-ansi-stream", false,
		"when set to true, will translate ansi stream to plain text before write to stdout/stderr, "+
			"when set to false, do nothing to the ansi stream, "+
			"when not set, will behavior as set to true if stdout/stderr is not a pty environment",
	)
	flags.BoolVar(&retainANSIStyle, "retain-ansi-style", retainANSIStyle,
		"when set to true, will retain ansi style when write to stdout/stderr, only effective "+
			"when ansi stream is going to be translated",
	)

	setupTaskCompletion(&appCtx, rootCmd)

	err := setupMatrixCompletion(&appCtx, rootCmd, "matrix")
	if err != nil {
		panic(fmt.Errorf("failed to setup matrix flag completion: %w", err))
	}

	rootCmd.AddCommand(render.NewRenderCmd(&appCtx))

	return rootCmd
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
