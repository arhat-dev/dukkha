package output

import (
	"fmt"
	"os"

	"github.com/muesli/termenv"
)

func WriteUsingExpiredCacheWarning(key string) {
	_, _ = fmt.Fprintln(os.Stderr,
		termenv.String(
			fmt.Sprintf("[WARNING] using expired local cache for %q\n", key),
		).Foreground(termenv.ANSIBrightYellow).String(),
	)
}
