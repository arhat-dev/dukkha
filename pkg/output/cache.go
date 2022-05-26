package output

import (
	"fmt"
	"io"

	"github.com/muesli/termenv"
)

func WriteUsingExpiredCacheWarning(stderr io.Writer, key string) {
	_, _ = fmt.Fprintln(stderr,
		termenv.String(
			fmt.Sprintf("[WARNING] using expired local cache for %q\n", key),
		).Foreground(termenv.ANSIBrightYellow).String(),
	)
}
