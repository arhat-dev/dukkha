package debug

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/tools"
	"gopkg.in/yaml.v3"
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
		mCtx, mExecOpts, err := tools.CreateTaskMatrixContext(&tools.TaskExecRequest{
			Context: tskCtx,
			Tool:    tool,
			Task:    tsk,
		}, ms, execOpts)
		_ = mExecOpts
		if err != nil {
			return nil, fmt.Errorf("failed to create task matrix context: %w", err)
		}

		err = tsk.DoAfterFieldsResolved(mCtx, -1, func() error {
			err := enc.Encode(tsk)
			if err != nil {
				return err
			}

			// handle go-yaml inline field issue with custom marshaler implementation
			// get the inline fields value first
			tskVal := reflect.ValueOf(tsk)
			tskTyp := tskVal.Type()
			for tskTyp.Kind() != reflect.Struct {
				tskTyp = tskTyp.Elem()
			}

			for tskVal.Kind() != reflect.Struct {
				tskVal = tskVal.Elem()
			}

			// var inlineValues []interface{}
			for i := 0; i < tskTyp.NumField(); i++ {
				f := tskTyp.Field(i)
				yTags := strings.Split(f.Tag.Get("yaml"), ",")
				for _, tg := range yTags {
					if tg == "inline" {
						data, err := yaml.Marshal(tskVal.Field(i).Interface())
						if err != nil {
							return fmt.Errorf("failed to marshal inline value of field %q: %w", f.Name, err)
						}

						_, err = os.Stdout.Write(data)
						if err != nil {
							return fmt.Errorf("failed to write inline value of field %q: %w", f.Name, err)
						}

						break
					}
				}
			}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate resolved yaml: %w", err)
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
