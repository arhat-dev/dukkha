package tools

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
)

var _ dukkha.Tool = (*BaseTool)(nil)

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
	stdoutIsTty       bool

	tasks map[dukkha.TaskKey]dukkha.Task

	mu sync.Mutex
}

func (t *BaseTool) ResolveFields(rc field.RenderingHandler, depth int, fieldName string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.BaseField.ResolveFields(rc, depth, fieldName)
}

// Init the tool, called when resolving tools config when dukkha start
//
// override it if the value of your tool kind is different from its
// default executable
func (t *BaseTool) Init(kind dukkha.ToolKind, cachdDir string) error {
	return t.InitBaseTool(kind, string(kind), cachdDir)
}

func (t *BaseTool) Kind() dukkha.ToolKind { return t.kind }
func (t *BaseTool) Name() dukkha.ToolName { return dukkha.ToolName(t.ToolName) }

func (t *BaseTool) UseShell() bool {
	return t.UsingShell
}

func (t *BaseTool) ShellName() string {
	return t.UsingShellName
}

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
func (t *BaseTool) InitBaseTool(kind dukkha.ToolKind, defaultExecutable, cacheDir string) error {
	t.kind = kind

	t.cacheDir = cacheDir
	t.defaultExecutable = defaultExecutable
	t.stdoutIsTty = term.IsTerminal(int(os.Stdout.Fd()))

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

	return runTask(taskCtx, t, tsk)
}

// GetExecSpec is a helper func for shells
func (t *BaseTool) GetExecSpec(toExec []string, isFilePath bool) (env, cmd []string, err error) {
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
