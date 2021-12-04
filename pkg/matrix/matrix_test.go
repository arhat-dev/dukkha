package matrix

import (
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
)

func TestMatrixConfig_GenerateEntries(t *testing.T) {
	tests := []struct {
		name string
		in   Spec
		// specs are sorted by name, put them in order
		expected []Entry

		matchFilter  map[string]*Vector
		ignoreFilter [][2]string
	}{
		{
			name: "basic",
			in: Spec{
				Kernel: NewVector("linux", "darwin"),
				Arch:   NewVector("amd64", "arm64"),
				Custom: map[string]*Vector{
					"foo": NewVector("a", "b"),
				},
			},
			expected: []Entry{
				// sort order: arch=amd64 foo=a,b, os=linux,darwin
				{"kernel": "linux", "arch": "amd64", "foo": "a"},
				{"kernel": "darwin", "arch": "amd64", "foo": "a"},
				{"kernel": "linux", "arch": "amd64", "foo": "b"},
				{"kernel": "darwin", "arch": "amd64", "foo": "b"},

				// sort order: arch=arm64 foo=a,b, os=linux,darwin
				{"kernel": "linux", "arch": "arm64", "foo": "a"},
				{"kernel": "darwin", "arch": "arm64", "foo": "a"},
				{"kernel": "linux", "arch": "arm64", "foo": "b"},
				{"kernel": "darwin", "arch": "arm64", "foo": "b"},
			},
		},
		{
			name: "basic + matchFilter",
			in: Spec{
				Kernel: NewVector("linux", "darwin"),
				Arch:   NewVector("amd64", "arm64"),
				Custom: map[string]*Vector{
					"foo": NewVector("a", "b"),
				},
			},
			matchFilter: map[string]*Vector{
				"kernel": NewVector("linux"),
				"arch":   NewVector("arm64"),
			},
			expected: []Entry{
				// sort order: foo=a,b
				{"kernel": "linux", "arch": "arm64", "foo": "a"},
				{"kernel": "linux", "arch": "arm64", "foo": "b"},
			},
		},
		{
			name: "basic + ignoreFilter",
			in: Spec{
				Kernel: NewVector("linux", "darwin"),
				Arch:   NewVector("amd64", "arm64"),
				Custom: map[string]*Vector{
					"foo": NewVector("a", "b"),
				},
			},
			ignoreFilter: [][2]string{
				{"kernel", "linux"},
				{"arch", "arm64"},
			},
			expected: []Entry{
				// sort order: arch=amd64 foo=a,b, os=linux,darwin
				{"kernel": "darwin", "arch": "amd64", "foo": "a"},
				{"kernel": "darwin", "arch": "amd64", "foo": "b"},
			},
		},
		{
			name: "include",
			in: Spec{
				Include: []*specItem{
					{
						Data: map[string]*Vector{
							"kernel": NewVector("windows"),
							"arch":   NewVector("arm64", "amd64"),
						},
					},
					{
						Data: map[string]*Vector{
							"kernel": NewVector("darwin"),
							"arch":   NewVector("arm64"),
						},
					},
				},
				Kernel: NewVector("linux"),
				Arch:   NewVector("amd64"),
			},
			expected: []Entry{
				{"kernel": "linux", "arch": "amd64"},
				{"kernel": "windows", "arch": "arm64"},
				{"kernel": "windows", "arch": "amd64"},
				{"kernel": "darwin", "arch": "arm64"},
			},
		},
		{
			name: "include + MatchFilter",
			in: Spec{
				Include: []*specItem{
					{
						Data: map[string]*Vector{
							"kernel": NewVector("aix"),
							"arch":   NewVector("ppc64le"),
						},
					},
					{
						Data: map[string]*Vector{
							"kernel": NewVector("darwin"),
							"arch":   NewVector("arm64"),
						},
					},
				},
				Kernel: NewVector("linux"),
				Arch:   NewVector("amd64", "arm64"),
			},
			matchFilter: map[string]*Vector{
				"arch": NewVector("amd64"),
			},
			expected: []Entry{
				{"kernel": "linux", "arch": "amd64"},
			},
		},
		{
			name: "exclude-all-full-match",
			in: Spec{
				Exclude: []*specItem{
					{
						Data: map[string]*Vector{
							"kernel": NewVector("linux"),
							"arch":   NewVector("arm64", "amd64"),
						},
					},
				},
				Kernel: NewVector("linux"),
				Arch:   NewVector("amd64", "arm64"),
			},
			expected: nil,
		},
		{
			name: "exclude-all-single-match",
			in: Spec{
				Exclude: []*specItem{
					{
						Data: map[string]*Vector{
							"kernel": NewVector("linux"),
						},
					},
				},
				Kernel: NewVector("linux"),
				Arch:   NewVector("amd64", "arm64"),
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(
				t,
				test.expected,
				test.in.GenerateEntries(&Filter{
					match:  test.matchFilter,
					ignore: test.ignoreFilter,
				}, "", ""),
			)
		})
	}
}

func TestSpec_GenerateEntries_Fixture(t *testing.T) {
	type testInputSpec struct {
		rs.BaseField

		MatchFilter  map[string]*Vector `yaml:"match_filter"`
		IgnoreFilter [][2]string        `yaml:"ignore_fitler"`
		Spec         Spec               `yaml:"spec"`
	}

	testhelper.TestFixtures(t, "./fixtures/gen-entries",
		func() interface{} { return rs.Init(&testInputSpec{}, nil).(*testInputSpec) },
		func() interface{} {
			var data []Entry
			return &data
		},
		func(t *testing.T, in, exp interface{}) {
			spec := in.(*testInputSpec)
			err := spec.ResolveFields(rs.RenderingHandleFunc(
				func(renderer string, rawData interface{}) (result []byte, err error) {
					data, err := rs.NormalizeRawData(rawData)
					if err != nil {
						return nil, err
					}

					return yamlhelper.ToYamlBytes(data)
				},
			), -1)
			assert.NoError(t, err)

			actual := spec.Spec.GenerateEntries(&Filter{
				match:  spec.MatchFilter,
				ignore: spec.IgnoreFilter,
			}, "", "")

			assert.EqualValues(t, exp, &actual)
		},
	)
}
