package sliceutils

import "sort"

func NewStrings(base []string, other ...string) []string {
	return append(append([]string(nil), base...), other...)
}

func FormatStringMap(m map[string]string, kvSep string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	ret := make([]string, 0, len(keys))
	for _, k := range keys {
		ret = append(ret, k+kvSep+m[k])
	}

	return ret
}
