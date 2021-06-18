package conf

import (
	"context"
	"reflect"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

type ToolsConfig struct {
	field.BaseField

	// map[tools.ToolKey]tools.ToolConfig
	Tools []tools.ToolConfig `yaml:",inline" dukkha:"other"`
}

func (c *ToolsConfig) Type() reflect.Type {
	return reflect.TypeOf(c)
}

func (c *ToolsConfig) Resolve(ctx context.Context) error {
	return nil
}
