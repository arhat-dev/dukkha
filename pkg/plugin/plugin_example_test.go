package plugin

import (
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

type fooRenderer struct {
	rs.BaseField

	name string
}

func (r *fooRenderer) Init(ctx dukkha.ConfigResolvingContext) error {
	return nil
}

func (r *fooRenderer) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) (result []byte, err error) {
	return []byte("HELLO " + r.name), nil
}

// NewRenderer_{renderer-default-name}
// nolint:revive
func NewRenderer_foo(name string) *fooRenderer {
	return &fooRenderer{name: name}
}

type fooTool struct {
	rs.BaseField

	tools.BaseTool
}

func (f *fooTool) Init(kind dukkha.ToolKind, cachdDir string) error {
	return f.BaseTool.InitBaseTool(kind, "foo", cachdDir, f)
}

// NewTool_{tool-kind}
// nolint:revive
func NewTool_foo_tool() *fooTool {
	return &fooTool{}
}

type fooTask struct {
	rs.BaseField

	tools.BaseTask
}

func (f *fooTask) GetExecSpecs(rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions) ([]dukkha.TaskExecSpec, error) {
	return nil, nil
}

// NewTask_{tool-kind}_{task-kind}
// nolint:revive
func NewTask_foo_tool_foo_task(name string) *fooTask {
	return &fooTask{}
}

type barTask fooTask

// NewTask_{tool-kind}_{task-kind}
// nolint:revive
func NewTask_foo_tool_bar_task(name string) *barTask {
	return &barTask{}
}

// NewTask_{tool-kind}_{task-kind}
// nolint:revive
func NewTask_bar_tool_bar_task(name string) *barTask {
	return &barTask{}
}
