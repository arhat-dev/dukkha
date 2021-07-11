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
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindBuild, t)
			return t
		},
	)
}

type TaskBuild buildah.TaskBud

// GetExecSpecs
// TODO: Handle manifests locally [#27](https://github.com/arhat-dev/dukkha/issues/27)
func (c *TaskBuild) GetExecSpecs(
	rc dukkha.TaskExecContext, options *dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		targets := c.ImageNames
		if len(targets) == 0 {
			targets = []buildah.ImageNameSpec{{
				Image:    c.TaskName,
				Manifest: "",
			}}
		}

		buildCmd := sliceutils.NewStrings(options.ToolCmd, "build")
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

		steps = append(steps, dukkha.TaskExecSpec{
			Command:   buildCmd,
			Env:       sliceutils.NewStrings(c.Env),
			UseShell:  options.UseShell,
			ShellName: options.ShellName,
		})

		return nil
	})

	return steps, err
}
