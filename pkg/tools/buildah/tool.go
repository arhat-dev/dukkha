package buildah

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "buildah"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Buildah struct{}

func (t *Buildah) DefaultExecutable() string { return "buildah" }
func (t *Buildah) Kind() dukkha.ToolKind     { return ToolKind }

type Tool struct {
	tools.BaseTool[Buildah, *Buildah]
}
