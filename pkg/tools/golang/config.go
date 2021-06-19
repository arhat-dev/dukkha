package golang

import (
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolName = "golang"

func init() {
	field.RegisterInterfaceField(
		tools.ToolConfigType,
		regexp.MustCompile("^golang$"),
		func() interface{} { return &Config{} },
	)
}

var _ tools.ToolConfig = (*Config)(nil)

type Config struct {
	field.BaseField

	Name string `yaml:"name"`
}

func (c *Config) Kind() string { return ToolName }
