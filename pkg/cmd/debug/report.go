package debug

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/tools"
)

type reportOptions struct {
	cliOptions

	matrixFilter *matrix.Filter
}

func (ropts *reportOptions) generateTaskReport(
	appCtx dukkha.Context,
	tool dukkha.Tool,
	tsk dukkha.Task,
) (*taskReport, error) {
	// TODO: implement
	appCtx = appCtx.DeriveNew()

	appCtx.SetMatrixFilter(ropts.matrixFilter)

	matrixSpecs, err := tsk.GetMatrixSpecs(appCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get task matrix specs: %w", err)
	}

	execOpts := dukkha.CreateTaskExecOptions(0, len(matrixSpecs))
	tskCtx := appCtx.DeriveNew()
	tskCtx.SetTask(tool.Key(), tsk.Key())

	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	defer func() { _ = enc.Close() }()

	for _, ms := range matrixSpecs {
		mCtx, mExecOpts, err2 := tools.CreateTaskMatrixContext(&tools.TaskExecRequest{
			Context: tskCtx,
			Tool:    tool,
			Task:    tsk,
		}, ms, execOpts)
		_ = mExecOpts
		if err2 != nil {
			return nil, fmt.Errorf("failed to create task matrix context: %w", err2)
		}

		err2 = tsk.DoAfterFieldsResolved(mCtx, -1, func() error {
			return enc.Encode(tsk)
		})
		if err2 != nil {
			return nil, fmt.Errorf("failed to generate resolved yaml: %w", err2)
		}
	}

	return &taskReport{
		Name: tsk.Name(),
	}, nil
}

type taskReport struct {
	Name dukkha.TaskName `yaml:"name"`

	MatrixRun []struct {
		Spec []map[string]string `yaml:"-"`
	} ``
}

// func (r *taskReport) MarshalYAML() (interface{}, error) {
// 	entYaml := "{ " + strings.Join(sliceutils.FormatStringMap(ent, ": ", false), ", ") + " }"
// }
