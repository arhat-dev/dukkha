package output

import (
	"os"

	"github.com/fatih/color"
)

func WriteUsingExpiredCacheWarning(key string) {
	_, _ = color.New(color.FgHiYellow).Fprintf(os.Stderr,
		"[WARNING] using expired local cache for %q\n", key,
	)
}
