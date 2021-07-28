package dukkha_test

import (
	"context"

	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
)

var _ dukkha.Renderer = (*echoRenderer)(nil)

type echoRenderer struct {
	field.BaseField
}

func (r *echoRenderer) Init(ctx dukkha.ConfigResolvingContext) error {
	ctx.AddRenderer("echo", r)
	return nil
}

func (*echoRenderer) RenderYaml(rc dukkha.RenderingContext, rawData interface{}) ([]byte, error) {
	switch t := rawData.(type) {
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	default:
		data, err := yaml.Marshal(rawData)
		return data, err
	}
}

func testBootstrapExec(toExec []string, isFilePath bool) (env []string, cmd []string, err error) {
	return []string{"DUKKHA_TEST=true"}, []string{}, nil
}

func NewTestContext(ctx context.Context) dukkha.ConfigResolvingContext {
	return NewTestContextWithGlobalEnv(ctx, nil)
}

func NewTestContextWithGlobalEnv(
	ctx context.Context,
	globalEnv map[string]string,
) dukkha.ConfigResolvingContext {
	d := dukkha.NewConfigResolvingContext(
		ctx,
		dukkha.GlobalInterfaceTypeHandler,
		true,
		false, // turn off color output
		1,
	)

	d.AddRenderer("echo", &echoRenderer{})

	return d
}
