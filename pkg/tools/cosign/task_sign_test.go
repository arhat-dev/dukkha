package cosign

import (
	"testing"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools/tests"
	"arhat.dev/rs"
)

func TestTaskSign(t *testing.T) {
	type Expected struct {
		rs.BaseField

		Signature string `yaml:"signature"`
	}

	t.SkipNow()

	tests.TestTask(t, "./fixtures/sign", &Tool{},
		func() dukkha.Task { return newTaskSign("test") },
		func() rs.Field { return &Expected{} },
		func(e, a rs.Field) {
			exp, actual := e.(*Expected), a.(*Expected)
			_, _ = exp, actual
			// assert.EqualValues(t, exp.Signature, actual.Signature)
		},
	)
}
