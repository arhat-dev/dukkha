package sliceutils

import "sort"

func NewStrings(base []string, other ...string) []string {
	return append(append([]string(nil), base...), other...)
}

func FormatStringMap(m map[string]string, kvSep string, omitKey bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	SortByKernelCmdArchLibcOther(keys, nil)

	ret := make([]string, 0, len(keys))
	for _, k := range keys {
		if omitKey {
			ret = append(ret, m[k])
			continue
		}

		ret = append(ret, k+kvSep+m[k])
	}

	return ret
}

var _ sort.Interface = (*biasedSorter[string, struct{}])(nil)

// biasedSorter is a sort.Interface implementation
//
// priority: high to low
// 	- kernel
//  - cmd
//  - arch
//  - libc
//  ...other
type biasedSorter[K ~string, V any] struct {
	keys   []K
	values [][]V
}

// Len implements sort.Interface
func (s *biasedSorter[K, V]) Len() int {
	return len(s.keys)
}

// Less implements sort.Interface
func (s *biasedSorter[K, V]) Less(i int, j int) bool {
	switch {
	case s.keys[i] == "kernel":
		return s.keys[j] != "kernel"
	case s.keys[j] == "kernel":
		return false
	case s.keys[i] == "cmd":
		return s.keys[j] != "cmd"
	case s.keys[j] == "cmd":
		return false
	case s.keys[i] == "arch":
		return s.keys[j] != "arch"
	case s.keys[j] == "arch":
		return false
	case s.keys[i] == "libc":
		return s.keys[j] != "libc"
	case s.keys[j] == "libc":
		return false
	default:
		return s.keys[i] < s.keys[j]
	}
}

// Swap implements sort.Interface
func (s *biasedSorter[K, V]) Swap(i int, j int) {
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
	if s.values != nil {
		s.values[i], s.values[j] = s.values[j], s.values[i]
	}
}

// SortByKernelCmdArchLibcOther
//
// mat can be nil, but when mat is not nil, its length should be the same of names's
func SortByKernelCmdArchLibcOther(names []string, mat [][]string) {
	s := biasedSorter[string, string]{
		keys:   names,
		values: mat,
	}

	sort.Stable(&s)
}
