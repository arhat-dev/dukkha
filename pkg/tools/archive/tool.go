package archive

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "archive"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Archive struct{}

func (t *Archive) DefaultExecutable() string { return "" }
func (t *Archive) Kind() dukkha.ToolKind     { return ToolKind }

type Tool struct {
	tools.BaseTool[Archive, *Archive]
}
