package plugin

import (
	"fmt"
	"reflect"

	"arhat.dev/rs"
	"github.com/traefik/yaegi/interp"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
)

var rendererMethodNames []string

func init() {
	rendererType := reflect.TypeOf((*dukkha.Renderer)(nil)).Elem()
	for i := 0; i < rendererType.NumMethod(); i++ {
		method := rendererType.Method(i)
		rendererMethodNames = append(rendererMethodNames, method.Name)
	}
}

func NewReflectRenderer(
	interp *interp.Interpreter,
	factoryFuncName, rendererName string,
) dukkha.Renderer {
	methods, err := evalObjectMethods(
		interp,
		fmt.Sprintf(`%s("%s")`, factoryFuncName, rendererName),
		rendererMethodNames,
	)
	if err != nil {
		panic(err)
	}

	return &reflectRenderer{methods: methods}
}

var _ dukkha.Renderer = (*reflectRenderer)(nil)

type reflectRenderer struct {
	methods map[string]reflect.Value
}

func (r *reflectRenderer) Init(ctx dukkha.ConfigResolvingContext) error {
	return r.methods["Init"].Call([]reflect.Value{reflect.ValueOf(ctx)})[0].Interface().(error)
}

func (r *reflectRenderer) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) (result []byte, err error) {
	rawDataVal := reflect.ValueOf(rawData)
	if !rawDataVal.IsValid() || rawDataVal.IsZero() {
		// TODO: set nil value and do not cause
		// 		 `panic: reflect: call of reflect.Value.IsZero on zero Value`
		rawDataVal = reflect.ValueOf("")
	}

	ret := r.methods["RenderYaml"].Call([]reflect.Value{
		reflect.ValueOf(rc), rawDataVal,
	})

	res := ret[0].Interface().([]byte)
	if ret[1].IsZero() || ret[1].IsNil() {
		return res, nil
	}

	return res, ret[1].Interface().(error)
}

func (r *reflectRenderer) UnmarshalYAML(n *yaml.Node) error {
	return r.methods["UnmarshalYAML"].Call([]reflect.Value{
		reflect.ValueOf(n),
	})[0].Interface().(error)
}

func (r *reflectRenderer) ResolveFields(rc rs.RenderingHandler, depth int, fieldNames ...string) error {
	return r.methods["ResolveFields"].Call([]reflect.Value{
		reflect.ValueOf(rc), reflect.ValueOf(depth), reflect.ValueOf(fieldNames),
	})[0].Interface().(error)
}
