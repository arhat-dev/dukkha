package utils

import (
	"io"

	"github.com/muesli/termenv"
)

type prefixWriter struct {
	writePrefix func() error
	writeOutput func(p []byte) (int, error)

	_w io.Writer
}

func (p *prefixWriter) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return p._w.Write(data)
	}

	//
	// 	lines := bytes.SplitAfter(data, []byte{'\n'})
	//
	// 	err = p.writePrefix()
	// 	if err != nil {
	// 		return 0, err
	// 	}
	//
	// 	var lineN int
	// 	for _, line := range lines {
	// 		if len(line) == 0 {
	// 			continue
	// 		}
	//
	// 		if lineN != 0 {
	// 			err = p.writePrefix()
	// 			if err != nil {
	// 				return
	// 			}
	// 		}
	//
	// 		lineN, err = p.writeOutput(line)
	// 		n += lineN
	// 		if err != nil {
	// 			return
	// 		}
	// 	}

	n, err = p.writeOutput(data)
	if n > len(data) {
		n = len(data)
	}

	return n, err
}

func TermWriter(
	prefix string,
	useColor bool,
	prefixColor, outputColor termenv.Color,
	w io.Writer,
) io.Writer {
	prefixBytes := []byte(prefix)
	writePrefix := func() error {
		_, err := w.Write(prefixBytes)
		return err
	}

	writeOutput := func(p []byte) (int, error) {
		return w.Write(p)
	}

	if useColor {
		if prefixColor != nil {
			style := termenv.Style{}.Foreground(prefixColor)
			writePrefix = func() error {
				_, err := w.Write([]byte(style.Styled(prefix)))
				return err
			}
		}

		if outputColor != nil {
			style := termenv.Style{}.Foreground(outputColor)
			writeOutput = func(p []byte) (int, error) {
				return w.Write([]byte(style.Styled(string(p))))
			}
		}
	}

	return &prefixWriter{
		writePrefix: writePrefix,
		writeOutput: writeOutput,

		_w: w,
	}
}
