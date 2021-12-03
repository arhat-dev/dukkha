package tools

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
)

var _ dukkha.Tool = (*ShellTool)(nil)

type ShellTool struct {
	rs.BaseField `yaml:"-"`

	ToolName dukkha.ToolName `yaml:"name"`

	BaseTool `yaml:",inline"`
}

func (t *ShellTool) Init(cacheFS *fshelper.OSFS) error {
	return t.InitBaseTool("", cacheFS, t)
}

func (t *ShellTool) Name() dukkha.ToolName { return t.ToolName }
func (t *ShellTool) Kind() dukkha.ToolKind { return "shell" }
func (t *ShellTool) Key() dukkha.ToolKey {
	return dukkha.ToolKey{Kind: t.Kind(), Name: t.Name()}
}

// GetExecSpec is a helper func for shells
func (t *ShellTool) GetExecSpec(
	toExec []string, isFilePath bool,
) (env dukkha.Env, cmd []string, err error) {
	if len(toExec) == 0 {
		return nil, nil, fmt.Errorf("invalid empty exec spec")
	}

	scriptPath := ""
	if !isFilePath {
		scriptPath, err = GetScriptCache(t.CacheFS, strings.Join(toExec, " "))
		if err != nil {
			return nil, nil, fmt.Errorf("tools: failed to ensure script cache: %w", err)
		}
	} else {
		scriptPath = toExec[0]
	}

	cmd = sliceutils.NewStrings(t.Cmd)
	if len(cmd) == 0 {
		cmd = append(cmd, t.defaultExecutable)
	}

	return t.Env, append(cmd, scriptPath), nil
}

type BaseTool struct {
	rs.BaseField `yaml:"-"`

	Env dukkha.Env `yaml:"env"`
	Cmd []string   `yaml:"cmd"`

	CacheFS *fshelper.OSFS `yaml:"-"`

	defaultExecutable string

	impl  dukkha.Tool
	tasks map[dukkha.TaskKey]dukkha.Task

	fieldsToResolve []string

	mu sync.Mutex
}

func (t *BaseTool) GetCmd() []string {
	toolCmd := sliceutils.NewStrings(t.Cmd)
	if len(toolCmd) == 0 && len(t.defaultExecutable) != 0 {
		toolCmd = append(toolCmd, t.defaultExecutable)
	}

	return toolCmd
}

func (t *BaseTool) GetTask(k dukkha.TaskKey) (dukkha.Task, bool) {
	tsk, ok := t.tasks[k]
	return tsk, ok
}

func (t *BaseTool) AllTasks() map[dukkha.TaskKey]dukkha.Task { return t.tasks }
func (t *BaseTool) GetEnv() dukkha.Env                       { return t.Env }

// InitBaseTool must be called in your own version of Init()
// with correct defaultExecutable name
//
// MUST be called when in Init
func (t *BaseTool) InitBaseTool(
	defaultExecutable string,
	cacheFS *fshelper.OSFS,
	impl dukkha.Tool,
) error {
	t.CacheFS = cacheFS
	t.defaultExecutable = defaultExecutable

	t.impl = impl
	t.tasks = make(map[dukkha.TaskKey]dukkha.Task)

	typ := reflect.TypeOf(impl).Elem()
	t.fieldsToResolve = getTagNamesToResolve(typ)

	return nil
}

// AddTasks accepts all tasks, override this function if your tool need
// different handling of tasks
func (t *BaseTool) AddTasks(tasks []dukkha.Task) error {
	for i, tsk := range tasks {
		t.tasks[dukkha.TaskKey{Kind: tsk.Kind(), Name: tsk.Name()}] = tasks[i]
	}

	return nil
}

// Run task
func (t *BaseTool) Run(ctx dukkha.TaskExecContext, key dukkha.TaskKey) error {
	tsk, ok := t.tasks[key]
	if !ok {
		return fmt.Errorf("task %q not found", key)
	}

	return RunTask(&TaskExecRequest{
		Context: ctx,
		Tool:    t.impl,
		Task:    tsk,
	})
}

func (t *BaseTool) DoAfterFieldsResolved(
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
		err := t.impl.ResolveFields(ctx, depth, t.fieldsToResolve...)
		if err != nil {
			return fmt.Errorf("failed to resolve tool fields: %w", err)
		}
	} else {
		forBase, forImpl := separateBaseAndImpl("BaseTool.", tagNames)
		if len(forBase) != 0 {
			err := t.ResolveFields(ctx, depth, forBase...)
			if err != nil {
				return fmt.Errorf("failed to resolve requested BaseTool fields: %w", err)
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

func separateBaseAndImpl(basePrefix string, names []string) (forBase, forImpl []string) {
	for _, name := range names {
		if strings.HasPrefix(name, basePrefix) {
			forBase = append(forBase, strings.TrimPrefix(name, basePrefix))
		} else {
			forImpl = append(forImpl, name)
		}
	}

	return
}
