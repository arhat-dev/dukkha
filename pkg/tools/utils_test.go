package tools

import (
	"reflect"
	"testing"

	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
)

func TestGetTagNamesToResolve(t *testing.T) {
	t.Parallel()

	t.Run("tool", func(t *testing.T) {
		expectedTagNames := []string{
			"name",
			"env",
			"cmd",

			"a",
			"bar",
			"NestedInlineMap",
			"c",
			"non_anonymous_task",
			"inner_field",
		}

		actualTagNames := getTagNamesToResolve(reflect.TypeOf(&TestTool{}).Elem())
		assert.EqualValues(t, expectedTagNames, actualTagNames)

		tool := rs.Init(&TestTool{}, nil).(*TestTool)
		assert.NoError(t, tool.ResolveFields(nil, -1, expectedTagNames...))
		assert.Error(t, tool.ResolveFields(nil, -1, "non-existing"))
	})

	t.Run("task", func(t *testing.T) {
		expectedTagNames := []string{
			"name",
			"env",
			"matrix",
			"hooks",
			"continue_on_error",

			"a",
			"bar",
			"NestedInlineMap",
			"c",
			"non_anonymous_task",
			"inner_field",
		}

		actualTagNames := getTagNamesToResolve(reflect.TypeOf(&TestTask{}).Elem())
		assert.EqualValues(t, expectedTagNames, actualTagNames)

		task := rs.Init(&TestTask{}, nil).(*TestTask)
		assert.NoError(t, task.ResolveFields(nil, -1, expectedTagNames...))
		assert.Error(t, task.ResolveFields(nil, -1, "non-existing"))
	})
}
