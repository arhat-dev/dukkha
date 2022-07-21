// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// item represents a token or text string returned from the scanner.
type item struct {
	typ  itemType // The type of this item.
	pos  Pos      // The starting position, in bytes, of this item in the input string.
	val  string   // The value of this item.
	line int      // The line number at the start of this item.

	notEmpty bool
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError        itemType = iota // error occurred; value is text of error
	itemBool                         // boolean constant
	itemChar                         // printable ASCII character; grab bag for comma etc.
	itemCharConstant                 // character constant
	itemComment                      // comment text
	itemComplex                      // complex constant (1+2i); imaginary is just a number
	itemAssign                       // equals ('=') introducing an assignment
	itemDeclare                      // colon-equals (':=') introducing a declaration
	itemEOF
	itemField      // alphanumeric identifier starting with '.'
	itemIdentifier // alphanumeric identifier not starting with '.'
	itemLeftDelim  // left action delimiter
	itemLeftParen  // '(' inside action
	itemNumber     // simple number, including imaginary
	itemPipe       // pipe symbol
	itemRawString  // raw quoted string (includes quotes)
	itemRightDelim // right action delimiter
	itemRightParen // ')' inside action
	itemSpace      // run of spaces separating arguments
	itemString     // quoted string (includes quotes)
	// itemText       // plain text
	itemVariable // variable starting with '$', such as '$' or  '$1' or '$hello'
	// Keywords appear after all the rest.
	itemKeyword  // used only to delimit the keywords
	itemBlock    // block keyword
	itemBreak    // break keyword
	itemContinue // continue keyword
	itemDot      // the cursor, spelled '.'
	itemDefine   // define keyword
	itemElse     // else keyword
	itemEnd      // end keyword
	itemIf       // if keyword
	itemNil      // the untyped nil constant, easiest to treat as a keyword
	itemRange    // range keyword
	itemTemplate // template keyword
	itemWith     // with keyword
)

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) (ret item, next stateFn)

