package helm

import (
	"fmt"
	"path/filepath"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindIndex = "index"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindIndex,
		func(toolName string) dukkha.Task {
			t := &TaskIndex{}
			t.InitBaseTask(ToolKind, dukkha.ToolName(toolName), TaskKindIndex)
			return t
		},
	)
}

type TaskIndex struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	RepoURL     string `yaml:"repo_url"`
	PackagesDir string `yaml:"packages_dir"`
	Merge       string `yaml:"merge"`
}

func (c *TaskIndex) GetExecSpecs(
	rc dukkha.RenderingContext,
	useShell bool,
	shellName string,
	helmCmd []string,
) ([]dukkha.TaskExecSpec, error) {
	indexCmd := sliceutils.NewStrings(helmCmd, "repo", "index")

	if len(c.RepoURL) != 0 {
		indexCmd = append(indexCmd, "--url", c.RepoURL)
	}

	dukkhaWorkingDir := rc.WorkingDir()
	if len(c.PackagesDir) != 0 {
		pkgDir, err2 := filepath.Abs(c.PackagesDir)
		if err2 != nil {
			return nil, fmt.Errorf("failed to determine absolute path of package_dir: %w", err2)
		}

		indexCmd = append(indexCmd, pkgDir)
	} else {
		indexCmd = append(indexCmd, dukkhaWorkingDir)
	}

	if len(c.Merge) != 0 {
		indexCmd = append(indexCmd, "--merge", c.Merge)
	}

	return []dukkha.TaskExecSpec{
		{
			Command:   indexCmd,
			UseShell:  useShell,
			ShellName: shellName,
		},
	}, nil
}
