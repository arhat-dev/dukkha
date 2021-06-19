package docker

import (
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindPush = "docker:push"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^docker(:\w+)?:push$`),
		func() interface{} { return &TaskPush{} },
	)
}

var _ tools.Task = (*TaskPush)(nil)

type TaskPush struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	ImageName    string   `yaml:"image_name"`
	ManifestName string   `yaml:"manifest_name"`
	ExtraArgs    []string `yaml:"extraArgs"`
}

func (c *TaskPush) Kind() string { return TaskKindPush }

func (c *TaskPush) Inherit(bc *TaskBuild) {
	if bc == nil {
		return
	}

	if c.Matrix == nil {
		c.Matrix = bc.Matrix
	}

	if len(c.ImageName) == 0 {
		c.ImageName = bc.ImageName
	}

	if len(c.ManifestName) == 0 {
		c.ManifestName = bc.ManifestName
	}
}
