package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseShellEval(t *testing.T) {
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
			result, err := ParseShellEval(test.toExpand)
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestParseTaskReference(t *testing.T) {
	tests := []struct {
		name string

		input string

		expectedToolKind   string
		expectedToolName   string
		expectedTaskKind   string
		expectedTaskName   string
		expectedMatrixSpec map[string]string

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

			expectedToolKind: "foo",
			expectedToolName: "",
			expectedTaskKind: "bar",
			expectedTaskName: "something",
		},
		{
			name:  "Valid With ToolName",
			input: "foo:tool:bar(something)",

			expectedToolKind: "foo",
			expectedToolName: "tool",
			expectedTaskKind: "bar",
			expectedTaskName: "something",
		},
		{
			name:  "Valid Custom Matrix",
			input: "foo:bar(something, {foo: bar, bar: foo})",

			expectedToolKind:   "foo",
			expectedToolName:   "",
			expectedTaskKind:   "bar",
			expectedTaskName:   "something",
			expectedMatrixSpec: map[string]string{"foo": "bar", "bar": "foo"},
		},
		{
			name:  "Valid Custom Matrix With ToolName",
			input: "foo:tool:bar(something, {foo: bar, bar: foo})",

			expectedToolKind:   "foo",
			expectedToolName:   "tool",
			expectedTaskKind:   "bar",
			expectedTaskName:   "something",
			expectedMatrixSpec: map[string]string{"foo": "bar", "bar": "foo"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			toolKind, toolName, taskKind, taskName, ms, err := ParseTaskReference(test.input)
			if test.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expectedToolKind, toolKind)
			assert.Equal(t, test.expectedToolName, toolName)
			assert.Equal(t, test.expectedTaskKind, taskKind)
			assert.Equal(t, test.expectedTaskName, taskName)
			assert.EqualValues(t, test.expectedMatrixSpec, ms)
		})
	}
}
