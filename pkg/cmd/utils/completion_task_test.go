package utils

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestHandleTaskCompletion(t *testing.T) {
	t.Parallel()

	type Result struct {
		candidates []string
		directive  cobra.ShellCompDirective
	}

	for _, test := range []struct {
		name string

		args       []string
		toComplete string

		expected Result
	}{
		{
			name:       "Tool Kind",
			args:       []string{},
			toComplete: "",
			expected: Result{
				candidates: []string{"workflow"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Partial Tool Kind",
			args:       nil,
			toComplete: "work",
			expected: Result{
				candidates: []string{"workflow"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Complete Tool Kind",
			args:       []string{},
			toComplete: "workflow",
			expected: Result{
				candidates: []string{"workflow"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Non Existing Tool Kind",
			args:       []string{},
			toComplete: "FOO",
			expected: Result{
				candidates: nil,
				directive:  cobra.ShellCompDirectiveNoSpace,
			},
		},

		{
			name:       "Tool Name",
			args:       []string{"workflow"},
			toComplete: "",
			expected: Result{
				candidates: []string{"local"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Partial Tool Name",
			args:       []string{"workflow"},
			toComplete: "l",
			expected: Result{
				candidates: []string{"local"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Complete Tool Name",
			args:       []string{"workflow"},
			toComplete: "local",
			expected: Result{
				candidates: []string{"local"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Non Existing Tool Name",
			args:       []string{"workflow"},
			toComplete: "FOO",
			expected: Result{
				candidates: nil,
				directive:  cobra.ShellCompDirectiveNoSpace,
			},
		},

		{
			name:       "Task Kind",
			args:       []string{"workflow", "local"},
			toComplete: "",
			expected: Result{
				candidates: []string{"run"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Partial Task Kind",
			args:       []string{"workflow", "local"},
			toComplete: "r",
			expected: Result{
				candidates: []string{"run"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Complete Task Kind",
			args:       []string{"workflow", "local"},
			toComplete: "run",
			expected: Result{
				candidates: []string{"run"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Non Existing Task Kind",
			args:       []string{"workflow", "local"},
			toComplete: "FOO",
			expected: Result{
				candidates: nil,
				directive:  cobra.ShellCompDirectiveNoSpace,
			},
		},

		{
			name:       "Task Name",
			args:       []string{"workflow", "local", "run"},
			toComplete: "",
			expected: Result{
				candidates: []string{"test"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Partial Task Name",
			args:       []string{"workflow", "local", "run"},
			toComplete: "te",
			expected: Result{
				candidates: []string{"test"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Complete Task Name",
			args:       []string{"workflow", "local", "run"},
			toComplete: "test",
			expected: Result{
				candidates: []string{"test"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name:       "Non Existing Task Name",
			args:       []string{"workflow", "local", "run"},
			toComplete: "FOO",
			expected: Result{
				candidates: nil,
				directive:  cobra.ShellCompDirectiveNoSpace,
			},
		},

		{
			name:       "Suggest Matrix",
			args:       []string{"workflow", "local", "run", "test"},
			toComplete: "",
			expected: Result{
				candidates: []string{"-m"},
				directive:  cobra.ShellCompDirectiveNoFileComp,
			},
		},
	} {
		ctx := newCompletionContext(t)

		t.Run(test.name, func(t *testing.T) {
			actualCandidates, directive := handleTaskCompletion(
				ctx, test.args, test.toComplete,
			)
			assert.EqualValues(t, test.expected.candidates, actualCandidates)
			assert.EqualValues(t, test.expected.directive, directive)
		})
	}
}
