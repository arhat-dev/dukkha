package buildah

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/rs"
)

type stepRun struct {
	rs.BaseField

	Script *string  `yaml:"script"`
	Cmd    []string `yaml:"cmd"`

	// Shell overrides default shell
	Shell []string `yaml:"shell"`
}

func (s *stepRun) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	stepCtx *xbuildContext,
) ([]dukkha.TaskExecSpec, error) {
	return nil, nil
}
