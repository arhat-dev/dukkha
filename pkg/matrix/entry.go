package matrix

import (
	"strings"

	"arhat.dev/dukkha/pkg/sliceutils"
)

// An Entry represents a set of all key value pairs of a matrix operation
type Entry map[string]string

// String formats the Entry as
//
// 	key1: value1, key2: value2, ...
func (m Entry) String() string {
	return strings.Join(sliceutils.FormatStringMap(m, ": ", false), ", ")
}

// BriefString return all values concatenated with slash
func (m Entry) BriefString() string {
	return strings.Join(sliceutils.FormatStringMap(m, "", true), "/")
}

// Contains returns true when x is a subset of the Entry
func (m Entry) Contains(x map[string]string) bool {
	if len(x) == 0 {
		return len(m) == 0
	}

	for k, v := range x {
		if m[k] != v {
			return false
		}
	}

	return true
}

// MatchKV returns true when it possesses the same key value pair
func (m Entry) MatchKV(key, value string) bool {
	return m[key] == value
}

// Equals returns true all entries in x are the same with all entries in m
func (m Entry) Equals(x map[string]string) bool {
	if m == nil {
		return x == nil
	}

	if len(x) != len(m) {
		return false
	}

	for k, v := range x {
		mv, ok := m[k]
		if !ok || mv != v {
			return false
		}
	}

	return true
}
