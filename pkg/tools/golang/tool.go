package golang

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "golang"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Tool struct {
	rs.BaseField

	tools.BaseTool `yaml:",inline"`
}

func (t *Tool) Init(kind dukkha.ToolKind, cachdDir string) error {
	return t.BaseTool.InitBaseTool(ToolKind, "go", cachdDir, t)
}
