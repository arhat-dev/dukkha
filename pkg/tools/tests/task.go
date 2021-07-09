package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
)

type ExecSpecGenerationTestCase struct {
	Name     string
	Prepare  func() error
	Finalize func()

	Options   dukkha.TaskExecOptions
	Task      dukkha.Task
	Expected  []dukkha.TaskExecSpec
	ExpectErr bool
}

func RunTaskExecSpecGenerationTests(
	t *testing.T,
	taskCtx dukkha.TaskExecContext,
	tests []ExecSpecGenerationTestCase,
) {
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
				_, err := test.Task.GetExecSpecs(taskCtx, test.Options)
				assert.Error(t, err)
				return
			}

			specs, err := test.Task.GetExecSpecs(taskCtx, test.Options)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, test.Expected, specs)
		})
	}
}
