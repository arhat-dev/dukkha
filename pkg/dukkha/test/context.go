package dukkha_test

import (
	"context"
	"os"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/utils"
)

func NewTestContext(ctx context.Context) dukkha.ConfigResolvingContext {
	return NewTestContextWithGlobalEnv(ctx, make(map[string]utils.LazyValue))
}

func NewTestContextWithGlobalEnv(
	ctx context.Context,
	globalEnv map[string]utils.LazyValue,
) dukkha.ConfigResolvingContext {
	d := dukkha.NewConfigResolvingContext(
		ctx,
		dukkha.GlobalInterfaceTypeHandler,
		globalEnv,
	)

	if len(d.WorkDir()) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		d.(di.WorkDirOverrider).OverrideWorkDir(cwd)
	}

	d.SetRuntimeOptions(dukkha.RuntimeOptions{
		FailFast:            true,
		ColorOutput:         false,
		TranslateANSIStream: false,
		RetainANSIStyle:     false,
		Workers:             1,
	})

	return d
}
