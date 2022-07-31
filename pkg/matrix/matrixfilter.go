package matrix

import (
	"strings"
)

// ParseMatrixFilter parses text form of Filter
// 	- a "key=value" is considered as a match rule
// 	- a "key!=value" is considered as an ignore rule
func ParseMatrixFilter(arr []string) (ret Filter) {
	for _, v := range arr {
		key, value, found := strings.Cut(v, "!=")
		if found {
			ret.AddIgnore(key, value)
			continue
		}

		key, value, found = strings.Cut(v, "=")
		if found {
			ret.AddMatch(key, value)
			continue
		}
	}

	return ret
}
