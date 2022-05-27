package templateutils

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Masterminds/goutils"
	"github.com/gosimple/slug"
	"github.com/huandu/xstrings"
	"mvdan.cc/sh/v3/syntax"
)

// TODO: support typed string
type stringsNS struct{}

func handleText1[T string](s String, do func(string) T) (_ T, err error) {
	str, err := toString(s)
	if err != nil {
		return
	}

	return do(str), nil
}

func handleText2[T string | bool | []string](s1, s2 String, do func(s1, s2 string) T) (_ T, err error) {
	str1, err := toString(s1)
	if err != nil {
		return
	}

	str2, err := toString(s2)
	if err != nil {
		return
	}

	return do(str1, str2), nil
}

func handleText3[T string](s1, s2, s3 String, do func(s1, s2, s3 string) T) (_ T, err error) {
	str1, err := toString(s1)
	if err != nil {
		return
	}

	str2, err := toString(s2)
	if err != nil {
		return
	}

	str3, err := toString(s3)
	if err != nil {
		return
	}

	return do(str1, str2, str3), nil
}

func (stringsNS) Abbrev(args ...any) (_ string, err error) {
	var (
		maxWidth, offset int

		str string
	)
	if len(args) < 2 {
		return "", fmt.Errorf("abbrev requires a 'maxWidth' and 'input' argument")
	}
	if len(args) == 2 {
		maxWidth = toIntegerOrPanic[int](args[0])
		str, err = toString(args[1])
		if err != nil {
			return
		}
	}
	if len(args) == 3 {
		offset = toIntegerOrPanic[int](args[0])
		maxWidth = toIntegerOrPanic[int](args[1])
		str, err = toString(args[2])
		if err != nil {
			return
		}
	}
	if len(str) <= maxWidth {
		return str, nil
	}

	return goutils.AbbreviateFull(str, offset, maxWidth)
}

func (stringsNS) Initials(s String) (string, error) {
	return handleText1(s, func(s string) string {
		return goutils.Initials(s)
	})
}

func (stringsNS) ReplaceAll(old, new, s String) (string, error) {
	return handleText3(s, old, new, strings.ReplaceAll)
}

func (stringsNS) Contains(substr, s String) (bool, error) {
	return handleText2(s, substr, strings.Contains)
}

func (stringsNS) ContainsAny(charset, s String) (bool, error) {
	return handleText2(s, charset, strings.ContainsAny)
}

func (stringsNS) Join(sep, strSlice any) (_ string, err error) {
	sepStr, err := toString(sep)
	if err != nil {
		return
	}

	parts, err := anyToStrings(strSlice)
	if err != nil {
		return
	}

	return strings.Join(parts, sepStr), nil
}

func (stringsNS) HasPrefix(prefix, s String) (bool, error) {
	return handleText2(s, prefix, strings.HasPrefix)
}

func (stringsNS) HasSuffix(suffix, s String) (bool, error) {
	return handleText2(s, suffix, strings.HasSuffix)
}

func (stringsNS) Repeat(count Number, s String) (_ string, err error) {
	str, err := toString(s)
	if err != nil {
		return
	}

	cnt, err := parseInteger[int](count)
	if err != nil {
		return
	}

	if cnt < 0 {
		err = fmt.Errorf("repeat: unexpected negative repeat count %d", cnt)
		return
	}

	return strings.Repeat(str, cnt), nil
}

func (stringsNS) Split(sep, s String) ([]string, error) {
	return handleText2(s, sep, strings.Split)
}

func (stringsNS) SplitN(sep String, n Number, s String) (_ []string, err error) {
	str, err := toString(s)
	if err != nil {
		return
	}

	sepStr, err := toString(sep)
	if err != nil {
		return
	}

	ni, err := parseInteger[int](n)
	if err != nil {
		return
	}

	return strings.SplitN(str, sepStr, ni), nil
}

