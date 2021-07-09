package docker

import (
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools/buildah"
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

type TaskLogin buildah.TaskLogin

func (c *TaskLogin) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		loginCmd := sliceutils.NewStrings(
			options.ToolCmd, "login",
			"--username", c.Username,
			"--password-stdin",
		)

		password := c.Password + "\n"

		steps = append(steps, dukkha.TaskExecSpec{
			Stdin:       strings.NewReader(password),
			Command:     append(loginCmd, c.Registry),
			IgnoreError: options.ContinueOnError,
			UseShell:    options.UseShell,
			ShellName:   options.ShellName,
		})

		return nil
	})

	return steps, err
}
