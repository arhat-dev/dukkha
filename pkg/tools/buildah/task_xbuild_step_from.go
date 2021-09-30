package buildah

import (
	"fmt"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/rs"
)

type stepFrom struct {
	rs.BaseField

	// Image as base rootfs
	Image *fromImageSpec `yaml:"image"`

	// Step as base rootfs
	Step *fromStepSpec `yaml:"step"`

	// TODO
	Mount []*mountSpec `yaml:"mount"`
}

type mountSpec struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`

	// Options of bind mount
	// 	ro, rw, z, Z, O
	// 	shared, slave, private, unbindable
	//  rshared, rslave, rprivate, runbindable
	Options []string `yaml:"options"`

	// FixUser adds `U` option to the mount, buildah will set correct uid and gid
	FixUser bool

	// AsOverlay adds `O` option to the mount, build will mount it using overlayfs
	AsOverlay bool
}

func (s *stepFrom) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	stepCtx *xbuildContext,
) ([]dukkha.TaskExecSpec, error) {
	switch {
	case s.Image != nil:
		return s.Image.genSpec(rc, options, stepCtx)
	case s.Step != nil:
		return s.Step.genSpec(rc, options, stepCtx)
	default:
		return nil, fmt.Errorf("invalid empty from spec")
	}
}

type fromImageSpec struct {
	rs.BaseField

	Name   *string `yaml:"name"`
	Digest *string `yaml:"digest"`

	OS   string `yaml:"os"`
	Arch string `yaml:"arch"`

	AlwaysPull bool `yaml:"always_pull"`
	NeverPull  bool `yaml:"never_pull"`
}

func (s *fromImageSpec) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	stepCtx *xbuildContext,
) ([]dukkha.TaskExecSpec, error) {
	switch {
	case s.Name != nil:
	case s.Digest != nil:
		fromCmd := sliceutils.NewStrings(options.ToolCmd(), "from", "--cidfile")
		_ = fromCmd
	default:
		return nil, fmt.Errorf("invalid no image reference")
	}

	return nil, nil
}

type fromStepSpec struct {
	rs.BaseField

	ID string `yaml:"id"`
}

func (s *fromStepSpec) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	stepCtx *xbuildContext,
) ([]dukkha.TaskExecSpec, error) {
	refCtx := stepCtx.Steps[s.ID]
	_ = refCtx
	return nil, nil
}