func (stringsNS) Trim(cutset, s String) (string, error) { return handleText2(s, cutset, strings.Trim) }

func (stringsNS) TrimLeft(cutset, s String) (string, error) {
	return handleText2(s, cutset, strings.TrimLeft)
}

func (stringsNS) TrimRight(cutset, s String) (string, error) {
	return handleText2(s, cutset, strings.TrimRight)
}

func (stringsNS) TrimSpace(s String) (string, error) { return handleText1(s, strings.TrimSpace) }

func (stringsNS) TrimPrefix(prefix, s String) (string, error) {
	return handleText2(s, prefix, strings.TrimPrefix)
}

func (stringsNS) TrimSuffix(suffix, s String) (string, error) {
	return handleText2(s, suffix, strings.TrimSuffix)
}

// NoSpace removes all whitespaces in s
func (stringsNS) NoSpace(s String) (_ string, err error) {
	str, err := toString(s)
	if err != nil {
		return
	}

	return RemoveMatchedRunesCopy(str, unicode.IsSpace), nil
}

func (stringsNS) Title(s String) (string, error) {
	return handleText1(s, strings.Title)
}

func (stringsNS) Untitle(s String) (string, error) {
	return handleText1(s, func(s string) string {
		return goutils.Uncapitalize(s)
	})
}

// Substr creates a substring of the string (the last argument), start, end are rune index
//
// Substr(start Number, s String)
//
// Substr(start, end Number, s String)
func (stringsNS) Substr(args ...any) (_ string, err error) {
	n := len(args)
	if n < 2 {
		err = fmt.Errorf("at least 2 args expected, got %d", n)
		return
	}

	str, err := toString(args[n-1])
	if err != nil {
		return
	}

	low, err := parseInteger[int](args[0])
	if err != nil {
		return
	}

	high := math.MaxInt
	if n > 2 {
		high, err = parseInteger[int](args[1])
		if err != nil {
			return
		}
	}

	sz := utf8.RuneCountInString(str)
	low, high, _ = validSliceArgs(low, high, math.MaxInt, sz, sz)

	runeIdx := 0
	start, end := -1, -1
	for byteIdx := range str {
		if runeIdx == low {
			start = byteIdx
		}

		if runeIdx == high {
			end = byteIdx
		}

		runeIdx++
	}

	if start == -1 {
		start = len(str)
	}

	if end == -1 {
		end = len(str)
	}

	return str[start:end], nil
}

func (stringsNS) Upper(s String) (string, error) {
	return handleText1(s, strings.ToUpper)
}

func (stringsNS) Lower(s String) (string, error) {
	return handleText1(s, strings.ToLower)
}

func (stringsNS) Unquote(s String) (ret string, err error) {
	var err2 error
	ret, err = handleText1(s, func(s string) string {
		ret, err2 = strconv.Unquote(s)
		return ret
	})

	if err != nil {
		return
	}

	err = err2
	return
}

// Indent is a wrapper of {{ text.AddPrefix (text.Repeat n $indent) $data }}
//
// Indent(n Integer, data String): add n spaces as prefix to each line of data
//
// Indent(n Integer, indent, data String): add n * indent as prefix to each line of data
func (stringsNS) Indent(args ...any) (_ string, err error) { return indent(args, false) }

// NIndent is like Indent but prepends a newline ("\n") to the return value
func (stringsNS) NIndent(args ...any) (_ string, err error) { return indent(args, true) }

func indent(args []any, prependNewline bool) (_ string, err error) {
	n := len(args)
	if n < 2 {
		err = fmt.Errorf("at least 2 args expected, got %d", n)
		return
	}

	count, err := parseInteger[int](args[0])
	if err != nil {
		return
	}

	if count < 0 {
		err = fmt.Errorf("invalid negative count of indentation: %d", count)
		return
	}

	indent := " "
	if n > 2 {
		indent, err = toString(args[1])
		if err != nil {
			return
		}
	}

	data, err := toString(args[n-1])
	if err != nil {
		return
	}

	var sb strings.Builder

	if prependNewline {
		sb.WriteString("\n")
	}

	addPrefixW(&sb, data, strings.Repeat(indent, count), "\n")

	return sb.String(), nil
}

