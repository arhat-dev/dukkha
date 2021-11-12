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
	"arhat.dev/pkg/rshelper"
	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/renderer/echo"
	"arhat.dev/dukkha/pkg/renderer/env"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/renderer/transform"
	"arhat.dev/dukkha/pkg/tools"
)

func NewConfig() *Config {
	return rshelper.InitAll(&Config{}, &rs.Options{
		InterfaceTypeHandler: dukkha.GlobalInterfaceTypeHandler,
	}).(*Config)
}

type GlobalConfig struct {
	rs.BaseField `yaml:"-"`

	// CacheDir set DUKKHA_CACHE_DIR to store script file and intermediate
	// task execution data
	CacheDir string `yaml:"cache_dir"`

	// DefaultGitBranch set GIT_DEFAULT_BRANCH, useful when dukkha can not
	// detect branch name of origin/HEAD (e.g. github ci environment)
	//
	// If your have multiple definitions of this option in different config
	// file, only the first occurrence of the option is used.
	DefaultGitBranch string `yaml:"default_git_branch"`

	// Env add global environment variables for all working parts in dukkha
	Env dukkha.Env `yaml:"env"`

	// Values is the global store of runtime values
	//
	// accessible from renderer template `{{ values.YOUR_VAL_KEY }}`
	// and renderer env/shell `${VALUES.YOUR_VAL_KEY}`
	Values rs.AnyObjectMap `yaml:"values"`
}

func (g *GlobalConfig) Merge(a *GlobalConfig) error {
	err := g.BaseField.Inherit(&a.BaseField)
	if err != nil {
		return fmt.Errorf("failed to inherit other global config: %w", err)
	}

	g.Env = append(g.Env, a.Env...)
	if len(a.CacheDir) != 0 {
		g.CacheDir = a.CacheDir
	}

	if len(a.DefaultGitBranch) != 0 {
		g.DefaultGitBranch = a.DefaultGitBranch
	}

	err = g.Values.Inherit(&a.Values.BaseField)
	if err != nil {
		return fmt.Errorf("failed to merge global values: %w", err)
	}

	if len(a.Values.Data) != 0 {
		if g.Values.Data == nil {
			g.Values.Data = a.Values.Data
		} else {
			for k, v := range a.Values.Data {
				g.Values.Data[k] = v
			}
		}
	}

	return nil
}

func (g *GlobalConfig) ResolveAllButValues(rc dukkha.ConfigResolvingContext) error {
	err := dukkha.ResolveEnv(g, rc, "Env", "env")
	if err != nil {
		return fmt.Errorf("failed to resolve global env: %w", err)
	}

	err = g.ResolveFields(rc, -1, "cache_dir")
	if err != nil {
		return fmt.Errorf("failed to resolve cache dir: %w", err)
	}

	err = g.ResolveFields(rc, -1, "default_git_branch")
	if err != nil {
		return fmt.Errorf("failed to resolve default git branch: %w", err)
	}

	return nil
}

type Config struct {
	rs.BaseField `yaml:"-"`

	// Global options only have limited rendering suffix support
	Global GlobalConfig `yaml:"global"`

	// Include other files using path relative to DUKKHA_WORKING_DIR
	// only local path (and path glob) is supported.
	//
	// You should always use relative path unless you do not want to
	// maintain compatibility with other environments.
	//
	// no rendering suffix support for this field
	Include []string `yaml:"include"`

	// Shells for command execution
	Shells []*tools.BaseToolWithInit `yaml:"shells"`

	// Renderers config options
	Renderers map[string]dukkha.Renderer `yaml:"renderers"`

	// Tools config options for registered tools
	Tools Tools `yaml:"tools"`

	Tasks map[string][]dukkha.Task `rs:"other"`
}

var _ yaml.Unmarshaler = (*Tools)(nil)

type Tools struct {
	rs.BaseField `yaml:"-"`

	Data map[string][]dukkha.Tool `rs:"other"`
}

func (m *Tools) Merge(a *Tools) error {
	err := m.BaseField.Inherit(&a.BaseField)
	if err != nil {
		return fmt.Errorf("failed to inherit other tools config: %w", err)
	}

	if len(a.Data) != 0 {
		if m.Data == nil {
			m.Data = make(map[string][]dukkha.Tool)
		}

		for k := range a.Data {
			m.Data[k] = append(m.Data[k], a.Data[k]...)
		}
	}

	return nil
}

func (c *Config) Merge(a *Config) error {
	err := c.BaseField.Inherit(&a.BaseField)
	if err != nil {
		return fmt.Errorf("failed to inherit other top level base field: %w", err)
	}

	err = c.Global.Merge(&a.Global)
	if err != nil {
		return err
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

// Resolve resolves all top level dukkha config
// to gain an overview of all tools and tasks
// nolint:gocyclo
func (c *Config) Resolve(appCtx dukkha.ConfigResolvingContext, needTasks bool) error {
	logger := log.Log.WithName("config")

	// step 1: create essential renderers to initialize renderers
	{
		logger.V("creating essential renderers")

		// TODO: let user decide what renderers to use
		// 		 resolve renderers first?
		appCtx.AddRenderer(echo.DefaultName, echo.NewDefault(""))
		appCtx.AddRenderer(env.DefaultName, env.NewDefault(""))
		appCtx.AddRenderer(shell.DefaultName, shell.NewDefault(""))
		appCtx.AddRenderer(template.DefaultName, template.NewDefault(""))
		appCtx.AddRenderer(file.DefaultName, file.NewDefault(""))
		appCtx.AddRenderer(transform.DefaultName, transform.NewDefault(""))

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

	// step 2: resolve global config (except Values), ensure cache dir exists
	{
		logger.D("resolving global config")
		err := c.ResolveFields(appCtx, 1, "global")
		if err != nil {
			return fmt.Errorf("failed to get global config overview: %w", err)
		}

		err = c.Global.ResolveAllButValues(appCtx)
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

	// NOTE: resolving renderers requires dukkha cache dir to been set
	// step 3: resolve renderers
	{
		logger.D("resolving renderers config overview")
		err := c.ResolveFields(appCtx, 1, "renderers")
		if err != nil {
			return fmt.Errorf("failed to get renderers config overview: %w", err)
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

	// step 4: resolve global Values
	{
		logger.D("resolving global values")

		err := c.Global.ResolveFields(appCtx, -1, "values")
		if err != nil {
			return fmt.Errorf("failed to resolve global values: %w", err)
		}

		logger.V("resolved global values", log.Any("values", c.Global.Values))

		logger.D("adding global values")
		values := c.Global.Values.NormalizedValue()
		if err != nil {
			return fmt.Errorf("failed to normalize global values: %w", err)
		}

		err = appCtx.AddValues(values)
		if err != nil {
			return fmt.Errorf("failed to add global values: %w", err)
		}
	}

	// step 5: resolve shells
	{
		logger.D("resolving shell config overview")

		err := c.ResolveFields(appCtx, 1, "shells")
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

	// save some time if the command is not interacting with tasks
	if !needTasks {
		return nil
	}

	// step 6: resolve tools and tasks
	logger.D("resolving top level config")
	err := c.ResolveFields(appCtx, 1, "tools", "Tasks")
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

	logger.V("resolving tools", log.Int("count", len(c.Tools.Data)))
	for tk, toolSet := range c.Tools.Data {
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
			err = t.ResolveFields(appCtx, 1)
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

			appCtx.AddTool(key, c.Tools.Data[string(toolKind)][i])
		}
	}

	return nil
}
