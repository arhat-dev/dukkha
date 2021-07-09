package helm

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "helm"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool {
		return &Tool{}
	})
}

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`
}
