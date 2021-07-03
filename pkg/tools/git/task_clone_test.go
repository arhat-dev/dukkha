package git_test

import (
	"strings"
	"testing"

	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/git"
	"arhat.dev/dukkha/pkg/tools/tests"
	"github.com/stretchr/testify/assert"
)

func TestTaskClone_GetExecSpecs(t *testing.T) {
	toolCmd := []string{"git"}
	testCases := []tests.ExecSpecGenerationTestCase{
		{
			Name:      "Invalid Empty",
			Task:      &git.TaskClone{},
			ExpectErr: true,
		},
		{
			Name: "Valid Clone Using Default Branch",
			Task: &git.TaskClone{URL: "example/foo.git"},
			Expected: []tools.TaskExecSpec{
				{
					Command: strings.Split("git clone --no-checkout --origin origin example/foo.git", " "),
				},
				{
					OutputAsReplace: "<DEFAULT_BRANCH>",
					Chdir:           "foo",
					Command:         strings.Split("git symbolic-ref refs/remotes/origin/HEAD", " "),
				},
				{
					Chdir:   "foo",
					Command: strings.Split("git checkout -b <DEFAULT_BRANCH> origin/<DEFAULT_BRANCH>", " "),
				},
			},
		},
		{
			Name: "Valid Clone Changing Remote Name",
			Task: &git.TaskClone{URL: "example/foo", RemoteName: "bar"},
			Expected: []tools.TaskExecSpec{
				{
					Command: strings.Split("git clone --no-checkout --origin bar example/foo", " "),
				},
				{
					OutputAsReplace: "<DEFAULT_BRANCH>",
					Chdir:           "foo",
					Command:         strings.Split("git symbolic-ref refs/remotes/bar/HEAD", " "),
				},
				{
					Chdir:   "foo",
					Command: strings.Split("git checkout -b <DEFAULT_BRANCH> bar/<DEFAULT_BRANCH>", " "),
				},
			},
		},
	}

	tests.RunTaskExecSpecGenerationTests(t, nil, toolCmd, testCases)

	assert.EqualValues(t, []string{"git"}, toolCmd)
}