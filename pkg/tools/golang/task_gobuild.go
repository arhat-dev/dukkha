package golang

import (
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindBuild = "build"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindBuild, tools.NewTask[TaskBuild, *TaskBuild])
}

type TaskBuild struct {
	tools.BaseTask[GolangBuild, *GolangBuild]
}

type GolangBuild struct {
	// Chdir into a different dir when running go command while keep `dukkha.WorkDir` unchanged
	// this can be helpful when you are managing multiple go modules in one workspace
	Chdir string `yaml:"chdir"`

	// Path is the import path of the source code to be built
	// it will be the last argument of this go command execution
	//
	// see `go help packages` for more information about import path
	Path string `yaml:"path"`

	// Outputs of the go build command, when multiple entries specified, will build multiple times with same
	// arguments
	Outputs []string `yaml:"outputs"`

	BuildOptions buildOptions `yaml:",inline"`

	// CGo options
	CGo CGOSepc `yaml:"cgo"`

	// ExtraArgs for go build (inserted before `Path`)
	ExtraArgs []string `yaml:"extra_args"`

	parent tools.BaseTaskType
}

func (w *GolangBuild) ToolKind() dukkha.ToolKind       { return ToolKind }
func (w *GolangBuild) Kind() dukkha.TaskKind           { return TaskKindBuild }
func (w *GolangBuild) LinkParent(p tools.BaseTaskType) { w.parent = p }

func (c *GolangBuild) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var buildSteps []dukkha.TaskExecSpec

	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
		outputs := sliceutils.NewStrings(c.Outputs)
		if len(outputs) == 0 {
			outputs = []string{string(c.parent.Name())}
		}

		buildEnv := createBuildEnv(rc, c.BuildOptions, c.CGo)
		for _, output := range outputs {
			spec := &dukkha.TaskExecSpec{
				Chdir: c.Chdir,

				// put generated env first, so user can override them
				EnvSuggest: buildEnv,
				Command:    []string{constant.DUKKHA_TOOL_CMD, "build", "-o", output},
			}

			spec.Command = append(spec.Command, c.BuildOptions.generateArgs()...)
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
