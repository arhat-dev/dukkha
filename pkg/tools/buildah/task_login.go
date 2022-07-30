package buildah

import (
	"strconv"
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
	tools.BaseTask[BuildahLogin, *BuildahLogin]
}

type BuildahLogin struct {
	Registry string `yaml:"registry"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	TLSSkipVerify *bool `yaml:"tls_skip_verify"`

	parent tools.BaseTaskType
}

func (w *BuildahLogin) ToolKind() dukkha.ToolKind       { return ToolKind }
func (w *BuildahLogin) Kind() dukkha.TaskKind           { return TaskKindLogin }
func (w *BuildahLogin) LinkParent(p tools.BaseTaskType) { w.parent = p }

func (c *BuildahLogin) GetExecSpecs(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec

	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
		loginCmd := []string{constant.DUKKHA_TOOL_CMD, "login",
			"--username", c.Username,
			"--password-stdin",
		}

		if c.TLSSkipVerify != nil {
			loginCmd = append(
				loginCmd,
				"--tls-verify", strconv.FormatBool(!*c.TLSSkipVerify),
			)
		}

		steps = append(steps, dukkha.TaskExecSpec{
			Stdin:       strings.NewReader(c.Password),
			Command:     append(loginCmd, c.Registry),
			IgnoreError: false,
		})

		return nil
	})

	return steps, err
}
