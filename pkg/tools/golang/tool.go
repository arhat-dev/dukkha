package golang

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "golang"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Golang struct{}

func (t *Golang) DefaultExecutable() string { return "go" }
func (t *Golang) Kind() dukkha.ToolKind     { return ToolKind }

type Tool struct {
	tools.BaseTool[Golang, *Golang]
}
