package buildah

import (
	"fmt"
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindBud = "bud"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^buildah(:.+){0,1}:bud$`),
		func(params []string) interface{} {
			t := &TaskBud{}
			if len(params) != 0 {
				t.SetToolName(strings.TrimPrefix(params[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskBud)(nil)

type TaskBud struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Context    string          `yaml:"context"`
	ImageNames []ImageNameSpec `yaml:"image_names"`
	Dockerfile string          `yaml:"dockerfile"`
	ExtraArgs  []string        `yaml:"extraArgs"`
}

type ImageNameSpec struct {
	Image    string `yaml:"image"`
	Manifest string `yaml:"manifest"`
}

func (c *TaskBud) ToolKind() string { return ToolKind }
func (c *TaskBud) TaskKind() string { return TaskKindBud }

func (c *TaskBud) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	budCmd := sliceutils.NewStringSlice(toolCmd, "bud")
	if len(c.Dockerfile) != 0 {
		budCmd = append(budCmd, "-f", c.Dockerfile)
	}

	budCmd = append(budCmd, c.ExtraArgs...)

	context := c.Context
	if len(context) == 0 {
		context = "."
	}

	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []ImageNameSpec{{
			Image:    c.Name,
			Manifest: "",
		}}
	}
	var result []tools.TaskExecSpec
	for _, spec := range targets {
		if len(spec.Manifest) == 0 {
			continue
		}

		// TODO: buildah only allows one --manifest for each bud run
		// 	    so we have to build multiple times for the same image with different
		// 		manifest name
		// 		need to find a way to use buildah manifest create & add correctly
		result = append(result, tools.TaskExecSpec{
			Command: sliceutils.NewStringSlice(budCmd,
				"-t", spec.Image,
				"--manifest", spec.Manifest, context,
			),
			IgnoreError: false,
		})

		// NOTE: buildah will treat --os and --arch values to bud as pull target
		// 	     which is not desierd in most use cases, especially when cross compiling
		//
		// so we MUST update manifest os/arch/variant after bud
		annotateCmd := sliceutils.NewStringSlice(toolCmd, "manifest", "annotate")
		mArch := ctx.Values().Env[constant.ENV_MATRIX_ARCH]
		annotateCmd = append(annotateCmd,
			"--os", constant.GetOciOS(ctx.Values().Env[constant.ENV_MATRIX_KERNEL]),
			"--arch", constant.GetOciArch(mArch),
		)

		variant := constant.GetOciArchVariant(mArch)
		if len(variant) != 0 {
			annotateCmd = append(annotateCmd, "--variant", variant)
		}

		annotateCmd = append(annotateCmd, spec.Manifest)
		annotateCmd = append(annotateCmd,
			fmt.Sprintf("$(%s)",
				strings.Join(
					sliceutils.NewStringSlice(
						toolCmd,
						"inspect", "--type", "image",
						"--format", `"{{ .FromImageDigest }}"`,
						spec.Image,
					),
					" ",
				),
			),
		)

		result = append(result, tools.TaskExecSpec{
			Command:     annotateCmd,
			IgnoreError: false,
		})
	}

	if len(result) == 0 {
		// no manifest set, build image without handling manifests
		for _, spec := range targets {
			budCmd = append(budCmd, "-t", spec.Image)
		}

		result = append(result, tools.TaskExecSpec{
			Command:     append(budCmd, context),
			IgnoreError: false,
		})
	}

	return result, nil
}
