package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestSize_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
	}{
		{
			name:     "B",
			input:    `128B`,
			expected: 128,
		},
		{
			name:     "KB",
			input:    `128KB`,
			expected: 128 * 1024,
		},
		{
			name:     "MB",
			input:    `128M`,
			expected: 128 * 1024 * 1024,
		},
		{
			name:     "GB",
			input:    `128G`,
			expected: 128 * 1024 * 1024 * 1024,
		},
		{
			name:     "TB",
			input:    `128T`,
			expected: 128 * 1024 * 1024 * 1024 * 1024,
		},
		{
			name:     "PB",
			input:    `128PB`,
			expected: 128 * 1024 * 1024 * 1024 * 1024 * 1024,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var s Size

			err := yaml.Unmarshal([]byte(test.input), &s)
			assert.NoError(t, err)

			assert.EqualValues(t, test.expected, s)
		})
	}
}
