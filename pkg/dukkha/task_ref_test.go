package dukkha

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/matrix"
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
				MatrixFilter: matrix.NewFilter(map[string][]string{
					"foo": {"bar"},
					"bar": {"foo", "bar"},
				}),
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
				MatrixFilter: matrix.NewFilter(map[string][]string{
					"foo": {"bar"},
					"bar": {"foo", "bar"},
				}),
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
