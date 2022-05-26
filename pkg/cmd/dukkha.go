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
	"io"
	"os"
	"os/signal"
	"strings"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/log"
	"arhat.dev/pkg/versionhelper"
	"github.com/spf13/cobra"

	"arhat.dev/dukkha/pkg/cmd/completion"
	"arhat.dev/dukkha/pkg/cmd/debug"
	"arhat.dev/dukkha/pkg/cmd/diff"
	"arhat.dev/dukkha/pkg/cmd/render"
	"arhat.dev/dukkha/pkg/cmd/run"
	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer/echo"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/tpl"
	"arhat.dev/dukkha/pkg/renderer/transform"
)

// NewRootCmd creates the dukkha command with all sub commands added
func NewRootCmd(prevCtx dukkha.Context) *cobra.Command {
	var (
		stdout io.Writer

		logConfig = new(log.Config)

		configPaths []string
		// merged config
		config = conf.NewConfig()

		appCtx                     = prevCtx
		appBaseCtx context.Context = prevCtx
	)

	if appBaseCtx == nil {
		var cancel context.CancelFunc

		stdout = os.Stdout
		appBaseCtx, cancel = context.WithCancel(context.Background())

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		go func() {
			for range sigCh {
				cancel()
			}
		}()
	} else {
		stdout = appCtx.Stdout()
	}

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

		// PersistentPreRunE resolves all config for dukkha
		//
		// NOTE: this fucntion is used in pkg/templateutils/dukkhaNS.Self
		//       make sure we only have config loading happening here
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
				return fmt.Errorf("initializing logger: %w", err)
			}

			logger := log.Log.WithName("pre-run")

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("check working dir: %w", err)
			}

			_appCtx := dukkha.NewConfigResolvingContext(
				appBaseCtx, dukkha.GlobalInterfaceTypeHandler,
				createGlobalEnv(appBaseCtx, cwd),
			)
			_appCtx.AddListEnv(os.Environ()...)

			// add essential renderers for bootstraping
			{
				logger.V("creating essential renderers")

				_appCtx.AddRenderer(echo.DefaultName, echo.NewDefault(echo.DefaultName))
				_appCtx.AddRenderer(env.DefaultName, env.NewDefault(env.DefaultName))
				_appCtx.AddRenderer(shell.DefaultName, shell.NewDefault(shell.DefaultName))
				_appCtx.AddRenderer(tpl.DefaultName, tpl.NewDefault(tpl.DefaultName))
				_appCtx.AddRenderer(file.DefaultName, file.NewDefault(file.DefaultName))
				_appCtx.AddRenderer(transform.DefaultName, transform.NewDefault(transform.DefaultName))

				essentialRenderers := _appCtx.AllRenderers()
				for name, r := range essentialRenderers {
					// using default config, no need to resolve fields
					err = r.Init(_appCtx.RendererCacheFS(name))
					if err != nil {
						return fmt.Errorf("initialize essential renderer %q: %w", name, err)
					}
				}
			}

			// read all configration files
			visitedPaths := make(map[string]struct{})
			err = conf.Read(
				_appCtx,
				fshelper.NewOSFS(false, os.Getwd),
				configPaths,
				!cmd.PersistentFlags().Changed("config"),
				&visitedPaths,
				config,
			)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			logger.V("initializing dukkha", log.Any("raw_config", config))

			// here we always have tasks resolved to make template function
			// `dukkha.Self` work under all circumstances (e.g. `dukkha.Self run` used in `dukkha render`)
			err = config.Resolve(_appCtx, true /* need tasks */)
			if err != nil {
				return err
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
		// version
		versionhelper.NewVersionCmd(stdout),
		// completion
		completion.NewCompletionCmd(&appCtx),
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
