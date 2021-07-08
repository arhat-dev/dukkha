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
			t.SetToolName(toolName)
			return t
		},
	)
}

var _ dukkha.Task = (*TaskLogin)(nil)

type TaskLogin buildah.TaskLogin

func (c *TaskLogin) ToolKind() dukkha.ToolKind { return ToolKind }
func (c *TaskLogin) Kind() dukkha.TaskKind     { return TaskKindLogin }

func (c *TaskLogin) GetExecSpecs(rc dukkha.RenderingContext, dockerCmd []string) ([]dukkha.TaskExecSpec, error) {
	loginCmd := sliceutils.NewStrings(
		dockerCmd, "login",
		"--username", c.Username,
		"--password-stdin",
	)

	password := c.Password + "\n"
	return []dukkha.TaskExecSpec{{
		Stdin:       strings.NewReader(password),
		Command:     append(loginCmd, c.Registry),
		IgnoreError: false,
	}}, nil
}
