package golang

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindTest = "test"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindTest,
		func(toolName string) dukkha.Task {
			t := &TaskTest{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindTest, t)
			return t
		},
	)
}

type TaskTest struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Chdir string `yaml:"chdir"`
}

func (c *TaskTest) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		steps = append(steps, dukkha.TaskExecSpec{
			Env:       sliceutils.NewStrings(c.Env),
			Chdir:     c.Chdir,
			Command:   sliceutils.NewStrings(options.ToolCmd, "test", "./..."),
			UseShell:  options.UseShell,
			ShellName: options.ShellName,
		})
		return nil
	})

	return steps, err
}
