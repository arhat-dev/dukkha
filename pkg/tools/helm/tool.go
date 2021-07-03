package helm

import (
	"context"
	"fmt"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "helm"

func init() {
	field.RegisterInterfaceField(
		tools.ToolType,
		regexp.MustCompile("^helm$"),
		func(_ []string) interface{} { return &Tool{} },
	)
}

var _ tools.Tool = (*Tool)(nil)

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`

	packageTasks map[string]*TaskPackage
	indexTasks   map[string]*TaskIndex
}

func (t *Tool) ToolKind() string { return ToolKind }

func (t *Tool) Init(
	cacheDir string,
	rf field.RenderingFunc,
	getBaseExecSpec field.ExecSpecGetFunc,
) error {
	err := t.BaseTool.InitBaseTool(cacheDir, "helm", rf, getBaseExecSpec)
	if err != nil {
		return fmt.Errorf("helm: failed to init tool base: %w", err)
	}

	t.packageTasks = make(map[string]*TaskPackage)
	t.indexTasks = make(map[string]*TaskIndex)

	return nil
}

func (t *Tool) ResolveTasks(tasks []tools.Task) error {
	for i, tsk := range tasks {
		switch typ := tasks[i].(type) {
		case *TaskPackage:
			t.packageTasks[tsk.TaskName()] = typ
		case *TaskIndex:
			t.indexTasks[tsk.TaskName()] = typ
		default:
			return fmt.Errorf("helm: unknown task type %T with name %q", tsk, tsk.TaskName())
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
	case TaskKindPackage:
		task, ok = t.packageTasks[taskName]
	case TaskKindIndex:
		task, ok = t.indexTasks[taskName]
	default:
		return fmt.Errorf("helm: unknown task kind %q", taskKind)
	}

	if !ok {
		return fmt.Errorf("helm: %s task %q not found", taskKind, taskName)
	}

	return t.BaseTool.RunTask(ctx, t, allTools, allShells, task)
}
