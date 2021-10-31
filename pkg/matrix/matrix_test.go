package matrix

import (
	"testing"

	"arhat.dev/pkg/testhelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	_ "embed"
)

func TestMatrixConfig_GenerateEntries(t *testing.T) {
	tests := []struct {
		name string
		in   Spec
		// specs are sorted by name, put them in order
		expected []Entry

		matchFilter  map[string][]string
		ignoreFilter [][2]string
	}{
		{
			name: "basic",
			in: Spec{
				Kernel: []string{"linux", "darwin"},
				Arch:   []string{"amd64", "arm64"},
				Custom: map[string][]string{
					"foo": {"a", "b"},
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
				Kernel: []string{"linux", "darwin"},
				Arch:   []string{"amd64", "arm64"},
				Custom: map[string][]string{
					"foo": {"a", "b"},
				},
			},
			matchFilter: map[string][]string{
				"kernel": {"linux"},
				"arch":   {"arm64"},
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
				Kernel: []string{"linux", "darwin"},
				Arch:   []string{"amd64", "arm64"},
				Custom: map[string][]string{
					"foo": {"a", "b"},
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
			matchFilter: map[string][]string{
				"arch": {"amd64"},
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

		MatchFilter  map[string][]string `yaml:"match_filter"`
		IgnoreFilter [][2]string         `yaml:"ignore_fitler"`
		Spec         Spec                `yaml:"spec"`
	}

	testhelper.TestFixtures(t, "./_fixtures/gen-entries",
		func() interface{} { return rs.Init(&testInputSpec{}, nil).(*testInputSpec) },
		func() interface{} {
			var data []Entry
			return &data
		},
		func(t *testing.T, in, exp interface{}) {
			spec := in.(*testInputSpec)
			actual := spec.Spec.GenerateEntries(&Filter{
				match:  spec.MatchFilter,
				ignore: spec.IgnoreFilter,
			}, "", "")

			assert.EqualValues(t, exp, &actual)
		},
	)
}
