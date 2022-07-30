package tools

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/rs"
)

var _ dukkha.Tool = (*BaseTool[testToolImpl, *testToolImpl])(nil)

type testToolImpl struct {
	A            string
	B            string `yaml:"bar"`
	InlineStruct struct {
		rs.BaseField

		NestedInlineMap    map[string]int `yaml:",inline"`
		NestedInlineStruct struct {
			rs.BaseField

			C string
		} `yaml:",inline"`
	} `yaml:",inline"`

	Ignored string `yaml:"-"`

	NonAnonymousTask BaseTool[struct{}, *testToolImpl] `yaml:"non_anonymous_task"`

	InnerField struct {
		// we do not care about fields inside inner field
		BaseTool[struct{}, *testToolImpl]
		Foo string `yaml:"foo"`
	} `yaml:"inner_field"`
}

func (t *testToolImpl) Kind() dukkha.ToolKind     { return "" }
func (t *testToolImpl) DefaultExecutable() string { return "" }

type TestTool = BaseTool[testToolImpl, *testToolImpl]
