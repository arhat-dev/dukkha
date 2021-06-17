package tools

import "strings"

func joinReplaceEmpty(sep string, defaults []string, s ...string) string {
	var v []string
	for i, str := range s {
		if len(str) == 0 {
			v = append(v, defaults[i])
		}

		v = append(v, str)
	}

	return strings.Join(v, sep)
}
