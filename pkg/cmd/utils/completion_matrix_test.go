package utils

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	_ "arhat.dev/dukkha/cmd/dukkha/addon"
)

func TestHandleTaskMatrixCompletion(t *testing.T) {
	type Result struct {
		candidates []string
		directive  cobra.ShellCompDirective
	}

	for _, test := range []struct {
		name string

		existing   []string
		args       []string
		toComplete string

		expected Result
	}{
		{
			name:       "All",
			existing:   []string{"-m"},
			args:       []string{"workflow", "local", "run", "test"},
			toComplete: "",
			expected: Result{
				candidates: []string{
					"a=a1", "a=a2", "b=b", "b=c",
				},
				directive: cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Prefix a",
			existing:   []string{"-m"},
			args:       []string{"workflow", "local", "run", "test"},
			toComplete: "a",
			expected: Result{
				candidates: []string{
					"a=a1", "a=a2",
				},
				directive: cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Non Existing",
			existing:   []string{"-m"},
			args:       []string{"workflow", "local", "run", "test"},
			toComplete: "non-existing",
			expected: Result{
				candidates: nil,
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Dedup",
			existing:   []string{"-m", "a=a1"},
			args:       []string{"workflow", "local", "run", "test"},
			toComplete: "",
			expected: Result{
				candidates: []string{"a=a2", "b=b", "b=c"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
	} {
		ctx := newCompletionContext(t)

		t.Run(test.name, func(t *testing.T) {
			actualCandidates, directive := handleTaskMatrixCompletion(
				ctx, test.existing, test.args, test.toComplete,
			)
			assert.EqualValues(t, test.expected.candidates, actualCandidates)
			assert.EqualValues(t, test.expected.directive, directive)
		})
	}
}
