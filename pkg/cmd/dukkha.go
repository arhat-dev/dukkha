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
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"arhat.dev/pkg/log"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/shell_file"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/renderer/template_file"
	"arhat.dev/dukkha/pkg/tools"
)

func NewRootCmd() *cobra.Command {
	var (
		appCtx      context.Context
		configPaths []string
		logConfig   = new(log.Config)
		config      = conf.NewConfig()

		renderingMgr = renderer.NewManager()
	)

	rootCmd := &cobra.Command{
		Use:           "dukkha",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "version" {
				return nil
			}

			config.Log = logConfig

			err := log.SetDefaultLogger(log.ConfigSet{*logConfig})
			if err != nil {
				return err
			}

			err = readConfig(
				configPaths,
				cmd.PersistentFlags().Changed("config"),
				config,
			)
			if err != nil {
				return fmt.Errorf("failed to read config: %w", err)
			}

			appCtx = context.Background()

			err = config.Bootstrap.Resolve()
			if err != nil {
				return fmt.Errorf("failed to resolve bootstrap config: %w", err)
			}

			// bootstrap config was resolved when unmarshaling
			if config.Bootstrap.Shell == "" {
				return fmt.Errorf("unable to get a shell name, please set bootstrap.shell manually")
			}

			// create a renderer manager with essential renderers
			err = multierr.Combine(err,
				renderingMgr.Add(
					&shell.Config{ExecFunc: config.Bootstrap.Exec},
					shell.DefaultName,
				),
				renderingMgr.Add(
					&shell_file.Config{ExecFunc: config.Bootstrap.Exec},
					shell_file.DefaultName,
				),
				renderingMgr.Add(&template.Config{}, template.DefaultName),
				renderingMgr.Add(&template_file.Config{}, template_file.DefaultName),
				renderingMgr.Add(&file.Config{}, file.DefaultName),
			)

			if err != nil {
				return fmt.Errorf("failed to create essential renderers: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(appCtx, renderingMgr, config)
		},
	}

	globalFlags := rootCmd.PersistentFlags()

	globalFlags.StringSliceVarP(&configPaths, "config", "c", []string{".dukkha", ".dukkha.yaml"}, "")

	// logging
	globalFlags.StringVarP(&logConfig.Level, "log.level", "v", "info",
		"log level, one of [verbose, debug, info, error, silent]")
	globalFlags.StringVar(&logConfig.Format, "log.format", "console",
		"log output format, one of [console, json]")
	globalFlags.StringVar(&logConfig.File, "log.file", "stderr",
		"log output to this file")

	return rootCmd
}

func run(appCtx context.Context, renderingMgr *renderer.Manager, config *conf.Config) error {
	logger := log.Log.WithName("app")

	// ensure all top-level config resolved using basic renderers
	logger.V("resolving top-level config")
	err := config.ResolveFields(field.WithRenderingValues(appCtx, nil), renderingMgr.Render, 1)
	if err != nil {
		return fmt.Errorf("failed to resolve config: %w", err)
	}

	logger.V("top-level config resolved", log.Any("resolved_config", config))

	// resolve all shells, add them as shell & shell_file renderers

	for i, v := range config.Shells {
		logger.V("resolving shell config",
			log.String("shell", v.ToolName()),
			log.Int("index", i),
		)

		// resolve all config
		err = v.ResolveFields(field.WithRenderingValues(appCtx, nil), renderingMgr.Render, -1)
		if err != nil {
			return fmt.Errorf("failed to resolve config for shell %q #%d", v.ToolName(), i)
		}

		if i == 0 {
			err = multierr.Combine(err,
				renderingMgr.Add(&shell.Config{ExecFunc: v.RenderingExec}, shell.DefaultName),
				renderingMgr.Add(&shell_file.Config{ExecFunc: v.RenderingExec}, shell_file.DefaultName),
			)
		}

		err = multierr.Combine(err,
			renderingMgr.Add(&shell.Config{ExecFunc: v.RenderingExec}, shell.DefaultName+":"+v.ToolName()),
			renderingMgr.Add(&shell_file.Config{ExecFunc: v.RenderingExec}, shell_file.DefaultName+":"+v.ToolName()),
		)

		if err != nil {
			return fmt.Errorf("failed to add shell renderer %q", v.ToolName())
		}
	}

	// gather tasks for tools
	toolSpecificTasks := make(map[string][]tools.Task)

	for _, tasks := range config.Tasks {
		if len(tasks) == 0 {
			continue
		}

		toolKind := tasks[0].ToolKind()

		toolSpecificTasks[toolKind] = append(
			toolSpecificTasks[toolKind], tasks...,
		)
	}

	for _, tools := range config.Tools {
		for i, t := range tools {
			logger := logger.WithFields(
				log.String("tool", t.ToolKind()),
				log.Int("index", i),
				log.String("name", t.ToolName()),
			)

			toolID := t.ToolKind() + "#" + strconv.FormatInt(int64(i), 10)
			if len(t.ToolName()) != 0 {
				toolID = t.ToolKind() + ":" + t.ToolName()
			}

			logger.V("resolving tool config")

			err = t.ResolveFields(field.WithRenderingValues(appCtx, nil), renderingMgr.Render, -1)
			if err != nil {
				return fmt.Errorf(
					"failed to resolve config for tool %q: %w",
					toolID, err,
				)
			}

			logger.V("initializing tool")
			err = t.Init(renderingMgr.Render)
			if err != nil {
				return fmt.Errorf(
					"failed to initialize tool %q: %w",
					toolID, err,
				)
			}

			logger.V("resolving tool tasks")

			err = t.ResolveTasks(toolSpecificTasks[t.ToolKind()])
			if err != nil {
				return fmt.Errorf(
					"failed to resolve tasks for tool %q: %w",
					toolID, err,
				)
			}
		}
	}

	logger.D("application configured", log.Any("config", config))

	return nil
}

func readConfig(configPaths []string, failOnFileNotFoundError bool, mergedConfig *conf.Config) error {
	readAndMergeConfigFile := func(path string) error {
		configBytes, err2 := os.ReadFile(path)
		if err2 != nil {
			return fmt.Errorf("failed to read config file %q: %w", path, err2)
		}

		current := conf.NewConfig()
		err2 = yaml.Unmarshal(configBytes, &current)
		if err2 != nil {
			return fmt.Errorf("failed to unmarshal config file %q: %w", path, err2)
		}

		log.Log.V("config unmarshaled", log.String("file", path), log.Any("config", current))

		mergedConfig.Merge(current)

		return err2
	}

	for _, path := range configPaths {
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				if !failOnFileNotFoundError {
					continue
				}
			}

			return err
		}

		if !info.IsDir() {
			err = readAndMergeConfigFile(path)
			if err != nil {
				return err
			}

			continue
		}

		err = fs.WalkDir(os.DirFS(path), ".", func(pathInDir string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			switch filepath.Ext(pathInDir) {
			case ".yml", ".yaml":
			default:
				return nil
			}

			return readAndMergeConfigFile(filepath.Join(path, pathInDir))
		})

		if err != nil {
			return err
		}
	}

	return nil
}
