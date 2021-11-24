package dukkha

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mvdan.cc/sh/v3/expand"
)

func TestGenEnvForValues(t *testing.T) {
	tests := []struct {
		name     string
		values   map[string]interface{}
		expected map[string]expand.Variable
	}{
		{
			name:   "Simple",
			values: map[string]interface{}{"foo": "bar"},
			expected: map[string]expand.Variable{
				"values.foo": createVariable(`bar`),
			},
		},
		{
			name: "Nested_1",
			values: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "woo",
				},
			},
			expected: map[string]expand.Variable{
				"values.foo":     createVariable(`{"bar":"woo"}`),
				"values.foo.bar": createVariable(`woo`),
			},
		},
		{
			name: "Nested_2",
			values: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"woo": "few",
					},
				},
			},
			expected: map[string]expand.Variable{
				"values.foo":         createVariable(`{"bar":{"woo":"few"}}`),
				"values.foo.bar":     createVariable(`{"woo":"few"}`),
				"values.foo.bar.woo": createVariable(`few`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values, err := genEnvForValues(test.values)
			assert.NoError(t, err)
			assert.EqualValues(t, test.expected, values)
		})
	}
}
