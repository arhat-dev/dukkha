package shell_file

import (
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/shell"
)

var _ renderer.Config = (*Config)(nil)

type Config shell.Config
