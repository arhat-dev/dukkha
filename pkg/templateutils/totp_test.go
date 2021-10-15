package templateutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTOTPCode(t *testing.T) {
	tests := []struct {
		name string

		token  string
		time   time.Time
		length int

		expected string
	}{
		{
			name:     "Default Length",
			token:    "JBSWY3DPEHPK3PXP",
			time:     time.Unix(1634290334, 0),
			length:   0,
			expected: "569116",
		},
		{
			name:     "Set Length 6",
			token:    "JBSWY3DPEHPK3PXP",
			time:     time.Unix(1634290334, 0),
			length:   6,
			expected: "569116",
		},
		{
			name:     "Set Length 8",
			token:    "JBSWY3DPEHPK3PXP",
			time:     time.Unix(1634290428, 0),
			length:   8,
			expected: "16382746",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, err := GenerateTOTPCode(test.token, test.time, test.length)
			assert.NoError(t, err)
			assert.EqualValues(t, test.expected, code)
		})
	}
}
