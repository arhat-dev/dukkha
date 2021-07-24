package docker

import (
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools/buildah"
)

const TaskKindPush = "push"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindPush,
		func(toolName string) dukkha.Task {
			t := &TaskPush{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindPush, t)
			return t
		},
	)
}

type TaskPush buildah.TaskPush

func (c *TaskPush) GetExecSpecs(
	rc dukkha.TaskExecContext, options *dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var result []dukkha.TaskExecSpec

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		targets := c.ImageNames
		if len(targets) == 0 {
			targets = []buildah.ImageNameSpec{{
				Image:    c.TaskName,
				Manifest: "",
			}}
		}

		var (
			manifestCmd = sliceutils.NewStrings(options.ToolCmd, "manifest")
		)

		for _, spec := range targets {
			if len(spec.Image) == 0 {
				continue
			}

			imageName := buildah.SetDefaultImageTagIfNoTagSet(rc, spec.Image)
			// docker push <image-name>
			if buildah.ImageOrManifestHasFQDN(imageName) {
				result = append(result, dukkha.TaskExecSpec{
					Env: sliceutils.NewStrings(c.Env),
					Command: sliceutils.NewStrings(
						options.ToolCmd, "push", imageName,
					),
					IgnoreError: false,
					UseShell:    options.UseShell,
					ShellName:   options.ShellName,
				})
			}

			if len(spec.Manifest) == 0 {
				continue
			}

			manifestName := buildah.SetDefaultManifestTagIfNoTagSet(rc, spec.Manifest)
			result = append(result,
				// ensure manifest exists
				dukkha.TaskExecSpec{
					Env: sliceutils.NewStrings(c.Env),
					Command: sliceutils.NewStrings(
						manifestCmd, "create", manifestName, imageName,
					),
					// may already exists
					IgnoreError: true,
					UseShell:    options.UseShell,
					ShellName:   options.ShellName,
				},
				// link manifest and image
				dukkha.TaskExecSpec{
					Command: sliceutils.NewStrings(
						manifestCmd, "create", manifestName,
						"--amend", imageName,
					),
					IgnoreError: false,
					UseShell:    options.UseShell,
					ShellName:   options.ShellName,
				},
			)

			// docker manifest annotate \
			// 		<manifest-list-name> <image-name> \
			// 		--os <arch> --arch <arch> {--variant <variant>}
			mArch := rc.MatrixArch()
			os, _ := constant.GetDockerOS(rc.MatrixKernel())
			arch, _ := constant.GetDockerArch(mArch)
			annotateCmd := sliceutils.NewStrings(
				manifestCmd, "annotate", manifestName, imageName,
				"--os", os, "--arch", arch,
			)

			variant, _ := constant.GetDockerArchVariant(mArch)
			if len(variant) != 0 {
				annotateCmd = append(annotateCmd, "--variant", variant)
			}

			result = append(result, dukkha.TaskExecSpec{
				Env:         sliceutils.NewStrings(c.Env),
				Command:     annotateCmd,
				IgnoreError: false,
				UseShell:    options.UseShell,
				ShellName:   options.ShellName,
			})

			// docker manifest push <manifest-list-name>
			if buildah.ImageOrManifestHasFQDN(manifestName) {
				result = append(result, dukkha.TaskExecSpec{
					Env:         sliceutils.NewStrings(c.Env),
					Command:     sliceutils.NewStrings(options.ToolCmd, "manifest", "push", spec.Manifest),
					IgnoreError: false,
					UseShell:    options.UseShell,
					ShellName:   options.ShellName,
				})
			}
		}

		return nil
	})

	return result, err
}
