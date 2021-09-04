package plugin

import (
	"fmt"
	"reflect"

	"arhat.dev/rs"
	"github.com/traefik/yaegi/interp"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
)

var requiredToolMethods []string

func init() {
	toolType := reflect.TypeOf((*dukkha.Tool)(nil)).Elem()
	for i := 0; i < toolType.NumMethod(); i++ {
		method := toolType.Method(i)
		requiredToolMethods = append(requiredToolMethods, method.Name)
	}
}

func NewReflectTool(
	interp *interp.Interpreter,
	factoryFuncName string,
) dukkha.Tool {
	// TODO: evaluate functions for renderer
	methods, err := evalObjectMethods(
		interp,
		fmt.Sprintf(`%s()`, factoryFuncName),
		requiredToolMethods,
	)
	if err != nil {
		panic(err)
	}

	return &reflectTool{methods: methods}
}

var _ dukkha.Tool = (*reflectTool)(nil)

type reflectTool struct {
	methods map[string]reflect.Value
}

// Kind of the tool managing this tool (e.g. docker)
func (r *reflectTool) Kind() dukkha.ToolKind {
	return r.methods["Kind"].Call(nil)[0].Interface().(dukkha.ToolKind)
}

// Name of the tool managing this tool (e.g. my-tool)
func (r *reflectTool) Name() dukkha.ToolName {
	return r.methods["ToolName"].Call(nil)[0].Interface().(dukkha.ToolName)
}

// Key of this tool
func (r *reflectTool) Key() dukkha.ToolKey {
	return r.methods["Key"].Call(nil)[0].Interface().(dukkha.ToolKey)
}

func (r *reflectTool) GetCmd() []string {
	return r.methods["GetCmd"].Call(nil)[0].Interface().([]string)
}

func (r *reflectTool) GetEnv() dukkha.Env {
	return r.methods["GetEnv"].Call(nil)[0].Interface().(dukkha.Env)
}

func (r *reflectTool) UseShell() bool {
	return r.methods["UseShell"].Call(nil)[0].Bool()
}

func (r *reflectTool) ShellName() string {
	return r.methods["ShellName"].Call(nil)[0].String()
}

func (r *reflectTool) GetTask(dukkha.TaskKey) (dukkha.Task, bool) {
	ret := r.methods["GetTask"].Call(nil)
	return ret[0].Interface().(dukkha.Task), ret[1].Bool()
}

func (r *reflectTool) Init(kind dukkha.ToolKind, cachdDir string) error {
	return r.methods["Init"].Call([]reflect.Value{
		reflect.ValueOf(kind), reflect.ValueOf(cachdDir),
	})[0].Interface().(error)
}

func (r *reflectTool) ResolveTasks(tasks []dukkha.Task) error {
	return r.methods["ResolveTasks"].Call([]reflect.Value{
		reflect.ValueOf(tasks),
	})[0].Interface().(error)
}

func (r *reflectTool) Run(taskCtx dukkha.TaskExecContext) error {
	return r.methods["Run"].Call([]reflect.Value{
		reflect.ValueOf(taskCtx),
	})[0].Interface().(error)
}

func (r *reflectTool) DoAfterFieldsResolved(
	rc dukkha.TaskExecContext,
	depth int,
	do func() error,
	fieldNames ...string,
) error {
	return r.methods["DoAfterFieldsResolved"].Call([]reflect.Value{
		reflect.ValueOf(rc), reflect.ValueOf(depth),
		reflect.ValueOf(do), reflect.ValueOf(fieldNames),
	})[0].Interface().(error)
}

func (r *reflectTool) UnmarshalYAML(n *yaml.Node) error {
	return r.methods["UnmarshalYAML"].Call([]reflect.Value{
		reflect.ValueOf(n),
	})[0].Interface().(error)
}

func (r *reflectTool) ResolveFields(rc rs.RenderingHandler, depth int, fieldNames ...string) error {
	return r.methods["ResolveFields"].Call([]reflect.Value{
		reflect.ValueOf(rc), reflect.ValueOf(depth), reflect.ValueOf(fieldNames),
	})[0].Interface().(error)
}
