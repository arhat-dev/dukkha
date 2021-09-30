package buildah

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/rs"
)

type step struct {
	rs.BaseField

	// ID of this step, if not set, will be the array index of this step
	ID string `yaml:"id"`

	// Workdir overrides default workdir settings
	Workdir string `yaml:"workdir"`

	// Commit this step as a new layer after this step finished
	Commit *bool `yaml:"commit"`

	// Set default options for all following steps
	Set *stepSet `yaml:"set"`

	// From some rootfs
	From *stepFrom `yaml:"from"`

	// Run some command
	Run *stepRun `yaml:"run"`

	// Copy files to somewhere
	Copy *stepCopy `yaml:"copy"`
}

func (s *step) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	stepCtx *xbuildContext,
) ([]dukkha.TaskExecSpec, error) {
	switch {
	case s.Set != nil:
		return nil, nil
	case s.From != nil:
		return s.From.genSpec(rc, options, stepCtx)
	case s.Run != nil:
		return s.Run.genSpec(rc, options, stepCtx)
	case s.Copy != nil:
		return s.Copy.genSpec(rc, options, stepCtx)
	default:
		return nil, fmt.Errorf("invalid empty step")
	}
}
