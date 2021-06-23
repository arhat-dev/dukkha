package golang

import (
	"fmt"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindBuild = "build"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^golang(:.+)?:build$`),
		func(params []string) interface{} {
			t := &TaskBuild{}
			if len(params) != 0 {
				t.SetToolName(params[0])
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskBuild)(nil)

type TaskBuild struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`
}

func (c *TaskBuild) ToolKind() string { return ToolKind }
func (c *TaskBuild) TaskKind() string { return TaskKindBuild }

func (c *TaskBuild) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	return nil, fmt.Errorf("unimplemented")
}
