package buildah

import (
	"regexp"
	"strconv"
	"strings"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindLogin = "login"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^buildah(:.+){0,1}:login$`),
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

type TaskLogin struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Registry string `yaml:"registry"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	TLSSkipVerify *bool `yaml:"tls_skip_verify"`
}

func (c *TaskLogin) ToolKind() string { return ToolKind }
func (c *TaskLogin) TaskKind() string { return TaskKindLogin }

func (c *TaskLogin) GetExecSpecs(ctx *field.RenderingContext, buildahCmd []string) ([]tools.TaskExecSpec, error) {
	loginCmd := sliceutils.NewStrings(
		buildahCmd, "login",
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
	return []tools.TaskExecSpec{{
		Stdin:       strings.NewReader(password),
		Command:     append(loginCmd, c.Registry),
		IgnoreError: false,
	}}, nil
}
