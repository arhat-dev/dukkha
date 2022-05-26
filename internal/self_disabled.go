//go:build noself

package dukkha_internal

import (
	"io"

	"arhat.dev/dukkha/pkg/dukkha"
)

func RunSelf(
	ctx dukkha.Context, stdin io.Reader, stdout, stderr io.Writer, args ...string,
) error {
	return nil
}
