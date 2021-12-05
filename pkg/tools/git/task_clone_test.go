package tool_git_test

import (
	"context"
	"strings"
	"testing"

	di "arhat.dev/dukkha/internal"
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	tool_git "arhat.dev/dukkha/pkg/tools/git"
	"arhat.dev/dukkha/pkg/tools/tests"
)

func TestTaskClone_GetExecSpecs(t *testing.T) {
	testCases := []tests.ExecSpecGenerationTestCase{
		{
			Name: "Invalid Empty",
			Task: &tool_git.TaskClone{
				TaskName: "foo",
			},
			ExpectErr: true,
			Options:   dt.CreateTaskMatrixExecOptions(),
		},
		{
			Name: "Valid Clone Using Default Branch",
			Task: &tool_git.TaskClone{
				TaskName: "foo",
				URL:      "example/foo.git",
			},
			Options: dt.CreateTaskMatrixExecOptions(),
			Expected: []dukkha.TaskExecSpec{
				{
					Command: strings.Split(constant.DUKKHA_TOOL_CMD+" clone --no-checkout --origin origin example/foo.git", " "),
				},
				{
					StdoutAsReplace: "<DEFAULT_BRANCH>",
					Chdir:           "foo",
					Command:         strings.Split(constant.DUKKHA_TOOL_CMD+" symbolic-ref refs/remotes/origin/HEAD", " "),
				},
				{
					Chdir:   "foo",
					Command: strings.Split(constant.DUKKHA_TOOL_CMD+" checkout -b <DEFAULT_BRANCH> origin/<DEFAULT_BRANCH>", " "),
				},
			},
		},
		{
			Name: "Valid Clone Changing Remote Name",
			Task: &tool_git.TaskClone{
				TaskName:   "foo",
				URL:        "example/foo",
				RemoteName: "bar",
			},
			Options: dt.CreateTaskMatrixExecOptions(),
			Expected: []dukkha.TaskExecSpec{
				{
					Command: strings.Split(constant.DUKKHA_TOOL_CMD+" clone --no-checkout --origin bar example/foo", " "),
				},
				{
					StdoutAsReplace: "<DEFAULT_BRANCH>",
					Chdir:           "foo",
					Command:         strings.Split(constant.DUKKHA_TOOL_CMD+" symbolic-ref refs/remotes/bar/HEAD", " "),
				},
				{
					Chdir:   "foo",
					Command: strings.Split(constant.DUKKHA_TOOL_CMD+" checkout -b <DEFAULT_BRANCH> bar/<DEFAULT_BRANCH>", " "),
				},
			},
		},
	}

	ctx := dt.NewTestContext(context.TODO())
	ctx.(di.CacheDirSetter).SetCacheDir(t.TempDir())

	tests.RunTaskExecSpecGenerationTests(
		t, ctx, testCases,
	)
}
