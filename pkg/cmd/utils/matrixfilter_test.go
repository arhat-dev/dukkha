package utils

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"arhat.dev/dukkha/pkg/matrix"
)

func TestParseMatrixFilter(t *testing.T) {
	for _, test := range []struct {
		name    string
		filters []string

		match  map[string][]string
		ignore [][2]string
	}{
		{
			name: "Match",
			filters: []string{
				"a=b",
			},

			match: map[string][]string{"a": {"b"}},
		},
		{
			name: "Match Multiple",
			filters: []string{
				"a=b", "a=c",
			},

			match: map[string][]string{"a": {"b", "c"}},
		},
		{
			name: "Ignore",
			filters: []string{
				"a!=b",
			},

			match:  map[string][]string{},
			ignore: [][2]string{{"a", "b"}},
		},
		{
			name: "Ignore Multiple",
			filters: []string{
				"a!=b", "a!=c",
			},

			match:  map[string][]string{},
			ignore: [][2]string{{"a", "b"}, {"a", "c"}},
		},
		{
			name: "Match And Ignore",
			filters: []string{
				"a!=b", "b=c",
			},

			match:  map[string][]string{"b": {"c"}},
			ignore: [][2]string{{"a", "b"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mf := ParseMatrixFilter(test.filters)

			expected := struct {
				match  map[string][]string
				ignore [][2]string
			}{
				match:  test.match,
				ignore: test.ignore,
			}
			assert.EqualValues(t, (*matrix.Filter)(unsafe.Pointer(&expected)), mf)
		})
	}
}
