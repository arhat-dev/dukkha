package dukkha_test

import (
	"context"

	"arhat.dev/rs"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
)

var _ dukkha.Renderer = (*echoRenderer)(nil)

type echoRenderer struct {
	rs.BaseField
}

func (r *echoRenderer) Init(ctx dukkha.ConfigResolvingContext) error {
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

func NewTestContext(ctx context.Context) dukkha.ConfigResolvingContext {
	return NewTestContextWithGlobalEnv(ctx, nil)
}

func NewTestContextWithGlobalEnv(
	ctx context.Context,
	globalEnv map[string]string,
) dukkha.ConfigResolvingContext {
	d := dukkha.NewConfigResolvingContext(
		ctx,
		dukkha.ContextOptions{
			InterfaceTypeHandler: dukkha.GlobalInterfaceTypeHandler,
			FailFast:             true,
			ColorOutput:          false,
			TranslateANSIStream:  false,
			RetainANSIStyle:      false,
			Workers:              1,
			GlobalEnv:            globalEnv,
		},
	)

	d.AddRenderer("echo", &echoRenderer{})

	return d
}
