package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestANSIWriter(t *testing.T) {
	const (
		fakeProgress = "" +
			"#       10.2%\r" +
			"##      16.2%\r" +
			"###     22.2%\r" +
			"####    81.8%\r" +
			"#####   98.3%\r" +
			"###### 100.0%\r\n"
		fakeProgressCompacted = "###### 100.0%\n"
	)

	tests := []struct {
		name string

		retainANSI bool
		input      string
		expected   string
	}{
		{
			name: "Plain",

			retainANSI: false,
			input:      "foo",
			expected:   "foo\n",
		},
		{
			name: "Style Stripped",

			retainANSI: false,
			input:      "\x1b[30mfoo\x1b[0m",
			expected:   "foo\n",
		},
		{
			name: "Style Retained",

			retainANSI: true,
			input:      "\x1b[30mfoo\x1b[0m",
			expected:   "\x1b[30mfoo\x1b[0m\n",
		},
		{
			name: "Lines Compacted",

			retainANSI: true,
			input:      strings.Repeat(fakeProgress, 10),
			expected:   strings.Repeat(fakeProgressCompacted, 10),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := NewANSIWriter(buf, test.retainANSI)
			_, err := w.Write([]byte(test.input))
			assert.NoError(t, err)
			_, err = w.Flush()

			assert.NoError(t, err)
			assert.Equal(t, test.expected, buf.String())
		})
	}
}
