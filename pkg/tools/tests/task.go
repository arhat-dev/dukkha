package tests

import (
	"testing"

	"arhat.dev/dukkha/pkg/field"
	"arhat.dev/dukkha/pkg/sliceutils"
	"arhat.dev/dukkha/pkg/tools"
	"github.com/stretchr/testify/assert"
)

type ExecSpecGenerationTestCase struct {
	Name     string
	Prepare  func() error
	Finalize func()

	Task      tools.Task
	Expected  []tools.TaskExecSpec
	ExpectErr bool
}

func RunTaskExecSpecGenerationTests(
	t *testing.T,
	rc *field.RenderingContext,
	toolCmd []string,
	tests []ExecSpecGenerationTestCase,
) {
	originalToolCmd := sliceutils.NewStringSlice(toolCmd)
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.Finalize != nil {
				defer test.Finalize()
			}

			if test.Prepare != nil {
				if !assert.NoError(t, test.Prepare(), "failed to prepare test environment") {
					return
				}
			}

			if test.ExpectErr {
				_, err := test.Task.GetExecSpecs(rc, toolCmd)
				assert.EqualValues(t, originalToolCmd, toolCmd, "task is not allowed to changed tool cmd")
				assert.Error(t, err)
				return
			}

			specs, err := test.Task.GetExecSpecs(rc, toolCmd)
			assert.EqualValues(t, originalToolCmd, toolCmd, "task is not allowed to changed tool cmd")
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, test.Expected, specs)
		})
	}
}
