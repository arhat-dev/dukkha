package github

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindRelease = "release"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindRelease,
		func(toolName string) dukkha.Task {
			t := &TaskRelease{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), t)
			return t
		},
	)
}

type TaskRelease struct {
	rs.BaseField `yaml:"-"`

	TaskName string `yaml:"name"`

	tools.BaseTask `yaml:",inline"`

	Tag        string `yaml:"tag"`
	Draft      bool   `yaml:"draft"`
	PreRelease bool   `yaml:"pre_release"`
	Title      string `yaml:"title"`
	Notes      string `yaml:"notes"`

	Files []ReleaseFileSpec `yaml:"files"`
}

type ReleaseFileSpec struct {
	rs.BaseField `yaml:"-"`

	// path to the file, can use glob
	Path string `yaml:"path"`
	// the display label as noted in gh docs
	// https://cli.github.com/manual/gh_release_create
	Label string `yaml:"label"`
}

func (c *TaskRelease) Kind() dukkha.TaskKind { return TaskKindRelease }
func (c *TaskRelease) Name() dukkha.TaskName { return dukkha.TaskName(c.TaskName) }
func (c *TaskRelease) Key() dukkha.TaskKey {
	return dukkha.TaskKey{Kind: c.Kind(), Name: c.Name()}
}

func (c *TaskRelease) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {

	var steps []dukkha.TaskExecSpec
	err := c.DoAfterFieldsResolved(rc, -1, true, func() error {
		createCmd := []string{constant.DUKKHA_TOOL_CMD, "release", "create", c.Tag}

		if c.Draft {
			createCmd = append(createCmd, "--draft")
		}

		if c.PreRelease {
			createCmd = append(createCmd, "--prerelease")
		}

		if len(c.Title) != 0 {
			createCmd = append(createCmd, "--title", c.Title)
		}

		if len(c.Notes) != 0 {
			f, err := os.CreateTemp(rc.CacheDir(), "github-release-note-*")
			if err != nil {
				return fmt.Errorf("creating temporary release note file: %w", err)
			}

			noteFile := f.Name()
			_, err = f.Write([]byte(c.Notes))
			_ = f.Close()
			if err != nil {
				return fmt.Errorf("writing release note: %w", err)
			}

			createCmd = append(createCmd, "--notes-file", noteFile)
		}

		for _, spec := range c.Files {
			matches, err := rc.FS().Glob(spec.Path)
			if err != nil {
				matches = []string{spec.Path}
			}

			for i, file := range matches {
				file, err = rc.FS().Abs(file)
				if err != nil {
					return err
				}

				var arg string
				if len(spec.Label) != 0 {
					arg += file + `#` + spec.Label
					if i != 0 {
						arg += " " + strconv.FormatInt(int64(i), 10)
					}
				} else {
					arg += file + `#` + filepath.Base(file)
				}

				createCmd = append(createCmd, arg)
			}
		}

		steps = append(steps, dukkha.TaskExecSpec{
			Command: createCmd,
		})

		return nil
	})

	return steps, err
}
