package templateutils

import (
	"strings"

	"arhat.dev/pkg/pathhelper"
	"mvdan.cc/sh/v3/expand"
)

func isSlash(c byte) bool { return pathhelper.IsWindowsSlash(c) }

func pathExts(env expand.Environ) []string {
	pathext := env.Get("PATHEXT").String()
	if pathext == "" {
		return []string{".com", ".exe", ".bat", ".cmd", ""}
	}

	var exts []string
	for _, e := range strings.Split(strings.ToLower(pathext), `;`) {
		if e == "" {
			continue
		}
		if e[0] != '.' {
			e = "." + e
		}
		exts = append(exts, e)
	}
	return exts
}
