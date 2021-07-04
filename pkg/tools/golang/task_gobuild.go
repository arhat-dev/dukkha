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

	Chdir     string   `yaml:"chdir"`
	Path      string   `yaml:"path"`
	Env       []string `yaml:"env"`
	LDFlags   []string `yaml:"ldflags"`
	Tags      []string `yaml:"tags"`
	ExtraArgs []string `yaml:"extra_args"`
	Outputs   []string `yaml:"outputs"`

	CGO CGOSepc `yaml:"cgo"`
}

type CGOSepc struct {
	field.BaseField

	Enabled bool     `yaml:"enabled"`
	CFlags  []string `yaml:"cflags"`
	LDFlags []string `yaml:"ldflags"`

	HostCC  string `yaml:"hostCC"`
	HostCXX string `yaml:"hostCXX"`

	TargetCC  string `yaml:"targetCC"`
	TargetCXX string `yaml:"targetCXX"`
}

func (c *TaskBuild) ToolKind() string { return ToolKind }
func (c *TaskBuild) TaskKind() string { return TaskKindBuild }

func (c *TaskBuild) GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]tools.TaskExecSpec, error) {
	mKernel := ctx.Values().Env[constant.ENV_MATRIX_KERNEL]
	mArch := ctx.Values().Env[constant.ENV_MATRIX_ARCH]
	doingCrossCompiling := ctx.Values().Env[constant.ENV_HOST_KERNEL] != mKernel ||
		ctx.Values().Env[constant.ENV_HOST_ARCH] != mArch

	env := sliceutils.NewStringSlice(c.Env, c.CGO.getEnv(
		doingCrossCompiling, mKernel, mArch,
		ctx.Values().Env[constant.ENV_HOST_OS],
		ctx.Values().Env[constant.ENV_MATRIX_LIBC],
	)...)

	env = append(env, "GOOS="+constant.GetGolangOS(mKernel))
	env = append(env, "GOARCH="+constant.GetGolangArch(mArch))

	if envGOMIPS := c.getGOMIPS(mArch); len(envGOMIPS) != 0 {
		env = append(env, "GOMIPS="+envGOMIPS)
	}

	if envGOARM := c.getGOARM(mArch); len(envGOARM) != 0 {
		env = append(env, "GOARM="+envGOARM)
	}

	outputs := c.Outputs
	if len(outputs) == 0 {
		outputs = []string{c.Name}
	}

	var buildSteps []tools.TaskExecSpec
	for _, output := range outputs {
		spec := &tools.TaskExecSpec{
			Chdir: c.Chdir,

			Env:     sliceutils.NewStringSlice(env),
			Command: sliceutils.NewStringSlice(toolCmd, "build", "-o", output),
		}

		spec.Command = append(spec.Command, c.ExtraArgs...)

		if len(c.LDFlags) != 0 {
			spec.Command = append(spec.Command,
				fmt.Sprintf("-ldflags=\"%s\"", strings.Join(c.LDFlags, " ")),
			)
		}

		if len(c.Tags) != 0 {
			spec.Command = append(spec.Command,
				fmt.Sprintf("-tags=\"%s\"", strings.Join(c.Tags, " ")),
			)
		}

		if len(c.Path) != 0 {
			spec.Command = append(spec.Command, c.Path)
		} else {
			spec.Command = append(spec.Command, ".")
		}

		buildSteps = append(buildSteps, *spec)
	}

	return buildSteps, nil
}

func (c *TaskBuild) getGOARM(mArch string) string {
	if strings.HasPrefix(mArch, "armv") {
		return strings.TrimPrefix(mArch, "armv")
	}

	return ""
}

func (c *TaskBuild) getGOMIPS(mArch string) string {
	if !strings.HasPrefix(mArch, "mips") {
		return ""
	}

	if strings.HasSuffix(mArch, "hf") {
		return "hardfloat"
	}

	return "softfloat"
}

func (c *CGOSepc) getEnv(doingCrossCompiling bool, mKernel, mArch, hostOS, targetLibc string) []string {
	if !c.Enabled {
		return []string{"CGO_ENABLED=0"}
	}

	var ret []string
	ret = append(ret, "CGO_ENABLED=1")

	if len(c.CFlags) != 0 {
		ret = append(ret, fmt.Sprintf("CGO_CFLAGS=%s", strings.Join(c.CFlags, " ")))
	}

	if len(c.LDFlags) != 0 {
		ret = append(ret, fmt.Sprintf("CGO_LDFLAGS=%s", strings.Join(c.LDFlags, " ")))
	}

	if len(c.HostCC) != 0 {
		ret = append(ret, "CC="+c.HostCC)
	}

	if len(c.HostCXX) != 0 {
		ret = append(ret, "CXX="+c.HostCXX)
	}

	targetCC := "gcc"
	targetCXX := "g++"
	if doingCrossCompiling {
		switch hostOS {
		case constant.OS_DEBIAN,
			constant.OS_UBUNTU:
			var tripleName string
			switch mKernel {
			case constant.KERNEL_LINUX:
				tripleName = constant.GetDebianTripleName(mArch, mKernel, targetLibc)
			case constant.KERNEL_DARWIN:
				// TODO: set darwin version
				tripleName = constant.GetAppleTripleName(mArch, "")
			case constant.KERNEL_WINDOWS:
				tripleName = constant.GetDebianTripleName(mArch, mKernel, targetLibc)
			default:
			}

			targetCC = tripleName + "-gcc"
			targetCXX = tripleName + "-g++"
		case constant.OS_ALPINE:
			tripleName := constant.GetAlpineTripleName(mArch)
			targetCC = tripleName + "-gcc"
			targetCXX = tripleName + "-g++"
		case constant.OS_MACOS:
			targetCC = "clang"
			targetCXX = "clang++"
		}
	}

	if len(c.TargetCC) != 0 {
		ret = append(ret, "CC_FOR_TARGET="+c.TargetCC)
	} else if doingCrossCompiling {
		ret = append(ret, "CC_FOR_TARGET="+targetCC)
	}

	if len(c.TargetCXX) != 0 {
		ret = append(ret, "CXX_FOR_TARGET="+c.TargetCXX)
	} else if doingCrossCompiling {
		ret = append(ret, "CC_FOR_TARGET="+targetCXX)
	}

	return ret
}
