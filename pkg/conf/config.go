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
	"arhat.dev/pkg/log"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

func NewConfig() *Config {
	return field.Init(&Config{}).(*Config)
}

type Config struct {
	field.BaseField

	// no rendering suffix support
	Log       *log.Config     `yaml:"log"`
	Bootstrap BootstrapConfig `yaml:"bootstrap"`

	// Shells for rendering and command execution
	Shells []tools.BaseTool `yaml:"shells"`

	// Language or tool specific tools
	Tools map[string][]tools.Tool `yaml:"tools"`
	Tasks map[string][]tools.Task `dukkha:"other"`
}

func (c *Config) Merge(a *Config) {
	if a.Log != nil {
		c.Log = a.Log
	}

	if len(a.Bootstrap.ScriptCmd) != 0 {
		c.Bootstrap = a.Bootstrap
	}

	c.Shells = append(c.Shells, a.Shells...)

	if len(a.Tools) != 0 {
		if c.Tools == nil {
			c.Tools = a.Tools
		} else {
			for k := range a.Tools {
				c.Tools[k] = a.Tools[k]
			}
		}
	}

	if len(a.Tasks) != 0 {
		if c.Tasks == nil {
			c.Tasks = a.Tasks
		} else {
			for k := range a.Tasks {
				c.Tasks[k] = a.Tasks[k]
			}
		}
	}
}
