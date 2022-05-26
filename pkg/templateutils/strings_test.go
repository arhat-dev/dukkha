package templateutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringsNS_Substr(t *testing.T) {
	for _, test := range []struct {
		args     []any
		expected string
	}{
		{
			args:     []any{2, "一二三四"},
			expected: "三四",
		},
		{
			args:     []any{1, 2, "一二三四"},
			expected: "二",
		},
		{
			args:     []any{2, -1, "一二三四"},
			expected: "三四",
		},
		{
			args:     []any{-2, "一二三四"},
			expected: "四",
		},
		{
			args:     []any{-1, "一二三四"},
			expected: "",
		},
	} {
		var ns stringsNS

		ret, err := ns.Substr(test.args...)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, ret)
	}
}

func TestStringsNS_NoSpace(t *testing.T) {
	for _, test := range []struct {
		in       string
		expected string
	}{
		{"abc", "abc"},   // no space
		{"     ", ""},    // all space
		{"  abc", "abc"}, // space at start
		{"abc  ", "abc"}, // space at end
		{"a  bc", "abc"}, // space in middle

		{"  \t一  二 \v三", "一二三"}, // space everywhere
	} {
		ret, err := stringsNS{}.NoSpace(test.in)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, ret)
	}
}

func TestStringsNS_AddPrefix(t *testing.T) {
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
			assert.Equal(t, test.expected, must(stringsNS{}.AddPrefix(test.prefix, test.sep, test.origin)))
		})
	}
}

func TestStringsNS_RemovePrefix(t *testing.T) {
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
			assert.Equal(t, test.expected, must(stringsNS{}.RemovePrefix(test.prefix, test.sep, test.origin)))
		})
	}
}

func TestStringsNS_AddSuffix(t *testing.T) {
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
			assert.Equal(t, test.expected, must(stringsNS{}.AddSuffix(test.prefix, test.sep, test.origin)))
		})
	}
}

func TestStringsNS_RemoveSuffix(t *testing.T) {
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
			assert.Equal(t, test.expected, must(stringsNS{}.RemoveSuffix(test.suffix, test.sep, test.origin)))
		})
	}
}
