package tools

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
)

type TaskImpl interface {
	ToolKind() dukkha.ToolKind
	Kind() dukkha.TaskKind
	LinkParent(p BaseTaskType)

	GetExecSpecs(
		rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
	) ([]dukkha.TaskExecSpec, error)
}

type BaseTaskType interface {
	dukkha.Task
	SetToolName(string)
	CacheFS() *fshelper.OSFS
}

func NewTask[V any, T BaseTaskType](toolName string) dukkha.Task {
	val := new(V)
	ret := *(*T)(unsafe.Pointer(&val))
	ret.SetToolName(toolName)
	return ret
}

// BaseTask is the helper to wrap plain old task spec as dukkha.Task
//
// NOTE: V MUST be a struct type, T MUST be *V
type BaseTask[V any, T TaskImpl] struct {
	rs.BaseField `yaml:"-"`

	TaskName            dukkha.TaskName      `yaml:"name"`
	Env                 dukkha.NameValueList `yaml:"env"`
	Matrix              matrix.Spec          `yaml:"matrix"`
	Hooks               TaskHooks            `yaml:"hooks,omitempty"`
	ContinueOnErrorFlag bool                 `yaml:"continue_on_error"`

	Impl V `yaml:",inline"`

	// fields managed by BaseTask

	cacheFS       *fshelper.OSFS
	toolName      dukkha.ToolName
	tagsToResolve []string
	lock          sync.Mutex
}

func (t *BaseTask[V, T]) SetToolName(name string) {
	t.toolName = dukkha.ToolName(name)
}

func (t *BaseTask[V, T]) getTaskImpl() T {
	ptr := &t.Impl
	return *(*T)(unsafe.Pointer(&ptr))
}

// Kind implements dukkha.Task
func (t *BaseTask[V, T]) Kind() dukkha.TaskKind {
	return t.getTaskImpl().Kind()
}

// Name implements dukkha.Task
func (t *BaseTask[V, T]) Name() dukkha.TaskName {
	return t.TaskName
}

// Key implements dukkha.Task
func (t *BaseTask[V, T]) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: t.getTaskImpl().Kind(), Name: t.TaskName}
}

// GetExecSpecs implements dukkha.Task
func (t *BaseTask[V, T]) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	return t.getTaskImpl().GetExecSpecs(rc, options)
}

func (t *BaseTask[V, T]) Init(cacheFS *fshelper.OSFS) error {
	t.cacheFS = cacheFS
	t.getTaskImpl().LinkParent(t)

	typ := reflect.TypeOf(t.Impl)
	t.tagsToResolve = getTagNamesToResolve(typ)
	return nil
}

func (t *BaseTask[V, T]) CacheFS() *fshelper.OSFS {
	return t.cacheFS
}

func (t *BaseTask[V, T]) DoAfterFieldsResolved(
	ctx dukkha.RenderingContext,
	depth int,
	resolveEnv bool,
	do func() error,
	tagNames ...string,
) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if resolveEnv {
		err = dukkha.ResolveAndAddEnv(ctx, t, "Env", "env")
		if err != nil {
			return
		}
	}

	if len(tagNames) == 0 {
		// resolve all fields of the real task type
		err = t.ResolveFields(ctx, depth, t.tagsToResolve...)
		if err != nil {
			return fmt.Errorf("resolving tool fields: %w", err)
		}
	} else {
		err = t.ResolveFields(ctx, depth, tagNames...)
		if err != nil {
			return fmt.Errorf("resolving requested fields: %w", err)
		}
	}

	return do()
}

func (t *BaseTask[V, T]) ToolKind() dukkha.ToolKind { return t.getTaskImpl().ToolKind() }
func (t *BaseTask[V, T]) ToolName() dukkha.ToolName { return t.toolName }

func (t *BaseTask[V, T]) ContinueOnError() bool { return t.ContinueOnErrorFlag }

func (t *BaseTask[V, T]) GetHookExecSpecs(
	taskCtx dukkha.TaskExecContext,
	stage dukkha.TaskExecStage,
) ([]dukkha.TaskExecSpec, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// hooks may have reference to env defined in task scope

	err := dukkha.ResolveAndAddEnv(taskCtx, t, "Env", "env")
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

func (t *BaseTask[V, T]) GetMatrixSpecs(rc dukkha.RenderingContext) (ret []matrix.Entry, err error) {
	err = t.DoAfterFieldsResolved(rc, -1, true, func() error {
		if t.Matrix.IsEmpty() {
			ret = []matrix.Entry{
				{
					"kernel": rc.HostKernel(),
					"arch":   rc.HostArch(),
				},
			}
		} else {
			ret = t.Matrix.GenerateEntries(rc.MatrixFilter())
		}

		return nil
	},
		"matrix",
	)

	return
}
