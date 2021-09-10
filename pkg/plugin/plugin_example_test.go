package plugin

import (
	"reflect"

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
	ret := &fooRenderer{}
	rs.InitRecursively(reflect.ValueOf(ret), dukkha.GlobalInterfaceTypeHandler)
	return ret
}

type toolFoo struct {
	// rs.BaseField

	tools.BaseTool `yaml:",inline"`
}

func (f *toolFoo) Init(kind dukkha.ToolKind, cachdDir string) error {
	return f.BaseTool.InitBaseTool(kind, "foo", cachdDir, f)
}

// NewTool_{tool-kind}
// nolint:revive
func NewTool_tool_foo() *toolFoo {
	ret := &toolFoo{}
	rs.InitRecursively(reflect.ValueOf(ret), dukkha.GlobalInterfaceTypeHandler)
	return ret
}

type taskFoo struct {
	// rs.BaseField

	tools.BaseTask `yaml:",inline"`

	Echo string `yaml:"echo"`
}

func (f *taskFoo) GetExecSpecs(rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions) ([]dukkha.TaskExecSpec, error) {
	println("TASK_FOO SPEC: " + f.Echo)
	return nil, nil
}

// NewTask_{tool-kind}_{task-kind}
// nolint:revive
func NewTask_tool_foo_task_foo(name string) *taskFoo {
	ret := &taskFoo{}
	rs.InitRecursively(reflect.ValueOf(ret), dukkha.GlobalInterfaceTypeHandler)
	return ret
}

type barTask = taskFoo

// NewTask_{tool-kind}_{task-kind}
// nolint:revive
func NewTask_tool_foo_task_bar(name string) *barTask {
	ret := &barTask{}
	rs.InitRecursively(reflect.ValueOf(ret), dukkha.GlobalInterfaceTypeHandler)
	return ret
}

// NewTask_{tool-kind}_{task-kind}
// nolint:revive
func NewTask_bar_tool_task_bar(name string) *barTask {
	ret := &barTask{}
	rs.InitRecursively(reflect.ValueOf(ret), dukkha.GlobalInterfaceTypeHandler)
	return ret
}
