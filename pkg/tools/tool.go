package tools

import (
	"fmt"
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

func (t *BaseTool) DoAfterFieldsResolved(mCtx dukkha.TaskExecContext, do func() error) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	err := t.impl.ResolveFields(mCtx, -1, "")
	if err != nil {
		return fmt.Errorf("failed to resolve tool fields: %w", err)
	}

	return do()
}
