package renderer

import (
	"path/filepath"
)

func FormatCacheDir(dukkhaCacheDir, rendererName string) string {
	return filepath.Join(dukkhaCacheDir, "renderer-"+rendererName)
}
