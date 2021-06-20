package shell

import (
	"arhat.dev/pkg/exechelper"

	"arhat.dev/dukkha/pkg/renderer"
)

type ExecFunc func(script string, spec *exechelper.Spec) (exitCode int, err error)

var _ renderer.Config = (*Config)(nil)

type Config struct {
	ExecFunc ExecFunc
}
