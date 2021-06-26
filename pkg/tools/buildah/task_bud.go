package buildah

import (
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
	targets := c.ImageNames
	if len(targets) == 0 {
		targets = []ImageNameSpec{{
			Image:    c.Name,
			Manifest: "",
		}}
	}

	var (
		budCmd      = sliceutils.NewStringSlice(toolCmd, "bud")
		manifestCmd = sliceutils.NewStringSlice(toolCmd, "manifest")

		taskExecAfterBudCmd []tools.TaskExecSpec
	)

	for _, spec := range targets {
		if len(spec.Image) == 0 {
			continue
		}

		budCmd = append(budCmd, "-t", spec.Image)

		if len(spec.Manifest) == 0 {
			// no manifest handling
			continue
		}

		taskExecAfterBudCmd = append(taskExecAfterBudCmd,
			// ensure manifest exists
			tools.TaskExecSpec{
				Command: sliceutils.NewStringSlice(manifestCmd, "create", spec.Manifest, spec.Image),
				// can already exist
				IgnoreError: true,
			},
			// add image to manifest
			tools.TaskExecSpec{
				Command: sliceutils.NewStringSlice(manifestCmd, "add", spec.Manifest, spec.Image),
				// can already added
				IgnoreError: true,
			},
		)

		// buildah manifest annotate --os <arch> --arch <arch> {--variant <variant>} \
		// 		<manifest-list-name> <image-name>
		mArch := ctx.Values().Env[constant.ENV_MATRIX_ARCH]
		annotateCmd := sliceutils.NewStringSlice(
			manifestCmd, "annotate",
			"--os", ctx.Values().Env[constant.ENV_MATRIX_KERNEL],
			"--arch", constant.GetOCIArch(mArch),
		)

		variant := constant.GetOCIArchVariant(mArch)
		if len(variant) != 0 {
			annotateCmd = append(annotateCmd, "--variant", variant)
		}

		taskExecAfterBudCmd = append(taskExecAfterBudCmd, tools.TaskExecSpec{
			Command:     append(annotateCmd, spec.Manifest, spec.Image),
			IgnoreError: false,
		})
	}

	if len(c.Dockerfile) != 0 {
		budCmd = append(budCmd, "-f", c.Dockerfile)
	}

	budCmd = append(budCmd, c.ExtraArgs...)

	if len(c.Context) == 0 {
		budCmd = append(budCmd, ".")
	} else {
		budCmd = append(budCmd, c.Context)
	}

	return append([]tools.TaskExecSpec{
		{
			Command:     budCmd,
			IgnoreError: false,
		},
	}, taskExecAfterBudCmd...), nil
}
