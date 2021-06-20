package docker

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindPush = "push"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^docker(:.+)?:push$`),
		func(params []string) interface{} {
			t := &TaskPush{}
			if len(params) != 0 {
				t.SetToolName(strings.TrimPrefix(params[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskPush)(nil)

type TaskPush struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	ImageNames []ImageNameSpec `yaml:"image_names"`
	ExtraArgs  []string        `yaml:"extraArgs"`
}

func (c *TaskPush) ToolKind() string { return ToolKind }
func (c *TaskPush) TaskKind() string { return TaskKindPush }

func (c *TaskPush) Inherit(bc *TaskBuild) {
	if bc == nil {
		return
	}

	if c.Matrix == nil {
		c.Matrix = bc.Matrix
	}

	if len(c.ImageNames) == 0 {
		c.ImageNames = bc.ImageNames
	}
}