// lexer holds the state of the scanner.
type lexer struct {
	name  string // the name of the input; used only for error reports
	input string // the string being scanned
	// leftDelim   string // start of action
	emitComment bool // emit itemComment tokens.
	pos         Pos  // current position in the input
	start       Pos  // start position of this item
	width       Pos  // width of last rune read from input
	parenDepth  int  // nesting depth of ( ) exprs
	line        int  // 1+number of newlines seen
	startLine   int  // start line of this item

	nextState stateFn
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	for _, r := range l.input[l.pos:] {
		return r
	}

	return eof
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) (ret item) {
	ret = item{t, l.start, l.input[l.start:l.pos], l.startLine, true}
	l.start = l.pos
	l.startLine = l.line
	return
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...any) item {
	return item{itemError, l.start, fmt.Sprintf(format, args...), l.startLine, true}
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() (ret item) {
	for l.nextState != nil {
		ret, l.nextState = l.nextState(l)
		if l.nextState == nil {
			// is the last state, return unconditionally
			return
		}

		if ret.notEmpty {
			// got a meaningful item
			return
		}

		// empty item, scan next
	}

	// fake EOF
	return l.emit(itemEOF)
}

// lex creates a new scanner for the input string.
func lex(name, input string, emitComment bool) *lexer {
	l := &lexer{
		name:        name,
		input:       input,
		emitComment: emitComment,
		line:        1,
		startLine:   1,

		nextState: lexWhitespace,
	}
	return l
}

// state functions

// lexWhitespace eats all whitespace prefix, when there is non-whitespace found, return
// itemLeftDelim
func lexWhitespace(l *lexer) (ret item, next stateFn) {
	var (
		i int
		r rune = eof
	)

	data := l.input[l.pos:]
	for i, r = range data {
		switch r {
		case ' ', '\r', '\t':
			continue
		case '\n':
			l.line++
			continue
		}

		// not a whitespace
		break
	}

	if r == eof {
		l.width = 0
	} else {
		if i != 0 {
			l.pos += Pos(i)
			l.width = 1 // all whitespaces are 1 byte in width
		}
	}

	l.start = l.pos
	l.startLine = l.line
	switch r {
	case ' ', '\r', '\t', '\n', eof:
		// when r end up being whitespace, we MUST have reached EOF
		return l.emit(itemEOF), nil
	case '#':
		return lexComment(l)
	}

	return l.emit(itemLeftDelim), lexInsideAction
}

// lexComment scans a comment line with prefix '#'
func lexComment(l *lexer) (ret item, next stateFn) {
	i := strings.IndexByte(l.input[l.pos:], '\n')

	if i < 0 {
		l.pos = Pos(len(l.input))
	} else {
		l.width = 1 // '\n'
		l.pos += Pos(i + 1)
		l.line++
	}

	if l.emitComment {
		ret = l.emit(itemComment)
		ret.pos++
		ret.val = ret.val[1:] // trim '#'
	} else {
		l.start = l.pos
		l.startLine = l.line
	}

	return ret, lexWhitespace
}

// lexInsideAction scans the elements inside action delimiters.
func lexInsideAction(l *lexer) (ret item, next stateFn) {
	var (
		i int
		r rune
	)

	data := l.input[l.pos:]

	if len(data) == 0 {
		if l.parenDepth > 0 {
			return l.errorf("unclosed left paren"), nil
		}

		// schedule lexWhitespace as next to emit EOF
		return l.emit(itemRightDelim), lexWhitespace
	}

	for i, r = range data {
		break
	}

	// fast path for identifiers and constants (template funcs)
	switch {
	case r >= '0' && r <= '9':
		return lexNumber(l)
	case isAlphaNumeric(r):
		return lexIdentifier(l)
	}

	switch r {
	case ' ', '\n', '\t', '\r':
		return lexInActionSpace(l)
	case '.':
		if i == len(data)-1 {
			// at the end of the input
			l.width = 1
			l.pos += 1
			return l.emit(itemDot), lexInsideAction
		}

		if data[i+1] < '0' || data[i+1] > '9' {
			l.width = 1
			l.pos += 1
			return lexField(l)
		}

		// .[0-9]
		return lexNumber(l)
	case '|':
		l.width = 1
		l.pos += 1
		return l.emit(itemPipe), lexInsideAction
	case '=':
		l.width = 1
		l.pos += 1
		return l.emit(itemAssign), lexInsideAction
	case ':':
		if i == len(data)-1 || data[i+1] != '=' {
			return l.errorf("expected :="), nil
		}

		l.width = 1
		l.pos += 2
		return l.emit(itemDeclare), lexInsideAction
	case '"':
		l.width = 1
		l.pos += 1
		return lexQuote(l)
	case '`':
		l.width = 1
		l.pos += 1
		return lexRawQuote(l)
	case '$':
		l.width = 1
		l.pos += 1
		return lexVariable(l)
	case '\'':
		l.width = 1
		l.pos += 1
		return lexChar(l)
	case '(':
		l.width = 1
		l.pos += 1
		ret = l.emit(itemLeftParen)
		l.parenDepth++
		return ret, lexInsideAction
	case ')':
		l.width = 1
		l.pos += 1
		ret = l.emit(itemRightParen)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r), nil
		}

		return ret, lexInsideAction
	case '+', '-':
		return lexNumber(l)
	case ';':
		l.width = 1
		l.pos += 1

		l.start = l.pos
		return l.emit(itemRightDelim), lexWhitespace
	default:
		if r <= unicode.MaxASCII && unicode.IsPrint(r) {
			// punctuations

			l.width = 1 // ascii
			l.pos += 1
			return l.emit(itemChar), lexInActionSpace
		}

		return l.errorf("unrecognized character in action: %#U", r), nil
	}
}

