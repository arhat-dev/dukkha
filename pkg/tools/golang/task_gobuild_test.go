package golang

import (
	"context"
	"testing"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/tests"
)

func TestTaskBuild_GetExecSpecs(t *testing.T) {
	t.Parallel()

	testCases := []tests.ExecSpecGenerationTestCase{
		{
			Name: "Default Build Task",
			Task: &TaskBuild{
				TaskName: "foo",
				BaseTask: tools.BaseTask{},
			},
			Options: dukkha_test.CreateTaskMatrixExecOptions(),
			Expected: []dukkha.TaskExecSpec{
				{
					EnvSuggest: dukkha.Env{
						{
							Name:  "CGO_ENABLED",
							Value: "0",
						},
					},
					Command: []string{constant.DUKKHA_TOOL_CMD, "build", "-o", "foo", "./"},
				},
			},
		},
	}

	ctx := dukkha_test.NewTestContext(context.TODO(), t.TempDir())

	tests.RunTaskExecSpecGenerationTests(t, ctx, testCases)
}
