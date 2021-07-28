package golang

import (
	"strings"

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
	ExtraArgs []string `yaml:"extra_args"`
	Outputs   []string `yaml:"outputs"`

	BuildOptions buildOptions `yaml:",inline"`

	CGO CGOSepc `yaml:"cgo"`
}

func (c *TaskBuild) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var buildSteps []dukkha.TaskExecSpec

	err := c.DoAfterFieldsResolved(rc, -1, func() error {
		mKernel := rc.MatrixKernel()
		mArch := rc.MatrixArch()

		buildEnv := append(c.CGO.getEnv(
			rc.HostKernel() != mKernel || rc.HostArch() != mArch,
			mKernel, mArch,
			rc.HostOS(),
			rc.MatrixLibc(),
		), createBuildEnv(mKernel, mArch)...)

		outputs := sliceutils.NewStrings(c.Outputs)
		if len(outputs) == 0 {
			outputs = []string{c.BaseTask.TaskName}
		}

		for _, output := range outputs {
			spec := &dukkha.TaskExecSpec{
				Chdir: c.Chdir,

				// put generated env first, so user can override them
				EnvSuggest: buildEnv,
				Command:    sliceutils.NewStrings(options.ToolCmd(), "build", "-o", output),
				UseShell:   options.UseShell(),
				ShellName:  options.ShellName(),
			}

			spec.Command = append(spec.Command, c.BuildOptions.generateArgs(options.UseShell())...)
			spec.Command = append(spec.Command, c.ExtraArgs...)

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

func formatArgs(args []string, useShell bool) string {
	ret := strings.Join(args, " ")
	if useShell {
		ret = `"` + ret + `"`
	}

	return ret
}
