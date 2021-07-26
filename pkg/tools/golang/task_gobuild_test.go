package golang

import (
	"context"
	"testing"

	"arhat.dev/dukkha/pkg/dukkha"
	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/tests"
)

func TestTaskRelease_GetExecSpecs(t *testing.T) {
	toolCmd := []string{"go"}
	testCases := []tests.ExecSpecGenerationTestCase{
		{
			Name: "Default Build Task",
			Task: &TaskBuild{
				BaseTask: tools.BaseTask{
					TaskName: "foo",
				},
			},
			Options: dukkha_test.CreateTaskMatrixExecOptions(toolCmd),
			Expected: []dukkha.TaskExecSpec{
				{
					Env:     []string{"CGO_ENABLED=0"},
					Command: []string{"go", "build", "-o", "foo", "./"},
				},
			},
		},
	}

	ctx := dukkha_test.NewTestContext(context.TODO())

	tests.RunTaskExecSpecGenerationTests(t, ctx, testCases)
}