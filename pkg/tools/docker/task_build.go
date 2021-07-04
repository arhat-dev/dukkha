package docker

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/buildah"
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

type TaskBuild buildah.TaskBud

func (c *TaskBuild) ToolKind() string { return ToolKind }
func (c *TaskBuild) TaskKind() string { return TaskKindBuild }

// GetExecSpecs
// TODO: Handle manifests locally [#27](https://github.com/arhat-dev/dukkha/issues/27)
func (c *TaskBuild) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []buildah.ImageNameSpec{{
			Image:    c.Name,
			Manifest: "",
		}}
	}

	var (
		buildCmd = sliceutils.NewStrings(toolCmd, "build")
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
