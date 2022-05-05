package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBrackets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		toExpand  string
		expected  string
		expectErr bool
	}{
		{
			name:     "Valid Simple",
			toExpand: "foo)(",
			expected: "foo",
		},
		{
			name:     "Valid Simple 2",
			toExpand: "foo)))))",
			expected: "foo",
		},
		{
			name:     "Valid Empty",
			toExpand: "))",
			expected: "",
		},
		{
			name:     "Valid One Pair",
			toExpand: "foo())",
			expected: "foo()",
		},
		{
			name:     "Valid Many Pairs",
			toExpand: "foo()()()())",
			expected: "foo()()()()",
		},
		{
			name:      "Invalid",
			toExpand:  "foo",
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ParseBrackets(test.toExpand)
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}
