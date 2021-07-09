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
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindLogin)
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
	rc dukkha.RenderingContext,
	useShell bool,
	shellName string,
	buildahCmd []string,
) ([]dukkha.TaskExecSpec, error) {
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
	return []dukkha.TaskExecSpec{
		{
			Stdin:       strings.NewReader(password),
			Command:     append(loginCmd, c.Registry),
			IgnoreError: false,
			UseShell:    useShell,
			ShellName:   shellName,
		},
	}, nil
}
