package template_file

import (
	"arhat.dev/dukkha/pkg/renderer"
	"arhat.dev/dukkha/pkg/renderer/template"
)

var _ renderer.Config = (*Config)(nil)

type Config template.Config
