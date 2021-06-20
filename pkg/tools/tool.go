package tools

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/field"
)

// ToolType for interface type registration
var ToolType = reflect.TypeOf((*Tool)(nil)).Elem()

// nolint:revive
type Tool interface {
	field.Interface

	// Kind of the tool, e.g. golang, docker
	ToolKind() string

	ToolName() string

	Init(rf field.RenderingFunc) error

	ResolveTasks(tasks []Task) error

	Run(ctx context.Context, taskKind, taskName string) error
}

type BaseTool struct {
	field.BaseField

	Name string   `yaml:"name"`
	Path string   `yaml:"path"`
	Env  []string `yaml:"env"`

	GlobalArgs []string `yaml:"args"`

	RenderingFunc field.RenderingFunc `json:"-" yaml:"-"`
}

func (t *BaseTool) Init(rf field.RenderingFunc) error {
	t.RenderingFunc = rf
	return nil
}

func (t *BaseTool) ToolName() string { return t.Name }

func (t *BaseTool) RunTask(ctx context.Context, toolKind string, task Task) error {
	specs, err := task.GetMatrixSpecs(
		field.WithRenderingValues(ctx, t.Env), t.RenderingFunc,
	)
	if err != nil {
		return fmt.Errorf("failed to create build matrix: %w", err)
	}

	for _, s := range specs {
		fmt.Println("---", task.TaskKind(), "with {", s.String(), "}")

		var matrixEnv []string
		for k, v := range s {
			matrixEnv = append(matrixEnv, "MATRIX_"+strings.ToUpper(k)+"="+v)
		}

		values := field.WithRenderingValues(ctx, append(matrixEnv, t.Env...))
		err = task.ResolveFields(values, t.RenderingFunc, -1)
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

		cmd = append(cmd, t.GlobalArgs...)
		cmd = append(cmd, args...)

		fmt.Println(">>>", toolKind, "[", strings.Join(cmd, " "), "]")
		p, err := exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: cmd,
			Env:     values.Values().Env,
			Stdin:   os.Stdin,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		})
		if err != nil {
			return fmt.Errorf("failed to execute command [ %s ]: %w", strings.Join(cmd, " "), err)
		}

		code, err := p.Wait()
		if err != nil {
			return fmt.Errorf("command exited with code %d: %w", code, err)
		}
	}

	return nil
}

// RenderingExec is a helper func for shell renderer
func (t *BaseTool) RenderingExec(script string, spec *exechelper.Spec) (int, error) {
	// TODO
	_, _ = exechelper.Do(*spec)

	return -1, nil
}
