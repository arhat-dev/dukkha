package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuneKeyFunc(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected []string
	}{
		{
			name:     "Empty",
			key:      "",
			expected: nil,
		},
		{
			name:     "ASCII",
			key:      "AbC",
			expected: []string{"A", "b", "C"},
		},
		{
			name:     "Unicode",
			key:      "测🙂试",
			expected: []string{"测", "🙂", "试"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var keys []string
			var (
				k    string
				next int
			)

			for k, next = RuneKeyFunc(test.key, 0); next > 0; k, next = RuneKeyFunc(test.key, next) {
				keys = append(keys, k)
			}

			assert.EqualValues(t, test.expected, keys)
			assert.Equal(t, -1, next)
		})
	}
}
