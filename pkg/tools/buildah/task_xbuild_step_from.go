package buildah

import (
	"bytes"
	"strings"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

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

type stepFrom struct {
	rs.BaseField `yaml:"-"`

	// Ref image
	Ref string `yaml:"ref"`

	Kernel string `yaml:"kernel"`
	Arch   string `yaml:"arch"`

	ExtraPullArgs []string `yaml:"extra_pull_args"`

	// TODO: implement
	Mount []mountSpec `yaml:"mount"`

	ExtraArgs []string `yaml:"extra_args"`
}

func (s *stepFrom) genSpec(
	rc dukkha.TaskExecContext,
	_ dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	_ = rc

	platformArgs := generatePlatformArgs(s.Kernel, s.Arch)

	var steps []dukkha.TaskExecSpec

	const (
		replace_XBUILD_FROM_IMAGE_REF = "<XBUILD_FROM_IMAGE_REF>"
	)

	var imageRef = replace_XBUILD_FROM_IMAGE_REF
	switch strings.ToLower(s.Ref) {
	case "scratch":
		imageRef = s.Ref
	default:
		// pull with os/arch/variant
		//
		// flag option --platform can override it

		// buildah pull [--policy always|never|missing] [--os OCI_OS] [--arch OCI_ARCH] [--variant OCI_VARIANT]
		pullCmd := []string{constant.DUKKHA_TOOL_CMD, "pull"}

		pullCmd = append(pullCmd, platformArgs...)
		pullCmd = append(pullCmd, s.ExtraPullArgs...)
		pullCmd = append(pullCmd, s.Ref)

		steps = append(steps, dukkha.TaskExecSpec{
			StdoutAsReplace:          replace_XBUILD_FROM_IMAGE_REF,
			FixStdoutValueForReplace: bytes.TrimSpace,

			ShowStdout:  true,
			IgnoreError: false,
			Command:     pullCmd,
		})
	}

	fromCmd := []string{constant.DUKKHA_TOOL_CMD, "from"}
	fromCmd = append(fromCmd, platformArgs...)
	fromCmd = append(fromCmd, s.ExtraArgs...)
	fromCmd = append(fromCmd, imageRef)

	const (
		replace_XBUILD_CURRENT_CONTAINER_NAME = "<XBUILD_CURRENT_CONTAINER_NAME>"
	)

	// produce container name
	steps = append(steps, dukkha.TaskExecSpec{
		StdoutAsReplace:          replace_XBUILD_CURRENT_CONTAINER_NAME,
		FixStdoutValueForReplace: bytes.TrimSpace,

		ShowStdout:  true,
		IgnoreError: false,
		Command:     fromCmd,
	})

	// retrieve container id
	steps = append(steps, dukkha.TaskExecSpec{
		StdoutAsReplace:          replace_XBUILD_CURRENT_CONTAINER_ID,
		FixStdoutValueForReplace: bytes.TrimSpace,

		ShowStdout:  true,
		IgnoreError: false,
		Command: []string{constant.DUKKHA_TOOL_CMD,
			"inspect",
			"--type", "container",
			"--format", "{{ .ContainerID }}",
			replace_XBUILD_CURRENT_CONTAINER_NAME,
		},
	})

	return steps, nil
}

func generatePlatformArgs(kernel, arch string) []string {
	var platformArgs []string

	if len(kernel) != 0 {
		ociOS, ok := constant.GetOciOS(kernel)
		if !ok {
			ociOS = kernel
		}

		platformArgs = append(platformArgs, "--os", ociOS)
	}

	if len(arch) != 0 {
		ociArch, ok := constant.GetOciArch(arch)
		if !ok {
			ociArch = arch
		}

		platformArgs = append(platformArgs, "--arch", ociArch)

		ociVariant, ok := constant.GetOciArchVariant(arch)
		if ok && len(ociVariant) != 0 {
			platformArgs = append(platformArgs, "--variant", ociVariant)
		}
	}

	return platformArgs
}
