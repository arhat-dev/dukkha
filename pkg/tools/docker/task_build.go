package docker

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const TaskKindBuild = "build"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindBuild,
		func(toolName string) dukkha.Task {
			t := &TaskBuild{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindBuild)
			return t
		},
	)
}

type TaskBuild buildah.TaskBud

// GetExecSpecs
// TODO: Handle manifests locally [#27](https://github.com/arhat-dev/dukkha/issues/27)
func (c *TaskBuild) GetExecSpecs(
	rc dukkha.TaskExecContext,
	useShell bool,
	shellName string,
	toolCmd []string,
) ([]dukkha.TaskExecSpec, error) {
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
			UseShell:    useShell,
			ShellName:   shellName,
		},
	}, nil
}
