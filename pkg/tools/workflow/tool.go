package workflow

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "workflow"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Workflow struct{}

func (t *Workflow) DefaultExecutable() string { return "" }
func (t *Workflow) Kind() dukkha.ToolKind     { return ToolKind }

type Tool struct {
	tools.BaseTool[Workflow, *Workflow]
}
