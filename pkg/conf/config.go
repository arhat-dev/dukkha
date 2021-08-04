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
	"os"
	"path/filepath"

	"arhat.dev/pkg/log"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer/echo"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/tools"
)

func NewConfig() *Config {
	cfg := field.Init(&Config{}, dukkha.GlobalInterfaceTypeHandler).(*Config)
	_ = field.Init(&cfg.Global, dukkha.GlobalInterfaceTypeHandler)
	return cfg
}

type GlobalConfig struct {
	field.BaseField

	// CacheDir to store script file and temporary task execution data
	CacheDir string `yaml:"cache_dir"`

	DefaultGitBranch string `yaml:"default_git_branch"`

	// Env
	Env dukkha.Env `yaml:"env"`
}

func (g *GlobalConfig) Resolve(rc dukkha.ConfigResolvingContext) error {
	err := dukkha.ResolveEnv(g, rc, "Env")
	if err != nil {
		return fmt.Errorf("failed to resolve global env: %w", err)
	}

	err = g.ResolveFields(rc, -1, "CacheDir")
	if err != nil {
		return fmt.Errorf("failed to resolve cache dir: %w", err)
	}

	err = g.ResolveFields(rc, -1, "DefaultGitBranch")
	if err != nil {
		return fmt.Errorf("failed to resolve default git branch: %w", err)
	}

	return nil
}

type Config struct {
	field.BaseField

	// Global has no rendering suffix support
	Global GlobalConfig `yaml:"global"`

	// Include other files using path relative to this config
	// also no rendering suffix support
	Include []string `yaml:"include"`

	// Shells for rendering and command execution
	//
	// this option is host specific and do not support
	// renderers like `http`
	Shells []*tools.BaseToolWithInit `yaml:"shells"`

	Renderers map[string]dukkha.Renderer `yaml:"renderers"`

	// Language or tool specific tools
	Tools map[string][]dukkha.Tool `yaml:"tools"`

	Tasks map[string][]dukkha.Task `dukkha:"other"`
}

func (c *Config) Merge(a *Config) {
	err := c.BaseField.Inherit(&a.BaseField)
	if err != nil {
		panic(fmt.Errorf("failed to inherit other top level base field: %w", err))
	}

	c.Global.Env = append(c.Global.Env, a.Global.Env...)
	if len(a.Global.CacheDir) != 0 {
		c.Global.CacheDir = a.Global.CacheDir
	}

	if len(a.Global.DefaultGitBranch) != 0 {
		c.Global.DefaultGitBranch = a.Global.DefaultGitBranch
	}

	err = c.Global.BaseField.Inherit(&a.Global.BaseField)
	if err != nil {
		panic(fmt.Errorf("failed to inherit other global config: %w", err))
	}

	c.Shells = append(c.Shells, a.Shells...)

	if c.Renderers == nil {
		c.Renderers = a.Renderers
	} else {
		// TODO: handle duplicated renderers
		for k, v := range a.Renderers {
			c.Renderers[k] = v
		}
	}

	if len(a.Tools) != 0 {
		if c.Tools == nil {
			c.Tools = make(map[string][]dukkha.Tool)
		}

		for k := range a.Tools {
			c.Tools[k] = append(c.Tools[k], a.Tools[k]...)
		}
	}

	if len(a.Tasks) != 0 {
		if c.Tasks == nil {
			c.Tasks = make(map[string][]dukkha.Task)
		}

		for k := range a.Tasks {
			c.Tasks[k] = append(c.Tasks[k], a.Tasks[k]...)
		}
	}
}

