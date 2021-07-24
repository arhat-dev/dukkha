package golang

import (
	"fmt"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindBuild = "build"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindBuild,
		func(toolName string) dukkha.Task {
			t := &TaskBuild{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindBuild, t)
			return t
		},
	)
}

type TaskBuild struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Chdir     string   `yaml:"chdir"`
	Path      string   `yaml:"path"`
	LDFlags   []string `yaml:"ldflags"`
	Tags      []string `yaml:"tags"`
	ExtraArgs []string `yaml:"extra_args"`
	Outputs   []string `yaml:"outputs"`

	CGO CGOSepc `yaml:"cgo"`
}

func (c *TaskBuild) GetExecSpecs(
	rc dukkha.TaskExecContext, options *dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	mKernel := rc.MatrixKernel()
	mArch := rc.MatrixArch()

	var buildSteps []dukkha.TaskExecSpec

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		env := sliceutils.NewStrings(
			c.CGO.getEnv(
				rc.HostKernel() != mKernel || rc.HostArch() != mArch,
				mKernel, mArch,
				rc.HostOS(),
				rc.MatrixLibc(),
			),
		)

		goos := constant.GetGolangOS(mKernel)
		if len(goos) != 0 {
			env = append(env, "GOOS="+goos)
		}

		goarch := constant.GetGolangArch(mArch)
		if len(goarch) != 0 {
			env = append(env, "GOARCH="+goarch)
		}

		if envGOMIPS := c.getGOMIPS(mArch); len(envGOMIPS) != 0 {
			env = append(env, "GOMIPS="+envGOMIPS, "GOMIPS64="+envGOMIPS)
		}

		if envGOARM := c.getGOARM(mArch); len(envGOARM) != 0 {
			env = append(env, "GOARM="+envGOARM)
		}

		outputs := sliceutils.NewStrings(c.Outputs)
		if len(outputs) == 0 {
			outputs = []string{c.BaseTask.TaskName}
		}

		for _, output := range outputs {
			spec := &dukkha.TaskExecSpec{
				Chdir: c.Chdir,

				// put generated env first, so user can override them
				Env:       sliceutils.NewStrings(env, c.Env...),
				Command:   sliceutils.NewStrings(options.ToolCmd, "build", "-o", output),
				UseShell:  options.UseShell,
				ShellName: options.ShellName,
			}

			spec.Command = append(spec.Command, c.ExtraArgs...)

			if len(c.LDFlags) != 0 {
				spec.Command = append(
					spec.Command,
					"-ldflags",
					formatArgs(c.LDFlags, options.UseShell),
				)
			}

			if len(c.Tags) != 0 {
				spec.Command = append(
					spec.Command,
					"-tags",
					// ref: https://golang.org/doc/go1.13#go-command
					// The go build flag -tags now takes a comma-separated list of build tags,
					// to allow for multiple tags in GOFLAGS. The space-separated form is
					// deprecated but still recognized and will be maintained.
					strings.Join(c.Tags, ","),
				)
			}

			if len(c.Path) != 0 {
				spec.Command = append(spec.Command, c.Path)
			} else {
				spec.Command = append(spec.Command, "./")
			}

			buildSteps = append(buildSteps, *spec)
		}
		return nil
	})

	return buildSteps, err
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

	if strings.HasSuffix(mArch, "sf") {
		return "softfloat"
	}

	return "hardfloat"
}

type CGOSepc struct {
	field.BaseField

	// Enable cgo
	Enabled bool `yaml:"enabled"`

	// CGO_CPPFLAGS
	CPPFlags []string `yaml:"cppflags"`

	// CGO_CFLAGS
	CFlags []string `yaml:"cflags"`

	// CGO_CXXFLAGS
	CXXFlags []string `yaml:"cxxflags"`

	// CGO_FFLAGS
	FFlags []string `yaml:"fflags"`

	// CGO_LDFLAGS
	LDFlags []string `yaml:"ldflags"`

	CC  string `yaml:"cc"`
	CXX string `yaml:"cxx"`
}

func (c *CGOSepc) getEnv(
	doingCrossCompiling bool,
	mKernel, mArch, hostOS, targetLibc string,
) []string {
	if !c.Enabled {
		return []string{"CGO_ENABLED=0"}
	}

	var ret []string
	ret = append(ret, "CGO_ENABLED=1")

	appendListEnv := func(name string, values, defVals []string) {
		actual := values
		if len(values) == 0 {
			actual = defVals
		}

		if len(actual) != 0 {
			ret = append(ret, fmt.Sprintf("%s=%s", name, strings.Join(actual, " ")))
		}
	}

	appendEnv := func(name, value, defVal string) {
		actual := value
		if len(value) == 0 {
			actual = defVal
		}

		if len(actual) != 0 {
			ret = append(ret, fmt.Sprintf("%s=%s", name, actual))
		}
	}

	var (
		cppflags []string
		cflags   []string
		cxxflags []string
		fflags   []string
		ldflags  []string

		cc  = "gcc"
		cxx = "g++"
	)

	if hostOS == constant.OS_MACOS {
		cc = "clang"
		cxx = "clang++"
	}

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

			cc = tripleName + "-gcc"
			cxx = tripleName + "-g++"
		case constant.OS_ALPINE:
			tripleName := constant.GetAlpineTripleName(mArch)
			cc = tripleName + "-gcc"
			cxx = tripleName + "-g++"
		case constant.OS_MACOS:
			cc = "clang"
			cxx = "clang++"
		}
	}

	// TODO: generate suitable flags
	appendListEnv("CGO_CPPFLAGS", c.CPPFlags, cppflags)
	appendListEnv("CGO_CFLAGS", c.CFlags, cflags)
	appendListEnv("CGO_CXXFLAGS", c.CXXFlags, cxxflags)
	appendListEnv("CGO_FFLAGS", c.FFlags, fflags)
	appendListEnv("CGO_LDFLAGS", c.LDFlags, ldflags)

	appendEnv("CC", c.CC, cc)
	appendEnv("CXX", c.CXX, cxx)

	return ret
}

func formatArgs(args []string, useShell bool) string {
	ret := strings.Join(args, " ")
	if useShell {
		ret = `"` + ret + `"`
	}

	return ret
}
