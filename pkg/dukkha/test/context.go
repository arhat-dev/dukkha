package dukkha_test

import (
	"context"
	"os"

	"arhat.dev/tlang"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

func NewTestContext(ctx context.Context, cacheDir string) dukkha.ConfigResolvingContext {
	return NewTestContextWithGlobalEnv(ctx, &dukkha.GlobalEnvSet{
		constant.GlobalEnv_DUKKHA_CACHE_DIR: tlang.ImmediateString(cacheDir),
	})
}

func NewTestContextWithGlobalEnv(
	ctx context.Context,
	globalEnv *dukkha.GlobalEnvSet,
) dukkha.ConfigResolvingContext {
	if globalEnv == nil {
		globalEnv = &dukkha.GlobalEnvSet{}
	}

	for i, v := range globalEnv {
		if v == nil {
			globalEnv[i] = tlang.ImmediateString("")
		}
	}

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
