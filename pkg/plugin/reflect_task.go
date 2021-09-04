package plugin

import (
	"fmt"
	"reflect"

	"arhat.dev/rs"
	"github.com/traefik/yaegi/interp"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/matrix"
)

var requiredTaskMethods []string

func init() {
	taskType := reflect.TypeOf((*dukkha.Task)(nil)).Elem()
	for i := 0; i < taskType.NumMethod(); i++ {
		method := taskType.Method(i)
		requiredTaskMethods = append(requiredTaskMethods, method.Name)
	}
}

func NewReflectTask(
	interp *interp.Interpreter,
	factoryFuncName, toolName string,
) dukkha.Task {
	// TODO: evaluate functions for renderer
	methods, err := evalObjectMethods(
		interp,
		fmt.Sprintf(`%s("%s")`, factoryFuncName, toolName),
		requiredTaskMethods,
	)
	if err != nil {
		panic(err)
	}

	return &reflectTask{methods: methods}
}

var _ dukkha.Task = (*reflectTask)(nil)

type reflectTask struct {
	methods map[string]reflect.Value
}

// Kind of the tool managing this task (e.g. docker)
func (r *reflectTask) ToolKind() dukkha.ToolKind {
	return r.methods["ToolKind"].Call(nil)[0].Interface().(dukkha.ToolKind)
}

// Name of the tool managing this task (e.g. my-tool)
func (r *reflectTask) ToolName() dukkha.ToolName {
	return r.methods["ToolName"].Call(nil)[0].Interface().(dukkha.ToolName)
}

// Kind of the task (e.g. build)
func (r *reflectTask) Kind() dukkha.TaskKind {
	return r.methods["Kind"].Call(nil)[0].Interface().(dukkha.TaskKind)
}

// Name of the task (e.g. foo)
func (r *reflectTask) Name() dukkha.TaskName {
	return r.methods["Name"].Call(nil)[0].Interface().(dukkha.TaskName)
}

// Key of this task
func (r *reflectTask) Key() dukkha.TaskKey {
	return r.methods["Key"].Call(nil)[0].Interface().(dukkha.TaskKey)
}

func (r *reflectTask) GetMatrixSpecs(
	rc dukkha.RenderingContext,
) ([]matrix.Entry, error) {
	ret := r.methods["GetMatrixSpecs"].Call([]reflect.Value{reflect.ValueOf(rc)})
	return ret[0].Interface().([]matrix.Entry), ret[1].Interface().(error)
}

func (r *reflectTask) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	ret := r.methods["GetExecSpecs"].Call([]reflect.Value{
		reflect.ValueOf(rc), reflect.ValueOf(options),
	})
	return ret[0].Interface().([]dukkha.TaskExecSpec), ret[1].Interface().(error)
}

func (r *reflectTask) GetHookExecSpecs(
	ctx dukkha.TaskExecContext, state dukkha.TaskExecStage,
) ([]dukkha.RunTaskOrRunCmd, error) {
	ret := r.methods["GetHookExecSpecs"].Call([]reflect.Value{
		reflect.ValueOf(ctx), reflect.ValueOf(state),
	})
	return ret[0].Interface().([]dukkha.RunTaskOrRunCmd), ret[1].Interface().(error)
}

func (r *reflectTask) DoAfterFieldsResolved(
	rc dukkha.RenderingContext,
	depth int,
	do func() error,
	fieldNames ...string,
) error {
	return r.methods["DoAfterFieldsResolved"].Call([]reflect.Value{
		reflect.ValueOf(rc), reflect.ValueOf(depth),
		reflect.ValueOf(do), reflect.ValueOf(fieldNames),
	})[0].Interface().(error)
}

func (r *reflectTask) ContinueOnError() bool {
	return r.methods["ContinueOnError"].Call(nil)[0].Interface().(bool)
}

func (r *reflectTask) UnmarshalYAML(n *yaml.Node) error {
	return r.methods["UnmarshalYAML"].Call([]reflect.Value{
		reflect.ValueOf(n),
	})[0].Interface().(error)
}

func (r *reflectTask) ResolveFields(rc rs.RenderingHandler, depth int, fieldNames ...string) error {
	return r.methods["ResolveFields"].Call([]reflect.Value{
		reflect.ValueOf(rc), reflect.ValueOf(depth), reflect.ValueOf(fieldNames),
	})[0].Interface().(error)
}
