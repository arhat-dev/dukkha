package rs

import (
	"unicode"
	"unicode/utf8"
)

// isExported is the copy of go/token.isExported
func isExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}
