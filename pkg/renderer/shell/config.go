package shell

import (
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/shell_file"
)

var _ renderer.Config = (*Config)(nil)

type Config shell_file.Config
