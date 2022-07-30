package docker

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "docker"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Docker struct{}

func (t *Docker) DefaultExecutable() string { return "docker" }
func (t *Docker) Kind() dukkha.ToolKind     { return ToolKind }

type Tool struct {
	tools.BaseTool[Docker, *Docker]
}
