package templateutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		origin string
		prefix string
		sep    string

		expected string
	}{
		{
			name:     "Simple",
			origin:   "foo\nfoo\n",
			prefix:   "- ",
			sep:      "\n",
			expected: "- foo\n- foo\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, stringsNS{}.AddPrefix(test.origin, test.prefix, test.sep))
		})
	}
}

func TestRemovePrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		origin string
		prefix string
		sep    string

		expected string
	}{
		{
			name:     "Simple",
			origin:   "barfoo\nbarfoo\n",
			prefix:   "bar",
			sep:      "\n",
			expected: "foo\nfoo\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, stringsNS{}.RemovePrefix(test.origin, test.prefix, test.sep))
		})
	}
}

func TestAddSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		origin string
		prefix string
		sep    string

		expected string
	}{
		{
			name:     "Simple",
			origin:   "foo\nfoo\n",
			prefix:   "-suffix",
			sep:      "\n",
			expected: "foo-suffix\nfoo-suffix\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, stringsNS{}.AddSuffix(test.origin, test.prefix, test.sep))
		})
	}
}

func TestRemoveSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		origin string
		suffix string
		sep    string

		expected string
	}{
		{
			name:     "Simple",
			origin:   "barfoo\nbarfoo\n",
			suffix:   "foo",
			sep:      "\n",
			expected: "bar\nbar\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, stringsNS{}.RemoveSuffix(test.origin, test.suffix, test.sep))
		})
	}
}
