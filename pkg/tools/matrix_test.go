package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatrixConfig_GetSpecs(t *testing.T) {
	tests := []struct {
		name string
		in   MatrixConfig
		// specs are sorted by name, put them in order
		expected []MatrixSpec
	}{
		{
			name: "normal",
			in: MatrixConfig{
				OS:   []string{"linux", "windows", "darwin"},
				Arch: []string{"amd64", "arm64"},
				Custom: map[string][]string{
					"foo": {"a", "b"},
				},
			},
			expected: []MatrixSpec{
				// sort order: arch=amd64 foo=a,b, os=linux,windows,darwin
				{
					"os":   "linux",
					"arch": "amd64",
					"foo":  "a",
				},
				{
					"os":   "windows",
					"arch": "amd64",
					"foo":  "a",
				},
				{
					"os":   "darwin",
					"arch": "amd64",
					"foo":  "a",
				},
				{
					"os":   "linux",
					"arch": "amd64",
					"foo":  "b",
				},
				{
					"os":   "windows",
					"arch": "amd64",
					"foo":  "b",
				},
				{
					"os":   "darwin",
					"arch": "amd64",
					"foo":  "b",
				},

				// sort order: arch=arm64 foo=a,b, os=linux,windows,darwin

				{
					"os":   "linux",
					"arch": "arm64",
					"foo":  "a",
				},
				{
					"os":   "windows",
					"arch": "arm64",
					"foo":  "a",
				},
				{
					"os":   "darwin",
					"arch": "arm64",
					"foo":  "a",
				},

				{
					"os":   "linux",
					"arch": "arm64",
					"foo":  "b",
				},
				{
					"os":   "windows",
					"arch": "arm64",
					"foo":  "b",
				},
				{
					"os":   "darwin",
					"arch": "arm64",
					"foo":  "b",
				},
			},
		},
		{
			name: "include",
			in: MatrixConfig{
				Include: []map[string][]string{
					{
						"os":   []string{"windows"},
						"arch": []string{"arm64", "amd64"},
					},
					{
						"os":   []string{"darwin"},
						"arch": []string{"arm64"},
					},
				},
				OS:   []string{"linux"},
				Arch: []string{"amd64"},
			},
			expected: []MatrixSpec{
				{
					"os":   "linux",
					"arch": "amd64",
				},
				{
					"os":   "windows",
					"arch": "arm64",
				},
				{
					"os":   "windows",
					"arch": "amd64",
				},
				{
					"os":   "darwin",
					"arch": "arm64",
				},
			},
		},
		{
			name: "exclude",
			in: MatrixConfig{
				Exclude: []map[string][]string{
					{
						"os":   []string{"linux"},
						"arch": []string{"arm64", "amd64"},
					},
				},
				OS:   []string{"linux"},
				Arch: []string{"amd64", "arm64"},
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(t, test.expected, test.in.GetSpecs())
		})
	}
}
