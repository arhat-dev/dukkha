package buildah

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/md5helper"
	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

type stepRun struct {
	rs.BaseField `yaml:"-"`

	// Script
	//
	// helpful when you need to run remote script, use:
	//
	// 	run:
	// 	  script@http: http://some-script.company
	//
	Script string `yaml:"script"`

	// ScriptArgs for the script
	ScriptArgs []string `yaml:"script_args"`

	// ExecutableFile path in local fs, run it in container
	//
	// Will copy the executable to container and remove it after executaion
	//
	// helpful when your executable is too large to be loaded as `script`
	ExecutableFile string `yaml:"executable_file"`

	// ExecutableArgs args to the ExecutableFile
	ExecutableArgs []string `yaml:"executable_args"`

	// Cmd as bare exec
	Cmd []string `yaml:"cmd"`

	ExtraArgs []string `yaml:"extra_args"`
}

const (
	execPathInContainer = "/tmp/xbuild-run"
)

func (s *stepRun) genSpec(
	rc dukkha.TaskExecContext,
	cacheFS *fshelper.OSFS,
	record bool,
) ([]dukkha.TaskExecSpec, error) {
	runCmd := []string{constant.DUKKHA_TOOL_CMD, "run"}
	if record {
		runCmd = append(runCmd, "--add-history")
	}
	runCmd = append(runCmd, s.ExtraArgs...)
	runCmd = append(runCmd, replace_XBUILD_CURRENT_CONTAINER_ID, "--")

	var steps []dukkha.TaskExecSpec

	switch {
	case len(s.Cmd) != 0:
		steps = append(steps, dukkha.TaskExecSpec{
			IgnoreError: false,
			Command:     append(runCmd, s.Cmd...),
		})
	case len(s.ExecutableFile) != 0:
		localExecutablePath := s.ExecutableFile

		const (
			replace_XBUILD_RUN_EXECUTABLE_SRC_REDACTED_PATH = "<XBUILD_RUN_EXECUTABLE_SRC_REDACTED_PATH>"
		)
		steps = append(steps,
			// write redacted file
			dukkha.TaskExecSpec{
				StdoutAsReplace:          replace_XBUILD_RUN_EXECUTABLE_SRC_REDACTED_PATH,
				FixStdoutValueForReplace: bytes.TrimSpace,

				AlterExecFunc: func(
					replace dukkha.ReplaceEntries,
					stdin io.Reader,
					stdout, stderr io.Writer,
				) (dukkha.RunTaskOrRunCmd, error) {
					file := "run-executable-" + hex.EncodeToString(
						md5helper.Sum([]byte(localExecutablePath)),
					) + "-redacted"

					err := cacheFS.WriteFile(file, []byte(""), 0644)
					if err != nil {
						return nil, err
					}

					srcFile, err := cacheFS.Abs(file)
					if err != nil {
						return nil, err
					}

					// TODO: remove additional \n for ansi translation flush
					_, err = stdout.Write([]byte(srcFile + "\n"))
					if err != nil {
						return nil, fmt.Errorf("write redacted executable file cache to stdout: %q", srcFile)
					}

					return nil, rc.FS().WriteFile(srcFile, []byte(""), 0644)
				},
			},
			// copy executable to container
			dukkha.TaskExecSpec{
				Command: []string{constant.DUKKHA_TOOL_CMD, "copy", "--chmod", "0755",
					replace_XBUILD_CURRENT_CONTAINER_ID,
					localExecutablePath, execPathInContainer,
				},
			},
			dukkha.TaskExecSpec{
				IgnoreError: false,
				Command:     append(append(runCmd, execPathInContainer), s.ExecutableArgs...),
			},
			// override that executable
			dukkha.TaskExecSpec{
				IgnoreError: false,
				Command: []string{constant.DUKKHA_TOOL_CMD, "copy", "--chmod", "0644",
					replace_XBUILD_CURRENT_CONTAINER_ID,
					replace_XBUILD_RUN_EXECUTABLE_SRC_REDACTED_PATH, execPathInContainer,
				},
			},
		)
	case len(s.Script) != 0:
		// copy this script to container
		const (
			replace_XBUILD_RUN_SCRIPT_SRC_PATH          = "<XBUILD_RUN_SCRIPT_SRC_PATH>"
			replace_XBUILD_RUN_SCRIPT_SRC_REDACTED_PATH = "<XBUILD_RUN_SCRIPT_SRC_REDACTED_PATH>"
		)

		script := s.Script
		steps = append(steps,
			// write script to local cache
			dukkha.TaskExecSpec{
				StdoutAsReplace:          replace_XBUILD_RUN_SCRIPT_SRC_PATH,
				FixStdoutValueForReplace: bytes.TrimSpace,

				ShowStdout: true,
				AlterExecFunc: func(
					replace dukkha.ReplaceEntries,
					stdin io.Reader,
					stdout, stderr io.Writer,
				) (dukkha.RunTaskOrRunCmd, error) {
					file := "run-script-" + hex.EncodeToString(md5helper.Sum([]byte(script)))
					err := cacheFS.WriteFile(file, []byte(script), 0644)
					if err != nil {
						return nil, err
					}

					srcFile, err := cacheFS.Abs(file)
					if err != nil {
						return nil, err
					}

					// TODO: remove additional \n for ansi translation flush
					_, err = stdout.Write([]byte(srcFile + "\n"))
					if err != nil {
						return nil, fmt.Errorf("write script cache path to stdout: %q", srcFile)
					}

					return nil, nil
				},
			},
			// write redacted file
			dukkha.TaskExecSpec{
				StdoutAsReplace:          replace_XBUILD_RUN_SCRIPT_SRC_REDACTED_PATH,
				FixStdoutValueForReplace: bytes.TrimSpace,

				AlterExecFunc: func(
					replace dukkha.ReplaceEntries,
					stdin io.Reader,
					stdout, stderr io.Writer,
				) (dukkha.RunTaskOrRunCmd, error) {
					v, ok := replace[replace_XBUILD_RUN_SCRIPT_SRC_PATH]
					if !ok {
						return nil, fmt.Errorf("unexpected script path not found")
					}

					srcFile := string(v.Data)
					redactedSrcFile := srcFile + "-redacted"

					// TODO: remove additional \n for ansi translation flush
					_, err := stdout.Write([]byte(redactedSrcFile + "\n"))
					if err != nil {
						return nil, fmt.Errorf("write redacted file path to stdout: %w", err)
					}

					return nil, rc.FS().WriteFile(redactedSrcFile, []byte(""), 0644)
				},
			},
			// copy script to container
			dukkha.TaskExecSpec{
				Command: []string{constant.DUKKHA_TOOL_CMD, "copy", "--chmod", "0755",
					replace_XBUILD_CURRENT_CONTAINER_ID,
					replace_XBUILD_RUN_SCRIPT_SRC_PATH, execPathInContainer,
				},
			},
			// run the script
			dukkha.TaskExecSpec{
				IgnoreError: false,
				Command:     append(append(runCmd, execPathInContainer), s.ScriptArgs...),
			},
			// override that script
			dukkha.TaskExecSpec{
				IgnoreError: false,
				Command: []string{constant.DUKKHA_TOOL_CMD, "copy", "--chmod", "0644",
					replace_XBUILD_CURRENT_CONTAINER_ID,
					replace_XBUILD_RUN_SCRIPT_SRC_REDACTED_PATH, execPathInContainer,
				},
			},
		)
	default:
		return nil, fmt.Errorf("invalid empty run statement")
	}

	return steps, nil
}
