package archive

import (
	"testing"

	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/tools"
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
		func() *TaskCreate { return tools.NewTask[TaskCreate, *TaskCreate]("").(*TaskCreate) },
		func() *Check { return &Check{} },
		func(t *testing.T, expected, actual *Check) {
			assert.EqualValues(t, expected.Files, actual.Files)
		},
	)
}
