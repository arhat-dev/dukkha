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
	"reflect"

	"arhat.dev/pkg/log"
	"go.uber.org/multierr"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/file"
	"arhat.dev/dukkha/pkg/renderer/shell"
	"arhat.dev/dukkha/pkg/renderer/shell_file"
	"arhat.dev/dukkha/pkg/renderer/template"
	"arhat.dev/dukkha/pkg/renderer/template_file"
)

func NewConfig() *Config {
	return field.New(&Config{
		Shell: field.New(&ShellConfigList{}).(*ShellConfigList),
		Tools: field.New(&ToolsConfig{}).(*ToolsConfig),
		Tasks: field.New(&TasksConfig{}).(*TasksConfig),
	}).(*Config)
}

type Config struct {
	field.BaseField

	Log       log.Config      `yaml:"log"`
	Bootstrap BootstrapConfig `yaml:"bootstrap"`

	Shell *ShellConfigList `yaml:"shell"`
	Tools *ToolsConfig     `yaml:"tools"`

	// use inline for all tasks so it will get notified with all yaml nodes
	Tasks *TasksConfig `yaml:",inline" dukkha:"other"`
}

func (c *Config) Type() reflect.Type {
	return reflect.TypeOf(c)
}

func (c *Config) Resolve(ctx context.Context) (context.Context, error) {
	// bootstrap config was resolved when unmarshaling
	if c.Bootstrap.Shell == "" {
		return nil, fmt.Errorf("conf: unable to get a shell name, please set bootstrap.shell manually")
	}

	var err error

	// create a renderer manager with essential renderers
	mgr := renderer.NewManager()
	err = multierr.Append(err, mgr.Add(&shell.Config{ExecFunc: c.Bootstrap.Exec}, shell.DefaultName))
	err = multierr.Append(err, mgr.Add(&shell_file.Config{ExecFunc: c.Bootstrap.Exec}, shell_file.DefaultName))
	err = multierr.Append(err, mgr.Add(&template.Config{}, template.DefaultName))
	err = multierr.Append(err, mgr.Add(&template_file.Config{}, template_file.DefaultName))
	err = multierr.Append(err, mgr.Add(&file.Config{}, file.DefaultName))
	if err != nil {
		return nil, fmt.Errorf("conf: failed to create essential renderers: %w", err)
	}

	ctx = renderer.WithManager(ctx, mgr)

	// resolve shells to add shell renderers
	err = c.Shell.Resolve(ctx)
	if err != nil {
		return nil, fmt.Errorf("conf: unable to resolve shell config: %w", err)
	}

	//
	// resolve other configs using fully configured renderers
	//

	// resolve tools config
	err = c.Tools.Resolve(ctx)
	if err != nil {
		return nil, fmt.Errorf("conf: unable to resolve tools config: %w", err)
	}

	// resolve tasks at last
	err = c.Tasks.Resolve(ctx)
	if err != nil {
		return nil, fmt.Errorf("conf: unable to resolve tasks: %w", err)
	}

	return ctx, nil
}
