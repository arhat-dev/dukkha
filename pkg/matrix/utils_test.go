package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCartesianProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		spec     map[string][]string
		expected []map[string]string
	}{
		{
			name:     "Empty",
			spec:     map[string][]string{},
			expected: nil,
		},
		{
			name: "Content Empty",
			spec: map[string][]string{
				"foo": {},
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, CartesianProduct(test.spec))
		})
	}
}
