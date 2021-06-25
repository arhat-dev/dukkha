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

	pushCmd     []string
	manifestCmd []string

	ImageNames []ImageNameSpec `yaml:"image_names"`
	ExtraArgs  []string        `yaml:"extraArgs"`
}

func (c *TaskPush) ToolKind() string { return ToolKind }
func (c *TaskPush) TaskKind() string { return TaskKindPush }

func (c *TaskPush) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []ImageNameSpec{
			{
				Image:    c.Name,
				Manifest: "",
			},
		}
	}

	var result []tools.TaskExecSpec
	for _, spec := range targets {
		if len(spec.Image) == 0 {
			continue
		}

		pushCmd := sliceutils.NewStringSlice(toolCmd, c.pushCmd...)
		if len(pushCmd) == len(toolCmd) {
			pushCmd = append(pushCmd, "push")
		}

		result = append(result, tools.TaskExecSpec{
			Command:     sliceutils.NewStringSlice(pushCmd, spec.Image),
			IgnoreError: false,
		})

		if len(spec.Manifest) == 0 {
			continue
		}

		manifestCmd := sliceutils.NewStringSlice(toolCmd, c.manifestCmd...)
		if len(manifestCmd) == len(toolCmd) {
			manifestCmd = append(manifestCmd, "manifest")
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

		mArch := ctx.Values().Env[constant.ENV_MATRIX_ARCH]
		annotateCmd := sliceutils.NewStringSlice(
			manifestCmd, "annotate", spec.Manifest, spec.Image,
			"--os", c.getManifestOS(ctx.Values().Env[constant.ENV_MATRIX_KERNEL]),
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
	}

	return result, nil
}

func (c *TaskPush) getManifestOS(os string) string {
	return strings.ToLower(os)
}
