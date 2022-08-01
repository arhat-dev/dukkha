package tools

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
)

type ToolImpl interface {
	DefaultExecutable() string
	Kind() dukkha.ToolKind
}

// BaseTool is the helper to wrap plain old tool spec as dukkha.Tool
//
// NOTE: V MUST be a struct type, T MUST be *V
type BaseTool[V any, T ToolImpl] struct {
	rs.BaseField `yaml:"-"`

	ToolName dukkha.ToolName      `yaml:"name"`
	Env      dukkha.NameValueList `yaml:"env"`
	Cmd      []string             `yaml:"cmd"`

	Impl V `yaml:",inline"`

	CacheFS *fshelper.OSFS `yaml:"-"`

	tasks map[dukkha.TaskKey]dukkha.Task

	fieldsToResolve []string

	mu sync.Mutex
}

func (t *BaseTool[V, T]) getImpl() T {
	ptr := &t.Impl
	return *(*T)(unsafe.Pointer(&ptr))
}

func (t *BaseTool[V, T]) Kind() dukkha.ToolKind {
	return t.getImpl().Kind()
}

func (t *BaseTool[V, T]) Name() dukkha.ToolName {
	return t.ToolName
}

func (t *BaseTool[V, T]) Key() dukkha.ToolKey {
	return dukkha.ToolKey{Kind: t.getImpl().Kind(), Name: t.ToolName}
}

func (t *BaseTool[V, T]) GetCmd() []string {
	toolCmd := sliceutils.NewStrings(t.Cmd)
	if len(toolCmd) == 0 && len(t.getImpl().DefaultExecutable()) != 0 {
		toolCmd = append(toolCmd, t.getImpl().DefaultExecutable())
	}

	return toolCmd
}

func (t *BaseTool[V, T]) GetTask(k dukkha.TaskKey) (dukkha.Task, bool) {
	tsk, ok := t.tasks[k]
	return tsk, ok
}

func (t *BaseTool[V, T]) AllTasks() map[dukkha.TaskKey]dukkha.Task { return t.tasks }
func (t *BaseTool[V, T]) GetEnv() dukkha.NameValueList             { return t.Env }

func (t *BaseTool[V, T]) Init(cacheFS *fshelper.OSFS) error {
	t.CacheFS = cacheFS

	t.tasks = make(map[dukkha.TaskKey]dukkha.Task)

	typ := reflect.TypeOf(t.Impl)
	t.fieldsToResolve = append([]string{"name", "cmd"}, getTagNamesToResolve(typ)...)

	return nil
}

// AddTasks accepts all tasks, override this function if your tool need
// different handling of tasks
func (t *BaseTool[V, T]) AddTasks(tasks []dukkha.Task) error {
	for i, tsk := range tasks {
		t.tasks[dukkha.TaskKey{Kind: tsk.Kind(), Name: tsk.Name()}] = tasks[i]
	}

	return nil
}

// Run task
func (t *BaseTool[V, T]) Run(ctx dukkha.TaskExecContext, key dukkha.TaskKey) error {
	tsk, ok := t.tasks[key]
	if !ok {
		return fmt.Errorf("task %q not found", key)
	}

	return RunTask(&TaskExecRequest{
		Context: ctx,
		Tool:    t,
		Task:    tsk,
	})
}

func (t *BaseTool[V, T]) DoAfterFieldsResolved(
	ctx dukkha.RenderingContext,
	depth int,
	resolveEnv bool,
	do func() error,
	tagNames ...string,
) (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if resolveEnv {
		err = dukkha.ResolveAndAddEnv(ctx, t, "Env", "env")
		if err != nil {
			return
		}
	}

	if len(tagNames) == 0 {
		// resolve all fields of the real task type
		err = t.ResolveFields(ctx, depth, t.fieldsToResolve...)
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
