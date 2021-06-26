package docker

import (
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindPush = "push"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^docker(:.+){0,1}:push$`),
		func(params []string) interface{} {
			t := &TaskPush{}
			if len(params) != 0 {
				t.SetToolName(strings.TrimPrefix(params[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskPush)(nil)

type TaskPush struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	ImageNames []ImageNameSpec `yaml:"image_names"`
}

func (c *TaskPush) ToolKind() string { return ToolKind }
func (c *TaskPush) TaskKind() string { return TaskKindPush }

func (c *TaskPush) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []ImageNameSpec{{
			Image:    c.Name,
			Manifest: "",
		}}
	}

	var (
		result []tools.TaskExecSpec

		manifestCmd = sliceutils.NewStringSlice(toolCmd, "manifest")
	)

	for _, spec := range targets {
		if len(spec.Image) == 0 {
			continue
		}

		// docker push <image-name>
		result = append(result, tools.TaskExecSpec{
			Command:     sliceutils.NewStringSlice(toolCmd, "push", spec.Image),
			IgnoreError: false,
		})

		if len(spec.Manifest) == 0 {
			continue
		}

		result = append(result,
			// ensure manifest exists
			tools.TaskExecSpec{
				Command:     sliceutils.NewStringSlice(manifestCmd, "create", spec.Manifest, spec.Image),
				IgnoreError: true,
			},
			// link manifest and image
			tools.TaskExecSpec{
				Command:     sliceutils.NewStringSlice(manifestCmd, "create", spec.Manifest, "--amend", spec.Image),
				IgnoreError: false,
			},
		)

		// docker manifest annotate \
		// 		<manifest-list-name> <image-name> \
		// 		--os <arch> --arch <arch> {--variant <variant>}
		mArch := ctx.Values().Env[constant.ENV_MATRIX_ARCH]
		annotateCmd := sliceutils.NewStringSlice(
			manifestCmd, "annotate", spec.Manifest, spec.Image,
			"--os", constant.GetDockerOS(ctx.Values().Env[constant.ENV_MATRIX_KERNEL]),
			"--arch", constant.GetDockerArch(mArch),
		)

		variant := constant.GetDockerArchVariant(mArch)
		if len(variant) != 0 {
			annotateCmd = append(annotateCmd, "--variant", variant)
		}

		result = append(result, tools.TaskExecSpec{
			Command:     annotateCmd,
			IgnoreError: false,
		})

		// docker manifest push <manifest-list-name>
		result = append(result, tools.TaskExecSpec{
			Command:     sliceutils.NewStringSlice(toolCmd, "manifest", "push", spec.Manifest),
			IgnoreError: false,
		})
	}

	return result, nil
}
