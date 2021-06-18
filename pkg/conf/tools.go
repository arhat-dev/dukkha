package conf

import (
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

type ToolsConfig struct {
	field.BaseField

	// map[tools.ToolKey]tools.ToolConfig
	Tools []tools.ToolConfig `yaml:",inline" dukkha:"other"`
}