// Resolve resolves all top level dukkha config
// to gain an overview of all tools and tasks
func (c *Config) Resolve(appCtx dukkha.ConfigResolvingContext) error {
	logger := log.Log.WithName("config")

	// step 1: create essential renderers
	{
		logger.V("creating essential renderers")

		// TODO: let user decide what renderers to use
		// 		 resolve renderers first?
		appCtx.AddRenderer(echo.DefaultName, echo.NewDefault())
		appCtx.AddRenderer(env.DefaultName, env.NewDefault())
		appCtx.AddRenderer(shell.DefaultName, shell.NewDefault())
		appCtx.AddRenderer(template.DefaultName, template.NewDefault())
		appCtx.AddRenderer(file.DefaultName, file.NewDefault())

		essentialRenderers := appCtx.AllRenderers()
		logger.D("initializing essential renderers",
			log.Int("count", len(essentialRenderers)),
		)

		for name, r := range essentialRenderers {
			// using default config, no need to resolve fields

			err := r.Init(appCtx)
			if err != nil {
				return fmt.Errorf("failed to initialize essential renderer %q: %w", name, err)
			}
		}
	}

	// step 2: resolve global config, ensure cache dir exists
	{
		logger.D("resolving global config")
		err := c.ResolveFields(appCtx, 1, "Global")
		if err != nil {
			return fmt.Errorf("failed to get global config overview: %w", err)
		}

		err = c.Global.Resolve(appCtx)
		if err != nil {
			return fmt.Errorf("failed to resolve global config: %w", err)
		}

		logger.V("resolved global config", log.Any("result", c.Global))

		cacheDir := c.Global.CacheDir
		if len(cacheDir) == 0 {
			cacheDir = constant.DefaultCacheDir
		}

		cacheDir, err = filepath.Abs(cacheDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path of cache dir: %w", err)
		}

		err = os.MkdirAll(cacheDir, 0750)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to ensure cache dir: %w", err)
		}

		if len(c.Global.DefaultGitBranch) != 0 {
			appCtx.OverrideDefaultGitBranch(c.Global.DefaultGitBranch)
		}

		appCtx.SetCacheDir(cacheDir)
	}

	// step 3: resolve renderers
	{
		logger.D("resolving global config overview")
		err := c.ResolveFields(appCtx, 1, "Renderers")
		if err != nil {
			return fmt.Errorf("failed to get global config overview: %w", err)
		}

		logger.D("resolving user renderers", log.Int("count", len(c.Renderers)))
		for name, r := range c.Renderers {
			logger := logger.WithFields(
				log.Any("renderer", name),
			)

			logger.D("resolving renderer config fields")
			err = r.ResolveFields(appCtx, -1)
			if err != nil {
				return fmt.Errorf("failed to resolve renderer %q config: %w", name, err)
			}

			err = r.Init(appCtx)
			if err != nil {
				return fmt.Errorf("failed to initialize renderer %q: %w", name, err)
			}

			appCtx.AddRenderer(name, c.Renderers[name])
		}
		logger.D("resolved all renderers", log.Int("count", len(appCtx.AllRenderers())))

	}

	// step 4: resolve shells
	{
		logger.D("resolving shell config overview")

		err := c.ResolveFields(appCtx, 1, "Shells")
		if err != nil {
			return fmt.Errorf("failed to resolve shell config overview: %w", err)
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
				return fmt.Errorf("failed to resolve shell %q #%d config: %w", v.Name(), i, err)
			}

			err = v.InitBaseTool(
				"shell", string(v.Name()), appCtx.CacheDir(), v,
			)
			if err != nil {
				return fmt.Errorf("failed to initialize shell %q", v.Name())
			}

			logger.V("adding shell")
			appCtx.AddShell(string(v.Name()), c.Shells[i])
		}
	}

	// step 5: resolve tools and tasks
	logger.D("resolving top level config")
	err := c.ResolveFields(appCtx, 1)
	if err != nil {
		return fmt.Errorf("failed to resolve top-level config: %w", err)
	}
	logger.V("resolved top-level config", log.Any("result", c))

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

			key := dukkha.ToolKey{
				Kind: toolKind,
				Name: t.Name(),
			}

			logger := logger.WithFields(
				log.String("key", key.String()),
				log.Int("index", i),
			)

			logger.D("resolving tool config fields")
			err = t.ResolveFields(appCtx, -1)
			if err != nil {
				return fmt.Errorf(
					"failed to resolve tool %q config: %w",
					key, err,
				)
			}

			logger.V("initializing tool")
			err = t.Init(toolKind, appCtx.CacheDir())
			if err != nil {
				return fmt.Errorf(
					"failed to initialize tool %q: %w",
					key, err,
				)
			}

			// append tasks without tool name
			// they are meant for all tools in the same kind they belong to
			noToolNameTasks, _ := appCtx.GetToolSpecificTasks(
				dukkha.ToolKey{Kind: toolKind, Name: ""},
			)
			appCtx.AddToolSpecificTasks(
				toolKind, t.Name(),
				noToolNameTasks,
			)

			tasks, _ := appCtx.GetToolSpecificTasks(key)

			if logger.Enabled(log.LevelVerbose) {
				logger.D("resolving tool tasks", log.Any("tasks", tasks))
			} else {
				logger.D("resolving tool tasks")
			}

			err = t.ResolveTasks(tasks)
			if err != nil {
				return fmt.Errorf(
					"failed to resolve tasks for tool %q: %w",
					key, err,
				)
			}

			appCtx.AddTool(key, c.Tools[string(toolKind)][i])
		}
	}

	return nil
}
