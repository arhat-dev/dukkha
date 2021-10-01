package buildah

import (
	"fmt"
	"strconv"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindXBuild = "xbuild"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindXBuild,
		func(toolName string) dukkha.Task {
			t := &TaskXBuild{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindXBuild, t)
			return t
		},
	)
}

type TaskXBuild struct {
	rs.BaseField

	tools.BaseTask `yaml:",inline"`

	Context string  `yaml:"context"`
	Steps   []*step `yaml:"steps"`

	ImageNames []ImageNameSpec `yaml:"image_names"`
}

func (w *TaskXBuild) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var ret []dukkha.TaskExecSpec

	err := w.DoAfterFieldsResolved(rc, -1, func() error {
		globalCtx := &xbuildContext{
			ID: "",

			Steps: make(map[string]*xbuildContext),

			User:    "",
			Shell:   nil,
			WorkDir: "",
			Commit:  false,
		}

		for i, step := range w.Steps {
			stepCtx := globalCtx.clone()
			stepCtx.ID = step.ID
			if len(stepCtx.ID) == 0 {
				stepCtx.ID = strconv.FormatInt(int64(i), 10)
			}

			// add this step to global step index

			if _, ok := globalCtx.Steps[stepCtx.ID]; ok {
				return fmt.Errorf("invalid duplicate step id %q", stepCtx.ID)
			}

			globalCtx.Steps[stepCtx.ID] = stepCtx

			if step.Set != nil {
				// set global ctx
				globalCtx = step.Set.genCtx(stepCtx)
				globalCtx.Steps[stepCtx.ID] = stepCtx
			} else {
				if step.Workdir != nil {
					stepCtx.WorkDir = *step.Workdir
				}

				if step.Commit != nil {
					stepCtx.Commit = *step.Commit
				}

				if step.User != nil {
					stepCtx.User = *step.User
				}

				steps, err := step.genSpec(rc, options, stepCtx)
				if err != nil {
					return err
				}

				ret = append(ret, steps...)
			}

			// commit this container as image when
			// - at last step
			// - next step is a from step
			// - set commit=true explicitly
			if i == len(w.Steps)-1 || w.Steps[i+1].From != nil || stepCtx.Commit {

			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate build spec: %w", err)
	}

	return ret, nil
}

type xbuildContext struct {
	contextDir string

	ID    string
	Steps map[string]*xbuildContext

	User    string
	Shell   []string
	WorkDir string
	Commit  bool
}

func (c *xbuildContext) clone() *xbuildContext {
	return &xbuildContext{
		Steps: c.Steps,

		User:    c.User,
		Shell:   sliceutils.NewStrings(c.Shell),
		WorkDir: c.WorkDir,
		Commit:  c.Commit,
	}
}
