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

func resolveConfig(
	appCtx context.Context,
	renderingMgr *renderer.Manager,
	config *conf.Config,
	allShells *map[tools.ToolKey]*tools.BaseTool,
	allTools *map[tools.ToolKey]tools.Tool,
	toolSpecificTasks *map[tools.ToolKey][]tools.Task,
) error {
	logger := log.Log.WithName("config")

	err := config.Bootstrap.Resolve()
	if err != nil {
		return fmt.Errorf("failed to resolve bootstrap config: %w", err)
	}

	// bootstrap config was resolved when unmarshaling
	if len(config.Bootstrap.ScriptCmd) == 0 {
		return fmt.Errorf("bootstrap script_cmd not set")
	}

	for _, entry := range config.Bootstrap.Env {
		parts := strings.SplitN(entry, "=", 2)
		name, value := parts[0], ""
		if len(parts) == 2 {
			value = parts[1]
		}

		_ = os.Setenv(name, value)
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

	// ensure all top-level config resolved using basic renderers
	logger.V("resolving top-level config")
	err = config.ResolveFields(field.WithRenderingValues(appCtx, nil), renderingMgr.Render, 1)
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

		err = v.InitBaseTool(config.Bootstrap.CacheDir, v.ToolName(), nil, nil)
		if err != nil {
			return fmt.Errorf("failed to initialize shell %q", v.ToolName())
		}

		if i == 0 {
			(*allShells)[tools.ToolKey{ToolKind: "shell", ToolName: ""}] = config.Shells[i]
			err = multierr.Combine(err,
				renderingMgr.Add(&shell.Config{GetExecSpec: v.GetExecSpec}, shell.DefaultName),
				renderingMgr.Add(&shell_file.Config{GetExecSpec: v.GetExecSpec}, shell_file.DefaultName),
			)
		}

		(*allShells)[tools.ToolKey{ToolKind: "shell", ToolName: v.Name}] = config.Shells[i]
		err = multierr.Combine(err,
			renderingMgr.Add(&shell.Config{GetExecSpec: v.GetExecSpec}, shell.DefaultName+":"+v.ToolName()),
			renderingMgr.Add(&shell_file.Config{GetExecSpec: v.GetExecSpec}, shell_file.DefaultName+":"+v.ToolName()),
		)

		if err != nil {
			return fmt.Errorf("failed to add shell renderer %q", v.ToolName())
		}
	}

	for _, tasks := range config.Tasks {
		if len(tasks) == 0 {
			continue
		}

		key := tools.ToolKey{
			ToolKind: tasks[0].ToolKind(),
			ToolName: tasks[0].ToolName(),
		}

		(*toolSpecificTasks)[key] = append(
			(*toolSpecificTasks)[key], tasks...,
		)
	}

	for _, toolSet := range config.Tools {
		for i, t := range toolSet {
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

			fullToolKey := tools.ToolKey{
				ToolKind: t.ToolKind(),
				ToolName: t.ToolName(),
			}

			err = t.ResolveTasks((*toolSpecificTasks)[fullToolKey])
			if err != nil {
				return fmt.Errorf(
					"failed to resolve tasks for tool %q: %w",
					toolID, err,
				)
			}

			(*allTools)[fullToolKey] = t

			if i != 0 {
				continue
			}

			// setup default tasks
			if len(fullToolKey.ToolName) != 0 {
				// is default tool for this kind but using name before
				defaultToolKey := tools.ToolKey{
					ToolKind: t.ToolKind(),
					ToolName: "",
				}

				(*allTools)[defaultToolKey] = t

				tasksWithDefaultTool := (*toolSpecificTasks)[defaultToolKey]
				err = t.ResolveTasks(tasksWithDefaultTool)
				if err != nil {
					return fmt.Errorf(
						"failed to resolve tasks for default tool %q: %w",
						toolID, err,
					)
				}

				(*toolSpecificTasks)[fullToolKey] = tasksWithDefaultTool
			}
		}
	}

	return nil
}
