package golang

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindTest = "test"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^golang(:.+){0,1}:test$`),
		func(subMatches []string) interface{} {
			t := &TaskTest{}
			if len(subMatches) != 0 {
				t.SetToolName(strings.TrimPrefix(subMatches[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskTest)(nil)

type TaskTest struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Chdir string `yaml:"chdir"`
}

func (c *TaskTest) ToolKind() string { return ToolKind }
func (c *TaskTest) TaskKind() string { return TaskKindTest }

func (c *TaskTest) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	spec := &tools.TaskExecSpec{
		Chdir:   c.Chdir,
		Command: sliceutils.NewStringSlice(toolCmd, "test", "./..."),
	}
	return []tools.TaskExecSpec{*spec}, nil
}
