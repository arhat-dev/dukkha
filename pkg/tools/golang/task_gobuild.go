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
		func(subMatches []string) interface{} {
			t := &TaskBuild{}
			if len(subMatches) != 0 {
				t.SetToolName(strings.TrimPrefix(subMatches[0], ":"))
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

	CGO struct {
		Enabled bool     `yaml:"enabled"`
		CFlags  []string `yaml:"cflags"`
		LDFlags []string `yaml:"ldflags"`

		HostCC  string `yaml:"hostCC"`
		HostCXX string `yaml:"hostCXX"`

		TargetCC  string `yaml:"targetCC"`
		TargetCXX string `yaml:"targetCXX"`
	} `yaml:"cgo"`
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

		if c.CGO.Enabled {
			spec.Env = append(spec.Env, "CGO_ENABLED=1")

			if len(c.CGO.CFlags) != 0 {
				spec.Env = append(spec.Env, fmt.Sprintf("CGO_CFLAGS=%s", strings.Join(c.CGO.CFlags, " ")))
			}

			if len(c.CGO.LDFlags) != 0 {
				spec.Env = append(spec.Env, fmt.Sprintf("CGO_LDFLAGS=%s", strings.Join(c.CGO.LDFlags, " ")))
			}

			if len(c.CGO.HostCC) != 0 {
				spec.Env = append(spec.Env, "CC="+c.CGO.HostCC)
			}

			if len(c.CGO.HostCXX) != 0 {
				spec.Env = append(spec.Env, "CXX="+c.CGO.HostCXX)
			}

			if len(c.CGO.TargetCC) != 0 {
				spec.Env = append(spec.Env, "CC_FOR_TARGET="+c.CGO.TargetCC)
			}

			if len(c.CGO.TargetCXX) != 0 {
				spec.Env = append(spec.Env, "CXX_FOR_TARGET="+c.CGO.TargetCXX)
			}
		} else {
			spec.Env = append(spec.Env, "CGO_ENABLED=0")
		}

		envGOOS := c.getGOOS(ctx.Values().Env[constant.ENV_MATRIX_KERNEL])
		spec.Env = append(spec.Env, "GOOS="+envGOOS)

		mArch := ctx.Values().Env[constant.ENV_MATRIX_ARCH]
		spec.Env = append(spec.Env, "GOARCH="+constant.GetGolangArch(mArch))

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

func (c *TaskBuild) getGOOS(mKernel string) string {
	return mKernel
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
