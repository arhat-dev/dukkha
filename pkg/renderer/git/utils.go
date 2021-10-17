package git

import (
	"fmt"
	"io"
	"strconv"
)

func formatPktLines(data []string) string {
	var ret string
	for _, v := range data {
		ret += formatPktLine(v)
	}
	return ret
}

func formatPktLine(s string) string {
	var ret string
	size := uint64(len(s))
	for size != 0 {
		n := size
		if n > 65516 {
			n = 65516
		}

		ret += formatPktSize(n+4) + s[:n]
		size -= n
	}

	return ret
}

func readPktSize(r io.Reader, buf []byte) (int, pktStatus, error) {
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, pktInvalid, err
	}

	size := parsePktSize(buf)
	switch size {
	case 0:
		// 0000, flush pkt
		return 0, pktFlush, nil
	case 1:
		// 0001, delim-pkt
		// TODO: implement
		return 0, pktDelim, nil
	case 2:
		// 0002, response-end-pkt
		// TODO: implement
		return 0, pktResponseEnd, nil
	default:
		return size - 4, pktNromal, nil
	}
}

func formatPktSize(n uint64) string {
	sizeStr := strconv.FormatUint(n, 16)
	switch {
	case n < 0x0010:
		return "000" + sizeStr
	case n < 0x0100:
		return "00" + sizeStr
	case n < 0x1000:
		return "0" + sizeStr
	default:
		return sizeStr
	}
}

func parsePktSize(sizeBuf []byte) int {
	_ = sizeBuf[3]

	return parseHexDigit(sizeBuf[3]) |
		parseHexDigit(sizeBuf[2])<<4 |
		parseHexDigit(sizeBuf[1])<<8 |
		parseHexDigit(sizeBuf[0])<<12
}

func parseHexDigit(x byte) int {
	switch {
	case x >= 'A' && x <= 'F':
		return int(x - 'A' + 10)
	case x >= 'a' && x <= 'f':
		return int(x - 'a' + 10)
	case x >= '0' && x <= '9':
		return int(x - '0')
	default:
		panic(fmt.Errorf("invalid hex digit %q", x))
	}
}
