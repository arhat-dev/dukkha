package utils

import (
	"io"

	"github.com/fatih/color"
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

	return p.writeOutput(data)
}

func PrefixWriter(
	prefix string,
	useColor bool,
	prefixColor, outputColor *color.Color,
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
			writePrefix = func() error {
				_, err := prefixColor.Fprint(w, prefix)
				return err
			}
		}

		if outputColor != nil {
			writeOutput = func(p []byte) (int, error) {
				return outputColor.Fprint(w, string(p))
			}
		}
	}

	return &prefixWriter{
		writePrefix: writePrefix,
		writeOutput: writeOutput,

		_w: w,
	}
}
