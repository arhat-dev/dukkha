package utils

import (
	"strings"

	"arhat.dev/dukkha/pkg/matrix"
)

func ParseMatrixFilter(arr []string) *matrix.Filter {
	ret := matrix.NewFilter(make(map[string][]string))

	for _, v := range arr {
		if idx := strings.Index(v, "!="); idx > 0 {
			ret.AddIgnore(v[:idx], v[idx+2:])
			continue
		}

		if idx := strings.IndexByte(v, '='); idx > 0 {
			ret.AddMatch(v[:idx], v[idx+1:])
			continue
		}
	}

	return ret
}
