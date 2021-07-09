package github

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "github"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`
}

func (t *Tool) Init(kind dukkha.ToolKind, cachdDir string) error {
	return t.BaseTool.InitBaseTool(ToolKind, "gh", cachdDir)
}

func (t *Tool) ResolveFields(rc field.RenderingHandler, depth int, fieldName string) error {
	return t.BaseField.ResolveFields(rc, depth, fieldName)
}
