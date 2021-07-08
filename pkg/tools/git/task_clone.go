package git

import (
	"fmt"
	"path"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindClone = "clone"

func init() {
	dukkha.RegisterTask(
		ToolKind, TaskKindClone,
		func(toolName string) dukkha.Task {
			t := &TaskClone{}
			t.SetToolName(toolName)
			return t
		},
	)
}

var _ dukkha.Task = (*TaskClone)(nil)

type TaskClone struct {
	field.BaseField

	tools.BaseTask `yaml:",inline"`

	URL          string `yaml:"url"`
	Path         string `yaml:"path"`
	RemoteBranch string `yaml:"remote_branch"`
	LocalBranch  string `yaml:"local_branch"`
	RemoteName   string `yaml:"remote_name"`

	ExtraArgs []string `yaml:"extra_args"`
}

func (c *TaskClone) ToolKind() dukkha.ToolKind { return ToolKind }
func (c *TaskClone) Kind() dukkha.TaskKind     { return TaskKindClone }

func (c *TaskClone) GetExecSpecs(rc dukkha.RenderingContext, gitCmd []string) ([]dukkha.TaskExecSpec, error) {
	if len(c.URL) == 0 {
		return nil, fmt.Errorf("remote url not set")
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

	var steps []dukkha.TaskExecSpec
	cloneCmd := sliceutils.NewStrings(
		gitCmd,
		"clone", "--no-checkout", "--origin", remoteName,
	)

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
			OutputAsReplace: replaceTargetDefaultBranch,

			IgnoreError: false,
			Command: sliceutils.NewStrings(
				gitCmd,
				"symbolic-ref",
				fmt.Sprintf("refs/remotes/%s/HEAD", remoteName),
			),
		})
	}

	// checkout
	steps = append(steps, dukkha.TaskExecSpec{
		IgnoreError: false,
		Chdir:       localPath,
		Command: sliceutils.NewStrings(
			gitCmd, "checkout", "-b", localBranch,
			fmt.Sprintf("%s/%s", remoteName, remoteBranch),
		),
	})

	return steps, nil
}
