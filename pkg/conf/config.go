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
	"reflect"
	"sync"

	"arhat.dev/pkg/log"
	"arhat.dev/pkg/synchain"
	"arhat.dev/rs"
	"go.uber.org/multierr"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	tools_shell "arhat.dev/dukkha/pkg/tools/shell"
)

func NewConfig() *Config {
	cfg := new(Config)
	rs.InitRecursively(reflect.ValueOf(cfg), &rs.Options{
		InterfaceTypeHandler: dukkha.GlobalInterfaceTypeHandler,
	})
	return cfg
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
	Shells []*tools_shell.Tool `yaml:"shells"`

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
	for i, group := range c.Renderers {
		logger := logger.WithFields(log.Int("index", i))

		logger.D("resolving renderer group")

		// renderers in the same group should be resolved all at once
		// without knowning each other

		err = group.ResolveFields(appCtx, -1)
		if err != nil {
			return fmt.Errorf("resolving renderer group #%d: %w", i, err)
		}

		for name, r := range group.Renderers {
			err = r.ResolveFields(appCtx, -1)
			if err != nil {
				return fmt.Errorf("resolving renderer %q: %w", name, err)
			}

			logger.V("resoving renderer", log.String("name", name), log.String("alias", r.Alias()))

			err = r.Init(appCtx.RendererCacheFS(name))
			if err != nil {
				return fmt.Errorf("initializing renderer %q: %w", name, err)
			}

			appCtx.AddRenderer(name, group.Renderers[name])
			if len(r.Alias()) != 0 {
				appCtx.AddRenderer(r.Alias(), group.Renderers[name])
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

		err = v.InitWithName(string(v.Name()), appCtx.ToolCacheFS(v))
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
//
// nolint:gocyclo
func (c *Config) Resolve(appCtx dukkha.ConfigResolvingContext, flags ReadFlag) (err error) {
	logger := log.Log.WithName("config")

	if flags&ReadFlag_Global != 0 {
		// step 1: resolve global config (except Values)
		logger.D("resolving global config")
		err = c.ResolveFields(appCtx, 1, "global")
		if err != nil {
			return fmt.Errorf("get global config overview: %w", err)
		}

		err = c.Global.ResolveAllButValues(appCtx)
		if err != nil {
			return fmt.Errorf("resolve global config: %w", err)
		}

		logger.V("resolved global config", log.Any("result", c.Global))

		if len(c.Global.DefaultGitBranch) != 0 {
			appCtx.(di.DefaultGitBranchOverrider).OverrideDefaultGitBranch(c.Global.DefaultGitBranch)
		}

		// step 2: resolve global Values
		logger.D("resolving global values")

		err = c.Global.ResolveFields(appCtx, -1, "values")
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

	var sg synchain.Synchain
	sg.Init()

	// resolve tasks first, as tools need tasks
	if flags&ReadFlag_Task != 0 {
		sg.Go(func(t synchain.Ticket) error {
			logger.D("resolving tasks overview")
			err := c.ResolveFields(appCtx, 1, "Tasks" /* inline field, field name as tag name */)
			if err != nil {
				return fmt.Errorf("gain overview of tasks: %w", err)
			}
			logger.V("resolved tasks overview", log.Any("result", c.Tasks))

			var wg sync.WaitGroup
			logger.V("groupping tasks by tool")

			errCh := make(chan error, len(c.Tasks))
			wg.Add(len(c.Tasks))
			for _, tasks := range c.Tasks {
				go func(tasks []dukkha.Task) {
					defer wg.Done()

					for _, tsk := range tasks {
						err2 := tsk.ResolveFields(appCtx, -1, "name")
						if err2 != nil {
							select {
							case <-appCtx.Done():
							case errCh <- fmt.Errorf("reoslving task name: %w", err2):
							}

							return
						}

						// FIXME: task name is empty at this time
						err2 = tsk.Init(appCtx.TaskCacheFS(tsk))
						if err2 != nil {
							select {
							case <-appCtx.Done():
							case errCh <- fmt.Errorf("task init %q: %w", tsk.Name(), err2):
							}

							return
						}
					}
				}(tasks)
			}

			wg.Wait()

			for {
				select {
				case err2 := <-errCh:
					err = multierr.Append(err, err2)
					continue
				default:
				}

				break
			}

			for _, tasks := range c.Tasks {
				appCtx.AddToolSpecificTasks(
					tasks[0].ToolKind(),
					tasks[0].ToolName(),
					tasks,
				)
			}

			// TODO: do we need this?
			// sg.Lock(t)

			return nil
		})
	}

	// step 3: resolve tools and tasks
	if flags&ReadFlag_Tool != 0 {
		sg.Go(func(t synchain.Ticket) error {
			logger.D("resolving tools overview")
			err := c.ResolveFields(appCtx, 2, "tools")
			if err != nil {
				return fmt.Errorf("gain overview of tools: %w", err)
			}
			logger.V("resolved tools overview", log.Any("result", c.Tools))

			logger.V("resolving tools", log.Int("count", len(c.Tools.Tools)))
			for tk, toolSet := range c.Tools.Tools {
				visited := make(map[dukkha.ToolName]struct{})

				for i, t := range toolSet {
					err2 := t.ResolveFields(appCtx, -1, "name")
					if err2 != nil {
						return fmt.Errorf("resolve tool name: %w", err2)
					}

					key := dukkha.ToolKey{
						Kind: dukkha.ToolKind(tk),
						Name: t.Name(),
					}

					// do not allow empty name
					if len(key.Name) == 0 {
						return fmt.Errorf("invalid %q tool without name, index %d", key.Kind, i)
					}

					// ensure tool names are unique
					if _, ok := visited[key.Name]; ok {
						return fmt.Errorf("duplicate tool name %q of kind %q", t.Name(), key.Kind)
					}

					visited[key.Name] = struct{}{}

					logger := logger.WithFields(
						log.String("key", key.String()),
						log.Int("index", i),
					)

					logger.V("init tool")
					err2 = t.Init(appCtx.ToolCacheFS(t))
					if err2 != nil {
						return fmt.Errorf("init tool %q: %w", key, err2)
					}

					appCtx.AddTool(key, c.Tools.Tools[string(key.Kind)][i])
				}
			}

			// wait for tasks in other goroutine
			if !sg.Lock(t) {
				return nil
			}

			for tk, toolSet := range c.Tools.Tools {
				for _, t := range toolSet {

					key := dukkha.ToolKey{
						Kind: dukkha.ToolKind(tk),
						Name: t.Name(),
					}

					// append tasks without tool name
					// they are meant for all tools in the same kind they belong to
					noToolNameTasks, _ := appCtx.GetToolSpecificTasks(
						dukkha.ToolKey{Kind: key.Kind, Name: ""},
					)
					appCtx.AddToolSpecificTasks(
						key.Kind, key.Name, noToolNameTasks,
					)

					tasks, _ := appCtx.GetToolSpecificTasks(key)

					if logger.Enabled(log.LevelVerbose) {
						logger.V("resolving tool tasks", log.Any("tasks", tasks))
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
				}
			}

			return nil
		})
	}

	sg.Wait()
	return nil
}
