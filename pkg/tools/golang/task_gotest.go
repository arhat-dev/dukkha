package golang

import (
	"fmt"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindTest = "test"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^golang(:.+)?:test$`),
		func(params []string) interface{} {
			t := &TaskTest{}
			if len(params) != 0 {
				t.SetToolName(params[0])
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskTest)(nil)

type TaskTest struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`
}

func (c *TaskTest) ToolKind() string { return ToolKind }
func (c *TaskTest) TaskKind() string { return TaskKindBuild }

func (c *TaskTest) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	return nil, fmt.Errorf("unimplemented")
}
