package tools

import (
	"fmt"
	"reflect"
	"sync"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
)

var _ dukkha.Task = (*_baseTaskWithGetExecSpecs)(nil)

type _baseTaskWithGetExecSpecs struct{ BaseTask }

func (b *_baseTaskWithGetExecSpecs) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return nil, nil
}

type BaseTask struct {
	rs.BaseField `yaml:"-"`

	TaskName string      `yaml:"name"`
	Env      dukkha.Env  `yaml:"env"`
	Matrix   matrix.Spec `yaml:"matrix"`
	Hooks    TaskHooks   `yaml:"hooks,omitempty"`

	ContinueOnErrorFlag bool `yaml:"continue_on_error"`

	// fields managed by BaseTask

	toolName dukkha.ToolName `yaml:"-"`
	toolKind dukkha.ToolKind `yaml:"-"`
	taskKind dukkha.TaskKind `yaml:"-"`

	fieldsToResolve []string
	impl            dukkha.Task

	mu sync.Mutex
}

func (t *BaseTask) DoAfterFieldsResolved(
	ctx dukkha.RenderingContext,
	depth int,
	resolveEnv bool,
	do func() error,
	tagNames ...string,
) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if resolveEnv {
		err := dukkha.ResolveEnv(t, ctx, "Env", "env")
		if err != nil {
			return err
		}
	}

	if len(tagNames) == 0 {
		// resolve all fields of the real task type
		err := t.impl.ResolveFields(ctx, depth, t.fieldsToResolve...)
		if err != nil {
			return fmt.Errorf("failed to resolve tool fields: %w", err)
		}
	} else {
		forBase, forImpl := separateBaseAndImpl("BaseTask.", tagNames)
		if len(forBase) != 0 {
			err := t.ResolveFields(ctx, depth, forBase...)
			if err != nil {
				return fmt.Errorf("failed to resolve requested BaseTask fields: %w", err)
			}
		}

		if len(forImpl) != 0 {
			err := t.impl.ResolveFields(ctx, depth, forImpl...)
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
	t.fieldsToResolve = getTagNamesToResolve(typ)
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
) ([]dukkha.TaskExecSpec, error) {

	t.mu.Lock()
	defer t.mu.Unlock()

	// hooks may have reference to env defined in task scope

	err := dukkha.ResolveEnv(t, taskCtx, "Env", "env")
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
	err := t.DoAfterFieldsResolved(rc, -1, true, func() error {
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
		"BaseTask.matrix",
	)

	return ret, err
}
