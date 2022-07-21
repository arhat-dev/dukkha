package matrix

import (
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
)

func TestSpec_GenerateFilter(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name string
		in   *Spec
	}{} {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}

func TestSpec_GenerateEntries(t *testing.T) {
	t.Parallel()

	
	tests := []struct {
		name string
		in   *Spec
		// specs are sorted by name, put them in order
		expected []Entry

		matchFilter  map[string]*Vector
		ignoreFilter [][2]string
	}{
		{
			name:     "nil",
			in:       nil,
			expected: nil,
		},
		{
			name: "basic",
			in: &Spec{
				Data: map[string]*Vector{
					"foo":    NewVector("a", "b"),
					"kernel": NewVector("linux", "darwin"),
					"arch":   NewVector("amd64", "arm64"),
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
			in: &Spec{
				Data: map[string]*Vector{
					"foo":    NewVector("a", "b"),
					"kernel": NewVector("linux", "darwin"),
					"arch":   NewVector("amd64", "arm64"),
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
			in: &Spec{
				Data: map[string]*Vector{
					"foo":    NewVector("a", "b"),
					"kernel": NewVector("linux", "darwin"),
					"arch":   NewVector("amd64", "arm64"),
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
			in: &Spec{
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
				Data: map[string]*Vector{
					"kernel": NewVector("linux"),
					"arch":   NewVector("amd64"),
				},
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
			in: &Spec{
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
				Data: map[string]*Vector{
					"kernel": NewVector("linux"),
					"arch":   NewVector("amd64", "arm64"),
				},
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
			in: &Spec{
				Exclude: []*specItem{
					{
						Data: map[string]*Vector{
							"kernel": NewVector("linux"),
							"arch":   NewVector("arm64", "amd64"),
						},
					},
				},
				Data: map[string]*Vector{
					"kernel": NewVector("linux"),
					"arch":   NewVector("amd64", "arm64"),
				},
			},
			expected: nil,
		},
		{
			name: "exclude-all-single-match",
			in: &Spec{
				Exclude: []*specItem{
					{
						Data: map[string]*Vector{
							"kernel": NewVector("linux"),
						},
					},
				},
				Data: map[string]*Vector{
					"kernel": NewVector("linux"),
					"arch":   NewVector("amd64", "arm64"),
				},
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(
				t,
				test.expected,
				test.in.GenerateEntries(Filter{
					match:  test.matchFilter,
					ignore: test.ignoreFilter,
				}),
			)
		})
	}
}

func TestSpec_GenerateEntries_Fixture(t *testing.T) {
	t.Parallel()

	type testInputSpec struct {
		rs.BaseField

		MatchFilter  map[string]*Vector `yaml:"match_filter"`
		IgnoreFilter [][2]string        `yaml:"ignore_fitler"`
		Spec         Spec               `yaml:"spec"`
	}

	testhelper.TestFixtures(t, "./fixtures/gen-entries",
		func() any { return rs.InitAny(&testInputSpec{}, nil).(*testInputSpec) },
		func() any {
			var data []Entry
			return &data
		},
		func(t *testing.T, in, exp any) {
			spec := in.(*testInputSpec)
			err := spec.ResolveFields(rs.RenderingHandleFunc(
				func(renderer string, rawData any) (result []byte, err error) {
					data, err := rs.NormalizeRawData(rawData)
					if err != nil {
						return nil, err
					}

					return yamlhelper.ToYamlBytes(data)
				},
			), -1)
			assert.NoError(t, err)

			actual := spec.Spec.GenerateEntries(Filter{
				match:  spec.MatchFilter,
				ignore: spec.IgnoreFilter,
			})

			assert.EqualValues(t, exp, &actual)
		},
	)
}
