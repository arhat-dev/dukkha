package templateutils

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveMatchedRunes(t *testing.T) {
	for _, test := range []struct {
		name     string
		data     string
		expected string
	}{
		{"empty", "", ""},
		{"all match", ",,,,,,,,,,", ""},
		{"nomatch", "12345", "12345"},
		{"first match only", ",12345", "12345"},
		{"last match only", "12345,", "12345"},
		{"inner match only", "123,45", "12345"},
		{"multi first match only", ",,,12345", "12345"},
		{"multi last match only", "12345,,,", "12345"},
		{"multi inner match only", "12,,,345", "12345"},
		{"multi matches", "12,,,345", "12345"},
		{"all match cond", ",1,,,23,4,,,,,5,,,", "12345"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ret := RemoveMatchedRunesCopy(test.data, func(r rune) bool { return r == ',' })
			assert.Equal(t, test.expected, ret)

			data := []byte(test.data)
			sz := RemoveMatchedRunesInPlace(data, func(r rune) bool { return r == ',' })
			assert.Equal(t, len(test.expected), sz)
			assert.Equal(t, test.expected, string(data[:sz]))
		})
	}
}

func TestChunkedWriter(t *testing.T) {
	// nolint:revive
	const (
		TOTAL_SIZE = 103
		CHUNK_SIZE = 4

		EACH_WRITE_SIZE = 10

		CALL_AFTER = TOTAL_SIZE / CHUNK_SIZE
		REMAINDER  = TOTAL_SIZE % CHUNK_SIZE
	)

	// nolint:revive
	CALL_PRE := CALL_AFTER
	if REMAINDER != 0 {
		CALL_PRE++
	}

	src := strings.Repeat("1", TOTAL_SIZE)

	buf := &bytes.Buffer{}
	rd := strings.NewReader(src)

	pre, after := 0, 0
	cw := NewChunkedWriter(CHUNK_SIZE, buf,
		func() error {
			pre++
			return nil
		},
		func() error {
			after++
			return nil
		},
	)

	n, err := io.CopyBuffer(&cw, rd, make([]byte, EACH_WRITE_SIZE))
	assert.NoError(t, err)
	assert.Equal(t, int64(TOTAL_SIZE), n)
	assert.Equal(t, src, buf.String())

	assert.Equal(t, CALL_PRE, pre)
	assert.Equal(t, CALL_AFTER, after)
	assert.Equal(t, REMAINDER, cw.Remainder())
}
