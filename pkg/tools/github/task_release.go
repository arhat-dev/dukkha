package github

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindRelease = "release"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindRelease,
		func(toolName string) dukkha.Task {
			t := &TaskRelease{}
			t.SetToolName(toolName)
			return t
		},
	)
}

var _ dukkha.Task = (*TaskRelease)(nil)

type TaskRelease struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	Tag        string `yaml:"tag"`
	Draft      bool   `yaml:"draft"`
	PreRelease bool   `yaml:"pre_release"`
	Title      string `yaml:"title"`
	Notes      string `yaml:"notes"`

	Files []ReleaseFileSpec `yaml:"files"`
}

type ReleaseFileSpec struct {
	// path to the file, can use glob
	Path string `yaml:"path"`
	// the display label as noted in gh docs
	// https://cli.github.com/manual/gh_release_create
	Label string `yaml:"label"`
}

func (c *TaskRelease) ToolKind() dukkha.ToolKind { return ToolKind }
func (c *TaskRelease) Kind() dukkha.TaskKind     { return TaskKindRelease }

func (c *TaskRelease) GetExecSpecs(rc dukkha.RenderingContext, ghCmd []string) ([]dukkha.TaskExecSpec, error) {
	createCmd := sliceutils.NewStrings(
		ghCmd, "release", "create", c.Tag,
	)

	if c.Draft {
		createCmd = append(createCmd, "--draft")
	}

	if c.PreRelease {
		createCmd = append(createCmd, "--prerelease")
	}

	if len(c.Title) != 0 {
		createCmd = append(createCmd,
			"--title", fmt.Sprintf("%q", c.Title),
		)
	}

	if len(c.Notes) != 0 {
		cacheDir := rc.CacheDir()
		f, err := ioutil.TempFile(cacheDir, "github-release-note-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary release note file: %w", err)
		}

		noteFile := f.Name()
		_, err = f.Write([]byte(c.Notes))
		_ = f.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to write release note: %w", err)
		}

		createCmd = append(createCmd, "--notes-file", noteFile)
	}

	for _, spec := range c.Files {
		matches, err := filepath.Glob(spec.Path)
		if err != nil {
			matches = []string{spec.Path}
		}

		for i, file := range matches {
			var arg string
			if len(spec.Label) != 0 {
				arg = `'` + file + `#` + spec.Label
				if i != 0 {
					arg += " " + strconv.FormatInt(int64(i), 10)
				}

				arg += `'`
			} else {
				arg = `'` + file + `#` + filepath.Base(file) + `'`
			}

			createCmd = append(createCmd, arg)
		}
	}

	return []dukkha.TaskExecSpec{{
		Command: createCmd,
	}}, nil
}
