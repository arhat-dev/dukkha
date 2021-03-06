package tool_git_test

import (
	"context"
	"strings"
	"testing"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	dt "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/tools"
	tool_git "arhat.dev/dukkha/pkg/tools/git"
	"arhat.dev/dukkha/pkg/tools/tests"
)

func TestTaskClone_GetExecSpecs(t *testing.T) {
	t.Parallel()

	newTask := func(name string) *tool_git.TaskClone {
		return tools.NewTask[tool_git.TaskClone, *tool_git.TaskClone](name).(*tool_git.TaskClone)
	}

	testCases := []tests.ExecSpecGenerationTestCase{
		{
			Name:      "Invalid Empty",
			Task:      newTask("foo"),
			ExpectErr: true,
			Options:   dt.CreateTaskMatrixExecOptions(),
		},
		{
			Name: "Valid Clone Using Default Branch",
			Task: func() dukkha.Task {
				tsk := newTask("foo")
				tsk.Impl.URL = "example/foo.git"
				return tsk
			}(),
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
			Task: func() dukkha.Task {
				tsk := newTask("foo")
				tsk.Impl.URL = "example/foo"
				tsk.Impl.RemoteName = "bar"
				return tsk
			}(),
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

	ctx := dt.NewTestContext(context.TODO(), t.TempDir())

	tests.RunTaskExecSpecGenerationTests(
		t, ctx, testCases,
	)
}
