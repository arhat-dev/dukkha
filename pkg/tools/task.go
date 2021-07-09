package tools

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

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
	Env      []string    `yaml:"env"`
	Matrix   matrix.Spec `yaml:"matrix"`
	Hooks    TaskHooks   `yaml:"hooks"`

	toolName dukkha.ToolName `yaml:"-"`
	toolKind dukkha.ToolKind `yaml:"-"`
	taskKind dukkha.TaskKind `yaml:"-"`

	fieldsToResolve []string
	impl            dukkha.Task

	mu sync.Mutex
}

func (t *BaseTask) DoAfterFieldsResolved(
	mCtx dukkha.RenderingContext,
	depth int,
	do func() error,
	fieldNames ...string,
) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(fieldNames) == 0 {
		for _, name := range []string{"TaskName", "Env"} {
			err := t.ResolveFields(mCtx, depth, name)
			if err != nil {
				return fmt.Errorf(
					"failed to resolve basic task field %q: %w",
					name, err,
				)
			}
		}

		// all fields, including hooks
		for _, name := range t.fieldsToResolve {
			err := t.impl.ResolveFields(mCtx, depth, name)
			if err != nil {
				return fmt.Errorf(
					"failed to resolve task field %q from default list: %w",
					name, err,
				)
			}
		}
	} else {
		for _, name := range fieldNames {
			target := field.Field(t.impl)
			targetField := name
			if strings.HasPrefix(targetField, "BaseTask.") {
				targetField = strings.TrimPrefix(targetField, "BaseTask.")
				target = t
			}

			err := target.ResolveFields(mCtx, depth, targetField)
			if err != nil {
				return fmt.Errorf(
					"failed to resolve specified task field %q: %w",
					name, err,
				)
			}
		}
	}

	return do()
}

func (t *BaseTask) InitBaseTask(
	k dukkha.ToolKind,
	n dukkha.ToolName,
	tk dukkha.TaskKind,
	impl dukkha.Task,
) {
	t.toolKind = k
	t.toolName = n

	t.taskKind = tk

	t.impl = impl

	typ := reflect.TypeOf(impl).Elem()
	for i := 1; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if !(f.Name[0] >= 'A' && f.Name[0] <= 'Z') {
			// unexported, ignore
			continue
		}

		if f.Anonymous && f.Name == "BaseTask" {
			// it's me
			continue
		}

		val, ok := f.Tag.Lookup("yaml")
		if !ok {
			// no yaml field, ignore
			continue
		}

		if strings.Contains(val, "-") {
			// ignored explicitly
			continue
		}

		t.fieldsToResolve = append(t.fieldsToResolve, f.Name)
	}
}

func (t *BaseTask) ToolKind() dukkha.ToolKind { return t.toolKind }
func (t *BaseTask) ToolName() dukkha.ToolName { return t.toolName }
func (t *BaseTask) Kind() dukkha.TaskKind     { return t.taskKind }
func (t *BaseTask) Name() dukkha.TaskName     { return dukkha.TaskName(t.TaskName) }

func (t *BaseTask) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: t.taskKind, Name: dukkha.TaskName(t.TaskName)}
}

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
