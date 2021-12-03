package docker

import (
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const TaskKindLogin = "login"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindLogin,
		func(toolName string) dukkha.Task {
			t := &TaskLogin{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
			return t
		},
	)
}

type TaskLogin buildah.TaskLogin

func (w *TaskLogin) Kind() dukkha.TaskKind { return TaskKindLogin }
func (w *TaskLogin) Name() dukkha.TaskName { return dukkha.TaskName(w.TaskName) }
func (w *TaskLogin) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: w.Kind(), Name: w.Name()}
}

func (c *TaskLogin) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		loginCmd := []string{constant.DUKKHA_TOOL_CMD, "login",
			"--username", c.Username,
			"--password-stdin",
		}

		steps = append(steps, dukkha.TaskExecSpec{
			Stdin:   strings.NewReader(c.Password),
			Command: append(loginCmd, c.Registry),
		})

		return nil
	})

	return steps, err
}
