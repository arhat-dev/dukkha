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
	itemText       // plain text
	itemVariable   // variable starting with '$', such as '$' or  '$1' or '$hello'
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

func key(k string) itemType {
	switch k {
	case ".":
		return itemDot
	case "block":
		return itemBlock
	case "break":
		return itemBreak
	case "continue":
		return itemContinue
	case "define":
		return itemDefine
	case "else":
		return itemElse
	case "end":
		return itemEnd
	case "if":
		return itemIf
	case "range":
		return itemRange
	case "nil":
		return itemNil
	case "template":
		return itemTemplate
	case "with":
		return itemWith
	default:
		return 0
	}
}

const eof = -1

// Trimming spaces.
// If the action begins "{{- " rather than "{{", then all space/tab/newlines
// preceding the action are trimmed; conversely if it ends " -}}" the
// leading spaces are trimmed. This is done entirely in the lexer; the
// parser never sees it happen. We require an ASCII space (' ', \t, \r, \n)
// to be present to avoid ambiguity with things like "{{-3}}". It reads
// better with the space present anyway. For simplicity, only ASCII
// does the job.
const (
	spaceChars    = " \t\r\n"  // These are the space characters defined by Go itself.
	trimMarker    = '-'        // Attached to left/right delimiter, trims trailing spaces from preceding/following text.
	trimMarkerLen = Pos(1 + 1) // marker plus space before or after
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) (ret item, next stateFn)

