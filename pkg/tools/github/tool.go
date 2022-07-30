package github

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "github"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Github struct{}

func (t *Github) DefaultExecutable() string { return "gh" }
func (t *Github) Kind() dukkha.ToolKind     { return ToolKind }

type Tool struct {
	tools.BaseTool[Github, *Github]
}
