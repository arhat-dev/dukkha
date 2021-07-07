package docker

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools/buildah"
	"arhat.dev/dukkha/pkg/types"
)

const TaskKindBuild = "build"

func init() {
	field.RegisterInterfaceField(
		dukkha.TaskType,
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

var _ dukkha.Task = (*TaskBuild)(nil)

type TaskBuild buildah.TaskBud

func (c *TaskBuild) ToolKind() dukkha.ToolKind { return ToolKind }
func (c *TaskBuild) Kind() dukkha.TaskKind     { return TaskKindBuild }

// GetExecSpecs
// TODO: Handle manifests locally [#27](https://github.com/arhat-dev/dukkha/issues/27)
func (c *TaskBuild) GetExecSpecs(rc types.RenderingContext, toolCmd []string) ([]dukkha.TaskExecSpec, error) {
	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []buildah.ImageNameSpec{{
			Image:    c.TaskName,
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

	return []dukkha.TaskExecSpec{
		{
			Command:     buildCmd,
			IgnoreError: false,
		},
	}, nil
}
