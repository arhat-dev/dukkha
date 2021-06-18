package conf

import (
	"arhat.dev/dukkha/pkg/field"
)

type ShellConfig struct {
	field.BaseField

	Name    string   `yaml:"name"`
	Path    string   `yaml:"path"`
	Env     []string `yaml:"env"`
	Command []string `yaml:"command"`
	Args    []string `yaml:"args"`
}

type ShellConfigList struct {
	field.BaseField `yaml:"-"`

	ShellConfigs []ShellConfig `yaml:",inline"`
}
