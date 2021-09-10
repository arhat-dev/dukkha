package dukkha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTaskReference(t *testing.T) {
	tests := []struct {
		name string

		input    string
		expected TaskReference

		expectErr bool
	}{
		{
			name:      "Invalid Empty",
			input:     "",
			expectErr: true,
		},
		{
			name:      "Invalid No Task Call",
			input:     "foo:bar",
			expectErr: true,
		},
		{
			name:      "Invalid Non Terminated Task Call",
			input:     "foo:bar(something",
			expectErr: true,
		},
		{
			name:      "Invalid No Prefix",
			input:     "(something)",
			expectErr: true,
		},
		{
			name:      "Invalid Prefix Too Few Parts",
			input:     "foo(something)",
			expectErr: true,
		},
		{
			name:      "Invalid Prefix Too Many Parts",
			input:     "foo:bar:fooBar:extra(something)",
			expectErr: true,
		},
		{
			name:  "Valid Default Matrix",
			input: "foo:bar(something)",

			expected: TaskReference{
				ToolKind: "foo",
				ToolName: "",
				TaskKind: "bar",
				TaskName: "something",
			},
		},
		{
			name:  "Valid With ToolName",
			input: "foo:tool:bar(something)",

			expected: TaskReference{
				ToolKind: "foo",
				ToolName: "tool",
				TaskKind: "bar",
				TaskName: "something",
			},
		},
		{
			name:  "Valid Custom Matrix",
			input: "foo:bar(something, {foo: [bar], bar: [foo, bar]})",

			expected: TaskReference{
				ToolKind: "foo",
				ToolName: "",
				TaskKind: "bar",
				TaskName: "something",
				MatrixFilter: map[string][]string{
					"foo": {"bar"},
					"bar": {"foo", "bar"},
				},
			},
		},
		{
			name:  "Valid Custom Matrix With ToolName",
			input: "foo:tool:bar(something, {foo: [bar], bar: [foo, bar]})",

			expected: TaskReference{
				ToolKind: "foo",
				ToolName: "tool",
				TaskKind: "bar",
				TaskName: "something",
				MatrixFilter: map[string][]string{
					"foo": {"bar"},
					"bar": {"foo", "bar"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ref, err := ParseTaskReference(test.input, "")
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.EqualValues(t, &test.expected, ref)
		})
	}
}

func TestParseBrackets(t *testing.T) {
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
			result, err := parseBrackets(test.toExpand)
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}
