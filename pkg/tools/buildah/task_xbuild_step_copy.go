package buildah

import (
	"bytes"
	"fmt"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/sliceutils"
)

// stepCopy is structured `buildah copy`
type stepCopy struct {
	rs.BaseField

	From copyFromSpec `yaml:"from"`
	To   copyToSpec   `yaml:"to"`

	// ExtraArgs for buildah copy
	ExtraArgs []string `yaml:"extra_args"`
}

func (s *stepCopy) genSpec(
	rc dukkha.TaskExecContext,
	options dukkha.TaskMatrixExecOptions,
	record bool,
) ([]dukkha.TaskExecSpec, error) {
	_ = rc
	var steps []dukkha.TaskExecSpec

	copyCmd := sliceutils.NewStrings(options.ToolCmd(), "copy")
	if record {
		copyCmd = append(copyCmd, "--add-history")
	}
	copyCmd = append(copyCmd, s.ExtraArgs...)

	switch {
	case s.From.Local != nil:
		copyCmd = append(copyCmd, replace_XBUILD_CURRENT_CONTAINER_ID, s.From.Local.Path)
	case s.From.HTTP != nil:
		copyCmd = append(copyCmd, replace_XBUILD_CURRENT_CONTAINER_ID, s.From.HTTP.URL)
	case s.From.Image != nil:
		from := *s.From.Image
		const (
			replace_XBUILD_COPY_FROM_IMAGE_ID = "<XBUILD_COPY_FROM_IMAGE_ID>"
		)

		pullCmd := sliceutils.NewStrings(options.ToolCmd(), "pull")
		pullCmd = append(pullCmd, from.ExtraPullArgs...)
		pullCmd = append(pullCmd, from.Ref)

		steps = append(steps, dukkha.TaskExecSpec{
			StdoutAsReplace:          replace_XBUILD_COPY_FROM_IMAGE_ID,
			FixStdoutValueForReplace: bytes.TrimSpace,

			ShowStdout:  true,
			IgnoreError: true,
			Command:     pullCmd,
			UseShell:    options.UseShell(),
			ShellName:   options.ShellName(),
		})

		copyCmd = append(
			copyCmd,
			"--from", replace_XBUILD_COPY_FROM_IMAGE_ID,
			replace_XBUILD_CURRENT_CONTAINER_ID,
			from.Path,
		)
	case s.From.Step != nil:
		from := *s.From.Step

		copyCmd = append(
			copyCmd,
			"--from", replace_XBUILD_STEP_CONTAINER_ID(from.ID),
			replace_XBUILD_CURRENT_CONTAINER_ID,
			from.Path,
		)
	default:
		return nil, fmt.Errorf("invalid no copy source specified")
	}

	// if path not set, will copy to workingdir
	if len(s.To.Path) != 0 {
		copyCmd = append(copyCmd, s.To.Path)
	}

	steps = append(steps, dukkha.TaskExecSpec{
		IgnoreError: false,
		Command:     copyCmd,
		UseShell:    options.UseShell(),
		ShellName:   options.ShellName(),
	})

	return steps, nil
}

type copyFromSpec struct {
	rs.BaseField

	Local *copyFromLocalSpec `yaml:"local"`
	HTTP  *copyFromHTTPSpec  `yaml:"http"`
	Image *copyFromImageSpec `yaml:"image"`
	Step  *copyFromStepSpec  `yaml:"step"`
}

type copyFromLocalSpec struct {
	rs.BaseField

	Path string `yaml:"path"`
}

type copyFromHTTPSpec struct {
	rs.BaseField

	URL string `yaml:"url"`
}

type copyFromImageSpec struct {
	rs.BaseField

	Ref string `yaml:"ref"`

	ExtraPullArgs []string `yaml:"extra_pull_args"`

	Path string `yaml:"path"`
}

type copyFromStepSpec struct {
	rs.BaseField

	// ID of that step
	ID string `yaml:"id"`

	Path string `yaml:"path"`
}

type copyToSpec struct {
	rs.BaseField

	Path string `yaml:"path"`

	// TODO: implement
	// Chmod []chmodSpec `yaml:"chmod"`
	// Chown []chownSpec `yaml:"chown"`
}

// type chmodSpec struct {
// 	rs.BaseField
//
// 	// Match glob pattern to match files
// 	Match string `yaml:"match"`
//
// 	// Ignore glob pattern to ignore files
// 	Ignore string `yaml:"ignore"`
//
// 	// Value for chmod (e.g. a+x, 0755)
// 	Value string `yaml:"value"`
//
// 	// Recursive run chmod on matched files
// 	Recursive bool `yaml:"recursive"`
// }

// type chownSpec struct {
// 	rs.BaseField
//
// 	// Match glob pattern to match files
// 	Match string `yaml:"match"`
//
// 	// Ignore glob pattern to ignore files
// 	Ignore string `yaml:"ignore"`
//
// 	// Value for chown (e.g. user:group, user, uid, uid:gid)
// 	Value string `yaml:"value"`
//
// 	// Recursive run chown on matched files
// 	Recursive bool `yaml:"recursive"`
// }