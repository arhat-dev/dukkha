package kubectl

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindDelete = "delete"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindDelete,
		func(toolName string) dukkha.Task {
			t := &TaskDelete{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindDelete, t)
			return t
		},
	)
}

type TaskDelete struct {
	rs.BaseField `yaml:"-"`

	tools.BaseTask `yaml:",inline"`
}

func (c *TaskDelete) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	deleteCmd := sliceutils.NewStrings(options.ToolCmd(), "delete")

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		return nil
	})

	return []dukkha.TaskExecSpec{{
		Command:   deleteCmd,
		UseShell:  options.UseShell(),
		ShellName: options.ShellName(),
	}}, err
}
