package helm

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "helm"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool {
		return &Tool{}
	})
}

type Tool struct {
	rs.BaseField

	tools.BaseTool `yaml:",inline"`
}

func (t *Tool) Init(kind dukkha.ToolKind, cacheDir string) error {
	return t.InitBaseTool(kind, "helm", cacheDir, t)
}
