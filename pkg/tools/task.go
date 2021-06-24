package tools

import (
	"fmt"
	"reflect"

	"arhat.dev/dukkha/pkg/field"
)

// TaskType for interface type registration
var TaskType = reflect.TypeOf((*Task)(nil)).Elem()

type TaskExecSpec struct {
	Env         []string
	Command     []string
	IgnoreError bool
}

type Task interface {
	field.Interface

	// Kind of the tool managing this task (e.g. docker)
	ToolKind() string

	// Name of the tool managing this task (e.g. my-tool)
	ToolName() string

	// Kind of the task (e.g. build)
	TaskKind() string

	// Name of the task
	TaskName() string

	// GetMatrixSpecs for matrix build
	GetMatrixSpecs(ctx *field.RenderingContext, rf field.RenderingFunc) ([]MatrixSpec, error)

	// GetExecSpecs generate commands using current field values
	GetExecSpecs(ctx *field.RenderingContext, toolCmd []string) ([]TaskExecSpec, error)
}

type BaseTask struct {
	field.BaseField

	Name   string        `yaml:"name"`
	Matrix *MatrixConfig `yaml:"matrix"`

	toolName string `yaml:"-"`
}

func (t *BaseTask) ToolName() string        { return t.toolName }
func (t *BaseTask) SetToolName(name string) { t.toolName = name }
func (t *BaseTask) TaskName() string        { return t.Name }

func (t *BaseTask) GetMatrixSpecs(ctx *field.RenderingContext, rf field.RenderingFunc) ([]MatrixSpec, error) {
	// resolve matrix config first
	err := t.ResolveFields(ctx, rf, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve base task fields: %w", err)
	}

	err = t.Matrix.ResolveFields(ctx, rf, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve task matrix: %w", err)
	}

	return t.Matrix.GetSpecs(), nil
}
