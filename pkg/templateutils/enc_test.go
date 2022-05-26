package templateutils

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncNS_Hex(t *testing.T) {
	const (
		TestData = "foo bar"
	)

	expected := hex.EncodeToString([]byte(TestData))

	var ns encNS

	t.Run("Encode", func(t *testing.T) {
		t.Run("Data", func(t *testing.T) {
			ret, err := ns.Hex(TestData)
			assert.NoError(t, err)
			assert.Equal(t, expected, ret)
		})

		t.Run("Reader", func(t *testing.T) {
			ret, err := ns.Hex(strings.NewReader(TestData))
			assert.NoError(t, err)
			assert.Equal(t, expected, ret)
		})

		t.Run("Writer", func(t *testing.T) {
			buf := &bytes.Buffer{}
			ret, err := ns.Hex(buf, TestData)
			assert.NoError(t, err)
			assert.Equal(t, "", ret)
			assert.Equal(t, expected, buf.String())
		})

		t.Run("ReadWriter", func(t *testing.T) {
			buf := &bytes.Buffer{}
			ret, err := ns.Hex(buf, strings.NewReader(TestData))
			assert.NoError(t, err)
			assert.Equal(t, "", ret)
			assert.Equal(t, expected, buf.String())
		})
	})

	t.Run("Decode", func(t *testing.T) {
		for _, f := range []string{"-d", "--decode"} {
			t.Run("Data", func(t *testing.T) {
				ret, err := ns.Hex(f, expected)
				assert.NoError(t, err)
				assert.Equal(t, TestData, ret)
			})

			t.Run("Reader", func(t *testing.T) {
				ret, err := ns.Hex(f, strings.NewReader(expected))
				assert.NoError(t, err)
				assert.Equal(t, TestData, ret)
			})

			t.Run("Writer", func(t *testing.T) {
				buf := &bytes.Buffer{}
				ret, err := ns.Hex(f, buf, expected)
				assert.NoError(t, err)
				assert.Equal(t, "", ret)
				assert.Equal(t, TestData, buf.String())
			})

			t.Run("ReadWriter", func(t *testing.T) {
				buf := &bytes.Buffer{}
				ret, err := ns.Hex(f, buf, strings.NewReader(expected))
				assert.NoError(t, err)
				assert.Equal(t, "", ret)
				assert.Equal(t, TestData, buf.String())
			})

		}
	})
}

