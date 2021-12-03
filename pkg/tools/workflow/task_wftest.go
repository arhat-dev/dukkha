package workflow

import (
	"fmt"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindTest = "test"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindTest,
		func(toolName string) dukkha.Task {
			t := &TaskTest{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
			return t
		},
	)
}

type TaskTest struct {
	rs.BaseField `yaml:"-"`

	TaskName string `yaml:"name"`

	tools.BaseTask `yaml:",inline"`
}

func (w *TaskTest) Kind() dukkha.TaskKind { return TaskKindTest }
func (w *TaskTest) Name() dukkha.TaskName { return dukkha.TaskName(w.TaskName) }

func (w *TaskTest) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: w.Kind(), Name: w.Name()}
}

func (w *TaskTest) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}
