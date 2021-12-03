package archive

import (
	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "archive"

func init() {
	dukkha.RegisterTool(ToolKind, func() dukkha.Tool { return &Tool{} })
}

type Tool struct {
	rs.BaseField `yaml:"-"`

	ToolName dukkha.ToolName `yaml:"name"`

	tools.BaseTool `yaml:",inline"`
}

func (t *Tool) Init(cacheFS *fshelper.OSFS) error {
	return t.InitBaseTool("", cacheFS, t)
}

func (t *Tool) Name() dukkha.ToolName { return t.ToolName }
func (t *Tool) Kind() dukkha.ToolKind { return ToolKind }
func (t *Tool) Key() dukkha.ToolKey {
	return dukkha.ToolKey{Kind: t.Kind(), Name: t.Name()}
}
