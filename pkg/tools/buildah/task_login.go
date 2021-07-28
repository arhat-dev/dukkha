package buildah

import (
	"strconv"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
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
	field.BaseField

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

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		loginCmd := sliceutils.NewStrings(
			options.ToolCmd(), "login",
			"--username", c.Username,
			"--password-stdin",
		)

		if c.TLSSkipVerify != nil {
			loginCmd = append(
				loginCmd,
				"--tls-verify", strconv.FormatBool(!*c.TLSSkipVerify),
			)
		}

		password := c.Password + "\n"

		steps = append(steps, dukkha.TaskExecSpec{
			Stdin:       strings.NewReader(password),
			Command:     append(loginCmd, c.Registry),
			IgnoreError: false,
			UseShell:    options.UseShell(),
			ShellName:   options.ShellName(),
		})

		return nil
	})

	return steps, err
}
