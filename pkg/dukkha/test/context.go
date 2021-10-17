package dukkha_test

import (
	"context"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
)

var _ dukkha.Renderer = (*echoRenderer)(nil)

type echoRenderer struct {
	rs.BaseField `yaml:"-"`
}

func (r *echoRenderer) Init(ctx dukkha.ConfigResolvingContext) error { return nil }

func (*echoRenderer) RenderYaml(
	rc dukkha.RenderingContext, rawData interface{},
) (interface{}, error) {
	return rawData, nil
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
		globalEnv,
	)

	d.SetRuntimeOptions(dukkha.RuntimeOptions{
		FailFast:            true,
		ColorOutput:         false,
		TranslateANSIStream: false,
		RetainANSIStyle:     false,
		Workers:             1,
	})
	d.AddRenderer("echo", &echoRenderer{})

	return d
}
