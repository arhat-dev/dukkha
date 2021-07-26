package tests

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/field"
)

type ExecSpecGenerationTestCase struct {
	Name     string
	Prepare  func() error
	Finalize func()

	Options   dukkha.TaskMatrixExecOptions
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

			field.InitRecursively(reflect.ValueOf(test.Task), nil)

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
