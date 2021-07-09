package tools

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/matrix"
)

var _ dukkha.Task = (*_baseTaskWithGetExecSpecs)(nil)

type _baseTaskWithGetExecSpecs struct {
	BaseTask
}

func (b *_baseTaskWithGetExecSpecs) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return nil, nil
}

type BaseTask struct {
	field.BaseField

	TaskName string      `yaml:"name"`
	Matrix   matrix.Spec `yaml:"matrix"`
	Hooks    TaskHooks   `yaml:"hooks"`

	toolName dukkha.ToolName `yaml:"-"`
	toolKind dukkha.ToolKind `yaml:"-"`
	taskKind dukkha.TaskKind `yaml:"-"`
}

func (t *BaseTask) InitBaseTask(k dukkha.ToolKind, n dukkha.ToolName, tk dukkha.TaskKind) {
	t.toolKind = k
	t.toolName = n

	t.taskKind = tk
}

func (t *BaseTask) Kind() dukkha.TaskKind {
	return t.taskKind
}

func (t *BaseTask) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: t.taskKind, Name: dukkha.TaskName(t.TaskName)}
}

func (t *BaseTask) ToolName() dukkha.ToolName { return t.toolName }
func (t *BaseTask) ToolKind() dukkha.ToolKind { return t.toolKind }

func (t *BaseTask) Name() dukkha.TaskName { return dukkha.TaskName(t.TaskName) }

func (t *BaseTask) GetHookExecSpecs(
	taskCtx dukkha.TaskExecContext,
	stage dukkha.TaskExecStage,
	options dukkha.TaskExecOptions,
) ([]dukkha.RunTaskOrRunShell, error) {
	specs, err := t.Hooks.GenSpecs(taskCtx, stage, options)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate exec specs for hook %q: %w",
			stage.String(), err,
		)
	}

	return specs, nil
}

func (t *BaseTask) GetMatrixSpecs(rc dukkha.RenderingContext) ([]matrix.Entry, error) {
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
