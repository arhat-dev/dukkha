package helm

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "helm"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Helm struct{}

func (t *Helm) DefaultExecutable() string { return "helm" }
func (t *Helm) Kind() dukkha.ToolKind     { return ToolKind }

type Tool struct{ tools.BaseTool[Helm, *Helm] }
