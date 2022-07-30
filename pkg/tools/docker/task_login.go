package docker

import (
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindLogin = "login"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindLogin, tools.NewTask[TaskLogin, *TaskLogin])
}

type TaskLogin struct {
	tools.BaseTask[DockerLogin, *DockerLogin]
}

// nolint:revive
type DockerLogin struct {
	Registry string `yaml:"registry"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	TLSSkipVerify *bool `yaml:"tls_skip_verify"`

	parent tools.BaseTaskType
}

func (c *DockerLogin) ToolKind() dukkha.ToolKind       { return ToolKind }
func (c *DockerLogin) Kind() dukkha.TaskKind           { return TaskKindLogin }
func (c *DockerLogin) LinkParent(p tools.BaseTaskType) { c.parent = p }

func (c *DockerLogin) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
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
