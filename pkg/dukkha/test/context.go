package dukkha_test

import (
	"context"

	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/types"
)

var _ dukkha.Renderer = (*echoRenderer)(nil)

type echoRenderer struct{}

func (*echoRenderer) RenderYaml(rc types.RenderingContext, rawData interface{}) ([]byte, error) {
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

func NewTestContext(ctx context.Context) dukkha.Context {
	return NewTestContextWithGlobalEnv(ctx, nil)
}

func NewTestContextWithGlobalEnv(ctx context.Context, globalEnv map[string]string) dukkha.Context {
	d := dukkha.NewConfigResolvingContext(
		ctx, globalEnv, testBootstrapExec, true, 1,
	)
	_ = d.AddRenderer(&echoRenderer{}, "echo")

	return d
}
