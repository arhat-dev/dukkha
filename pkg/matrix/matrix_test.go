package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatrixConfig_GetSpecs(t *testing.T) {
	tests := []struct {
		name string
		in   Spec
		// specs are sorted by name, put them in order
		expected []Entry
	}{
		{
			name: "normal",
			in: Spec{
				Kernel: []string{"linux", "windows", "darwin"},
				Arch:   []string{"amd64", "arm64"},
				Custom: map[string][]string{
					"foo": {"a", "b"},
				},
			},
			expected: []Entry{
				// sort order: arch=amd64 foo=a,b, os=linux,windows,darwin
				{
					"kernel": "linux",
					"arch":   "amd64",
					"foo":    "a",
				},
				{
					"kernel": "windows",
					"arch":   "amd64",
					"foo":    "a",
				},
				{
					"kernel": "darwin",
					"arch":   "amd64",
					"foo":    "a",
				},
				{
					"kernel": "linux",
					"arch":   "amd64",
					"foo":    "b",
				},
				{
					"kernel": "windows",
					"arch":   "amd64",
					"foo":    "b",
				},
				{
					"kernel": "darwin",
					"arch":   "amd64",
					"foo":    "b",
				},

				// sort order: arch=arm64 foo=a,b, os=linux,windows,darwin

				{
					"kernel": "linux",
					"arch":   "arm64",
					"foo":    "a",
				},
				{
					"kernel": "windows",
					"arch":   "arm64",
					"foo":    "a",
				},
				{
					"kernel": "darwin",
					"arch":   "arm64",
					"foo":    "a",
				},

				{
					"kernel": "linux",
					"arch":   "arm64",
					"foo":    "b",
				},
				{
					"kernel": "windows",
					"arch":   "arm64",
					"foo":    "b",
				},
				{
					"kernel": "darwin",
					"arch":   "arm64",
					"foo":    "b",
				},
			},
		},
		{
			name: "include",
			in: Spec{
				Include: []map[string][]string{
					{
						"kernel": []string{"windows"},
						"arch":   []string{"arm64", "amd64"},
					},
					{
						"kernel": []string{"darwin"},
						"arch":   []string{"arm64"},
					},
				},
				Kernel: []string{"linux"},
				Arch:   []string{"amd64"},
			},
			expected: []Entry{
				{
					"kernel": "linux",
					"arch":   "amd64",
				},
				{
					"kernel": "windows",
					"arch":   "arm64",
				},
				{
					"kernel": "windows",
					"arch":   "amd64",
				},
				{
					"kernel": "darwin",
					"arch":   "arm64",
				},
			},
		},
		{
			name: "exclude",
			in: Spec{
				Exclude: []map[string][]string{
					{
						"kernel": []string{"linux"},
						"arch":   []string{"arm64", "amd64"},
					},
				},
				Kernel: []string{"linux"},
				Arch:   []string{"amd64", "arm64"},
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(t, test.expected, test.in.GetSpecs(nil, "", ""))
		})
	}
}
