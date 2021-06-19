package golang

import (
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindBuild = "golang:build"

func init() {
	field.RegisterInterfaceField(
		tools.TaskConfigType,
		regexp.MustCompile(`^golang(:\w+)?:build$`),
		func() interface{} { return &TaskBuildConfig{} },
	)
}

var _ tools.TaskConfig = (*TaskBuildConfig)(nil)

type TaskBuildConfig struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Matrix BuildMatrixConfig `yaml:"matrix"`
}

type BuildMatrixConfig struct {
	field.BaseField

	tools.BaseMatrixConfig `yaml:",inline"`
}

func (c *TaskBuildConfig) Kind() string { return TaskKindBuild }