func (stringsNS) Slug(s String) (string, error) {
	return handleText1(s, slug.Make)
}

func (stringsNS) DoubleQuote(s String) (string, error) {
	return handleText1(s, strconv.Quote)
}

func (stringsNS) ShellQuote(args ...String) (_ string, err error) {
	parts, err := toStrings(args)
	if err != nil {
		return
	}

	for i, p := range parts {
		parts[i], err = syntax.Quote(p, syntax.LangBash)
		if err != nil {
			return
		}
	}

	return strings.Join(parts, " "), nil
}

func (stringsNS) SingleQuote(s String) (string, error) {
	return handleText1(s, func(s string) string {
		var sb strings.Builder
		sb.WriteString(`'`)
		sb.WriteString(strings.ReplaceAll(s, `'`, `''`))
		sb.WriteString(`'`)
		return sb.String()
	})
}

func (stringsNS) SnakeCase(s String) (string, error) {
	return handleText1(s, func(s string) string {
		return xstrings.ToSnakeCase(ReplaceMatchedRunesCopy(s,
			func(r rune) bool {
				return unicode.IsPunct(r) || unicode.IsSymbol(r)
			},
			func(matched string) string {
				return "-"
			},
		))
	})
}

func (stringsNS) CamelCase(s String) (string, error) {
	return handleText1(s, func(s string) string {
		return xstrings.ToCamelCase(ReplaceMatchedRunesCopy(s,
			func(r rune) bool {
				return unicode.IsPunct(r) || unicode.IsSymbol(r)
			},
			func(matched string) string {
				return "_"
			},
		))
	})
}

func (stringsNS) KebabCase(s String) (string, error) {
	return handleText1(s, func(s string) string {
		return xstrings.ToKebabCase(ReplaceMatchedRunesCopy(s,
			func(r rune) bool {
				return unicode.IsPunct(r) || unicode.IsSymbol(r)
			},
			func(matched string) string {
				return "_"
			},
		))
	})
}

func (stringsNS) Shuffle(s String) (string, error)  { return handleText1(s, xstrings.Shuffle) }
func (stringsNS) SwapCase(s String) (string, error) { return handleText1(s, goutils.SwapCase) }

func (stringsNS) WordWrap(args ...any) (_ string, err error) {
	n := len(args)

	if n == 0 || n > 3 {
		return "", fmt.Errorf("expected 1, 2, or 3 args, got %d", n)
	}

	in, err := toString(args[n-1])
	if err != nil {
		return
	}

	var opts WordWrapOpts
	if len(args) == 2 {
		switch a := (args[0]).(type) {
		case string:
			opts.LBSeq = a
		default:
			opts.Width = toIntegerOrPanic[uint](a)
		}
	}

	if len(args) == 3 {
		opts.Width = toIntegerOrPanic[uint](args[0])
		opts.LBSeq, err = toString(args[1])
		if err != nil {
			return
		}
	}

	return WordWrap(in, opts), nil
}

// WordWrapOpts defines the options to apply to the WordWrap function
type WordWrapOpts struct {
	// The desired maximum line length in characters (defaults to 80)
	Width uint

	// Line-break sequence to insert (defaults to "\n")
	LBSeq string
}

// applies default options
func wwDefaults(opts WordWrapOpts) WordWrapOpts {
	if opts.Width == 0 {
		opts.Width = 80
	}
	if opts.LBSeq == "" {
		opts.LBSeq = "\n"
	}
	return opts
}

