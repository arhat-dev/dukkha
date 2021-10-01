package buildah

import (
	"fmt"

	"arhat.dev/dukkha/pkg/constant"
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

	// Name of the image
	Name *string `yaml:"name"`
	// Digest of the image
	Digest *string `yaml:"digest"`

	Kernel string `yaml:"kernel"`
	Arch   string `yaml:"arch"`

	AlwaysPull    bool     `yaml:"always_pull"`
	NeverPull     bool     `yaml:"never_pull"`
	ExtraPullArgs []string `yaml:"extra_pull_args"`
}

func (s *fromImageSpec) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	stepCtx *xbuildContext,
) ([]dukkha.TaskExecSpec, error) {
	// pull with os/arch/variant
	//
	// flag option --platform can override it
	switch {
	case s.AlwaysPull:
		pullCmd := sliceutils.NewStrings(options.ToolCmd(), "pull")

		if len(s.Kernel) != 0 {
			os, ok := constant.GetOciOS(s.Kernel)
			if !ok {
				os = s.Kernel
			}

			pullCmd = append(pullCmd, "--os", os)
		}

		if len(s.Arch) != 0 {
			arch, ok := constant.GetOciArch(s.Arch)
			if !ok {
				arch = s.Arch
			}
			pullCmd = append(pullCmd, "--arch", arch)

			variant, ok := constant.GetOciArchVariant(s.Arch)
			if ok {
				pullCmd = append(pullCmd, "--variant", variant)
			}
		}

		pullCmd = append(pullCmd, s.ExtraPullArgs...)
	case s.NeverPull:
	default:
		// pull when image does not exist
	}

	fromCmd := sliceutils.NewStrings(options.ToolCmd(), "from", "--cidfile")

	switch {
	case s.Name != nil:
		fromCmd = append(fromCmd, "docker://"+*s.Name)
	case s.Digest != nil:
		fromCmd = append(fromCmd, "")
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
