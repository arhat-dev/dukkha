package tools

import (
	"fmt"
	"reflect"
	"sync"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
)

var _ dukkha.Task = (*_baseTaskWithGetExecSpecs)(nil)

type _baseTaskWithGetExecSpecs struct{ BaseTask }

func (b *_baseTaskWithGetExecSpecs) Kind() dukkha.TaskKind { return "_" }
func (b *_baseTaskWithGetExecSpecs) Name() dukkha.TaskName { return "_" }
func (b *_baseTaskWithGetExecSpecs) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: b.Kind(), Name: b.Name()}
}

func (b *_baseTaskWithGetExecSpecs) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return nil, nil
}

type BaseTask struct {
	rs.BaseField `yaml:"-"`

	Env    dukkha.Env  `yaml:"env"`
	Matrix matrix.Spec `yaml:"matrix"`
	Hooks  TaskHooks   `yaml:"hooks,omitempty"`

	ContinueOnErrorFlag bool `yaml:"continue_on_error"`

	// fields managed by BaseTask

	CacheFS *fshelper.OSFS `yaml:"-"`

	toolName dukkha.ToolName `yaml:"-"`
	toolKind dukkha.ToolKind `yaml:"-"`

	tagsToResolve []string

	impl dukkha.Task

	mu sync.Mutex
}

func (t *BaseTask) Init(cacheFS *fshelper.OSFS) error {
	t.CacheFS = cacheFS
	return nil
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
		err := dukkha.ResolveEnv(ctx, t, "Env", "env")
		if err != nil {
			return err
		}
	}

	if len(tagNames) == 0 {
		// resolve all fields of the real task type
		err := t.impl.ResolveFields(ctx, depth, t.tagsToResolve...)
		if err != nil {
			return fmt.Errorf("resolving tool fields: %w", err)
		}
	} else {
		forBase, forImpl := separateBaseAndImpl("BaseTask.", tagNames)
		if len(forBase) != 0 {
			err := t.ResolveFields(ctx, depth, forBase...)
			if err != nil {
				return fmt.Errorf("resolving requested BaseTask fields: %w", err)
			}
		}

		if len(forImpl) != 0 {
			err := t.impl.ResolveFields(ctx, depth, forImpl...)
			if err != nil {
				return fmt.Errorf("resolving requested fields: %w", err)
			}
		}
	}

	return do()
}

func (t *BaseTask) InitBaseTask(
	k dukkha.ToolKind,
	n dukkha.ToolName,
	impl dukkha.Task,
) {
	t.toolKind = k
	t.toolName = n

	t.impl = impl

	typ := reflect.TypeOf(impl).Elem()
	t.tagsToResolve = getTagNamesToResolve(typ)
}

func (t *BaseTask) ToolKind() dukkha.ToolKind { return t.toolKind }
func (t *BaseTask) ToolName() dukkha.ToolName { return t.toolName }

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

	err := dukkha.ResolveEnv(taskCtx, t, "Env", "env")
	if err != nil {
		return nil, fmt.Errorf(
			"preparing env for hook %q: %w",
			stage.String(), err,
		)
	}

	err = t.ResolveFields(taskCtx, 1, "hooks")
	if err != nil {
		return nil, fmt.Errorf(
			"resolving hooks overview for hook %q: %w",
			stage.String(), err,
		)
	}

	specs, err := t.Hooks.GenSpecs(taskCtx, stage)
	if err != nil {
		return nil, fmt.Errorf(
			"generating exec specs for hook %q: %w",
			stage.String(), err,
		)
	}

	return specs, nil
}

func (t *BaseTask) GetMatrixSpecs(rc dukkha.RenderingContext) (ret []matrix.Entry, err error) {
	err = t.DoAfterFieldsResolved(rc, -1, true, func() error {
		if t.Matrix.HasUserValue() {
			ret = t.Matrix.GenerateEntries(rc.MatrixFilter())
		} else {
			ret = []matrix.Entry{
				{
					"kernel": rc.HostKernel(),
					"arch":   rc.HostArch(),
				},
			}
		}

		return nil
	},
		// t.DoAfterFieldsResolved is intended to serve
		// real task type, so we have to add the prefix
		// `BaseTask.`
		"BaseTask.matrix",
	)

	return
}
