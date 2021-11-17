package buildah

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"arhat.dev/pkg/md5helper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

// stepCopy is structured `buildah copy`
type stepCopy struct {
	rs.BaseField `yaml:"-"`

	From copyFromSpec `yaml:"from"`
	To   copyToSpec   `yaml:"to"`

	// ExtraArgs for buildah copy
	ExtraArgs []string `yaml:"extra_args"`
}

func (s *stepCopy) genSpec(
	rc dukkha.TaskExecContext,
	_ dukkha.TaskMatrixExecOptions,
	record bool,
) ([]dukkha.TaskExecSpec, error) {
	_ = rc
	var steps []dukkha.TaskExecSpec

	copyCmd := []string{constant.DUKKHA_TOOL_CMD, "copy"}
	if record {
		copyCmd = append(copyCmd, "--add-history")
	}
	copyCmd = append(copyCmd, s.ExtraArgs...)

	switch {
	case s.From.Text != nil:
		data := s.From.Text.Data

		const (
			replace_XBUILD_COPY_FROM_TEXT_DATA_SRC_PATH = "<XBUILD_COPY_FROM_TEXT_DATA_FILE>"
		)

		steps = append(steps, dukkha.TaskExecSpec{
			StdoutAsReplace:          replace_XBUILD_COPY_FROM_TEXT_DATA_SRC_PATH,
			FixStdoutValueForReplace: bytes.TrimSpace,

			AlterExecFunc: func(
				replace dukkha.ReplaceEntries,
				stdin io.Reader,
				stdout, stderr io.Writer,
			) (dukkha.RunTaskOrRunCmd, error) {
				srcFile := filepath.Join(
					rc.CacheDir(),
					"buildah", "xbuild",
					"copy-text-"+hex.EncodeToString(md5helper.Sum([]byte(data))),
				)
				err := os.MkdirAll(filepath.Dir(srcFile), 0755)
				if err != nil && !os.IsExist(err) {
					return nil, fmt.Errorf("failed to ensure text data cache dir: %w", err)
				}

				// TODO: remove additional \n for ansi translation flush
				_, err = stdout.Write([]byte(srcFile + "\n"))
				if err != nil {
					return nil, fmt.Errorf("failed to create text data cache: %q", srcFile)
				}

				return nil, os.WriteFile(srcFile, []byte(data), 0644)
			},
		})

		copyCmd = append(copyCmd,
			replace_XBUILD_CURRENT_CONTAINER_ID,
			replace_XBUILD_COPY_FROM_TEXT_DATA_SRC_PATH,
		)
	case s.From.Local != nil:
		copyCmd = append(copyCmd, replace_XBUILD_CURRENT_CONTAINER_ID, s.From.Local.Path)
	case s.From.HTTP != nil:
		copyCmd = append(copyCmd, replace_XBUILD_CURRENT_CONTAINER_ID, s.From.HTTP.URL)
	case s.From.Image != nil:
		from := *s.From.Image
		const (
			replace_XBUILD_COPY_FROM_IMAGE_ID = "<XBUILD_COPY_FROM_IMAGE_ID>"
		)

		pullCmd := []string{constant.DUKKHA_TOOL_CMD, "pull"}
		pullCmd = append(pullCmd, generatePlatformArgs(from.Kernel, from.Arch)...)
		pullCmd = append(pullCmd, from.ExtraPullArgs...)
		pullCmd = append(pullCmd, from.Ref)

		steps = append(steps, dukkha.TaskExecSpec{
			StdoutAsReplace:          replace_XBUILD_COPY_FROM_IMAGE_ID,
			FixStdoutValueForReplace: bytes.TrimSpace,

			ShowStdout:  true,
			IgnoreError: true,
			Command:     pullCmd,
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
	})

	return steps, nil
}

type copyFromSpec struct {
	rs.BaseField `yaml:"-"`

	Text  *copyFromTextSpec  `yaml:"text"`
	Local *copyFromLocalSpec `yaml:"local"`
	HTTP  *copyFromHTTPSpec  `yaml:"http"`
	Image *copyFromImageSpec `yaml:"image"`
	Step  *copyFromStepSpec  `yaml:"step"`
}

type copyFromTextSpec struct {
	rs.BaseField `yaml:"-"`

	Data string `yaml:"data"`
}

type copyFromLocalSpec struct {
	rs.BaseField `yaml:"-"`

	Path string `yaml:"path"`
}

type copyFromHTTPSpec struct {
	rs.BaseField `yaml:"-"`

	URL string `yaml:"url"`
}

type copyFromImageSpec struct {
	rs.BaseField `yaml:"-"`

	Ref string `yaml:"ref"`

	Kernel string `yaml:"kernel"`
	Arch   string `yaml:"arch"`

	ExtraPullArgs []string `yaml:"extra_pull_args"`

	Path string `yaml:"path"`
}

type copyFromStepSpec struct {
	rs.BaseField `yaml:"-"`

	// ID of that step
	ID string `yaml:"id"`

	Path string `yaml:"path"`
}

type copyToSpec struct {
	rs.BaseField `yaml:"-"`

	Path string `yaml:"path"`

	// TODO: implement
	// Chmod []chmodSpec `yaml:"chmod"`
	// Chown []chownSpec `yaml:"chown"`
}

// type chmodSpec struct {
// 	rs.BaseField `yaml:"-"`
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
// 	rs.BaseField `yaml:"-"`
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
