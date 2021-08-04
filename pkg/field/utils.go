package field

import (
	"unicode"
	"unicode/utf8"
)

// IsExported is the copy of go/token.IsExported
func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}
