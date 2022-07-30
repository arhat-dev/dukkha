package tools

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/rs"
)

var _ dukkha.Task = (*BaseTask[testTaskImpl, *testTaskImpl])(nil)

type testTaskImpl struct {
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

	NonAnonymousTask BaseTask[struct{}, *testTaskImpl] `yaml:"non_anonymous_task"`

	InnerField struct {
		// we do not care about fields inside inner field
		BaseTask[struct{}, *testTaskImpl]
		Foo string `yaml:"foo"`
	} `yaml:"inner_field"`
}

func (b *testTaskImpl) ToolKind() dukkha.ToolKind { return "" }
func (b *testTaskImpl) Kind() dukkha.TaskKind     { return "" }
func (b *testTaskImpl) LinkParent(p BaseTaskType) {}

func (b *testTaskImpl) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return nil, nil
}

type TestTask = BaseTask[testTaskImpl, *testTaskImpl]
