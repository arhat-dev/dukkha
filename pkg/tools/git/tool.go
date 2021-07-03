package git

import (
	"context"
	"fmt"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "git"

func init() {
	field.RegisterInterfaceField(
		tools.ToolType,
		regexp.MustCompile("^git$"),
		func(_ []string) interface{} { return &Tool{} },
	)
}

var _ tools.Tool = (*Tool)(nil)

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`

	cloneTasks map[string]*TaskClone
}

func (t *Tool) ToolKind() string { return ToolKind }

func (t *Tool) Init(
	cacheDir string,
	rf field.RenderingFunc,
	getBaseExecSpec field.ExecSpecGetFunc,
) error {
	err := t.BaseTool.InitBaseTool(cacheDir, "go", rf, getBaseExecSpec)
	if err != nil {
		return fmt.Errorf("git: failed to init tool base: %w", err)
	}

	t.cloneTasks = make(map[string]*TaskClone)

	return nil
}

func (t *Tool) ResolveTasks(tasks []tools.Task) error {
	for i, tsk := range tasks {
		switch typ := tasks[i].(type) {
		case *TaskClone:
			t.cloneTasks[tsk.TaskName()] = typ
		default:
			return fmt.Errorf("git: unknown task type %T with name %q", tsk, tsk.TaskName())
		}
	}

	return nil
}

func (t *Tool) Run(
	ctx context.Context,
	allTools map[tools.ToolKey]tools.Tool,
	allShells map[tools.ToolKey]*tools.BaseTool,
	taskKind, taskName string,
) error {
	var (
		task tools.Task
		ok   bool
	)

	switch taskKind {
	case TaskKindClone:
		task, ok = t.cloneTasks[taskName]
	default:
		return fmt.Errorf("git: unknown task kind %q", taskKind)
	}

	if !ok {
		return fmt.Errorf("git: %s task %q not found", taskKind, taskName)
	}

	return t.BaseTool.RunTask(ctx, t, allTools, allShells, task)
}