// lexer holds the state of the scanner.
type lexer struct {
	name        string // the name of the input; used only for error reports
	input       string // the string being scanned
	leftDelim   string // start of action
	rightDelim  string // end of action
	emitComment bool   // emit itemComment tokens.
	pos         Pos    // current position in the input
	start       Pos    // start position of this item
	width       Pos    // width of last rune read from input
	parenDepth  int    // nesting depth of ( ) exprs
	line        int    // 1+number of newlines seen
	startLine   int    // start line of this item
	breakOK     bool   // break keyword allowed
	continueOK  bool   // continue keyword allowed

	state stateFn
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
	r := l.next()
	l.backup()
	return r
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

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.line += strings.Count(l.input[l.start:l.pos], "\n")
	l.start = l.pos
	l.startLine = l.line
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
	for l.state != nil {
		ret, l.state = l.state(l)
		if l.state == nil {
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
func lex(name, input, left, right string, emitComment bool) *lexer {
	if left == "" {
		left = leftDelim
	}
	if right == "" {
		right = rightDelim
	}
	l := &lexer{
		name:        name,
		input:       input,
		leftDelim:   left,
		rightDelim:  right,
		emitComment: emitComment,
		line:        1,
		startLine:   1,

		state: lexText,
	}
	return l
}

// state functions

const (
	leftDelim    = "{{"
	rightDelim   = "}}"
	leftComment  = "/*"
	rightComment = "*/"
)

// lexText scans until an opening action delimiter, "{{".
func lexText(l *lexer) (ret item, next stateFn) {
	l.width = 0
	if x := strings.Index(l.input[l.pos:], l.leftDelim); x >= 0 {
		ldn := Pos(len(l.leftDelim))
		l.pos += Pos(x)
		trimLength := Pos(0)
		if hasLeftTrimMarker(l.input[l.pos+ldn:]) {
			trimLength = rightTrimLength(l.input[l.start:l.pos])
		}
		l.pos -= trimLength
		if l.pos > l.start {
			l.line += strings.Count(l.input[l.start:l.pos], "\n")
			ret = l.emit(itemText)
		}
		l.pos += trimLength
		l.ignore()
		next = lexLeftDelim
		return
	}
	l.pos = Pos(len(l.input))
	// Correctly reached EOF.
	if l.pos > l.start {
		l.line += strings.Count(l.input[l.start:l.pos], "\n")
		ret = l.emit(itemText)
		return
	}
	ret = l.emit(itemEOF)
	return
}

// rightTrimLength returns the length of the spaces at the end of the string.
func rightTrimLength(s string) Pos {
	return Pos(len(s) - len(strings.TrimRight(s, spaceChars)))
}

// atRightDelim reports whether the lexer is at a right delimiter, possibly preceded by a trim marker.
func (l *lexer) atRightDelim() (delim, trimSpaces bool) {
	if hasRightTrimMarker(l.input[l.pos:]) && strings.HasPrefix(l.input[l.pos+trimMarkerLen:], l.rightDelim) { // With trim marker.
		return true, true
	}
	if strings.HasPrefix(l.input[l.pos:], l.rightDelim) { // Without trim marker.
		return true, false
	}
	return false, false
}

// leftTrimLength returns the length of the spaces at the beginning of the string.
func leftTrimLength(s string) Pos {
	return Pos(len(s) - len(strings.TrimLeft(s, spaceChars)))
}

// lexLeftDelim scans the left delimiter, which is known to be present, possibly with a trim marker.
func lexLeftDelim(l *lexer) (ret item, next stateFn) {
	l.pos += Pos(len(l.leftDelim))
	trimSpace := hasLeftTrimMarker(l.input[l.pos:])
	afterMarker := Pos(0)
	if trimSpace {
		afterMarker = trimMarkerLen
	}
	if strings.HasPrefix(l.input[l.pos+afterMarker:], leftComment) {
		l.pos += afterMarker
		l.ignore()
		next = lexComment
		return
	}
	ret = l.emit(itemLeftDelim)
	next = lexInsideAction
	l.pos += afterMarker
	l.ignore()
	l.parenDepth = 0
	return
}

// lexComment scans a comment. The left comment marker is known to be present.
func lexComment(l *lexer) (ret item, next stateFn) {
	l.pos += Pos(len(leftComment))
	i := strings.Index(l.input[l.pos:], rightComment)
	if i < 0 {
		ret = l.errorf("unclosed comment")
		return
	}
	l.pos += Pos(i + len(rightComment))
	delim, trimSpace := l.atRightDelim()
	if !delim {
		ret = l.errorf("comment ends before closing delimiter")
		return
	}
	if l.emitComment {
		ret = l.emit(itemComment)
	}
	if trimSpace {
		l.pos += trimMarkerLen
	}
	l.pos += Pos(len(l.rightDelim))
	if trimSpace {
		l.pos += leftTrimLength(l.input[l.pos:])
	}
	l.ignore()
	next = lexText
	return
}

// lexRightDelim scans the right delimiter, which is known to be present, possibly with a trim marker.
func lexRightDelim(l *lexer) (ret item, next stateFn) {
	trimSpace := hasRightTrimMarker(l.input[l.pos:])
	if trimSpace {
		l.pos += trimMarkerLen
		l.ignore()
	}
	l.pos += Pos(len(l.rightDelim))
	ret = l.emit(itemRightDelim)
	if trimSpace {
		l.pos += leftTrimLength(l.input[l.pos:])
		l.ignore()
	}
	next = lexText
	return
}

// lexInsideAction scans the elements inside action delimiters.
func lexInsideAction(l *lexer) (ret item, next stateFn) {
	// Either number, quoted string, or identifier.
	// Spaces separate arguments; runs of spaces turn into itemSpace.
	// Pipe symbols separate and are emitted.
	delim, _ := l.atRightDelim()
	if delim {
		if l.parenDepth == 0 {
			next = lexRightDelim
			return
		}
		ret = l.errorf("unclosed left paren")
		return
	}
	switch r := l.next(); {
	case r == eof:
		ret = l.errorf("unclosed action")
		return
	case isSpace(r):
		l.backup() // Put space back in case we have " -}}".
		next = lexSpace
		return
	case r == '=':
		ret = l.emit(itemAssign)
	case r == ':':
		if l.next() != '=' {
			ret = l.errorf("expected :=")
			return
		}
		ret = l.emit(itemDeclare)
	case r == '|':
		ret = l.emit(itemPipe)
	case r == '"':
		next = lexQuote
		return
	case r == '`':
		next = lexRawQuote
		return
	case r == '$':
		next = lexVariable
		return
	case r == '\'':
		next = lexChar
		return
	case r == '.':
		// special look-ahead for ".field" so we don't break l.backup().
		if l.pos < Pos(len(l.input)) {
			r := l.input[l.pos]
			if r < '0' || '9' < r {
				next = lexField
				return
			}
		}
		fallthrough // '.' can start a number.
	case r == '+' || r == '-' || ('0' <= r && r <= '9'):
		l.backup()
		next = lexNumber
		return
	case isAlphaNumeric(r):
		l.backup()
		next = lexIdentifier
		return
	case r == '(':
		ret = l.emit(itemLeftParen)
		l.parenDepth++
	case r == ')':
		ret = l.emit(itemRightParen)
		l.parenDepth--
		if l.parenDepth < 0 {
			ret = l.errorf("unexpected right paren %#U", r)
			return
		}
	case r <= unicode.MaxASCII && unicode.IsPrint(r):
		ret = l.emit(itemChar)
	default:
		ret = l.errorf("unrecognized character in action: %#U", r)
		return
	}

	next = lexInsideAction
	return
}

// lexSpace scans a run of space characters.
// We have not consumed the first space, which is known to be present.
// Take care if there is a trim-marked right delimiter, which starts with a space.
func lexSpace(l *lexer) (ret item, next stateFn) {
	var r rune
	var numSpaces int
	for {
		r = l.peek()
		if !isSpace(r) {
			break
		}
		l.next()
		numSpaces++
	}
	// Be careful about a trim-marked closing delimiter, which has a minus
	// after a space. We know there is a space, so check for the '-' that might follow.
	if hasRightTrimMarker(l.input[l.pos-1:]) && strings.HasPrefix(l.input[l.pos-1+trimMarkerLen:], l.rightDelim) {
		l.backup() // Before the space.
		if numSpaces == 1 {
			next = lexRightDelim // On the delim, so go right to that.
			return
		}
	}

	return l.emit(itemSpace), lexInsideAction
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer) (ret item, next stateFn) {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atTerminator() {
				ret = l.errorf("bad character %#U", r)
				return
			}
			switch {
			case key(word) > itemKeyword:
				item := key(word)
				if item == itemBreak && !l.breakOK || item == itemContinue && !l.continueOK {
					ret = l.emit(itemIdentifier)
				} else {
					ret = l.emit(item)
				}
			case word[0] == '.':
				ret = l.emit(itemField)
			case word == "true", word == "false":
				ret = l.emit(itemBool)
			default:
				ret = l.emit(itemIdentifier)
			}
			break Loop
		}
	}

	next = lexInsideAction
	return
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
	var r rune
	for {
		r = l.next()
		if !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}
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
	r := l.peek()
	if isSpace(r) {
		return true
	}
	switch r {
	case eof, '.', ',', '|', ':', ')', '(':
		return true
	}
	// Does r start the delimiter? This can be ambiguous (with delim=="//", $x/2 will
	// succeed but should fail) but only in extremely rare cases caused by willfully
	// bad choice of delimiter.
	if rd, _ := utf8.DecodeRuneInString(l.rightDelim); rd == r {
		return true
	}
	return false
}

// lexChar scans a character constant. The initial quote is already
// scanned. Syntax checking is done by the parser.
func lexChar(l *lexer) (ret item, next stateFn) {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			ret = l.errorf("unterminated character constant")
			return
		case '\'':
			break Loop
		}
	}

	return l.emit(itemCharConstant), lexInsideAction
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
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			ret = l.errorf("unterminated quoted string")
			return
		case '"':
			break Loop
		}
	}

	return l.emit(itemString), lexInsideAction
}

// lexRawQuote scans a raw quoted string.
func lexRawQuote(l *lexer) (ret item, next stateFn) {
Loop:
	for {
		switch l.next() {
		case eof:
			ret = l.errorf("unterminated raw quoted string")
			return
		case '`':
			break Loop
		}
	}

	return l.emit(itemRawString), lexInsideAction
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func hasLeftTrimMarker(s string) bool {
	return len(s) >= 2 && s[0] == trimMarker && isSpace(rune(s[1]))
}

func hasRightTrimMarker(s string) bool {
	return len(s) >= 2 && isSpace(rune(s[0])) && s[1] == trimMarker
}
