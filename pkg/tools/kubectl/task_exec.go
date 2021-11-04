package kubectl

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindExec = "exec"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindExec,
		func(toolName string) dukkha.Task {
			t := &TaskExec{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindExec, t)
			return t
		},
	)
}

type TaskExec struct {
	rs.BaseField `yaml:"-"`

	tools.BaseTask `yaml:",inline"`
}

func (c *TaskExec) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	execCmd := sliceutils.NewStrings(options.ToolCmd(), "exec")

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		return nil
	})

	return []dukkha.TaskExecSpec{{
		Command:   execCmd,
		UseShell:  options.UseShell(),
		ShellName: options.ShellName(),
	}}, err
}
