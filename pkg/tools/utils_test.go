package tools

import (
	"reflect"
	"testing"

	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
)

func TestSeparateBaseAndImpl(t *testing.T) {
	t.Parallel()

	forBase, forImpl := separateBaseAndImpl("base.", []string{"base.foo", "bar"})

	assert.EqualValues(t, []string{"foo"}, forBase)
	assert.EqualValues(t, []string{"bar"}, forImpl)
}

func TestGetTagNamesToResolve(t *testing.T) {
	t.Parallel()

	expectedTagNames := []string{
		"a",
		"bar",
		"NestedInlineMap",
		"c",
		"non_anonymous_task",
		"inner_field",
	}

	t.Run("tool", func(t *testing.T) {
		type TestTool struct {
			rs.BaseField

			// top level anonymous BaseTask should be ignored
			BaseTask

			A            string
			B            string `yaml:"bar"`
			InlineStruct struct {
				rs.BaseField

				NestedInlineMap    map[string]int `yaml:",inline"`
				NestedInlineStruct struct {
					rs.BaseField

					C string
				} `yaml:",inline"`
			} `yaml:",inline"`

			Ignored string `yaml:"-"`

			NonAnonymousTask BaseTask `yaml:"non_anonymous_task"`

			InnerField struct {
				// we do not care about fields inside inner field
				BaseTask
				Foo string `yaml:"foo"`
			} `yaml:"inner_field"`
		}

		actualTagNames := getTagNamesToResolve(reflect.TypeOf(&TestTool{}).Elem())
		assert.EqualValues(t, expectedTagNames, actualTagNames)

		tool := rs.Init(&TestTool{}, nil).(*TestTool)
		assert.NoError(t, tool.ResolveFields(nil, -1, expectedTagNames...))
		assert.Error(t, tool.ResolveFields(nil, -1, "non-existing"))
	})

	t.Run("task", func(t *testing.T) {
		type TestTask struct {
			rs.BaseField

			// top level anonymous BaseTask should be ignored
			BaseTask

			A            string
			B            string `yaml:"bar"`
			InlineStruct struct {
				rs.BaseField

				NestedInlineMap    map[string]int `yaml:",inline"`
				NestedInlineStruct struct {
					rs.BaseField

					C string
				} `yaml:",inline"`
			} `yaml:",inline"`

			Ignored string `yaml:"-"`

			NonAnonymousTask BaseTask `yaml:"non_anonymous_task"`

			InnerField struct {
				// we do not care about fields inside inner field
				BaseTask
				Foo string `yaml:"foo"`
			} `yaml:"inner_field"`
		}

		actualTagNames := getTagNamesToResolve(reflect.TypeOf(&TestTask{}).Elem())
		assert.EqualValues(t, expectedTagNames, actualTagNames)

		task := rs.Init(&TestTask{}, nil).(*TestTask)
		assert.NoError(t, task.ResolveFields(nil, -1, expectedTagNames...))
		assert.Error(t, task.ResolveFields(nil, -1, "non-existing"))
	})
}
