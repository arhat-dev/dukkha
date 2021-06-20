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
		regexp.MustCompile(`^docker(:.+)?:build$`),
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

func (c *TaskBuild) ExecArgs() ([]string, error) {
	args := []string{
		"build",
	}

	if len(c.Dockerfile) != 0 {
		args = append(args, "-f", c.Dockerfile)
	}

	if len(c.ImageNames) == 0 {
		args = append(args, "-t", c.Name)
	} else {
		for _, imgName := range c.ImageNames {
			args = append(args, "-t", imgName.Image)
		}
	}

	args = append(args, c.ExtraArgs...)

	if len(c.BuildContext) == 0 {
		args = append(args, ".")
	} else {
		args = append(args, c.BuildContext)
	}

	return args, nil
}
