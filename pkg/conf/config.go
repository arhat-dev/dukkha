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

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/shell_file"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/renderer/template_file"
	"arhat.dev/dukkha/pkg/tools"
)

func NewConfig() *Config {
	return field.Init(
		&Config{},
		dukkha.GlobalInterfaceTypeHandler,
	).(*Config)
}

type Config struct {
	field.BaseField

	// no rendering suffix support
	Bootstrap BootstrapConfig `yaml:"bootstrap"`

	// Shells for rendering and command execution
	Shells []*tools.BaseTool `yaml:"shells"`

	Renderers map[string]dukkha.Renderer `yaml:"renderers"`

	// Language or tool specific tools
	Tools map[string][]dukkha.Tool `yaml:"tools"`

	Tasks map[string][]dukkha.Task `dukkha:"other"`
}

func (c *Config) Merge(a *Config) {
	c.Bootstrap.Env = append(c.Bootstrap.Env, a.Bootstrap.Env...)
	if len(a.Bootstrap.CacheDir) != 0 {
		c.Bootstrap.CacheDir = a.Bootstrap.CacheDir
	}

	// once changed script cmd, replace the whole bootstrap config
	if len(a.Bootstrap.ScriptCmd) != 0 {
		c.Bootstrap = a.Bootstrap
	}

	c.Shells = append(c.Shells, a.Shells...)

	if len(a.Tools) != 0 {
		if c.Tools == nil {
			c.Tools = a.Tools
		} else {
			for k := range a.Tools {
				c.Tools[k] = append(c.Tools[k], a.Tools[k]...)
			}
		}
	}

	if len(a.Tasks) != 0 {
		if c.Tasks == nil {
			c.Tasks = a.Tasks
		} else {
			for k := range a.Tasks {
				c.Tasks[k] = append(c.Tasks[k], a.Tasks[k]...)
			}
		}
	}
}

// ResolveAfterBootstrap resolves all top level dukkha config
// to gain a overview of all tools and tasks
//
// 1. create a rendering manager with all essential renderers
//
// 2. resolve shells config using essential renderers,
// 	  add them as shell renderers
//
// 3. resolve tools and their tasks
func (c *Config) ResolveAfterBootstrap(appCtx dukkha.ConfigResolvingContext) error {
	logger := log.Log.WithName("config")

	logger.V("creating essential renderers")
	appCtx.AddRenderer(
		shell_file.DefaultName,
		shell_file.NewDefault(appCtx.GetBootstrapExecSpec),
	)
	appCtx.AddRenderer(
		shell.DefaultName,
		shell.NewDefault(appCtx.GetBootstrapExecSpec),
	)
	appCtx.AddRenderer(
		env.DefaultName,
		env.NewDefault(appCtx.GetBootstrapExecSpec),
	)
	appCtx.AddRenderer(template.DefaultName, template.NewDefault())
	appCtx.AddRenderer(template_file.DefaultName, template_file.NewDefault())
	appCtx.AddRenderer(file.DefaultName, file.NewDefault())

	logger.D("resolving top level config")
	err := c.ResolveFields(appCtx, 1, "")
	if err != nil {
		return fmt.Errorf("failed to resolve config: %w", err)
	}
	logger.V("resolved top level config", log.Any("result", c))

	logger.D("resolving shells", log.Int("count", len(c.Shells)))
	for i, v := range c.Shells {
		logger := logger.WithFields(
			log.Any("shell", v.Name()),
			log.Int("index", i),
		)

		logger.D("resolving shell config fields")
		err = v.ResolveFields(appCtx, -1, "")
		if err != nil {
			return fmt.Errorf("failed to resolve config for shell %q #%d", v.Name(), i)
		}

		err = v.InitBaseTool("shell", string(v.Name()), appCtx.CacheDir())
		if err != nil {
			return fmt.Errorf("failed to initialize shell %q", v.Name())
		}

		if i == 0 {
			logger.V("adding default shell")
			appCtx.AddShell("", c.Shells[i])
		}

		logger.V("adding shell")
		appCtx.AddShell(string(v.Name()), c.Shells[i])
	}

	essentialRenderers := appCtx.AllRenderers()
	logger.D("initializing essential renderers",
		log.Int("count", len(essentialRenderers)),
	)
	for name, r := range essentialRenderers {
		// using default config, no need to resolve fields

		err = r.Init(appCtx)
		if err != nil {
			return fmt.Errorf("failed to initialize essential renderer %q: %w", name, err)
		}
	}

	logger.D("resolving user renderers", log.Int("count", len(c.Renderers)))
	for name, r := range c.Renderers {
		logger := logger.WithFields(
			log.Any("renderer", name),
		)

		logger.D("resolving renderer config fields")
		err = r.ResolveFields(appCtx, -1, "")
		if err != nil {
			return fmt.Errorf("failed to resolve config for renderer %q: %w", name, err)
		}
	}

	logger.D("initializing user renderers", log.Int("count", len(c.Renderers)))
	for name, r := range c.Renderers {
		err = r.Init(appCtx)
		if err != nil {
			return fmt.Errorf("failed to initialize renderer %q: %w", name, err)
		}
	}
	logger.D("resolved all renderers", log.Int("count", len(appCtx.AllRenderers())))

	logger.V("groupping tasks by tool")
	for _, tasks := range c.Tasks {
		if len(tasks) == 0 {
			continue
		}

		appCtx.AddToolSpecificTasks(
			tasks[0].ToolKind(),
			tasks[0].ToolName(),
			tasks,
		)
	}

	logger.V("resolving tools", log.Int("count", len(c.Tools)))
	for tk, toolSet := range c.Tools {
		toolKind := dukkha.ToolKind(tk)

		visited := make(map[dukkha.ToolName]struct{})

		for i, t := range toolSet {
			// do not allow empty name
			if len(t.Name()) == 0 {
				return fmt.Errorf("invalid %q tool without name, index %d", toolKind, i)
			}

			// ensure tool names are unique
			if _, ok := visited[t.Name()]; ok {
				return fmt.Errorf("invalid duplicate %q tool name %q", toolKind, t.Name())
			}

			visited[t.Name()] = struct{}{}

			logger := logger.WithFields(
				log.Any("kind", toolKind),
				log.Any("name", t.Name()),
				log.Int("index", i),
			)

			toolID := string(toolKind) + ":" + string(t.Name())

			logger.D("resolving tool config fields")
			err = t.ResolveFields(appCtx, -1, "")
			if err != nil {
				return fmt.Errorf(
					"failed to resolve config for tool %q: %w",
					toolID, err,
				)
			}

			logger.V("initializing tool")
			err = t.Init(toolKind, appCtx.CacheDir())
			if err != nil {
				return fmt.Errorf(
					"failed to initialize tool %q: %w",
					toolID, err,
				)
			}

			// append tasks without tool name
			// they are meant for all tools in the same kind they belong to
			noToolNameTasks, _ := appCtx.GetToolSpecificTasks(
				toolKind, "",
			)
			appCtx.AddToolSpecificTasks(
				toolKind, t.Name(),
				noToolNameTasks,
			)

			tasks, _ := appCtx.GetToolSpecificTasks(toolKind, t.Name())

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

			appCtx.AddTool(toolKind, t.Name(), c.Tools[string(toolKind)][i])
		}
	}

	return nil
}
