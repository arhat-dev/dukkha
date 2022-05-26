package templateutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashNS(t *testing.T) {
	const (
		TestData     = "foo bar"
		TestHMAC_KEY = "testkey"
	)

	var ns hashNS
	for _, test := range []struct {
		name string
		fn   func(...any) (string, error)
	}{
		{"alder32", ns.ADLER32},
		{"crc32", ns.CRC32},
		{"crc64", ns.CRC64},
		{"md4", ns.MD4},
		{"md5", ns.MD5},
		{"ripemd160", ns.RIPEMD160},
		{"sha1", ns.SHA1},
		{"sha224", ns.SHA224},
		{"sha256", ns.SHA256},
		{"sha384", ns.SHA384},
		{"sha512", ns.SHA512},
		{"sha512-224", ns.SHA512_224},
		{"sha512-256", ns.SHA512_256},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, cs := range []struct {
				name string
				args []any
			}{
				{"hex(default)", nil},
				{"hex(explicit)", []any{"--hex"}},
				{"base64", []any{"--base64"}},
				{"base32", []any{"--base32"}},
				{"raw", []any{"--raw"}},
			} {
				t.Run(cs.name, func(t *testing.T) {
					// string input
					args := append(append([]any{}, cs.args...), TestData)
					retDT, err := test.fn(args...)
					assert.NoError(t, err)
					assert.NotEmpty(t, retDT)
					t.Log(retDT)

					// reader input
					args = append(append([]any{}, cs.args...), strings.NewReader(TestData))
					retRD, err := test.fn(args...)
					assert.NoError(t, err)
					assert.Equal(t, retDT, retRD)

					// hmac
					t.Run("hmac", func(t *testing.T) {
						args := append(append([]any{"--hmac", TestHMAC_KEY}, cs.args...), TestData)
						retHMAC, err := test.fn(args...)
						assert.NoError(t, err)
						assert.NotEmpty(t, retHMAC)
						t.Log(retHMAC)

						assert.NotEqual(t, retDT, retHMAC)
					})
				})
			}
		})
	}
}
