package github_test

import (
	"testing"

	"arhat.dev/dukkha/pkg/tools/tests"
)

func TestTaskRelease_GetExecSpecs(t *testing.T) {
	// toolCmd := []string{"gh"}
	testCases := []tests.ExecSpecGenerationTestCase{}
	tests.RunTaskExecSpecGenerationTests(t, nil, testCases)
}
