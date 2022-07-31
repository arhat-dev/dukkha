package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifiableString_Ext(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"foo", ""},
		{"foo.json", ".json"},
		{"foo/json", ""},
		{"foo.无常", ".无常"},
	} {
		t.Run(test.input, func(t *testing.T) {
			assert.Equal(t, test.expected, IdentifiableString(test.input).Ext())
		})
	}
}
