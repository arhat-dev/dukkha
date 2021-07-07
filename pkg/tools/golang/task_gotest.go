package golang

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/types"
)

const TaskKindTest = "test"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindTest,
		func(toolName string) dukkha.Task {
			t := &TaskTest{}
			t.SetToolName(toolName)
			return t
		},
	)
}

var _ dukkha.Task = (*TaskTest)(nil)

type TaskTest struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Chdir string `yaml:"chdir"`
}

func (c *TaskTest) ToolKind() dukkha.ToolKind { return ToolKind }
func (c *TaskTest) Kind() dukkha.TaskKind     { return TaskKindTest }

func (c *TaskTest) GetExecSpecs(rc types.RenderingContext, toolCmd []string) ([]dukkha.TaskExecSpec, error) {
	spec := &dukkha.TaskExecSpec{
		Chdir:   c.Chdir,
		Command: sliceutils.NewStrings(toolCmd, "test", "./..."),
	}
	return []dukkha.TaskExecSpec{*spec}, nil
}
