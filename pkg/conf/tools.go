package conf

import (
	"context"

	"arhat.dev/dukkha/pkg/tools"
)

type ToolsConfig map[tools.ToolKey]tools.ToolConfig

func (c ToolsConfig) resolve(ctx context.Context, data interface{}) error {
	return nil
}
