package output

import (
	"bytes"
	"io"
)

type prefixWriter struct {
	prefix []byte
	w      io.Writer
}

func (p *prefixWriter) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return p.w.Write(data)
	}

	lines := bytes.SplitAfter(data, []byte{'\n'})

	var lineN int
	_, err = p.w.Write(p.prefix)
	if err != nil {
		return 0, err
	}

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		if lineN != 0 {
			_, err = p.w.Write(p.prefix)
			if err != nil {
				return 0, err
			}
		}

		lineN, err = p.w.Write(line)
		n += lineN
		if err != nil {
			return
		}
	}

	return
}

func PrefixWriter(prefix string, w io.Writer) io.Writer {
	return &prefixWriter{
		prefix: []byte(prefix),
		w:      w,
	}
}
