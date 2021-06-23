package docker

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
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

	BuildContext string          `yaml:"build_context"`
	ImageNames   []ImageNameSpec `yaml:"image_names"`
	Dockerfile   string          `yaml:"dockerfile"`
	ExtraArgs    []string        `yaml:"extraArgs"`
}

type ImageNameSpec struct {
	Image    string `yaml:"image"`
	Manifest string `yaml:"manifest"`
}

func (c *TaskBuild) ToolKind() string { return ToolKind }
func (c *TaskBuild) TaskKind() string { return TaskKindBuild }

func (c *TaskBuild) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	var cmd []string

	cmd = append(cmd, toolCmd...)
	cmd = append(cmd, "build")

	if len(c.Dockerfile) != 0 {
		cmd = append(cmd, "-f", c.Dockerfile)
	}

	if len(c.ImageNames) == 0 {
		cmd = append(cmd, "-t", c.Name)
	} else {
		for _, imgName := range c.ImageNames {
			cmd = append(cmd, "-t", imgName.Image)
		}
	}

	cmd = append(cmd, c.ExtraArgs...)

	if len(c.BuildContext) == 0 {
		cmd = append(cmd, ".")
	} else {
		cmd = append(cmd, c.BuildContext)
	}

	return []tools.TaskExecSpec{
		{
			Command:     cmd,
			IgnoreError: false,
		},
	}, nil
}
