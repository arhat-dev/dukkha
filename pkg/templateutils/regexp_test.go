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

	findFn := func(args ...any) (any, error) {
		return ns.Find(convertAnySlice[String](args)...)
	}

	findAllFn := func(args ...any) (any, error) {
		return ns.FindAll(convertAnySlice[String](args)...)
	}

	matchFn := func(args ...any) (any, error) {
		return ns.Match(convertAnySlice[String](args)...)
	}

	replaceFn := func(args ...any) (any, error) {
		return ns.Replace(convertAnySlice[String](args)...)
	}

	for _, test := range []struct {
		name string
		args []any

		fn func(...any) (any, error)

		expected any
	}{
		{
			name:     "Find dot-newline",
			args:     []any{"a.bc", "-s", "a\nbc"},
			fn:       findFn,
			expected: "a\nbc",
		},
		{
			name:     "FindAll",
			args:     []any{"abc", "abc,abc,abc,abc"},
			fn:       findAllFn,
			expected: []string{"abc", "abc", "abc", "abc"},
		},
		{
			name:     "FindAll Int N",
			args:     []any{"abc", 2, "abc,abc,abc,abc"},
			fn:       findAllFn,
			expected: []string{"abc", "abc"},
		},
		{
			name:     "Match case-insensitive",
			args:     []any{"^abc", "-i", "aBCdef"},
			fn:       matchFn,
			expected: true,
		},
		{
			name:     "Replace greedy",
			args:     []any{"a.*c", "foo", "abcdbc"},
			fn:       replaceFn,
			expected: "foo",
		},
		{
			name:     "Replace ungreedy",
			args:     []any{"a.*c", "-U", "foo", "abcdbc"},
			fn:       replaceFn,
			expected: "foodbc",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ret, err := test.fn(test.args...)
			assert.NoError(t, err)
			assert.EqualValues(t, test.expected, ret)
		})
	}
}
