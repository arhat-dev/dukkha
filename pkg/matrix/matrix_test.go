package matrix

import (
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/pkg/yamlhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"
)

func TestSpec_GenerateEntries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec *Spec

		matchFilter  map[string]*Vector
		ignoreFilter [][2]string

		expected []Entry
	}{
		{
			name:     "nil",
			spec:     nil,
			expected: nil,
		},
		{
			name: "Basic",
			spec: &Spec{
				Values: map[string]*Vector{
					"foo":    NewVector("a", "b"),
					"kernel": NewVector("linux", "darwin"),
					"arch":   NewVector("amd64", "arm64"),
				},
			},
			expected: []Entry{
				// sort order: os=linux arch=amd64,arm64 foo=a,b
				{"kernel": "linux", "arch": "amd64", "foo": "a"},
				{"kernel": "linux", "arch": "amd64", "foo": "b"},
				{"kernel": "linux", "arch": "arm64", "foo": "a"},
				{"kernel": "linux", "arch": "arm64", "foo": "b"},

				// sort order: os=darwin arch=amd64,arm64 foo=a,b
				{"kernel": "darwin", "arch": "amd64", "foo": "a"},
				{"kernel": "darwin", "arch": "amd64", "foo": "b"},
				{"kernel": "darwin", "arch": "arm64", "foo": "a"},
				{"kernel": "darwin", "arch": "arm64", "foo": "b"},
			},
		},
		{
			name: "Basic and MatchFilter",
			spec: &Spec{
				Values: map[string]*Vector{
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
			name: "Basic and IgnoreFilter",
			spec: &Spec{
				Values: map[string]*Vector{
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
				// sort order: os=linux,darwin arch=amd64 foo=a,b
				{"kernel": "darwin", "arch": "amd64", "foo": "a"},
				{"kernel": "darwin", "arch": "amd64", "foo": "b"},
			},
		},
		{
			name: "Include",
			spec: &Spec{
				Include: []*SpecItem{
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
				Values: map[string]*Vector{
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
			name: "Include And MatchFilter",
			spec: &Spec{
				Include: []*SpecItem{
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
				Values: map[string]*Vector{
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
			name: "Exclude All By Full Match",
			spec: &Spec{
				Exclude: []*SpecItem{
					{
						Data: map[string]*Vector{
							"kernel": NewVector("linux"),
							"arch":   NewVector("arm64", "amd64"),
						},
					},
				},
				Values: map[string]*Vector{
					"kernel": NewVector("linux"),
					"arch":   NewVector("amd64", "arm64"),
				},
			},
			expected: nil,
		},
		{
			name: "Exclude All By Single Match",
			spec: &Spec{
				Exclude: []*SpecItem{
					{
						Data: map[string]*Vector{
							"kernel": NewVector("linux"),
						},
					},
				},
				Values: map[string]*Vector{
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
				test.spec.GenerateEntries(Filter{
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
		func() *testInputSpec { return rs.Init(&testInputSpec{}, nil).(*testInputSpec) },
		func() *[]Entry {
			var data []Entry
			return &data
		},
		func(t *testing.T, spec *testInputSpec, exp *[]Entry) {
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

func TestSpec_AsFilter(t *testing.T) {
	t.Parallel()

	spec := Spec{
		Exclude: []*SpecItem{},
		Include: []*SpecItem{},
		Values: map[string]*Vector{
			"kernel": {
				Vec: []string{"k1", "k2", "k3"},
			},
			"arch": {
				Vec: []string{"a1", "a2", "a3"},
			},
		},
	}

	all := []Entry{
		{"kernel": "k1", "arch": "a1"},
		{"kernel": "k1", "arch": "a2"},
		{"kernel": "k1", "arch": "a3"},
		{"kernel": "k2", "arch": "a1"},
		{"kernel": "k2", "arch": "a2"},
		{"kernel": "k2", "arch": "a3"},
		{"kernel": "k3", "arch": "a1"},
		{"kernel": "k3", "arch": "a2"},
		{"kernel": "k3", "arch": "a3"},
	}

	assert.Equal(t, all, spec.GenerateEntries(Filter{}))

	for _, test := range []struct {
		name     string
		filter   *Spec
		expected []Entry
	}{
		// Type 1: empty filter
		{
			name:     "Nil Filter Spec Match All",
			filter:   nil,
			expected: all,
		},
		{
			name:     "Empty Filter Spec Match All",
			filter:   &Spec{},
			expected: all,
		},

		// Type 2: None Existing Key/Value
		{
			name: "Exclude None Existing Value Match All",
			filter: &Spec{
				Exclude: []*SpecItem{
					{
						Data: map[string]*Vector{
							"kernel": {Vec: []string{"non-exist"}},
						},
					},
				},
			},
			expected: all,
		},
		{
			name: "Exclude None Existing Key Match All",
			filter: &Spec{
				Exclude: []*SpecItem{
					{
						Data: map[string]*Vector{
							"non-exist": {Vec: []string{"k1"}},
						},
					},
				},
			},
			expected: all,
		},
		{
			name: "Include None Existing Value Match None",
			filter: &Spec{
				Include: []*SpecItem{
					{
						Data: map[string]*Vector{
							"kernel": {Vec: []string{"non-exist"}},
						},
					},
				},
			},
			expected: nil,
		},
		{
			name: "Include None Existing Key Match None",
			filter: &Spec{
				Include: []*SpecItem{
					{
						Data: map[string]*Vector{
							"non-exist": {Vec: []string{"k1"}},
						},
					},
				},
			},
			expected: nil,
		},
		{
			name: "Non Existing Value Match None",
			filter: &Spec{
				Values: map[string]*Vector{
					"kernel": {Vec: []string{"non-exist"}},
				},
			},
			expected: nil,
		},
		{
			name: "Non Existing Key Match None",
			filter: &Spec{
				Values: map[string]*Vector{
					"non-exist": {Vec: []string{"k1"}},
				},
			},
			expected: nil,
		},

		// Type 3: TBD
		{
			name: "Exclude All Value Of Single Key Match None",
			filter: &Spec{
				Exclude: []*SpecItem{
					{
						Data: map[string]*Vector{
							"kernel": {Vec: []string{"k1", "k2", "k3"}},
						},
					},
				},
			},
			expected: nil,
		},
		{
			name: "Exclude Some Value Of Single Key Match Some",
			filter: &Spec{
				Exclude: []*SpecItem{
					{
						Data: map[string]*Vector{
							"kernel": {Vec: []string{"k1", "k2"}},
						},
					},
				},
			},
			expected: []Entry{
				{"kernel": "k3", "arch": "a1"},
				{"kernel": "k3", "arch": "a2"},
				{"kernel": "k3", "arch": "a3"},
			},
		},
		{
			name: "Values All Key Single Value Match Single",
			filter: &Spec{
				Values: map[string]*Vector{
					"kernel": {Vec: []string{"k1"}},
					"arch":   {Vec: []string{"a1"}},
				},
			},
			expected: []Entry{
				{"kernel": "k1", "arch": "a1"},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			filter := test.filter.AsFilter()
			if !assert.EqualValues(t, test.expected, spec.GenerateEntries(filter)) {
				t.Log(filter)
			}
		})
	}
}
