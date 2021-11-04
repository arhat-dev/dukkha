package kubectl

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindApply = "apply"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindApply,
		func(toolName string) dukkha.Task {
			t := &TaskApply{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindApply, t)
			return t
		},
	)
}

type TaskApply struct {
	rs.BaseField `yaml:"-"`

	tools.BaseTask `yaml:",inline"`
}

func (c *TaskApply) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	applyCmd := sliceutils.NewStrings(options.ToolCmd(), "apply")

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		return nil
	})

	return []dukkha.TaskExecSpec{{
		Command:   applyCmd,
		UseShell:  options.UseShell(),
		ShellName: options.ShellName(),
	}}, err
}
