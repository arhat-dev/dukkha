package cosign

import (
	"testing"

	"arhat.dev/rs"

	"arhat.dev/dukkha/pkg/tools"
	"arhat.dev/dukkha/pkg/tools/tests"
)

func TestTaskSign(t *testing.T) {
	t.Parallel()

	type Expected struct {
		rs.BaseField

		Signature string `yaml:"signature"`
	}

	t.SkipNow()

	tests.TestTask(t, "./fixtures/sign", &Tool{},
		func() *TaskSign { return tools.NewTask[TaskSign, *TaskSign]("test").(*TaskSign) },
		func() *Expected { return &Expected{} },
		func(t *testing.T, exp, actual *Expected) {
			_, _ = exp, actual
			// assert.EqualValues(t, exp.Signature, actual.Signature)
		},
	)
}
