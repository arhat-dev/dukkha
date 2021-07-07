package docker

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools/buildah"
	"arhat.dev/dukkha/pkg/types"
)

const TaskKindLogin = "login"

func init() {
	field.RegisterInterfaceField(
		dukkha.TaskType,
		regexp.MustCompile(`^docker(:.+){0,1}:login$`),
		func(params []string) interface{} {
			t := &TaskLogin{}
			if len(params) != 0 {
				t.SetToolName(strings.TrimPrefix(params[0], ":"))
			}
			return t
		},
	)
}

var _ dukkha.Task = (*TaskLogin)(nil)

type TaskLogin buildah.TaskLogin

func (c *TaskLogin) ToolKind() dukkha.ToolKind { return ToolKind }
func (c *TaskLogin) Kind() dukkha.TaskKind     { return TaskKindLogin }

func (c *TaskLogin) GetExecSpecs(rc types.RenderingContext, dockerCmd []string) ([]dukkha.TaskExecSpec, error) {
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
