package archive

import (
	"testing"

	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/pkg/tools/tests"
)

func TestTaskCreate(t *testing.T) {
	t.Parallel()

	type Check struct {
		rs.BaseField

		Files map[string]string `yaml:",inline"`
	}

	tests.TestTask(t, "./fixtures/create",
		&Tool{},
		func() dukkha.Task { return newCreateTask("") },
		func() rs.Field { return &Check{} },
		func(t *testing.T, e, a rs.Field) {
			expected, actual := e.(*Check), a.(*Check)
			assert.EqualValues(t, expected.Files, actual.Files)
		},
	)
}
