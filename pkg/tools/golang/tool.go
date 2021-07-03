package golang

import (
	"context"
	"fmt"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "golang"

func init() {
	field.RegisterInterfaceField(
		tools.ToolType,
		regexp.MustCompile("^golang$"),
		func(_ []string) interface{} { return &Tool{} },
	)
}

var _ tools.Tool = (*Tool)(nil)

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`

	buildTasks map[string]*TaskBuild
	testTasks  map[string]*TaskTest
}

func (t *Tool) ToolKind() string { return ToolKind }

func (t *Tool) Init(
	cacheDir string,
	rf field.RenderingFunc,
	getBaseExecSpec field.ExecSpecGetFunc,
) error {
	err := t.BaseTool.InitBaseTool(cacheDir, "go", rf, getBaseExecSpec)
	if err != nil {
		return fmt.Errorf("golang: failed to init tool base: %w", err)
	}

	t.buildTasks = make(map[string]*TaskBuild)
	t.testTasks = make(map[string]*TaskTest)

	return nil
}

func (t *Tool) ResolveTasks(tasks []tools.Task) error {
	for i, tsk := range tasks {
		switch typ := tasks[i].(type) {
		case *TaskBuild:
			t.buildTasks[tsk.TaskName()] = typ
		case *TaskTest:
			t.testTasks[tsk.TaskName()] = typ
		default:
			return fmt.Errorf("golang: unknown task type %T with name %q", tsk, tsk.TaskName())
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
	case TaskKindBuild:
		task, ok = t.buildTasks[taskName]
	case TaskKindTest:
		task, ok = t.testTasks[taskName]
	default:
		return fmt.Errorf("golang: unknown task kind %q", taskKind)
	}

	if !ok {
		return fmt.Errorf("golang: %s task %q not found", taskKind, taskName)
	}

	return t.BaseTool.RunTask(ctx, t, allTools, allShells, task)
}
