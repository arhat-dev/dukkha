package docker

import (
	"context"
	"fmt"
	"regexp"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/tools"
)

const ToolKind = "docker"

func init() {
	field.RegisterInterfaceField(
		tools.ToolType,
		regexp.MustCompile("^docker$"),
		func() interface{} { return &Tool{} },
	)
}

var _ tools.Tool = (*Tool)(nil)

type Tool struct {
	field.BaseField

	tools.BaseTool `yaml:",inline"`

	mgr *renderer.Manager

	buildTasks map[string]*TaskBuild
	pushTasks  map[string]*TaskPush
}

func (t *Tool) Kind() string { return ToolKind }

func (t *Tool) ResolveTasks(tasks []tools.Task) error {
	// TODO
	_ = t.mgr
	return nil
}

func (t *Tool) Exec(ctx context.Context, taskKind, taskName string) error {
	switch taskKind {
	case TaskKindBuild:
		bt, ok := t.buildTasks[taskName]
		if !ok {
			return fmt.Errorf("docker: build task %q not found", taskName)
		}

		_ = bt
		return nil
	case TaskKindPush:
		pt, ok := t.pushTasks[taskName]
		if !ok {
			return fmt.Errorf("docker: push task %q not found", taskName)
		}

		pt.Inherit(t.buildTasks[taskName])

		return nil
	default:
		return fmt.Errorf("docker: unknown task kind %q", taskKind)
	}
}
