package kubectl

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindDiff = "diff"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindDiff,
		func(toolName string) dukkha.Task {
			t := &TaskDiff{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindDiff, t)
			return t
		},
	)
}

type TaskDiff struct {
	rs.BaseField `yaml:"-"`

	tools.BaseTask `yaml:",inline"`
}

func (c *TaskDiff) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	diffCmd := sliceutils.NewStrings(options.ToolCmd(), "diff")

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		return nil
	})

	return []dukkha.TaskExecSpec{{
		Command:   diffCmd,
		UseShell:  options.UseShell(),
		ShellName: options.ShellName(),
	}}, err
}
