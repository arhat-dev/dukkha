package dukkha_test

import (
	"context"
	"os"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/tlang"
)

func NewTestContext(ctx context.Context) dukkha.ConfigResolvingContext {
	return NewTestContextWithGlobalEnv(ctx, make(map[string]tlang.LazyValueType[string]))
}

func NewTestContextWithGlobalEnv(
	ctx context.Context,
	globalEnv map[string]tlang.LazyValueType[string],
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
