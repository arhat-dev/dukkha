package helm

import (
	"fmt"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindIndex = "index"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindIndex, tools.NewTask[TaskIndex, *TaskIndex])
}

type HelmIndex struct {
	RepoURL     string `yaml:"repo_url"`
	PackagesDir string `yaml:"packages_dir"`
	Merge       string `yaml:"merge"`

	parent tools.BaseTaskType
}

type TaskIndex struct {
	tools.BaseTask[HelmIndex, *HelmIndex]
}

func (w *HelmIndex) ToolKind() dukkha.ToolKind       { return ToolKind }
func (w *HelmIndex) Kind() dukkha.TaskKind           { return TaskKindIndex }
func (w *HelmIndex) LinkParent(p tools.BaseTaskType) { w.parent = p }

func (c *HelmIndex) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	indexCmd := []string{constant.DUKKHA_TOOL_CMD, "repo", "index"}

	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
		if len(c.RepoURL) != 0 {
			indexCmd = append(indexCmd, "--url", c.RepoURL)
		}

		if len(c.PackagesDir) != 0 {
			pkgDir, err := rc.FS().Abs(c.PackagesDir)
			if err != nil {
				return fmt.Errorf(
					"determine absolute path of package_dir %q: %w",
					pkgDir, err,
				)
			}

			indexCmd = append(indexCmd, pkgDir)
		} else {
			indexCmd = append(indexCmd, rc.WorkDir())
		}

		if len(c.Merge) != 0 {
			indexCmd = append(indexCmd, "--merge", c.Merge)
		}

		return nil
	})

	return []dukkha.TaskExecSpec{{
		Command: indexCmd,
	}}, err
}
