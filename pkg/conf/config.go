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
	"context"
	"fmt"
	"io"
	"strings"

	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/shell_file"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/renderer/template_file"
	"arhat.dev/dukkha/pkg/tools"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Bootstrap BootstrapConfig `yaml:"-"`
	Shell     ShellConfigList `yaml:"-"`

	Tools ToolsConfig `yaml:"-"`

	Tasks TasksConfig `yaml:"-"`
}

const (
	topLevelFieldBootstrap = "bootstrap"
	topLevelFieldShell     = "shell"
	topLevelFieldTools     = "tools"
)

func Unmarshal(r io.Reader, out interface{}) error {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	return dec.Decode(out)
}

func (c *Config) Decode(ctx context.Context, mergedConfig map[string]interface{}) error {
	if mergedConfig == nil {
		return nil
	}

	if c.Tools == nil {
		c.Tools = make(map[tools.ToolKey]tools.ToolConfig)
	}

	if c.Tasks == nil {
		c.Tasks = make(map[tools.TaskTypeKey][]tools.TaskConfig)
	}

	// resolve bootstrap config first
	err := c.Bootstrap.Resolve(ctx, mergedConfig[topLevelFieldBootstrap])
	if err != nil {
		return fmt.Errorf("conf: bootstrap config not resolved: %w", err)
	}

	if c.Bootstrap.Shell == "" {
		return fmt.Errorf("conf: unable to get a shell name, please set bootstrap.shell manually")
	}

	// create a renderer manager with essential renderers
	mgr := renderer.NewManager()
	err = multierr.Append(err, mgr.Add(&shell.Config{ExecFunc: c.Bootstrap.Exec}, shell.DefaultName))
	err = multierr.Append(err, mgr.Add(&shell_file.Config{ExecFunc: c.Bootstrap.Exec}, shell_file.DefaultName))
	err = multierr.Append(err, mgr.Add(&template.Config{}, template.DefaultName))
	err = multierr.Append(err, mgr.Add(&template_file.Config{}, template_file.DefaultName))
	err = multierr.Append(err, mgr.Add(&file.Config{}, file.DefaultName))
	ctx = renderer.WithManager(ctx, mgr)

	err = c.Shell.resolve(ctx, mergedConfig[topLevelFieldShell])
	if err != nil {
		return fmt.Errorf("conf: unable to resolve shell config: %w", err)
	}

	//
	// resolve other configs using fully configured renderers
	//

	// resolve tools config
	err = c.Tools.resolve(ctx, mergedConfig[topLevelFieldTools])
	if err != nil {
		return fmt.Errorf("conf: unable to resovle tools config: %w", err)
	}

	// resolve tasks
	for k, v := range mergedConfig {
		parts := strings.SplitN(k, "@", 2)

		f := &Field{
			Name: parts[0],
		}
		if len(parts) == 2 {
			// has rendering suffix
			f.Renderer = parts[1]
		}

		var err error
		switch f.Name {
		case topLevelFieldBootstrap, topLevelFieldShell, topLevelFieldTools:
			continue
		}

		err = c.resolveTasks(f, v)
		if err != nil {
			return fmt.Errorf("conf: invalid config: %w", err)
		}
	}

	return nil
}

func (c *Config) resolveTasks(taskField *Field, data interface{}) error {
	taskParts := strings.Split(taskField.Name, ":")

	var (
		toolName = taskParts[0]
		toolID   string
		taskType string
	)

	switch len(taskParts) {
	case 2:
		taskType = taskParts[1]
	case 3:
		toolID, taskType = taskParts[1], taskParts[2]
	default:
		return fmt.Errorf(
			"task: invalid task field %q, expecting 1 or 2 colon, got %d",
			taskField.Name, len(taskParts),
		)
	}

	key, err := tools.CreateTaskTypeKey(toolName, toolID, taskType)
	if err != nil {
		return fmt.Errorf("task: invalid task field: %w", err)
	}

	c.Tasks[*key] = nil
	if len(taskField.Renderer) != 0 {
		// requires extra rendering
		strVal, ok := data.(string)
		if !ok {
			return fmt.Errorf("task.%s: unexpected non string value", key.String())
		}

		// TODO: mark to be processed later
		_ = strVal
	} else {
		// can unmarshal childs directly
		tasksBytes, err := yaml.Marshal(data)
		if err != nil {
			return fmt.Errorf("task.%s: marhsal: %w", key.String(), err)
		}

		taskConfigType, err := tools.GetTaskConfigType(toolName, taskType)
		if err != nil {
			return fmt.Errorf("task.%s", key.String())
		}

		_ = tasksBytes
		_ = taskConfigType
	}

	return nil
}
