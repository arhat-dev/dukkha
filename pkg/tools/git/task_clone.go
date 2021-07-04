package git

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
)

const TaskKindClone = "clone"

func init() {
	field.RegisterInterfaceField(
		tools.TaskType,
		regexp.MustCompile(`^git(:.+)?:clone$`),
		func(subMatches []string) interface{} {
			t := &TaskClone{}
			if len(subMatches) != 0 {
				t.SetToolName(strings.TrimPrefix(subMatches[0], ":"))
			}
			return t
		},
	)
}

var _ tools.Task = (*TaskClone)(nil)

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

func (c *TaskClone) ToolKind() string { return ToolKind }
func (c *TaskClone) TaskKind() string { return TaskKindClone }

func (c *TaskClone) GetExecSpecs(ctx *field.RenderingContext, gitCmd []string) ([]tools.TaskExecSpec, error) {
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

	var steps []tools.TaskExecSpec
	cloneCmd := sliceutils.NewStringSlice(
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

	steps = append(steps, tools.TaskExecSpec{
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

		steps = append(steps, tools.TaskExecSpec{
			Chdir:           localPath,
			OutputAsReplace: replaceTargetDefaultBranch,

			IgnoreError: false,
			Command: sliceutils.NewStringSlice(
				gitCmd,
				"symbolic-ref",
				fmt.Sprintf("refs/remotes/%s/HEAD", remoteName),
			),
		})
	}

	// checkout
	steps = append(steps, tools.TaskExecSpec{
		IgnoreError: false,
		Chdir:       localPath,
		Command: sliceutils.NewStringSlice(
			gitCmd, "checkout", "-b", localBranch,
			fmt.Sprintf("%s/%s", remoteName, remoteBranch),
		),
	})

	return steps, nil
}
