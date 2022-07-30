package tool_git

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "git"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Git struct{}

func (t *Git) DefaultExecutable() string { return "git" }
func (t *Git) Kind() dukkha.ToolKind     { return ToolKind }

type Tool struct{ tools.BaseTool[Git, *Git] }
