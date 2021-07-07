package tools

import (
	"fmt"
	"sync"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/matrix"
	"arhat.dev/dukkha/pkg/types"
)

type BaseTask struct {
	field.BaseField

	TaskName string      `yaml:"name"`
	Matrix   matrix.Spec `yaml:"matrix"`
	Hooks    TaskHooks   `yaml:"hooks"`

	toolName dukkha.ToolName `yaml:"-"`

	hookMU sync.Mutex
}

func (t *BaseTask) ToolName() dukkha.ToolName { return t.toolName }
func (t *BaseTask) SetToolName(name string)   { t.toolName = dukkha.ToolName(name) }
func (t *BaseTask) Name() dukkha.TaskName     { return dukkha.TaskName(t.TaskName) }

func (t *BaseTask) RunHooks(taskCtx dukkha.TaskExecContext, stage dukkha.TaskExecStage) error {
	t.hookMU.Lock()
	defer t.hookMU.Unlock()

	err := t.Hooks.Run(taskCtx, stage)
	if err != nil {
		return fmt.Errorf("hook `%s` failed: %w", stage.String(), err)
	}

	return nil
}

func (t *BaseTask) GetMatrixSpecs(rc types.RenderingContext) ([]types.MatrixSpec, error) {
	err := t.ResolveFields(rc, -1, "Matrix")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve task matrix: %w", err)
	}

	return t.Matrix.GetSpecs(
		rc.MatrixFilter(),
		rc.HostKernel(),
		rc.HostArch(),
	), nil
}
