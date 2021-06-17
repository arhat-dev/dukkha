package shell

import (
	"arhat.dev/pkg/exechelper"
)

type ExecFunc func(script string, spec *exechelper.Spec) (exitCode int, err error)

type Config struct {
	ExecFunc ExecFunc
}
