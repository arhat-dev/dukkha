package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/output"
)

// ToolType for interface type registration
var ToolType = reflect.TypeOf((*Tool)(nil)).Elem()

// nolint:revive
type Tool interface {
	field.Interface

	// Kind of the tool, e.g. golang, docker
	ToolKind() string

	ToolName() string

	Init(
		cacheDir string,
		rf field.RenderingFunc,
		getBaseExecSpec field.ExecSpecGetFunc,
	) error

	ResolveTasks(tasks []Task) error

	Run(ctx context.Context, taskKind, taskName string) error
}

type BaseTool struct {
	field.BaseField

	Name string   `yaml:"name"`
	Path string   `yaml:"path"`
	Env  []string `yaml:"env"`

	Args []string `yaml:"args"`

	cacheDir             string                `json:"-" yaml:"-"`
	RenderingFunc        field.RenderingFunc   `json:"-" yaml:"-"`
	getBootstrapExecSpec field.ExecSpecGetFunc `json:"-" yaml:"-"`
}

func (t *BaseTool) Init(
	cacheDir string,
	rf field.RenderingFunc,
	getBaseExecSpec field.ExecSpecGetFunc,
) error {
	t.cacheDir = cacheDir
	t.RenderingFunc = rf
	t.getBootstrapExecSpec = getBaseExecSpec
	return nil
}

func (t *BaseTool) ToolName() string { return t.Name }

func (t *BaseTool) RunTask(ctx context.Context, toolKind string, task Task) error {
	baseCtx := field.WithRenderingValues(ctx, t.Env)

	specs, err := task.GetMatrixSpecs(baseCtx, t.RenderingFunc)
	if err != nil {
		return fmt.Errorf("failed to create build matrix: %w", err)
	}

	for _, s := range specs {
		taskCtx := baseCtx.Clone()

		output.WriteTaskStart(taskCtx.Context(),
			task.ToolKind(), task.ToolName(), task.TaskKind(), task.TaskName(),
			s.String(),
		)

		for k, v := range s {
			taskCtx.AddEnv("MATRIX_" + strings.ToUpper(k) + "=" + v)
		}

		err = task.ResolveFields(taskCtx, t.RenderingFunc, -1)
		if err != nil {
			return fmt.Errorf("failed to resolve task fields: %w", err)
		}

		// TODO: use generated args to execute tasks in parallel

		args, err := task.ExecArgs()
		if err != nil {
			return fmt.Errorf("failed to generate task args: %w", err)
		}

		var cmd []string
		if len(t.Path) != 0 {
			cmd = append(cmd, t.Path)
		} else {
			cmd = append(cmd, toolKind)
		}

		cmd = append(cmd, t.Args...)
		cmd = append(cmd, args...)

		_, runScriptCmd, err := t.getBootstrapExecSpec(strings.Join(cmd, " "), false)
		if err != nil {
			return fmt.Errorf("failed to get exec spec from bootstrap config: %w", err)
		}

		output.WriteExecStart(
			taskCtx.Context(),
			toolKind,
			cmd,
			filepath.Base(runScriptCmd[len(runScriptCmd)-1]),
		)

		p, err := exechelper.Do(exechelper.Spec{
			Context: taskCtx.Context(),
			Command: runScriptCmd,
			Env:     taskCtx.Values().Env,

			Stdin: os.Stdin,

			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})
		if err != nil {
			output.WriteExecFailure()

			return fmt.Errorf("failed to execute command [ %s ]: %w", strings.Join(cmd, " "), err)
		}

		code, err := p.Wait()
		if err != nil {
			return fmt.Errorf("command exited with code %d: %w", code, err)
		}

		output.WriteExecSuccess()
	}

	return nil
}

// GetExecSpec is a helper func for shell renderer
func (t *BaseTool) GetExecSpec(script string, isFilePath bool) (env, cmd []string, err error) {
	scriptPath := script
	if !isFilePath {
		scriptPath, err = GetScriptCache(t.cacheDir, script)
		if err != nil {
			return nil, nil, fmt.Errorf("tools: failed to ensure script cache: %w", err)
		}
	}

	if len(t.Path) != 0 {
		cmd = append(cmd, t.Path)
	} else {
		cmd = append(cmd, t.ToolName())
	}

	cmd = append(cmd, t.Args...)

	return t.Env, append(cmd, scriptPath), nil
}
