package constant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTriple(t *testing.T) {
	for _, test := range []struct {
		arch   string
		kernel string
		libc   string

		get func(arch, kernel, libc string) (string, bool)

		expected string
	}{
		{"armv7", "", "", GetZigTripleName, "arm-linux-musleabihf"},
		{"arm64", "", "", GetZigTripleName, "aarch64-linux-musl"},
	} {
		t.Run("", func(t *testing.T) {
			actual, ok := test.get(test.arch, test.kernel, test.libc)
			assert.True(t, ok)
			assert.Equal(t, test.expected, actual)
		})
	}
}
