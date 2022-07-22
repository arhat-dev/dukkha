package templateutils

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func convertAnySlice[T any](x []any) []T {
	return *(*[]T)(unsafe.Pointer(&x))
}

func TestRegexpNS(t *testing.T) {
	var ns regexpNS

	findFirstFn := func(args ...any) (any, error) {
		return ns.FindFirst(convertAnySlice[String](args)...)
	}

	findNFn := func(args ...any) (any, error) {
		return ns.FindN(convertAnySlice[String](args)...)
	}

	findAllFn := func(args ...any) (any, error) {
		return ns.FindAll(convertAnySlice[String](args)...)
	}

	matchFn := func(args ...any) (any, error) {
		return ns.Match(convertAnySlice[String](args)...)
	}

	replaceFirstFn := func(args ...any) (any, error) {
		return ns.ReplaceFirst(convertAnySlice[String](args)...)
	}

	replaceAllFn := func(args ...any) (any, error) {
		return ns.ReplaceAll(convertAnySlice[String](args)...)
	}

	for _, test := range []struct {
		name string
		args []any

		fn func(...any) (any, error)

		expected any
	}{
		{
			name:     "FindFirst dot match",
			args:     []any{".", "--dot-newline", "abc"},
			fn:       findFirstFn,
			expected: "a",
		},
		{
			name:     "FindN dot match",
			args:     []any{".", "--dot-newline", 2, "abc"},
			fn:       findNFn,
			expected: []string{"a", "b"},
		},
		{
			name:     "FindAll",
			args:     []any{"a", "abc,abc,abc,abc"},
			fn:       findAllFn,
			expected: []string{"a", "a", "a", "a"},
		},
		{
			name:     "Match case-insensitive",
			args:     []any{"^abc", "-i", "aBCdef"},
			fn:       matchFn,
			expected: true,
		},
		{
			name:     "Replace",
			args:     []any{"a", "foo", "aaa"},
			fn:       replaceFirstFn,
			expected: "fooaa",
		},
		{
			name:     "Replace greedy",
			args:     []any{"a.*c", "foo", "abcdbc"},
			fn:       replaceFirstFn,
			expected: "foo",
		},
		{
			name:     "Replace ungreedy",
			args:     []any{"a.*c", "-U", "foo", "abcdbc"},
			fn:       replaceFirstFn,
			expected: "foodbc",
		},
		{
			name:     "ReplaceAll",
			args:     []any{"a", "foo", "aaa"},
			fn:       replaceAllFn,
			expected: "foofoofoo",
		},
		{
			name:     `Replace including newline`,
			args:     []any{`//.*build.*\n`, "", "//go:build darwin && go1.18 && !go1.19\n// +build darwin,go1.18,!go1.19\n"},
			fn:       replaceAllFn,
			expected: "",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ret, err := test.fn(test.args...)
			assert.NoError(t, err)
			assert.EqualValues(t, test.expected, ret)
		})
	}
}
