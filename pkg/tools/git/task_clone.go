package tool_git

import (
	"fmt"
	"path"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindClone = "clone"

func init() {
	dukkha.RegisterTask(ToolKind, TaskKindClone, tools.NewTask[TaskClone, *TaskClone])
}

type TaskClone struct {
	tools.BaseTask[GitClone, *GitClone]
}

type GitClone struct {
	URL          string `yaml:"url"`
	Path         string `yaml:"path"`
	RemoteBranch string `yaml:"remote_branch"`
	LocalBranch  string `yaml:"local_branch"`
	RemoteName   string `yaml:"remote_name"`

	ExtraArgs []string `yaml:"extra_args"`

	parent tools.BaseTaskType
}

func (w *GitClone) ToolKind() dukkha.ToolKind       { return ToolKind }
func (w *GitClone) Kind() dukkha.TaskKind           { return TaskKindClone }
func (w *GitClone) LinkParent(p tools.BaseTaskType) { w.parent = p }

func (c *GitClone) GetExecSpecs(
	rc dukkha.TaskExecContext, options dukkha.TaskMatrixExecOptions,
) ([]dukkha.TaskExecSpec, error) {
	var steps []dukkha.TaskExecSpec

	err := c.parent.DoAfterFieldsResolved(rc, -1, true, func() error {
		if len(c.URL) == 0 {
			return fmt.Errorf("remote url not set")
		}

		// first determine the name of the remote
		remoteName := c.RemoteName
		if len(remoteName) == 0 {
			remoteName = "origin"
		}

		remoteBranch := c.RemoteBranch

		localBranch := c.LocalBranch
		if len(localBranch) == 0 {
			localBranch = remoteBranch
		}

		cloneCmd := []string{constant.DUKKHA_TOOL_CMD,
			"clone", "--no-checkout", "--origin", remoteName,
		}

		if len(remoteBranch) != 0 {
			cloneCmd = append(cloneCmd, "--branch", remoteBranch)
		}

		cloneCmd = append(cloneCmd, c.ExtraArgs...)
		cloneCmd = append(cloneCmd, c.URL)
		if len(c.Path) != 0 {
			cloneCmd = append(cloneCmd, c.Path)
		}

		steps = append(steps, dukkha.TaskExecSpec{
			Command:     cloneCmd,
			IgnoreError: false,
		})

		localPath := c.Path
		if len(localPath) == 0 {
			localPath = strings.TrimSuffix(path.Base(c.URL), ".git")
		}

		const replaceTargetDefaultBranch = "<DEFAULT_BRANCH>"
		if len(localBranch) == 0 {
			// local branch name not set
			// which means remote branch name is also not set

			localBranch = replaceTargetDefaultBranch
			remoteBranch = replaceTargetDefaultBranch

			steps = append(steps, dukkha.TaskExecSpec{
				Chdir:           localPath,
				StdoutAsReplace: replaceTargetDefaultBranch,

				IgnoreError: false,
				Command: []string{constant.DUKKHA_TOOL_CMD,
					"symbolic-ref",
					fmt.Sprintf("refs/remotes/%s/HEAD", remoteName),
				},
			})
		}

		// checkout
		steps = append(steps, dukkha.TaskExecSpec{
			IgnoreError: false,
			Chdir:       localPath,
			Command: []string{constant.DUKKHA_TOOL_CMD, "checkout", "-b", localBranch,
				fmt.Sprintf("%s/%s", remoteName, remoteBranch),
			},
		})

		return nil
	})

	return steps, err
}