// lexInActionSpace scans a run of space characters.
func lexInActionSpace(l *lexer) (ret item, next stateFn) {
	var (
		i int
		r rune = eof

		hasInlineBackslash bool
		emitRightDelim     bool
	)

	data := l.input[l.pos:]
	for i, r = range data {
		switch r {
		case ' ', '\t', '\r':
			continue
		case '\n':
			l.line++
			if hasInlineBackslash {
				// when we have backslash as last non-whitespace char in this line
				// the action continues
				hasInlineBackslash = false
				continue
			}

			emitRightDelim = true
			i++ // consume this newline
		case '\\':
			hasInlineBackslash = true
			continue
		}

		break
	}

	l.width = 1 // all whitespaces are 1 byte in width
	l.pos += Pos(i)

	switch r {
	case ' ', '\t', '\r':
		l.pos++
		fallthrough
	case '\n':
		if l.pos >= Pos(len(l.input)) {
			// reached EOF
			emitRightDelim = true
		}
	case ';':
		l.pos++ // include this one
		emitRightDelim = true
	case '#':
		emitRightDelim = true
	case eof:
		l.width = 0
		l.start = l.pos
		l.startLine = l.line
		return l.emit(itemEOF), nil
	}

	if emitRightDelim {
		l.start = l.pos
		l.startLine = l.line
		return l.emit(itemRightDelim), lexWhitespace
	}

	if i == 0 {
		return lexInsideAction(l)
	}

	return l.emit(itemSpace), lexInsideAction
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer) (ret item, next stateFn) {
	var (
		i int
		r rune
	)
	data := l.input[l.pos:]

	for i, r = range data {
		if !isAlphaNumeric(r) {
			break
		}
	}

	if isAlphaNumeric(r) { // r is the last rune
		i = len(data)
		l.width = Pos(utf8.RuneLen(r))
	} else {
		_, sz := utf8.DecodeLastRuneInString(data[:i])
		l.width = Pos(sz)
	}

	l.pos += Pos(i)
	if !l.atTerminator() {
		return l.errorf("bad character %#U", r), nil
	}

	switch data[:i] {
	case ".":
		return l.emit(itemDot), lexInsideAction
	case "block":
		return l.emit(itemBlock), lexInsideAction
	case "break":
		return l.emit(itemBreak), lexInsideAction
	case "continue":
		return l.emit(itemContinue), lexInsideAction
	case "define":
		return l.emit(itemDefine), lexInsideAction
	case "else":
		return l.emit(itemElse), lexInsideAction
	case "end":
		return l.emit(itemEnd), lexInsideAction
	case "if":
		return l.emit(itemIf), lexInsideAction
	case "range":
		return l.emit(itemRange), lexInsideAction
	case "nil":
		return l.emit(itemNil), lexInsideAction
	case "template":
		return l.emit(itemTemplate), lexInsideAction
	case "with":
		return l.emit(itemWith), lexInsideAction
	case "true", "false":
		return l.emit(itemBool), lexInsideAction
	default:
		if data[0] == '.' {
			return l.emit(itemField), lexInsideAction
		}

		return l.emit(itemIdentifier), lexInsideAction
	}
}

// lexField scans a field: .Alphanumeric.
// The . has been scanned.
func lexField(l *lexer) (ret item, next stateFn) {
	return lexFieldOrVariable(l, itemField)
}

// lexVariable scans a Variable: $Alphanumeric.
// The $ has been scanned.
func lexVariable(l *lexer) (ret item, next stateFn) {
	if l.atTerminator() { // Nothing interesting follows -> "$".
		return l.emit(itemVariable), lexInsideAction
	}
	return lexFieldOrVariable(l, itemVariable)
}

// lexVariable scans a field or variable: [.$]Alphanumeric.
// The . or $ has been scanned.
func lexFieldOrVariable(l *lexer, typ itemType) (ret item, next stateFn) {
	if l.atTerminator() { // Nothing interesting follows -> "." or "$".
		if typ == itemVariable {
			return l.emit(itemVariable), lexInsideAction
		}

		return l.emit(itemDot), lexInsideAction
	}
	var (
		i int
		r rune
	)

	data := l.input[l.pos:]
	for i, r = range data {
		if !isAlphaNumeric(r) {
			break
		}
	}

	if isAlphaNumeric(r) { // r is the last rune in data
		i = len(data)
		l.width = Pos(utf8.RuneLen(r))
	} else {
		_, sz := utf8.DecodeLastRuneInString(data[:i])
		l.width = Pos(sz)
	}

	l.pos += Pos(i)

	if !l.atTerminator() {
		ret = l.errorf("bad character %#U", r)
		return
	}

	return l.emit(typ), lexInsideAction
}

