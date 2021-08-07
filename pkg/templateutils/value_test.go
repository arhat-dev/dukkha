package templateutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenNewVal(t *testing.T) {
	tests := []struct {
		name string
		key  string

		value    interface{}
		expected map[string]interface{}

		expectErr bool
	}{
		{
			name:  "simple",
			key:   "foo",
			value: "bar",
			expected: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			name:  "quoted",
			key:   `"foo.bar"`,
			value: "woo",
			expected: map[string]interface{}{
				"foo.bar": "woo",
			},
		},
		{
			name:  "nested",
			key:   "foo.bar.woo",
			value: "koo",
			expected: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"woo": "koo",
					},
				},
			},
		},
		{
			name:  "nested quoted",
			key:   `foo."bar.woo"`,
			value: "koo",
			expected: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar.woo": "koo",
				},
			},
		},
		{
			name:  "invalid quoted",
			key:   `"foo`,
			value: "bar",

			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret := make(map[string]interface{})
			err := genNewVal(test.key, test.value, &ret)
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.EqualValues(t, test.expected, ret)
		})
	}
}
