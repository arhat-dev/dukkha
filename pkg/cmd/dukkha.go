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
	"strings"

	"arhat.dev/pkg/log"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/constant"
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
		appCtx      = context.Background()
		configPaths []string
		logConfig   = new(log.Config)
		config      = conf.NewConfig()
		workerCount int

		matrixFilter = make(map[string]string)
		renderingMgr = renderer.NewManager()
	)

	rootCmd := &cobra.Command{
		Use: "dukkha <tool-kind> {tool-name} <task-kind> <task-name>",
		Example: `dukkha docker build my-image
dukkha docker non-default-tool build my-image`,

		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.RangeArgs(3, 4),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "version" {
				return nil
			}

			populateGlobalEnv(appCtx)

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
			if len(config.Bootstrap.ScriptCmd) == 0 {
				return fmt.Errorf("bootstrap script_cmd not set")
			}

			// create a renderer manager with essential renderers
			err = multierr.Combine(err,
				renderingMgr.Add(
					&shell_file.Config{GetExecSpec: config.Bootstrap.GetExecSpec},
					shell_file.DefaultName,
				),
				renderingMgr.Add(
					&shell.Config{GetExecSpec: config.Bootstrap.GetExecSpec},
					shell.DefaultName,
				),
				renderingMgr.Add(&template.Config{}, template.DefaultName),
				renderingMgr.Add(&template_file.Config{}, template_file.DefaultName),
				renderingMgr.Add(&file.Config{}, file.DefaultName),
			)

			if err != nil {
				return fmt.Errorf("failed to create essential renderers: %w", err)
			}

			mf := make(map[string][]string)
			for k, v := range matrixFilter {
				mf[k] = strings.Split(v, ",")
			}

			appCtx = constant.WithWorkerCount(appCtx, workerCount)
			if len(mf) != 0 {
				appCtx = constant.WithMatrixFilter(appCtx, mf)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(appCtx, renderingMgr, config, args)
		},
	}

	globalFlags := rootCmd.PersistentFlags()

	globalFlags.StringSliceVarP(
		&configPaths, "config", "c",
		[]string{".dukkha", ".dukkha.yaml"},
		"path to your config files",
	)

	globalFlags.IntVarP(&workerCount, "workers", "j", 1, "set parallel worker count")

	globalFlags.StringToStringVarP(&matrixFilter, "matrix", "m", nil, "set matrix filter")

	// logging
	globalFlags.StringVarP(&logConfig.Level, "log.level", "v", "info",
		"log level, one of [verbose, debug, info, error, silent]")
	globalFlags.StringVar(&logConfig.Format, "log.format", "console",
		"log output format, one of [console, json]")
	globalFlags.StringVar(&logConfig.File, "log.file", "stderr",
		"log output to this file")

	return rootCmd
}

func run(
	appCtx context.Context,
	renderingMgr *renderer.Manager,
	config *conf.Config,
	args []string,
) error {
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
				renderingMgr.Add(&shell.Config{GetExecSpec: v.GetExecSpec}, shell.DefaultName),
				renderingMgr.Add(&shell_file.Config{GetExecSpec: v.GetExecSpec}, shell_file.DefaultName),
			)
		}

		err = multierr.Combine(err,
			renderingMgr.Add(&shell.Config{GetExecSpec: v.GetExecSpec}, shell.DefaultName+":"+v.ToolName()),
			renderingMgr.Add(&shell_file.Config{GetExecSpec: v.GetExecSpec}, shell_file.DefaultName+":"+v.ToolName()),
		)

		if err != nil {
			return fmt.Errorf("failed to add shell renderer %q", v.ToolName())
		}
	}

	// gather tasks for tools
	type toolKey struct {
		toolKind string
		toolName string
	}

	toolSpecificTasks := make(map[toolKey][]tools.Task)

	// Always initialize all tasks in case task dependencies

	for _, tasks := range config.Tasks {
		if len(tasks) == 0 {
			continue
		}

		key := toolKey{
			toolKind: tasks[0].ToolKind(),
			toolName: tasks[0].ToolName(),
		}

		toolSpecificTasks[key] = append(
			toolSpecificTasks[key], tasks...,
		)
	}

	allTools := make(map[toolKey]tools.Tool)

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
			err = t.Init(
				config.Bootstrap.CacheDir,
				renderingMgr.Render,
				config.Bootstrap.GetExecSpec,
			)
			if err != nil {
				return fmt.Errorf(
					"failed to initialize tool %q: %w",
					toolID, err,
				)
			}

			logger.V("resolving tool tasks")

			key := toolKey{
				toolKind: t.ToolKind(),
				toolName: t.ToolName(),
			}

			err = t.ResolveTasks(toolSpecificTasks[key])
			if err != nil {
				return fmt.Errorf(
					"failed to resolve tasks for tool %q: %w",
					toolID, err,
				)
			}

			allTools[key] = t

			if i == 0 && len(key.toolName) != 0 {
				// is default tool for this kind but using name before
				key = toolKey{
					toolKind: t.ToolKind(),
					toolName: "",
				}

				allTools[key] = t

				err = t.ResolveTasks(toolSpecificTasks[key])
				if err != nil {
					return fmt.Errorf(
						"failed to resolve tasks for default tool %q: %w",
						toolID, err,
					)
				}
			}
		}
	}

	logger.D("application configured", log.Any("config", config))

	type taskKey struct {
		taskKind string
		taskName string
	}

	var (
		targetTool toolKey
		targetTask taskKey
	)
	switch n := len(args); n {
	case 3:
		targetTool.toolKind, targetTool.toolName = args[0], ""
		targetTask.taskKind, targetTask.taskName = args[1], args[2]
	case 4:
		targetTool.toolKind, targetTool.toolName = args[0], args[1]
		targetTask.taskKind, targetTask.taskName = args[2], args[3]
	default:
		return fmt.Errorf("expecting 3 or 4 args, got %d", n)
	}

	tool, ok := allTools[targetTool]
	if !ok {
		return fmt.Errorf("tool %q with name %q not found", targetTool.toolKind, targetTool.toolName)
	}

	return tool.Run(appCtx, targetTask.taskKind, targetTask.taskName)
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
			case ".yaml":
				// leave .yml for customization
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
