package golang

import (
	"fmt"
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindBuild = "build"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^golang(:.+)?:build$`),
		func(params []string) interface{} {
			t := &TaskBuild{}
			if len(params) != 0 {
				t.SetToolName(params[0])
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskBuild)(nil)

type TaskBuild struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Path      string   `yaml:"path"`
	Env       []string `yaml:"env"`
	LDFlags   []string `yaml:"ldflags"`
	Tags      []string `yaml:"tags"`
	ExtraArgs []string `yaml:"extraArgs"`
	Outputs   []string `yaml:"outputs"`
}

func (c *TaskBuild) ToolKind() string { return ToolKind }
func (c *TaskBuild) TaskKind() string { return TaskKindBuild }

func (c *TaskBuild) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	outputs := c.Outputs
	if len(outputs) == 0 {
		outputs = []string{c.Name}
	}

	var ret []tools.TaskExecSpec
	for _, output := range outputs {
		spec := &tools.TaskExecSpec{}

		spec.Env = append(spec.Env, c.Env...)

		envGOOS := c.getGOOS(strings.ToLower(ctx.Values().Env[constant.ENV_MATRIX_OS]))
		spec.Env = append(spec.Env, "GOOS="+envGOOS)

		mArch := strings.ToLower(ctx.Values().Env[constant.ENV_MATRIX_ARCH])
		spec.Env = append(spec.Env, "GOARCH="+c.getGOARCH(mArch))

		switch {
		case strings.HasPrefix(mArch, "mips"):
			spec.Env = append(spec.Env, "GOMIPS="+c.getGOMIPS(mArch))
		case strings.HasPrefix(mArch, "armv"):
			spec.Env = append(spec.Env, "GOARM="+c.getGOARM(mArch))
		}

		spec.Command = append(spec.Command, c.ExtraArgs...)
		spec.Command = sliceutils.NewStringSlice(toolCmd, "build", "-o", output)

		if len(c.LDFlags) != 0 {
			spec.Command = append(spec.Command, fmt.Sprintf("-ldflags=\"%s\"", strings.Join(c.LDFlags, " ")))
		}

		if len(c.Tags) != 0 {
			spec.Command = append(spec.Command, fmt.Sprintf("-tags=\"%s\"", strings.Join(c.Tags, " ")))
		}

		if len(c.Path) != 0 {
			spec.Command = append(spec.Command, c.Path)
		} else {
			spec.Command = append(spec.Command, ".")
		}

		ret = append(ret, *spec)
	}

	return ret, nil
}

func (c *TaskBuild) getGOOS(mOS string) string {
	return mOS
}

func (c *TaskBuild) getGOARCH(mArch string) string {
	goArch := map[string]string{
		constant.ARCH_X86:          "386",
		constant.ARCH_AMD64:        "amd64",
		constant.ARCH_ARM64:        "arm64",
		constant.ARCH_ARM_V5:       "arm",
		constant.ARCH_ARM_V6:       "arm",
		constant.ARCH_ARM_V7:       "arm",
		constant.ARCH_MIPS:         "mips",
		constant.ARCH_MIPS_HF:      "mips",
		constant.ARCH_MIPS_LE:      "mipsle",
		constant.ARCH_MIPS_LE_HF:   "mipsle",
		constant.ARCH_MIPS64:       "mips64",
		constant.ARCH_MIPS64_HF:    "mips64",
		constant.ARCH_MIPS64_LE:    "mips64le",
		constant.ARCH_MIPS64_LE_HF: "mips64le",
		// constant.ARCH_PPC:          "ppc",
		constant.ARCH_PPC64:    "ppc64",
		constant.ARCH_PPC64LE:  "ppc64le",
		constant.ARCH_RISCV_64: "riscv64",
		constant.ARCH_S390X:    "s390x",
		constant.ARCH_IA64:     "ia64",
	}[mArch]

	return goArch
}

func (c *TaskBuild) getGOARM(mArch string) string {
	return strings.TrimPrefix(mArch, "armv")
}

func (c *TaskBuild) getGOMIPS(mArch string) string {
	if strings.HasSuffix(mArch, "hf") {
		return "hardfloat"
	}

	return "softfloat"
}
