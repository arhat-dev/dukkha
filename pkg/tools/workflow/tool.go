package workflow

import (
	"regexp"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "workflow"

func init() {
	field.RegisterInterfaceField(
		dukkha.ToolType,
		regexp.MustCompile("^workflow$"),
		func(_ []string) interface{} { return &Tool{} },
	)
}

var _ dukkha.Tool = (*Tool)(nil)

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`

	workflows map[string]*TaskRun
}

func (t *Tool) Kind() dukkha.ToolKind { return ToolKind }

func (t *Tool) Init(cachdDir string) error {
	return t.BaseTool.InitBaseTool("", cachdDir)
}
