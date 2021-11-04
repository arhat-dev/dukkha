package kubectl

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindLogs = "logs"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindLogs,
		func(toolName string) dukkha.Task {
			t := &TaskLogs{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindLogs, t)
			return t
		},
	)
}

type TaskLogs struct {
	rs.BaseField `yaml:"-"`

	tools.BaseTask `yaml:",inline"`

	TargetCluster kubeContext `yaml:"target"`
}

func (c *TaskLogs) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	logsCmd := sliceutils.NewStrings(options.ToolCmd(), "logs")

	err := c.DoAfterFieldsResolved(rc, -1, func() error {

		return nil
	})

	return []dukkha.TaskExecSpec{{
		Command:   logsCmd,
		UseShell:  options.UseShell(),
		ShellName: options.ShellName(),
	}}, err
}
