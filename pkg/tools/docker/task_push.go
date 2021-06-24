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

		pushCmd := sliceutils.NewStringSlice(toolCmd, "push")
		result = append(result, tools.TaskExecSpec{
			Command:     sliceutils.NewStringSlice(pushCmd, spec.Image),
			IgnoreError: false,
		})

		if len(spec.Manifest) == 0 {
			continue
		}

		manifestCmd := sliceutils.NewStringSlice(toolCmd, "manifest")
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
			"--arch", c.getManifestArch(mArch),
		)

		variant := c.getManifestArchVariant(mArch)
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

func (c *TaskPush) getManifestArch(arch string) string {
	arch = strings.ToLower(arch)

	mArch := map[string]string{
		constant.ARCH_X86:   "386",
		constant.ARCH_AMD64: "amd64",

		constant.ARCH_ARM_V5: "arm",
		constant.ARCH_ARM_V6: "arm",
		constant.ARCH_ARM_V7: "arm",

		constant.ARCH_ARM64: "arm64",

		constant.ARCH_MIPS64_LE:    "mips64le",
		constant.ARCH_MIPS64_LE_HF: "mips64le",

		constant.ARCH_PPC64_LE: "ppc64le",

		constant.ARCH_S390X: "s390x",
	}[arch]

	if len(mArch) == 0 {
		mArch = arch
	}

	return mArch
}

func (c *TaskPush) getManifestArchVariant(arch string) string {
	arch = strings.ToLower(arch)

	variant := map[string]string{
		constant.ARCH_ARM_V5: "v5",
		constant.ARCH_ARM_V6: "v6",
		constant.ARCH_ARM_V7: "v7",
		constant.ARCH_ARM64:  "v8",
	}[arch]

	return variant
}
