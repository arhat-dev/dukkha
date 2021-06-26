package docker

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindBuild = "build"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^docker(:.+){0,1}:build$`),
		func(params []string) interface{} {
			t := &TaskBuild{}
			if len(params) != 0 {
				t.SetToolName(strings.TrimPrefix(params[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskBuild)(nil)

type TaskBuild struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Context    string          `yaml:"context"`
	ImageNames []ImageNameSpec `yaml:"image_names"`
	Dockerfile string          `yaml:"dockerfile"`
	ExtraArgs  []string        `yaml:"extraArgs"`
}

type ImageNameSpec struct {
	Image    string `yaml:"image"`
	Manifest string `yaml:"manifest"`
}

func (c *TaskBuild) ToolKind() string { return ToolKind }
func (c *TaskBuild) TaskKind() string { return TaskKindBuild }

func (c *TaskBuild) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []ImageNameSpec{{
			Image:    c.Name,
			Manifest: "",
		}}
	}

	var (
		buildCmd = sliceutils.NewStringSlice(toolCmd, "build")
	)

	for _, spec := range targets {
		if len(spec.Image) == 0 {
			continue
		}

		buildCmd = append(buildCmd, "-t", spec.Image)
	}

	if len(c.Dockerfile) != 0 {
		buildCmd = append(buildCmd, "-f", c.Dockerfile)
	}

	buildCmd = append(buildCmd, c.ExtraArgs...)

	if len(c.Context) == 0 {
		buildCmd = append(buildCmd, ".")
	} else {
		buildCmd = append(buildCmd, c.Context)
	}

	return []tools.TaskExecSpec{
		{
			Command:     buildCmd,
			IgnoreError: false,
		},
	}, nil
}
