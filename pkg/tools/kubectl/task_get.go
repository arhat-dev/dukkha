package kubectl

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindGet = "get"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindGet,
		func(toolName string) dukkha.Task {
			t := &TaskGet{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindGet, t)
			return t
		},
	)
}

type TaskGet struct {
	rs.BaseField `yaml:"-"`

	tools.BaseTask `yaml:",inline"`
}

func (c *TaskGet) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	getCmd := sliceutils.NewStrings(options.ToolCmd(), "get")

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		return nil
	})

	return []dukkha.TaskExecSpec{{
		Command:   getCmd,
		UseShell:  options.UseShell(),
		ShellName: options.ShellName(),
	}}, err
}
