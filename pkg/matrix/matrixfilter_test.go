package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMatrixFilter(t *testing.T) {
	for _, test := range []struct {
		name    string
		filters []string

		match  map[string]*Vector
		ignore [][2]string
	}{
		{
			name: "Match",
			filters: []string{
				"a=b",
			},

			match: map[string]*Vector{
				"a": NewVector("b"),
			},
		},
		{
			name: "Match Multiple",
			filters: []string{
				"a=b", "a=c",
			},

			match: map[string]*Vector{
				"a": NewVector("b", "c"),
			},
		},
		{
			name: "Ignore",
			filters: []string{
				"a!=b",
			},

			match:  map[string]*Vector{},
			ignore: [][2]string{{"a", "b"}},
		},
		{
			name: "Ignore Multiple",
			filters: []string{
				"a!=b", "a!=c",
			},

			match:  map[string]*Vector{},
			ignore: [][2]string{{"a", "b"}, {"a", "c"}},
		},
		{
			name: "Match And Ignore",
			filters: []string{
				"a!=b", "b=c",
			},

			match: map[string]*Vector{
				"b": NewVector("c"),
			},
			ignore: [][2]string{{"a", "b"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mf := ParseMatrixFilter(test.filters)

			assert.Equal(t, len(test.match), len(mf.match))
			for k, v := range test.match {
				assert.EqualValues(t, v.Vector, mf.match[k].Vector)
			}

			assert.EqualValues(t, test.ignore, mf.ignore)
		})
	}
}
