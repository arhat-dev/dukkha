package tools

import (
	"context"
	"fmt"
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

func (t *BaseTool) RunTask(ctx context.Context, task Task) error {
	specs, err := task.GetMatrixSpec(
		field.WithRenderingValues(ctx, t.Env), t.RenderingFunc,
	)
	if err != nil {
		return fmt.Errorf("failed to create build matrix: %w", err)
	}

	for _, s := range specs {
		fmt.Println(task.TaskKind(), "using matrix:", s.String())

		var env []string
		for k, v := range s {
			env = append(env, "MATRIX_"+strings.ToUpper(k)+"="+v)
		}

		err = task.ResolveFields(
			field.WithRenderingValues(ctx, append(env, t.Env...)),
			t.RenderingFunc, -1,
		)
		if err != nil {
			return fmt.Errorf("failed to resolve task fields: %w", err)
		}

		// TODO: execute task
	}

	return nil
}

// RenderingExec is a helper func for shell renderer
func (t *BaseTool) RenderingExec(script string, spec *exechelper.Spec) (int, error) {
	// TODO
	_, _ = exechelper.Do(*spec)

	return -1, nil
}
