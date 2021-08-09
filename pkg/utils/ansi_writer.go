package utils

import (
	"io"
	"sync"

	"github.com/aoldershaw/ansi"
	"github.com/muesli/termenv"
)

var _ io.Writer = (*ANSIWriter)(nil)

func NewANSIWriter(w io.Writer, retainStyle bool) *ANSIWriter {
	lines := ansi.Lines{}
	p := &ANSIWriter{
		lines:       &lines,
		ansiWriter:  ansi.NewWriter(&lines),
		retainStyle: retainStyle,

		w: w,
	}

	return p
}

// ANSIWriter is a ansi escape sequence aware writer
type ANSIWriter struct {
	// ansi handling
	lines       *ansi.Lines
	ansiWriter  *ansi.Writer
	retainStyle bool

	w io.Writer

	mu sync.Mutex
}

func (p *ANSIWriter) Write(data []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	n, err := p.ansiWriter.Write(data)
	if err != nil {
		return int(n), err
	}

	// write lines if reached threshold
	if len(*p.lines) >= 16 {
		_, err = p.flushBufferredLines()
	}

	return int(n), err
}

func (p *ANSIWriter) Flush() (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.flushBufferredLines()
}

// lock the plainTextWriter before calling this method
func (p *ANSIWriter) flushBufferredLines() (int, error) {
	var (
		sum     int
		lastIdx = -1
	)

	for i, line := range *p.lines {
		if len(line) == 0 {
			continue
		}

		var lineBytes []byte
		for _, chk := range line {
			data := string(chk.Data)
			if p.retainStyle {
				data = restoreStyle(data, chk.Style)
			}

			lineBytes = append(lineBytes, data...)
		}

		if len(lineBytes) == 0 {
			continue
		}

		n, err := p.w.Write(append(lineBytes, '\n'))
		sum += n
		if err != nil {
			lines := (*p.lines)[lastIdx+1:]
			p.lines = &lines
			p.ansiWriter.Output = p.lines
			return sum, err
		}

		lastIdx = i
	}

	if lastIdx != -1 {
		// wrote all buffered lines
		lines := (*p.lines)[0:0]
		p.lines = &lines
		p.ansiWriter.Output = p.lines
	}

	return sum, nil
}

func restoreStyle(data string, s ansi.Style) string {
	fg := s.Foreground
	bg := s.Background
	Bold := s.Modifier&ansi.Bold != 0
	Faint := s.Modifier&ansi.Faint != 0
	Italic := s.Modifier&ansi.Italic != 0
	Underline := s.Modifier&ansi.Underline != 0
	Blink := s.Modifier&ansi.Blink != 0
	Inverted := s.Modifier&ansi.Inverted != 0
	Fraktur := s.Modifier&ansi.Fraktur != 0
	Framed := s.Modifier&ansi.Framed != 0

	style := termenv.Style{}
styleLoop:
	for i := 0; i < 100; i++ {
		switch {
		case fg != ansi.DefaultColor:
			// termenv.ANSIColor = ansi.Color - 1
			style = style.Foreground(termenv.ANSIColor(fg - 1))
			fg = ansi.DefaultColor
		case bg != ansi.DefaultColor:
			// termenv.ANSIColor = ansi.Color - 1
			style = style.Background(termenv.ANSIColor(bg - 1))
			bg = ansi.DefaultColor
		case Bold:
			style = style.Bold()
			Bold = false
		case Faint:
			style = style.Faint()
			Faint = false
		case Italic:
			style = style.Italic()
			Italic = false
		case Underline:
			style = style.Underline()
			Underline = false
		case Blink:
			style = style.Blink()
			Blink = false
		case Inverted:
			// invert fg/bg
			// https://terminalguide.namepad.de/attr/7/
			style = style.Reverse()
			Inverted = false
		case Fraktur:
			// set: 20, reset: 23
			// https://espterm.github.io/help.html
			// TODO: not implemented
			Fraktur = false
		case Framed:
			// unknown source
			// TODO: not implemented
			Framed = false
		default:
			break styleLoop
		}
	}

	return style.Styled(data)
}
