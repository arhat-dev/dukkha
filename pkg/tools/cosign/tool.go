package cosign

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "cosign"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Cosign struct{}

func (t *Cosign) DefaultExecutable() string { return "cosign" }
func (t *Cosign) Kind() dukkha.ToolKind     { return ToolKind }

type Tool struct {
	tools.BaseTool[Cosign, *Cosign]
}
