package tools

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
)

var _ dukkha.Tool = (*BaseToolWithInit)(nil)

type BaseToolWithInit struct {
	field.BaseField

	BaseTool `yaml:",inline"`
}

func (t *BaseToolWithInit) Init(kind dukkha.ToolKind, cacheDir string) error {
	return t.InitBaseTool(kind, "", cacheDir, t)
}

// GetExecSpec is a helper func for shells
func (t *BaseToolWithInit) GetExecSpec(toExec []string, isFilePath bool) (env, cmd []string, err error) {
	if len(toExec) == 0 {
		return nil, nil, fmt.Errorf("invalid empty exec spec")
	}

	scriptPath := ""
	if !isFilePath {
		scriptPath, err = GetScriptCache(t.cacheDir, strings.Join(toExec, " "))
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
	field.BaseField

	ToolName string   `yaml:"name"`
	Env      []string `yaml:"env"`
	Cmd      []string `yaml:"cmd"`

	// Whether to run this tool in shell and which shell to use
	UsingShell     bool   `yaml:"use_shell"`
	UsingShellName string `yaml:"shell_name"`

	kind dukkha.ToolKind

	cacheDir          string
	defaultExecutable string

	impl  dukkha.Tool
	tasks map[dukkha.TaskKey]dukkha.Task

	fieldsToResolve []string

	mu sync.Mutex
}

func (t *BaseTool) Kind() dukkha.ToolKind { return t.kind }
func (t *BaseTool) Name() dukkha.ToolName { return dukkha.ToolName(t.ToolName) }
func (t *BaseTool) UseShell() bool        { return t.UsingShell }
func (t *BaseTool) ShellName() string     { return t.UsingShellName }

func (t *BaseTool) GetCmd() []string {
	toolCmd := sliceutils.NewStrings(t.Cmd)
	if len(toolCmd) == 0 && len(t.defaultExecutable) != 0 {
		toolCmd = append(toolCmd, t.defaultExecutable)
	}

	return toolCmd
}

func (t *BaseTool) Key() dukkha.ToolKey {
	return dukkha.ToolKey{Kind: t.Kind(), Name: t.Name()}
}

func (t *BaseTool) GetTask(k dukkha.TaskKey) (dukkha.Task, bool) {
	tsk, ok := t.tasks[k]
	return tsk, ok
}

func (t *BaseTool) GetEnv() []string { return sliceutils.NewStrings(t.Env) }

// InitBaseTool must be called in your own version of Init()
// with correct defaultExecutable name
//
// MUST be called when in Init
func (t *BaseTool) InitBaseTool(
	kind dukkha.ToolKind,
	defaultExecutable,
	cacheDir string,
	impl dukkha.Tool,
) error {
	t.kind = kind

	t.cacheDir = cacheDir
	t.defaultExecutable = defaultExecutable

	t.impl = impl
	t.tasks = make(map[dukkha.TaskKey]dukkha.Task)

	typ := reflect.TypeOf(impl).Elem()
	t.fieldsToResolve = getFieldNamesToResolve(typ)

	return nil
}

// ResolveTasks accepts all tasks, override this function if your tool need
// different handling of tasks
func (t *BaseTool) ResolveTasks(tasks []dukkha.Task) error {
	for i, tsk := range tasks {
		t.tasks[dukkha.TaskKey{Kind: tsk.Kind(), Name: tsk.Name()}] = tasks[i]
	}

	return nil
}

// Run task
func (t *BaseTool) Run(taskCtx dukkha.TaskExecContext) error {
	tsk, ok := t.tasks[taskCtx.CurrentTask()]
	if !ok {
		return fmt.Errorf("task %q not found", taskCtx.CurrentTask())
	}

	return RunTask(taskCtx, t.impl, tsk)
}

func (t *BaseTool) DoAfterFieldsResolved(
	ctx dukkha.TaskExecContext,
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
		err = resolveFields(ctx, t.impl, depth, t.fieldsToResolve)
		if err != nil {
			return fmt.Errorf("failed to resolve tool fields: %w", err)
		}
	} else {
		forBase, forImpl := separateBaseAndImpl("BaseTool.", fieldNames)
		if len(forBase) != 0 {
			err := resolveFields(ctx, t, depth, forBase)
			if err != nil {
				return fmt.Errorf("failed to resolve requested BaseTool fields: %w", err)
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

func (t *BaseTool) resolveEssentialFieldsAndAddEnv(mCtx dukkha.RenderingContext) error {
	err := resolveFields(mCtx, t, -1, []string{"ToolName", "Env"})
	if err != nil {
		return fmt.Errorf("failed to resolve essential fields: %w", err)
	}

	mCtx.AddEnv(t.Env...)

	return nil
}

func separateBaseAndImpl(basePrefix string, fieldNames []string) (forBase, forImpl []string) {
	for _, name := range fieldNames {
		if strings.HasPrefix(name, basePrefix) {
			forBase = append(forBase, strings.TrimPrefix(name, basePrefix))
		} else {
			forImpl = append(forImpl, name)
		}
	}

	return
}

func resolveFields(rh field.RenderingHandler, f field.Field, depth int, fieldNames []string) error {
	for _, name := range fieldNames {
		err := f.ResolveFields(rh, depth, name)
		if err != nil {
			return fmt.Errorf("failed to resolve field %q: %w", name, err)
		}
	}

	return nil
}
