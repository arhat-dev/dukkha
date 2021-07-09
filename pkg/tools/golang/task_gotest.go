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
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindTest)
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
	rc dukkha.TaskExecContext,
	useShell bool,
	shellName string,
	toolCmd []string,
) ([]dukkha.TaskExecSpec, error) {
	spec := &dukkha.TaskExecSpec{
		Chdir:     c.Chdir,
		Command:   sliceutils.NewStrings(toolCmd, "test", "./..."),
		UseShell:  useShell,
		ShellName: shellName,
	}
	return []dukkha.TaskExecSpec{*spec}, nil
}
