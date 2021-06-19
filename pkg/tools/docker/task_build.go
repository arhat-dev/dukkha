package docker

import (
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindBuild = "docker:build"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^docker(:\w+)?:build$`),
		func() interface{} { return &TaskBuild{} },
	)
}

var _ tools.Task = (*TaskBuild)(nil)

type TaskBuild struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	ImageName    string   `yaml:"image_name"`
	ManifestName string   `yaml:"manifest_name"`
	Dockerfile   string   `yaml:"dockerfile"`
	ExtraArgs    []string `yaml:"extraArgs"`
}

func (c *TaskBuild) Kind() string { return TaskKindBuild }
