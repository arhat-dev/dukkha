package docker

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const TaskKindLogin = "login"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
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

var _ tools.Task = (*TaskLogin)(nil)

type TaskLogin buildah.TaskLogin

func (c *TaskLogin) ToolKind() string { return ToolKind }
func (c *TaskLogin) TaskKind() string { return TaskKindLogin }

func (c *TaskLogin) GetExecSpecs(ctx *field.RenderingContext, dockerCmd []string) ([]tools.TaskExecSpec, error) {
	loginCmd := sliceutils.NewStringSlice(
		dockerCmd, "login",
		"--username", c.Username,
		"--password-stdin",
	)

	password := c.Password + "\n"
	return []tools.TaskExecSpec{{
		Stdin:       strings.NewReader(password),
		Command:     loginCmd,
		IgnoreError: false,
	}}, nil
}
