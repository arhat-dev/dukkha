package dukkha_test

import (
	"context"

	"gopkg.in/yaml.v3"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
)

var _ dukkha.Renderer = (*echoRenderer)(nil)

type echoRenderer struct {
	rs.BaseField `yaml:"-"`
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
		dukkha.GlobalInterfaceTypeHandler,
		true,
		false, // turn off color output
		false, // do not translate ansi stream
		false, // retainANSIStyle is not used when translateANSIStream is disabled
		1,
		globalEnv,
	)

	d.AddRenderer("echo", &echoRenderer{})

	return d
}
