package golang

import (
	"fmt"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "golang"

func init() {
	field.RegisterInterfaceField(
		tools.ToolType,
		regexp.MustCompile("^golang$"),
		func() interface{} { return &Config{} },
	)
}

var _ tools.Tool = (*Config)(nil)

type Config struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`
}

func (c *Config) ToolKind() string { return ToolKind }

func (c *Config) ResolveTasks(tasks []tools.Task) error {
	return fmt.Errorf("unimplemented")
}
