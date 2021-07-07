package matrix

import (
	"sort"
	"strings"
)

type Spec map[string]string

func (m Spec) String() string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, k+": "+m[k])
	}

	return strings.Join(pairs, ", ")
}

// BriefString return all values concatenated with slash
func (m Spec) BriefString() string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, m[k])
	}

	return strings.Join(parts, "/")
}

func (m Spec) Match(a map[string]string) bool {
	if len(a) == 0 {
		return len(m) == 0
	}

	for k, v := range a {
		if m[k] != v {
			return false
		}
	}

	return true
}

func (m Spec) Equals(a map[string]string) bool {
	if a == nil {
		return m == nil
	}

	if len(a) != len(m) {
		return false
	}

	for k, v := range a {
		mv, ok := m[k]
		if !ok || mv != v {
			return false
		}
	}

	return true
}
