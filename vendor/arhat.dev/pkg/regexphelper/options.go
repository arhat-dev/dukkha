package regexphelper

import "strings"

// regexp options per regexp/syntax/doc.go
//
// i	case-insensitive (default false)
// m	multi-line mode: ^ and $ match begin/end line in addition to begin/end text (default false)
// s	let . match \n (default false)
// U	ungreedy: swap meaning of x* and x*?, x+ and x+?, etc (default false)
//
// see https://pkg.go.dev/regexp/syntax
type Options struct {
	CaseInsensitive bool
	Multiline       bool
	DotNewLine      bool
	Ungreedy        bool
}

func (opts Options) Wrap(expr string) string {
	const PREFIX = "(?"
	var sb strings.Builder

	sb.WriteString(PREFIX)

	if opts.CaseInsensitive {
		sb.WriteString("i")
	}

	if opts.Multiline {
		sb.WriteString("m")
	}

	if opts.DotNewLine {
		sb.WriteString("s")
	}

	if opts.Ungreedy {
		sb.WriteString("U")
	}

	if sb.Len() == len(PREFIX) {
		// no flags added
		return expr
	}

	// `(?flags:re)`
	sb.WriteString(":")
	sb.WriteString(expr)
	sb.WriteString(")")

	return sb.String()
}
