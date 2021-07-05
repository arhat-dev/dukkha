package env

import (
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/renderer"
)

var _ renderer.Config = (*Config)(nil)

type Config struct {
	GetExecSpec field.ExecSpecGetFunc
}