func TestEncNS_BaseX(t *testing.T) {
	const (
		TestData = "foo bar"
	)

	base64ExpectedStd := base64.StdEncoding.EncodeToString([]byte(TestData))
	base64ExpectedRawStd := base64.RawStdEncoding.EncodeToString([]byte(TestData))
	assert.NotEmpty(t, base64ExpectedStd)
	assert.NotEmpty(t, base64ExpectedRawStd)
	assert.NotEqual(t, base64ExpectedStd, base64ExpectedRawStd)

	base64ExpectedURL := base64.URLEncoding.EncodeToString([]byte(TestData))
	base64ExpectedRawURL := base64.RawURLEncoding.EncodeToString([]byte(TestData))
	assert.NotEmpty(t, base64ExpectedURL)
	assert.NotEmpty(t, base64ExpectedRawURL)
	assert.NotEqual(t, base64ExpectedURL, base64ExpectedRawURL)

	base32ExpectedStd := base32.StdEncoding.EncodeToString([]byte(TestData))
	base32ExpectedRawStd := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(TestData))
	assert.NotEmpty(t, base32ExpectedStd)
	assert.NotEmpty(t, base32ExpectedRawStd)
	assert.NotEqual(t, base32ExpectedStd, base32ExpectedRawStd)

	base32ExpectedHex := base32.HexEncoding.EncodeToString([]byte(TestData))
	base32ExpectedRawHex := base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(TestData))
	assert.NotEmpty(t, base32ExpectedHex)
	assert.NotEmpty(t, base32ExpectedRawHex)
	assert.NotEqual(t, base32ExpectedHex, base32ExpectedRawHex)

	for _, test := range []struct {
		name string
		fn   func(args ...any) (string, error)

		expectedStd    string
		expectedRawStd string
		expectedURL    string
		expectedRawURL string
		expectedHex    string
		expectedRawHex string
	}{
		{
			name: "Base64",
			fn:   encNS{}.Base64,

			expectedStd:    base64ExpectedStd,
			expectedRawStd: base64ExpectedRawStd,
			expectedURL:    base64ExpectedURL,
			expectedRawURL: base64ExpectedRawURL,
			expectedHex:    "",
			expectedRawHex: "",
		},
		{
			name: "Base32",
			fn:   encNS{}.Base32,

			expectedStd:    base32ExpectedStd,
			expectedRawStd: base32ExpectedRawStd,
			expectedURL:    "",
			expectedRawURL: "",
			expectedHex:    base32ExpectedHex,
			expectedRawHex: base32ExpectedRawHex,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("encode", func(t *testing.T) {
				t.Run("Std", func(t *testing.T) {
					ret, err := test.fn(TestData)
					assert.NoError(t, err)
					assert.Equal(t, test.expectedStd, ret)
				})

				t.Run("StdReader", func(t *testing.T) {
					ret, err := test.fn(strings.NewReader(TestData))
					assert.NoError(t, err)
					assert.Equal(t, test.expectedStd, ret)
				})

				t.Run("StdWriter", func(t *testing.T) {
					buf := &bytes.Buffer{}
					ret, err := test.fn(buf, TestData)
					assert.NoError(t, err)
					assert.Equal(t, "", ret)
					assert.Equal(t, test.expectedStd, buf.String())
				})

				t.Run("StdReadWriter", func(t *testing.T) {
					buf := &bytes.Buffer{}
					ret, err := test.fn(buf, strings.NewReader(TestData))
					assert.NoError(t, err)
					assert.Equal(t, "", ret)
					assert.Equal(t, test.expectedStd, buf.String())
				})

				t.Run("StdWrap", func(t *testing.T) {
					for _, f := range []string{"-w", "--wrap"} {
						ret, err := test.fn(f, 1, TestData)
						assert.NoError(t, err)
						assert.Equal(t, strings.Join(strings.Split(test.expectedStd, ""), "\n")+"\n", ret)
					}
				})

				t.Run("RawStd", func(t *testing.T) {
					for _, f := range []string{"-r", "--raw"} {
						ret, err := test.fn(f, TestData)
						if len(test.expectedRawStd) != 0 {
							assert.NoError(t, err)
							assert.Equal(t, test.expectedRawStd, ret)
						} else {
							assert.Error(t, err)
							assert.Equal(t, "", ret)
						}
					}
				})

				t.Run("URL", func(t *testing.T) {
					for _, f := range []string{"-u", "--url"} {
						ret, err := test.fn(f, TestData)
						if len(test.expectedURL) != 0 {
							assert.NoError(t, err)
							assert.Equal(t, test.expectedURL, ret)
						} else {
							assert.Error(t, err)
							assert.Equal(t, "", ret)
						}
					}
				})

				t.Run("RawURL", func(t *testing.T) {
					for _, f := range [][]string{
						{"-u", "-r"},
						{"--url", "--raw"},
						{"-u", "--raw"},
						{"--url", "-r"},
					} {
						args := make([]any, len(f)+1)
						for i, v := range f {
							args[i] = v
						}

						args[len(f)] = TestData

						ret, err := test.fn(args...)
						if len(test.expectedRawURL) != 0 {
							assert.NoError(t, err)
							assert.Equal(t, test.expectedRawURL, ret)
						} else {
							assert.Error(t, err)
							assert.Equal(t, "", ret)
						}
					}
				})

				t.Run("Hex", func(t *testing.T) {
					for _, f := range []string{"-h", "--hex"} {
						ret, err := test.fn(f, TestData)
						if len(test.expectedHex) != 0 {
							assert.NoError(t, err)
							assert.Equal(t, test.expectedHex, ret)
						} else {
							assert.Error(t, err)
							assert.Equal(t, "", ret)
						}
					}
				})

				t.Run("RawHex", func(t *testing.T) {
					for _, f := range [][]string{
						{"-h", "-r"},
						{"--hex", "--raw"},
						{"-h", "--raw"},
						{"--hex", "-r"},
					} {
						args := make([]any, len(f)+1)
						for i, v := range f {
							args[i] = v
						}

						args[len(f)] = TestData

						ret, err := test.fn(args...)
						if len(test.expectedRawHex) != 0 {
							assert.NoError(t, err)
							assert.Equal(t, test.expectedRawHex, ret)
						} else {
							assert.Error(t, err)
							assert.Equal(t, "", ret)
						}
					}
				})
			})

			t.Run("decode", func(t *testing.T) {
				for _, f := range []string{
					"-d", "--decode",
				} {
					ret, err := test.fn(f, test.expectedStd)
					assert.NoError(t, err)
					assert.Equal(t, TestData, ret)
				}

				// 				buf := &bytes.Buffer{}
				// 				ret, err = test.fn("--decode", "-u", buf, test.expectedURL)
				// 				assert.NoError(t, err)
				// 				assert.Equal(t, "", ret)
				// 				assert.Equal(t, TestData, buf.String())
			})
		})
	}
}
