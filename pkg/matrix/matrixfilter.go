package matrix

import (
	"strings"
)

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
