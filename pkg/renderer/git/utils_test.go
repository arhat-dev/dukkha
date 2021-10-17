package git

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPktSize(t *testing.T) {
	for i := 1; i < math.MaxUint16; i += 1 {
		prefix := ""
		switch {
		case i < 0x0010:
			prefix = "000"
		case i < 0x0100:
			prefix = "00"
		case i < 0x1000:
			prefix = "0"
		}

		assert.Equal(t, prefix+strconv.FormatInt(int64(i), 16), formatPktSize(uint64(i)))
	}
}

func TestParsePktSize(t *testing.T) {
	for i := 1; i < math.MaxUint16; i += 1 {
		src := formatPktSize(uint64(i))
		assert.Equal(t, i, parsePktSize([]byte(strings.ToLower(src))))
		assert.Equal(t, i, parsePktSize([]byte(strings.ToUpper(src))))
	}
}
