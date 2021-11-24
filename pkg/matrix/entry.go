package matrix

import (
	"strings"

	"arhat.dev/dukkha/pkg/sliceutils"
)

type Entry map[string]string

func (m Entry) String() string {
	return strings.Join(sliceutils.FormatStringMap(m, ": ", false), ", ")
}

// BriefString return all values concatenated with slash
func (m Entry) BriefString() string {
	return strings.Join(sliceutils.FormatStringMap(m, "", true), "/")
}

func (m Entry) Match(a map[string]string) bool {
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

func (m Entry) MatchKV(key, value string) bool {
	return m[key] == value
}

func (m Entry) Equals(a map[string]string) bool {
	if m == nil {
		return a == nil
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
