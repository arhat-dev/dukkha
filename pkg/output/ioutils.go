package output

import "io"

type prefixWriter struct {
	prefix []byte
	w      io.Writer
}

func (p *prefixWriter) Write(data []byte) (n int, err error) {
	_, err = p.w.Write(p.prefix)
	if err != nil {
		return 0, err
	}

	return p.w.Write(data)
}

func PrefixWriter(prefix string, w io.Writer) io.Writer {
	return &prefixWriter{
		prefix: []byte(prefix),
		w:      w,
	}
}