// atTerminator reports whether the input is at valid termination character to
// appear after an identifier. Breaks .X.Y into two pieces. Also catches cases
// like "$x+2" not being acceptable without a space, in case we decide one
// day to implement arithmetic.
func (l *lexer) atTerminator() bool {
	if l.pos >= Pos(len(l.input)) { // EOF
		return true
	}

	switch l.input[l.pos] {
	case '.', ',', '|', ':', ')', '(', ' ', '\t', '\r', '\n', ';':
		return true
	default:
		return false
	}
}

// lexChar scans a character constant. The initial quote is already
// scanned. Syntax checking is done by the parser.
func lexChar(l *lexer) (ret item, next stateFn) {
	var (
		i int
		r rune

		hasEnd bool
	)

	data := l.input[l.pos:]

Loop:
	for i, r = range data {
		switch r {
		case '\\':
			if i != len(data)-1 && data[i+1] != '\n' { // at the end of the data
				break
			}

			fallthrough
		case eof, '\n':
			goto ERR
		case '\'':
			if i == 1 && data[0] == '\\' {
				// we can only have one escaped single quote in a char constant
				continue
			}

			hasEnd = true
			break Loop
		}
	}

	if hasEnd {
		l.width = 1 // '\''
		l.pos += Pos(i) + 1

		return l.emit(itemCharConstant), lexInsideAction
	}

ERR:
	return l.errorf("unterminated character constant"), nil
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary. This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexNumber(l *lexer) (ret item, next stateFn) {
	if !l.scanNumber() {
		ret = l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		return
	}
	if sign := l.peek(); sign == '+' || sign == '-' {
		// Complex: 1+2i. No spaces, must end in 'i'.
		if !l.scanNumber() || l.input[l.pos-1] != 'i' {
			ret = l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
			return
		}

		return l.emit(itemComplex), lexInsideAction
	}

	return l.emit(itemNumber), lexInsideAction
}

func (l *lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789_"
	if l.accept("0") {
		// Note: Leading 0 does not mean octal in floats.
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			digits = "01234567_"
		} else if l.accept("bB") {
			digits = "01_"
		}
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	if len(digits) == 16+6+1 && l.accept("pP") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// lexQuote scans a quoted string.
func lexQuote(l *lexer) (ret item, next stateFn) {
	var (
		i int
		r rune

		hasEnd bool
	)

	data := l.input[l.pos:]
Loop:
	for i, r = range data {
		switch r {
		case '\\':
			if i != len(data)-1 && data[i+1] != '\n' { // at the end of the data
				break
			}

			fallthrough
		case eof, '\n':
			goto ERR
		case '"':
			if i != 0 && data[i-1] == '\\' {
				continue
			}
			hasEnd = true
			break Loop
		}
	}

	if hasEnd {
		l.width = 1 // '"'
		l.pos += Pos(i) + 1

		return l.emit(itemString), lexInsideAction
	}

ERR:
	return l.errorf("unterminated quoted string"), nil
}

// lexRawQuote scans a raw quoted string.
func lexRawQuote(l *lexer) (item, stateFn) {
	var (
		i int
		r rune

		hasEnd bool
	)

Loop:
	for i, r = range l.input[l.pos:] {
		switch r {
		case '\n':
			l.line++
		case '`':
			hasEnd = true
			break Loop
		}
	}

	if !hasEnd {
		return l.errorf("unterminated raw quoted string"), nil
	}

	l.width = 1
	l.pos += Pos(i) + 1
	return l.emit(itemRawString), lexInsideAction
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
