package funcs

import (
	"context"
	"sync"

	"arhat.dev/dukkha/third_party/gomplate/conv"
	"arhat.dev/dukkha/third_party/gomplate/env"
)

var (
	ef     *EnvFuncs
	efInit sync.Once
)

// EnvNS - the Env namespace
func EnvNS() *EnvFuncs {
	efInit.Do(func() { ef = &EnvFuncs{} })
	return ef
}

// AddEnvFuncs -
func AddEnvFuncs(f map[string]interface{}) {
	for k, v := range CreateEnvFuncs(context.Background()) {
		f[k] = v
	}
}

// CreateEnvFuncs -
func CreateEnvFuncs(ctx context.Context) map[string]interface{} {
	ns := EnvNS()
	ns.ctx = ctx
	return map[string]interface{}{
		"env":    EnvNS,
		"getenv": ns.Getenv,
	}
}

// EnvFuncs -
type EnvFuncs struct {
	ctx context.Context
}

// Getenv -
func (f *EnvFuncs) Getenv(key interface{}, def ...string) string {
	return env.Getenv(conv.ToString(key), def...)
}

// ExpandEnv -
func (f *EnvFuncs) ExpandEnv(s interface{}) string {
	return env.ExpandEnv(conv.ToString(s))
}
