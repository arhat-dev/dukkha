package tools

import (
	"fmt"
	"reflect"
	"sync"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/matrix"
)

var _ dukkha.Task = (*_baseTaskWithGetExecSpecs)(nil)

type _baseTaskWithGetExecSpecs struct{ BaseTask }

func (b *_baseTaskWithGetExecSpecs) GetExecSpecs(
	rc dukkha.TaskExecContext, options *dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return nil, nil
}

type BaseTask struct {
	field.BaseField

	TaskName string      `yaml:"name"`
	Env      []string    `yaml:"env"`
	Matrix   matrix.Spec `yaml:"matrix"`
	Hooks    TaskHooks   `yaml:"hooks"`

	ContinueOnErrorFlag bool `yaml:"continue_on_error"`

	// fields managed by BaseTask

	toolName dukkha.ToolName `yaml:"-"`
	toolKind dukkha.ToolKind `yaml:"-"`
	taskKind dukkha.TaskKind `yaml:"-"`

	fieldsToResolve []string
	impl            dukkha.Task

	mu sync.Mutex
}

func (t *BaseTask) resolveEssentialFieldsAndAddEnv(mCtx dukkha.RenderingContext) error {
	err := resolveFields(mCtx, t, -1, []string{"TaskName", "Env"})
	if err != nil {
		return fmt.Errorf("failed to resolve essential task fields: %w", err)
	}

	mCtx.AddEnv(t.Env...)

	return nil
}

func (t *BaseTask) DoAfterFieldsResolved(
	ctx dukkha.RenderingContext,
	depth int,
	do func() error,
	fieldNames ...string,
) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	err := t.resolveEssentialFieldsAndAddEnv(ctx)
	if err != nil {
		return err
	}

	if len(fieldNames) == 0 {
		// resolve all fields of the real task type
		err := resolveFields(ctx, t.impl, depth, t.fieldsToResolve)
		if err != nil {
			return fmt.Errorf("failed to resolve tool fields: %w", err)
		}
	} else {
		forBase, forImpl := separateBaseAndImpl("BaseTask.", fieldNames)
		if len(forBase) != 0 {
			err := resolveFields(ctx, t, depth, forBase)
			if err != nil {
				return fmt.Errorf("failed to resolve requested BaseTask fields: %w", err)
			}
		}

		if len(forImpl) != 0 {
			err := resolveFields(ctx, t.impl, depth, forImpl)
			if err != nil {
				return fmt.Errorf("failed to resolve requested fields: %w", err)
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
	t.fieldsToResolve = getFieldNamesToResolve(typ)
}

func (t *BaseTask) ToolKind() dukkha.ToolKind { return t.toolKind }
func (t *BaseTask) ToolName() dukkha.ToolName { return t.toolName }
func (t *BaseTask) Kind() dukkha.TaskKind     { return t.taskKind }
func (t *BaseTask) Name() dukkha.TaskName     { return dukkha.TaskName(t.TaskName) }

func (t *BaseTask) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: t.taskKind, Name: dukkha.TaskName(t.TaskName)}
}

func (t *BaseTask) ContinueOnError() bool {
	return t.ContinueOnErrorFlag
}

func (t *BaseTask) GetHookExecSpecs(
	taskCtx dukkha.TaskExecContext,
	stage dukkha.TaskExecStage,
) ([]dukkha.RunTaskOrRunCmd, error) {

	t.mu.Lock()
	defer t.mu.Unlock()

	// hooks may have reference to env defined in task scope

	err := t.resolveEssentialFieldsAndAddEnv(taskCtx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to prepare env for hook %q: %w",
			stage.String(), err,
		)
	}

	specs, err := t.Hooks.GenSpecs(taskCtx, stage)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate exec specs for hook %q: %w",
			stage.String(), err,
		)
	}

	return specs, nil
}

func (t *BaseTask) GetMatrixSpecs(rc dukkha.RenderingContext) ([]matrix.Entry, error) {
	var ret []matrix.Entry
	err := t.DoAfterFieldsResolved(rc, -1, func() error {
		ret = t.Matrix.GenerateEntries(
			rc.MatrixFilter(),
			rc.HostKernel(),
			rc.HostArch(),
		)

		return nil
	},
		// t.DoAfterFieldsResolved is intended to serve
		// real task type, so we have to add the prefix
		// `BaseTask.`
		"BaseTask.Matrix",
	)

	return ret, err
}
