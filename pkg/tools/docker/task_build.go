package docker

import (
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const TaskKindBuild = "build"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindBuild,
		func(toolName string) dukkha.Task {
			t := &TaskBuild{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindBuild, t)
			return t
		},
	)
}

type TaskBuild buildah.TaskBuild

// GetExecSpecs
// TODO: Handle manifests locally [#27](https://github.com/arhat-dev/dukkha/issues/27)
func (c *TaskBuild) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec

	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		targets := c.ImageNames
		if len(targets) == 0 {
			targets = []*buildah.ImageNameSpec{{
				Image:    c.TaskName,
				Manifest: "",
			}}
		}

		buildCmd := []string{constant.DUKKHA_TOOL_CMD, "build"}
		for _, spec := range targets {
			if len(spec.Image) == 0 {
				continue
			}

			buildCmd = append(buildCmd, "-t", spec.Image)
		}

		if len(c.File) != 0 {
			buildCmd = append(buildCmd, "-f", c.File)
		}

		for _, bArg := range c.BuildArgs {
			buildCmd = append(buildCmd, "--build-arg", bArg)
		}

		buildCmd = append(buildCmd, c.ExtraArgs...)

		if len(c.Context) == 0 {
			buildCmd = append(buildCmd, ".")
		} else {
			buildCmd = append(buildCmd, c.Context)
		}

		steps = append(steps, dukkha.TaskExecSpec{
			Command: buildCmd,
		})

		return nil
	})

	return steps, err
}
