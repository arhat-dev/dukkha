package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"arhat.dev/dukkha/pkg/field"

	_ "embed"
)

func TestMatrixConfig_GenerateEntries(t *testing.T) {
	tests := []struct {
		name string
		in   Spec
		// specs are sorted by name, put them in order
		expected []Entry

		filter map[string][]string
	}{
		{
			name: "normal",
			in: Spec{
				Kernel: []string{"linux", "darwin"},
				Arch:   []string{"amd64", "arm64"},
				Custom: map[string][]string{
					"foo": {"a", "b"},
				},
			},
			expected: []Entry{
				// sort order: arch=amd64 foo=a,b, os=linux,darwin
				{
					"kernel": "linux",
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
					"kernel": "darwin",
					"arch":   "amd64",
					"foo":    "b",
				},

				// sort order: arch=arm64 foo=a,b, os=linux,darwin

				{
					"kernel": "linux",
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
					"kernel": "darwin",
					"arch":   "arm64",
					"foo":    "b",
				},
			},
		},
		{
			name: "include",
			in: Spec{
				Include: []*specItem{
					{
						Data: map[string][]string{
							"kernel": {"windows"},
							"arch":   {"arm64", "amd64"},
						},
					},
					{
						Data: map[string][]string{
							"kernel": {"darwin"},
							"arch":   {"arm64"},
						},
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
			name: "include+filter",
			in: Spec{
				Include: []*specItem{
					{
						Data: map[string][]string{
							"kernel": {"aix"},
							"arch":   {"ppc64le"},
						},
					},
					{
						Data: map[string][]string{
							"kernel": {"darwin"},
							"arch":   {"arm64"},
						},
					},
				},
				Kernel: []string{"linux"},
				Arch:   []string{"amd64", "arm64"},
			},
			filter: map[string][]string{
				"arch": {"amd64"},
			},
			expected: []Entry{
				{
					"kernel": "linux",
					"arch":   "amd64",
				},
			},
		},
		{
			name: "exclude-all-full-match",
			in: Spec{
				Exclude: []*specItem{
					{
						Data: map[string][]string{
							"kernel": {"linux"},
							"arch":   {"arm64", "amd64"},
						},
					},
				},
				Kernel: []string{"linux"},
				Arch:   []string{"amd64", "arm64"},
			},
			expected: nil,
		},
		{
			name: "exclude-all-single-match",
			in: Spec{
				Exclude: []*specItem{
					{
						Data: map[string][]string{
							"kernel": {"linux"},
						},
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
			assert.EqualValues(
				t,
				test.expected,
				test.in.GenerateEntries(test.filter, "", ""),
			)
		})
	}
}

var (
	//go:embed _fixtures/001-filter-amd64-got-unwanted-aix.yaml
	fitlerAMD64GotUnwantedAIX []byte
)

func TestMatrixConfig_GenerateEntries_Fixture(t *testing.T) {
	tests := []struct {
		name           string
		yamlMatrixSpec []byte
		filter         map[string][]string
		expected       []Entry
	}{
		{
			name:           "001-filter-amd64-got-unwanted-aix",
			yamlMatrixSpec: fitlerAMD64GotUnwantedAIX,
			filter:         map[string][]string{"arch": {"amd64"}},
			expected: []Entry{
				{"arch": "amd64", "kernel": "linux"},
				{"arch": "amd64", "kernel": "darwin"},
				// {"arch": "amd64", "kernel": "aix"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spec := field.Init(&Spec{}, nil).(*Spec)
			if !assert.NoError(t, yaml.Unmarshal(test.yamlMatrixSpec, spec)) {
				return
			}

			entries := spec.GenerateEntries(test.filter, "", "")
			assert.EqualValues(t, test.expected, entries)
		})
	}
}
