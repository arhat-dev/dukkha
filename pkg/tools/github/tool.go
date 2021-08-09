package github

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "github"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Tool struct {
	rs.BaseField

	tools.BaseTool `yaml:",inline"`
}

func (t *Tool) Init(kind dukkha.ToolKind, cachdDir string) error {
	return t.BaseTool.InitBaseTool(ToolKind, "gh", cachdDir, t)
}
