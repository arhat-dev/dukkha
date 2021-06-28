package buildah

import (
	"context"
	"fmt"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "buildah"

func init() {
	field.RegisterInterfaceField(
		tools.ToolType,
		regexp.MustCompile("^buildah$"),
		func([]string) interface{} { return &Tool{} },
	)
}

var _ tools.Tool = (*Tool)(nil)

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`

	budTasks  map[string]*TaskBud
	pushTasks map[string]*TaskPush
}

func (t *Tool) ToolKind() string { return ToolKind }

func (t *Tool) Init(
	cacheDir string,
	rf field.RenderingFunc,
	getBaseExecSpec field.ExecSpecGetFunc,
) error {
	err := t.BaseTool.InitBaseTool(cacheDir, "buildah", rf, getBaseExecSpec)
	if err != nil {
		return fmt.Errorf("buildah: failed to init tool base: %w", err)
	}

	t.budTasks = make(map[string]*TaskBud)
	t.pushTasks = make(map[string]*TaskPush)

	return nil
}

func (t *Tool) ResolveTasks(tasks []tools.Task) error {
	for i, tsk := range tasks {
		switch typ := tasks[i].(type) {
		case *TaskBud:
			t.budTasks[tsk.TaskName()] = typ
		case *TaskPush:
			t.pushTasks[tsk.TaskName()] = typ
		default:
			return fmt.Errorf("unknown task type %T with name %q", tsk, tsk.TaskName())
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
	case TaskKindBud:
		task, ok = t.budTasks[taskName]
	case TaskKindPush:
		task, ok = t.pushTasks[taskName]
	default:
		return fmt.Errorf("buildah: unknown task kind %q", taskKind)
	}

	if !ok {
		return fmt.Errorf("buildah: %s task %q not found", taskKind, taskName)
	}

	return t.BaseTool.RunTask(ctx, t, allTools, allShells, task)
}
