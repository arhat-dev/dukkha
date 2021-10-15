// Copied from https://github.com/yitsushi/totp-cli/blob/main/internal/security/otp.go
// with modification to codeLength

/*
MIT License

Copyright (c)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package templateutils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

// args: token, optional<code length>, optional<time value>
func totpTemplateFunc(tokenVal interface{}, args ...interface{}) (string, error) {
	var (
		token  string
		t      time.Time
		length int
	)

	switch tt := tokenVal.(type) {
	case string:
		token = tt
	case []byte:
		token = string(tt)
	default:
		return "", fmt.Errorf("invalid token value type: %q", tt)
	}

	if len(args) != 0 {
		switch at := args[0].(type) {
		case string, []byte:
		case int, int8, int16, int32, int64:
			length = int(reflect.ValueOf(at).Int())
		case uint, uint8, uint16, uint32, uint64, uintptr:
			length = int(reflect.ValueOf(at).Uint())
		case float32, float64:
			length = int(reflect.ValueOf(at).Float())
		default:
			return "", fmt.Errorf(
				"invalid code length arg value type %T",
				at,
			)
		}
	}

	if len(args) > 1 {
		switch at := args[1].(type) {
		case string:
			var err error
			t, err = time.Parse(time.RFC3339, at)
			if err != nil {
				return "", fmt.Errorf(
					"invalid totp time arg value %q for format %q",
					at, time.RFC3339,
				)
			}
		case []byte:
			var err error
			t, err = time.Parse(time.RFC3339, string(at))
			if err != nil {
				return "", fmt.Errorf(
					"invalid totp time arg value %q for format %q",
					string(at), time.RFC3339,
				)
			}
		case time.Time:
			t = at
		case *time.Time:
			if at == nil {
				t = time.Now()
			}
		default:
			return "", fmt.Errorf("invalid time arg value type: %T", at)
		}
	} else {
		t = time.Now()
	}

	return GenerateTOTPCode(token, t, length)
}

func GenerateTOTPCode(token string, t time.Time, codeLength int) (string, error) {
	const (
		mask1              = 0xf
		mask2              = 0x7f
		mask3              = 0xff
		timeSplitInSeconds = 30
		shift24            = 24
		shift16            = 16
		shift8             = 8
		sumByteLength      = 8
		passwordHashLength = 32
	)

	if codeLength <= 0 {
		codeLength = 6
	}

	timer := uint64(math.Floor(float64(t.Unix()) / float64(timeSplitInSeconds)))
	// Remove spaces, some providers are giving us in a readable format
	// so they add spaces in there. If it's not removed while pasting in,
	// remove it now.
	token = strings.ReplaceAll(token, " ", "")

	// It should be uppercase always
	token = strings.ToUpper(token)

	secretBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("invalid totp token: %w", err)
	}

	buf := make([]byte, sumByteLength)
	mac := hmac.New(sha1.New, secretBytes)

	binary.BigEndian.PutUint64(buf, timer)
	_, _ = mac.Write(buf)
	sum := mac.Sum(nil)

	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & mask1
	value := int64(((int(sum[offset]) & mask2) << shift24) |
		((int(sum[offset+1] & mask3)) << shift16) |
		((int(sum[offset+2] & mask3)) << shift8) |
		(int(sum[offset+3]) & mask3))

	modulo := int32(value % int64(math.Pow10(codeLength)))

	format := fmt.Sprintf("%%0%dd", codeLength)

	return fmt.Sprintf(format, modulo), nil
}
