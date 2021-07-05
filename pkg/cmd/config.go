package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"arhat.dev/pkg/log"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/conf"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/shell_file"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/renderer/template_file"
	"arhat.dev/dukkha/pkg/tools"
)

func readConfigRecursively(
	configPaths []string,
	failOnFileNotFoundError bool,
	mergedConfig *conf.Config,
) error {
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

// resolveConfig step by step
//
// 1. resolve bootstrap config
// 2. create a rendering manager with all essential renderers
// 3. resolve shells config, add them as shell renderer
// 4. resolve tools and their tasks
func resolveConfig(
	appCtx context.Context,
	renderingMgr *renderer.Manager,
	config *conf.Config,
	allShells *map[tools.ToolKey]*tools.BaseTool,
	allTools *map[tools.ToolKey]tools.Tool,
	toolSpecificTasks *map[tools.ToolKey][]tools.Task,
) error {
	logger := log.Log.WithName("config")

	logger.V("resolving bootstrap config")
	err := config.Bootstrap.Resolve()
	if err != nil {
		return fmt.Errorf("failed to resolve bootstrap config: %w", err)
	}

	logger.D("ensuring cache dir exists", log.String("cache_dir", config.Bootstrap.CacheDir))
	err = os.MkdirAll(config.Bootstrap.CacheDir, 0750)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to ensure cache dir exists: %w", err)
	}

	logger.D("creating essential renderers")
	err = multierr.Combine(err,
		renderingMgr.Add(
			&shell_file.Config{GetExecSpec: config.Bootstrap.GetExecSpec},
			shell_file.DefaultName,
		),
		renderingMgr.Add(
			&shell.Config{GetExecSpec: config.Bootstrap.GetExecSpec},
			shell.DefaultName,
		),
		renderingMgr.Add(
			&env.Config{GetExecSpec: config.Bootstrap.GetExecSpec},
			env.DefaultName,
		),
		renderingMgr.Add(&template.Config{}, template.DefaultName),
		renderingMgr.Add(&template_file.Config{}, template_file.DefaultName),
		renderingMgr.Add(&file.Config{}, file.DefaultName),
	)
	if err != nil {
		return fmt.Errorf("failed to create essential renderers: %w", err)
	}

	// no need to pass config.Bootstrap.Env as extraEnv
	// they are already set with os.Setenv
	bootstrapCtx := field.WithRenderingValues(appCtx, os.Environ())

	logger.D("resolving top level config")
	err = config.ResolveFields(bootstrapCtx, renderingMgr.Render, 1, false)
	if err != nil {
		return fmt.Errorf("failed to resolve config: %w", err)
	}
	logger.V("resolved top level config", log.Any("result", config))

	logger.D("resolving shells", log.Int("count", len(config.Shells)))
	for i, v := range config.Shells {
		logger := logger.WithFields(
			log.String("shell", v.ToolName()),
			log.Int("index", i),
		)

		logger.D("resolving shell config fields")
		err = v.ResolveFields(bootstrapCtx, renderingMgr.Render, -1, false)
		if err != nil {
			return fmt.Errorf("failed to resolve config for shell %q #%d", v.ToolName(), i)
		}

		err = v.InitBaseTool(config.Bootstrap.CacheDir, v.ToolName(), nil, nil)
		if err != nil {
			return fmt.Errorf("failed to initialize shell %q", v.ToolName())
		}

		if i == 0 {
			logger.V("adding default shell")

			(*allShells)[tools.ToolKey{ToolKind: "shell", ToolName: ""}] = config.Shells[i]
			err = multierr.Combine(err,
				renderingMgr.Add(
					&shell.Config{GetExecSpec: v.GetExecSpec},
					shell.DefaultName,
				),
				renderingMgr.Add(
					&shell_file.Config{GetExecSpec: v.GetExecSpec},
					shell_file.DefaultName,
				),
				renderingMgr.Add(
					&env.Config{GetExecSpec: config.Bootstrap.GetExecSpec},
					env.DefaultName,
				),
			)
			if err != nil {
				return fmt.Errorf("failed to add default shell renderer %q", v.ToolName())
			}
		}

		logger.V("adding shell")
		(*allShells)[tools.ToolKey{ToolKind: "shell", ToolName: v.ToolName()}] = config.Shells[i]
		err = multierr.Combine(err,
			renderingMgr.Add(
				&shell.Config{GetExecSpec: v.GetExecSpec},
				shell.DefaultName+":"+v.ToolName(),
			),
			renderingMgr.Add(
				&shell_file.Config{GetExecSpec: v.GetExecSpec},
				shell_file.DefaultName+":"+v.ToolName(),
			),
			renderingMgr.Add(
				&env.Config{GetExecSpec: config.Bootstrap.GetExecSpec},
				env.DefaultName+":"+v.ToolName(),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to add shell renderer %q", v.ToolName())
		}
	}

	logger.V("groupping tasks by tool")
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

	logger.V("resolving tools", log.Int("count", len(config.Tools)))
	for toolKind, toolSet := range config.Tools {
		visited := make(map[string]struct{})

		for i, t := range toolSet {
			// do not allow empty name
			if len(t.ToolName()) == 0 {
				return fmt.Errorf("invalid %q tool without name, index %d", toolKind, i)
			}

			// ensure tool names are unique
			if _, ok := visited[t.ToolName()]; ok {
				return fmt.Errorf("invalid duplicate %q tool name %q", toolKind, t.ToolName())
			}

			visited[t.ToolName()] = struct{}{}

			logger := logger.WithFields(
				log.String("kind", toolKind),
				log.String("name", t.ToolName()),
				log.Int("index", i),
			)

			toolID := toolKind + ":" + t.ToolName()

			logger.D("resolving tool config fields")
			err = t.ResolveFields(bootstrapCtx, renderingMgr.Render, -1, false)
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

			fullToolKey := tools.ToolKey{
				ToolKind: toolKind,
				ToolName: t.ToolName(),
			}

			defaultToolKey := tools.ToolKey{
				ToolKind: toolKind,
				ToolName: "",
			}

			// append tasks without tool name
			//
			// this is also used by shell completion
			(*toolSpecificTasks)[fullToolKey] = append(
				(*toolSpecificTasks)[fullToolKey],
				(*toolSpecificTasks)[defaultToolKey]...,
			)

			tasks := (*toolSpecificTasks)[fullToolKey]

			if logger.Enabled(log.LevelVerbose) {
				logger.D("resolving tool tasks", log.Any("tasks", tasks))
			} else {
				logger.D("resolving tool tasks")
			}

			err = t.ResolveTasks(tasks)
			if err != nil {
				return fmt.Errorf(
					"failed to resolve tasks for tool %q: %w",
					toolID, err,
				)
			}

			(*allTools)[fullToolKey] = config.Tools[toolKind][i]
			if i == 0 {
				// is first tool, set default tool key
				//
				// this is used by shell completion
				(*allTools)[defaultToolKey] = config.Tools[toolKind][i]
			}
		}
	}

	return nil
}
