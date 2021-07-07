package buildah

import (
	"regexp"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "buildah"

func init() {
	field.RegisterInterfaceField(
		dukkha.ToolType,
		regexp.MustCompile("^buildah$"),
		func([]string) interface{} { return &Tool{} },
	)
}

var _ dukkha.Tool = (*Tool)(nil)

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`
}

func (t *Tool) Kind() dukkha.ToolKind { return ToolKind }
