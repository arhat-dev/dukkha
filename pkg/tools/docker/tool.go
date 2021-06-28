package docker

import (
	"context"
	"fmt"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "docker"

func init() {
	field.RegisterInterfaceField(
		tools.ToolType,
		regexp.MustCompile("^docker$"),
		func([]string) interface{} { return &Tool{} },
	)
}

var _ tools.Tool = (*Tool)(nil)

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`

	buildTasks map[string]*TaskBuild
	pushTasks  map[string]*TaskPush
}

func (t *Tool) ToolKind() string { return ToolKind }

func (t *Tool) Init(
	cacheDir string,
	rf field.RenderingFunc,
	getBaseExecSpec field.ExecSpecGetFunc,
) error {
	err := t.BaseTool.InitBaseTool(cacheDir, "docker", rf, getBaseExecSpec)
	if err != nil {
		return fmt.Errorf("docker: failed to init tool base: %w", err)
	}

	t.buildTasks = make(map[string]*TaskBuild)
	t.pushTasks = make(map[string]*TaskPush)

	return nil
}

func (t *Tool) ResolveTasks(tasks []tools.Task) error {
	for i, tsk := range tasks {
		switch typ := tasks[i].(type) {
		case *TaskBuild:
			t.buildTasks[tsk.TaskName()] = typ
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
	case TaskKindBuild:
		task, ok = t.buildTasks[taskName]
	case TaskKindPush:
		task, ok = t.pushTasks[taskName]
	default:
		return fmt.Errorf("docker: unknown task kind %q", taskKind)
	}

	if !ok {
		return fmt.Errorf("docker: %s task %q not found", taskKind, taskName)
	}

	return t.BaseTool.RunTask(ctx, t, allTools, allShells, task)
}