// WordWrap - insert line-breaks into the string, before it reaches the given
// width.
func WordWrap(in string, opts WordWrapOpts) string {
	opts = wwDefaults(opts)
	return goutils.WrapCustom(in, int(opts.Width), opts.LBSeq, false)
}

// RuneCount - like len(s), but for runes
func (stringsNS) RuneCount(args ...String) (_ int, err error) {
	var sb strings.Builder

	parts, err := toStrings(args)
	if err != nil {
		return
	}

	for _, p := range parts {
		_, err = sb.WriteString(p)
		if err != nil {
			return
		}
	}

	return utf8.RuneCountInString(sb.String()), nil
}

/*

	Start of Multi-Section string processing (usually multi-line)

*/

// TODO: support writer as the second last argument
func (stringsNS) AddPrefix(args ...String) (string, error) {
	return handleMultiSectionText_OPDATA_OptionalSep_DATA(args, AddPrefix)
}

// TODO: support writer as the second last argument
func (stringsNS) RemovePrefix(args ...String) (string, error) {
	return handleMultiSectionText_OPDATA_OptionalSep_DATA(args, RemovePrefix)
}

// TODO: support writer as the second last argument
func (stringsNS) AddSuffix(args ...String) (string, error) {
	return handleMultiSectionText_OPDATA_OptionalSep_DATA(args, AddSuffix)
}

// TODO: support writer as the second last argument
func (stringsNS) RemoveSuffix(args ...String) (string, error) {
	return handleMultiSectionText_OPDATA_OptionalSep_DATA(args, RemoveSuffix)
}

func forEachTextSection(s, sep string, do func(section string)) {
	n := len(s)
	szSep := len(sep)

	for start, i := 0, 0; start < n; start += i + szSep {
		i = strings.Index(s[start:], sep)
		if i == -1 {
			if start < n {
				do(s[start:])
			}

			break
		}

		do(s[start : start+i])
	}
}

// AddPrefix to each separated string elements
func AddPrefix(s, prefix, sep string) string {
	var sb strings.Builder
	addPrefixW(&sb, s, prefix, sep)
	return sb.String()
}

func addPrefixW(w io.StringWriter, s, prefix, sep string) {
	forEachTextSection(s, sep, func(section string) {
		w.WriteString(prefix)
		w.WriteString(section)
		w.WriteString(sep)
	})
}

// RemovePrefix of each separated string elements
func RemovePrefix(s, prefix, sep string) string {
	var sb strings.Builder

	forEachTextSection(s, sep, func(section string) {
		sb.WriteString(strings.TrimPrefix(section, prefix))
		sb.WriteString(sep)
	})

	return sb.String()
}

// AddSuffix to each separated string elements
func AddSuffix(s, suffix, sep string) string {
	var sb strings.Builder

	forEachTextSection(s, sep, func(section string) {
		sb.WriteString(section)
		sb.WriteString(suffix)
		sb.WriteString(sep)
	})

	return sb.String()
}

func RemoveSuffix(s, suffix, sep string) string {
	var sb strings.Builder

	forEachTextSection(s, sep, func(section string) {
		sb.WriteString(strings.TrimSuffix(section, suffix))
		sb.WriteString(sep)
	})

	return sb.String()
}

func handleMultiSectionText_OPDATA_OptionalSep_DATA(args []String, do func(str, op, sep string) string) (_ string, err error) {
	var (
		opData any
		data   any
		sep    any
	)

	switch n := len(args); n {
	case 0, 1:
		err = fmt.Errorf("invalid args: at least 2 args expected, got %d", n)
		return
	case 2:
		opData, sep, data = args[0], "\n", args[1]
	default:
		opData, sep, data = args[0], args[1], args[n-1]
	}

	op, err := toString(opData)
	if err != nil {
		return
	}

	sepStr, err := toString(sep)
	if err != nil {
		return
	}

	str, err := toString(data)
	if err != nil {
		return
	}

	return do(str, op, sepStr), nil
}

/*

	End of Multi-Section string processing

*/
