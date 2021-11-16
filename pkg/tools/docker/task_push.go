package docker

import (
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/templateutils"
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
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var result []dukkha.TaskExecSpec

	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		targets := c.ImageNames
		if len(targets) == 0 {
			targets = []buildah.ImageNameSpec{{
				Image:    c.TaskName,
				Manifest: "",
			}}
		}

		var (
			manifestCmd = []string{constant.DUKKHA_TOOL_CMD, "manifest"}
		)

		for _, spec := range targets {
			if len(spec.Image) == 0 {
				continue
			}

			imageName := templateutils.SetDefaultImageTagIfNoTagSet(
				rc, spec.Image, false,
			)
			// docker push <image-name>
			if imageOrManifestHasFQDN(imageName) {
				result = append(result, dukkha.TaskExecSpec{
					Command:     []string{constant.DUKKHA_TOOL_CMD, "push", imageName},
					IgnoreError: false,
				})
			}

			if len(spec.Manifest) == 0 {
				continue
			}

			manifestName := templateutils.SetDefaultManifestTagIfNoTagSet(
				rc, spec.Manifest,
			)
			result = append(result,
				// ensure manifest exists
				dukkha.TaskExecSpec{
					Command: sliceutils.NewStrings(
						manifestCmd, "create", manifestName, imageName,
					),
					// may already exists
					IgnoreError: true,
				},
				// link manifest and image
				dukkha.TaskExecSpec{
					Command: sliceutils.NewStrings(
						manifestCmd, "create", manifestName,
						"--amend", imageName,
					),
					IgnoreError: false,
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
				Command:     annotateCmd,
				IgnoreError: false,
			})

			// docker manifest push <manifest-list-name>
			if imageOrManifestHasFQDN(manifestName) {
				result = append(result, dukkha.TaskExecSpec{
					Command:     []string{constant.DUKKHA_TOOL_CMD, "manifest", "push", spec.Manifest},
					IgnoreError: false,
				})
			}
		}

		return nil
	})

	return result, err
}

func imageOrManifestHasFQDN(s string) bool {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) == 1 {
		return false
	}

	return strings.Contains(parts[0], ".")
}
