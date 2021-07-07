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
	"go.uber.org/multierr"

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
	return field.Init(&Config{}).(*Config)
}

type Config struct {
	field.BaseField

	// no rendering suffix support
	Bootstrap BootstrapConfig `yaml:"bootstrap"`

	// Shells for rendering and command execution
	Shells []*tools.BaseTool `yaml:"shells"`

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

	logger.D("creating essential renderers")
	err := multierr.Combine(
		appCtx.AddRenderer(
			shell_file.New(appCtx.GetBootstrapExecSpec),
			shell_file.DefaultName,
		),
		appCtx.AddRenderer(
			shell.New(appCtx.GetBootstrapExecSpec),
			shell.DefaultName,
		),
		appCtx.AddRenderer(
			env.New(appCtx.GetBootstrapExecSpec),
			env.DefaultName,
		),
		appCtx.AddRenderer(template.New(), template.DefaultName),
		appCtx.AddRenderer(template_file.New(), template_file.DefaultName),
		appCtx.AddRenderer(file.New(), file.DefaultName),
	)
	if err != nil {
		return fmt.Errorf("failed to create essential renderers: %w", err)
	}

	logger.D("resolving top level config")
	err = c.ResolveFields(appCtx, 1, "")
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

		err = v.InitBaseTool(string(v.Name()), appCtx.CacheDir())
		if err != nil {
			return fmt.Errorf("failed to initialize shell %q", v.Name())
		}

		if i == 0 {
			logger.V("adding default shell")

			appCtx.AddShell("", c.Shells[i])

			err = multierr.Combine(err,
				appCtx.AddRenderer(
					shell.New(c.Shells[i].GetExecSpec),
					shell.DefaultName,
				),
				appCtx.AddRenderer(
					shell_file.New(c.Shells[i].GetExecSpec),
					shell_file.DefaultName,
				),
				appCtx.AddRenderer(
					env.New(c.Shells[i].GetExecSpec),
					env.DefaultName,
				),
			)
			if err != nil {
				return fmt.Errorf("failed to add default shell renderer %q", v.Name())
			}
		}

		logger.V("adding shell")

		appCtx.AddShell(string(v.Name()), c.Shells[i])
		err = multierr.Combine(err,
			appCtx.AddRenderer(
				shell.New(c.Shells[i].GetExecSpec),
				shell.DefaultName+":"+string(v.Name()),
			),
			appCtx.AddRenderer(
				shell_file.New(c.Shells[i].GetExecSpec),
				shell_file.DefaultName+":"+string(v.Name()),
			),
			appCtx.AddRenderer(
				env.New(c.Shells[i].GetExecSpec),
				env.DefaultName+":"+string(v.Name()),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to add shell renderer %q", v.Name())
		}
	}

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
			err = t.Init(appCtx.CacheDir())
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
