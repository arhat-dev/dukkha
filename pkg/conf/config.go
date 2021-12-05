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

package conf

import (
	"fmt"

	"arhat.dev/pkg/log"
	"arhat.dev/pkg/rshelper"
	"arhat.dev/rs"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

func NewConfig() *Config {
	return rshelper.InitAll(&Config{}, &rs.Options{
		InterfaceTypeHandler: dukkha.GlobalInterfaceTypeHandler,
	}).(*Config)
}

type Config struct {
	rs.BaseField `yaml:"-"`

	// Global options only have limited rendering suffix support
	Global GlobalConfig `yaml:"global"`

	// Include other files using path relative to current file
	// only local path is supported
	//
	// With path glob pattern '*' and '**' support
	Include []*IncludeEntry `yaml:"include"`

	// Shells for command execution
	Shells []*tools.ShellTool `yaml:"shells"`

	// Renderers config options
	Renderers []*RendererGroup `yaml:"renderers"`

	// Tools config options for registered tools
	Tools Tools `yaml:"tools"`

	Tasks map[string][]dukkha.Task `yaml:",inline"`
}

func (c *Config) Merge(a *Config) error {
	err := c.BaseField.Inherit(&a.BaseField)
	if err != nil {
		return fmt.Errorf("inherit top level config: %w", err)
	}

	err = c.Global.Merge(&a.Global)
	if err != nil {
		return err
	}

	err = c.Tools.Merge(&a.Tools)
	if err != nil {
		return err
	}

	if len(a.Tasks) != 0 {
		if c.Tasks == nil {
			c.Tasks = make(map[string][]dukkha.Task)
		}

		for k := range a.Tasks {
			c.Tasks[k] = append(c.Tasks[k], a.Tasks[k]...)
		}
	}

	return nil
}

func (c *Config) resolveRenderers(appCtx dukkha.ConfigResolvingContext) error {
	logger := log.Log.WithName("config")

	logger.D("resolving to gain renderers overview")
	err := c.ResolveFields(appCtx, 2, "renderers")
	if err != nil {
		return fmt.Errorf("gain overview of renderers: %w", err)
	}

	logger.D("resolving user renderers", log.Int("count", len(c.Renderers)))
	for i, g := range c.Renderers {
		logger := logger.WithFields(log.Int("index", i))

		logger.D("resolving renderer group")

		// renderers in the same group should shall be resolved all at once
		// without knowning each other

		err = g.ResolveFields(appCtx, -1)
		if err != nil {
			return fmt.Errorf("resolving renderer group #%d: %w", i, err)
		}

		for name, r := range g.Renderers {
			err = r.ResolveFields(appCtx, -1)
			if err != nil {
				return fmt.Errorf("resolving renderer %q: %w", name, err)
			}

			logger.V("resoving renderer", log.String("name", name), log.String("alias", r.Alias()))

			err = r.Init(appCtx.RendererCacheFS(name))
			if err != nil {
				return fmt.Errorf("initializing renderer %q: %w", name, err)
			}

			appCtx.AddRenderer(name, g.Renderers[name])
			if len(r.Alias()) != 0 {
				appCtx.AddRenderer(r.Alias(), g.Renderers[name])
			}
		}
	}

	return nil
}

func (c *Config) resolveShells(appCtx dukkha.ConfigResolvingContext) error {
	logger := log.Log.WithName("config")

	err := c.ResolveFields(appCtx, 1, "shells")
	if err != nil {
		return fmt.Errorf("resolving overview of shells: %w", err)
	}
	logger.V("resolved shell config overview", log.Any("result", c.Shells))

	logger.D("resolving shells", log.Int("count", len(c.Shells)))
	for i, v := range c.Shells {
		logger := logger.WithFields(
			log.Any("shell", v.Name()),
			log.Int("index", i),
		)

		logger.D("resolving shell config fields")
		err := v.ResolveFields(appCtx, -1)
		if err != nil {
			return fmt.Errorf("resolving shell %q #%d config: %w", v.Name(), i, err)
		}

		err = v.InitBaseTool(string(v.Name()), appCtx.ToolCacheFS(v), v)
		if err != nil {
			return fmt.Errorf("initializing shell %q", v.Name())
		}

		logger.V("adding shell")
		appCtx.AddShell(string(v.Name()), c.Shells[i])
	}

	return nil
}

// Resolve resolves all top level dukkha config
// to gain an overview of all tools and tasks
func (c *Config) Resolve(appCtx dukkha.ConfigResolvingContext, needTasks bool) error {
	logger := log.Log.WithName("config")

	// step 1: resolve global config (except Values)
	{
		logger.D("resolving global config")
		err := c.ResolveFields(appCtx, 1, "global")
		if err != nil {
			return fmt.Errorf("get global config overview: %w", err)
		}

		err = c.Global.ResolveAllButValues(appCtx)
		if err != nil {
			return fmt.Errorf("resolve global config: %w", err)
		}

		logger.V("resolved global config", log.Any("result", c.Global))

		cacheDir := c.Global.CacheDir
		if len(cacheDir) == 0 {
			cacheDir = constant.DefaultCacheDir
		}

		cacheDir, err = appCtx.FS().Abs(cacheDir)
		if err != nil {
			return fmt.Errorf("get absolute path of cache dir: %w", err)
		}

		if len(c.Global.DefaultGitBranch) != 0 {
			appCtx.(di.DefaultGitBranchOverrider).OverrideDefaultGitBranch(c.Global.DefaultGitBranch)
		}

		appCtx.(di.CacheDirSetter).SetCacheDir(cacheDir)
	}

	// step 2: resolve global Values
	{
		logger.D("resolving global values")

		err := c.Global.ResolveFields(appCtx, -1, "values")
		if err != nil {
			return fmt.Errorf("resolving global values: %w", err)
		}

		logger.V("resolved global values", log.Any("values", c.Global.Values))

		logger.D("adding global values")
		values := c.Global.Values.NormalizedValue()
		if err != nil {
			return fmt.Errorf("normalizing global values: %w", err)
		}

		err = appCtx.AddValues(values)
		if err != nil {
			return fmt.Errorf("adding global values: %w", err)
		}
	}

	// save some time if the command is not interacting with tasks
	if !needTasks {
		return nil
	}

	// step 3: resolve tools and tasks
	logger.D("resolving tools overview")
	err := c.ResolveFields(appCtx, 2, "tools")
	if err != nil {
		return fmt.Errorf("gain overview of tools: %w", err)
	}
	logger.V("resolved tools overview", log.Any("result", c.Tools))

	logger.D("resolving tasks overview")
	err = c.ResolveFields(appCtx, 1, "Tasks")
	if err != nil {
		return fmt.Errorf("gain overview of tasks: %w", err)
	}
	logger.V("resolved tasks overview", log.Any("result", c.Tasks))

	logger.V("groupping tasks by tool")
	for _, tasks := range c.Tasks {
		for _, tsk := range tasks {
			err = tsk.ResolveFields(appCtx, -1, "name")
			if err != nil {
				return fmt.Errorf("reoslving task name: %w", err)
			}

			// FIXME: task name is empty at this time
			err = tsk.Init(appCtx.TaskCacheFS(tsk))
			if err != nil {
				return fmt.Errorf("task init: %w", err)
			}
		}

		if len(tasks) == 0 {
			continue
		}

		appCtx.AddToolSpecificTasks(
			tasks[0].ToolKind(),
			tasks[0].ToolName(),
			tasks,
		)
	}

	logger.V("resolving tools", log.Int("count", len(c.Tools.Tools)))
	for tk, toolSet := range c.Tools.Tools {
		toolKind := dukkha.ToolKind(tk)

		visited := make(map[dukkha.ToolName]struct{})

		for i, t := range toolSet {
			err = t.ResolveFields(appCtx, -1, "name")
			if err != nil {
				return fmt.Errorf("resolve tool name: %w", err)
			}

			// do not allow empty name
			name := t.Name()
			if len(name) == 0 {
				return fmt.Errorf("invalid %q tool without name, index %d", toolKind, i)
			}

			// ensure tool names are unique
			if _, ok := visited[name]; ok {
				return fmt.Errorf("duplicate tool name %q of kind %q", t.Name(), toolKind)
			}

			visited[name] = struct{}{}

			key := dukkha.ToolKey{
				Kind: toolKind,
				Name: name,
			}

			logger := logger.WithFields(
				log.String("key", key.String()),
				log.Int("index", i),
			)

			logger.V("initializing tool")
			err = t.Init(appCtx.ToolCacheFS(t))
			if err != nil {
				return fmt.Errorf(
					"initializing tool %q: %w",
					key, err,
				)
			}

			// append tasks without tool name
			// they are meant for all tools in the same kind they belong to
			noToolNameTasks, _ := appCtx.GetToolSpecificTasks(
				dukkha.ToolKey{Kind: toolKind, Name: ""},
			)
			appCtx.AddToolSpecificTasks(
				toolKind, name, noToolNameTasks,
			)

			tasks, _ := appCtx.GetToolSpecificTasks(key)

			if logger.Enabled(log.LevelVerbose) {
				logger.D("resolving tool tasks", log.Any("tasks", tasks))
			} else {
				logger.D("resolving tool tasks")
			}

			err = t.AddTasks(tasks)
			if err != nil {
				return fmt.Errorf(
					"admitting tasks to tool %q: %w",
					key, err,
				)
			}

			appCtx.AddTool(key, c.Tools.Tools[string(toolKind)][i])
		}
	}

	return nil
}
