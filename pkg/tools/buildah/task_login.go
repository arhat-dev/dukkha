package buildah

import (
	"strconv"
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindLogin = "login"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindLogin,
		func(toolName string) dukkha.Task {
			t := &TaskLogin{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindLogin, t)
			return t
		},
	)
}

type TaskLogin struct {
	rs.BaseField `yaml:"-"`

	tools.BaseTask `yaml:",inline"`

	Registry string `yaml:"registry"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	TLSSkipVerify *bool `yaml:"tls_skip_verify"`
}

func (c *TaskLogin) GetExecSpecs(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec

	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
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
