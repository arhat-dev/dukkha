package templateutils

import (
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/goutils"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"

	gompstrings "arhat.dev/dukkha/third_party/gomplate/strings"
)

type stringsNS struct{}

func (stringsNS) Abbrev(args ...interface{}) (string, error) {
	var (
		maxWidth, offset int

		str string
	)
	if len(args) < 2 {
		return "", errors.Errorf("abbrev requires a 'maxWidth' and 'input' argument")
	}
	if len(args) == 2 {
		maxWidth = int(toUint64(args[0]))
		str = toString(args[1])
	}
	if len(args) == 3 {
		offset = int(toUint64(args[0]))
		maxWidth = int(toUint64(args[1]))
		str = toString(args[2])
	}
	if len(str) <= maxWidth {
		return str, nil
	}

	return goutils.AbbreviateFull(str, offset, maxWidth)
}

func (stringsNS) ReplaceAll(old, new, s String) string {
	return strings.ReplaceAll(toString(s), toString(old), toString(new))
}

func (stringsNS) Contains(substr, s String) bool {
	return strings.Contains(toString(s), toString(substr))
}

func (stringsNS) HasPrefix(prefix, s String) bool {
	return strings.HasPrefix(toString(s), toString(prefix))
}

func (stringsNS) HasSuffix(suffix, s String) bool {
	return strings.HasSuffix(toString(s), toString(suffix))
}

func (stringsNS) Repeat(cnt Integer, s String) (string, error) {
	count := int(toUint64(cnt))
	if count < 0 {
		return "", errors.Errorf("negative count %d", count)
	}
	str := toString(s)
	if count > 0 && len(str)*count/count != len(str) {
		return "", errors.Errorf("count %d too long: causes overflow", count)
	}
	return strings.Repeat(str, count), nil
}

func (stringsNS) Split(sep, s String) []string {
	return strings.Split(toString(s), toString(sep))
}

func (stringsNS) SplitN(sep String, n Integer, s String) []string {
	return strings.SplitN(toString(s), toString(sep), int(toUint64(n)))
}

func (stringsNS) Trim(cutset, s String) string {
	return strings.Trim(toString(s), toString(cutset))
}

func (stringsNS) TrimPrefix(cutset, s String) string {
	return strings.TrimPrefix(toString(s), toString(cutset))
}

func (stringsNS) TrimSuffix(cutset, s String) string {
	return strings.TrimSuffix(toString(s), toString(cutset))
}

func (stringsNS) Title(s String) string {
	return strings.Title(toString(s))
}

func (stringsNS) ToUpper(s String) string {
	return strings.ToUpper(toString(s))
}

func (stringsNS) ToLower(s String) string {
	return strings.ToLower(toString(s))
}

func (stringsNS) TrimSpace(s String) string {
	return strings.TrimSpace(toString(s))
}

func (stringsNS) Trunc(length Integer, s String) string {
	return gompstrings.Trunc(int(toUint64(length)), toString(s))
}

func (stringsNS) Indent(args ...interface{}) (string, error) {
	input := toString(args[len(args)-1])
	indent := " "
	width := 1
	var ok bool
	switch len(args) {
	case 2:
		indent, ok = args[0].(string)
		if !ok {
			width, ok = args[0].(int)
			if !ok {
				return "", errors.New("indent: invalid arguments")
			}
			indent = " "
		}
	case 3:
		width, ok = args[0].(int)
		if !ok {
			return "", errors.New("indent: invalid arguments")
		}
		indent, ok = args[1].(string)
		if !ok {
			return "", errors.New("indent: invalid arguments")
		}
	}

	return gompstrings.Indent(width, indent, input), nil
}

func (stringsNS) Slug(s String) string {
	return slug.Make(toString(s))
}

func (stringsNS) Quote(s String) string {
	return strconv.Quote(toString(s))
}

func (stringsNS) ShellQuote(s String) string {
	val := reflect.ValueOf(s)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		var sb strings.Builder
		max := val.Len()
		for n := 0; n < max; n++ {
			sb.WriteString(gompstrings.ShellQuote(toString(val.Index(n))))
			if n+1 != max {
				sb.WriteRune(' ')
			}
		}
		return sb.String()
	}
	return gompstrings.ShellQuote(toString(s))
}

func (stringsNS) Squote(s String) string {
	var sb strings.Builder
	sb.WriteString(`'`)
	sb.WriteString(strings.ReplaceAll(toString(s), `'`, `''`))
	sb.WriteString(`'`)
	return sb.String()
}

func (stringsNS) SnakeCase(s String) string {
	return gompstrings.SnakeCase(toString(s))
}

func (stringsNS) CamelCase(s String) string {
	return gompstrings.CamelCase(toString(s))
}

func (stringsNS) KebabCase(s String) string {
	str := toString(s)
	if len(str) == 0 {
		return str
	}

	return gompstrings.KebabCase(str)
}

func (stringsNS) WordWrap(args ...interface{}) (string, error) {
	if len(args) == 0 || len(args) > 3 {
		return "", errors.Errorf("expected 1, 2, or 3 args, got %d", len(args))
	}
	in := toString(args[len(args)-1])

	opts := gompstrings.WordWrapOpts{}
	if len(args) == 2 {
		switch a := (args[0]).(type) {
		case string:
			opts.LBSeq = a
		default:
			opts.Width = uint(toUint64(a))
		}
	}
	if len(args) == 3 {
		opts.Width = uint(toUint64(args[0]))
		opts.LBSeq = toString(args[1])
	}
	return gompstrings.WordWrap(in, opts), nil
}

// RuneCount - like len(s), but for runes
func (stringsNS) RuneCount(args ...String) (int, error) {
	s := ""
	for _, arg := range args {
		s += toString(arg)
	}
	return utf8.RuneCountInString(s), nil
}

/*

	Start of Multi-Section string processing

*/

// AddPrefix to each separated string elements
func (stringsNS) AddPrefix(s, prefix, sep string) string {
	parts := strings.Split(s, sep)
	for i, p := range parts {
		if len(p) == 0 && i == len(parts)-1 {
			continue
		}

		parts[i] = prefix + p
	}

	return strings.Join(parts, sep)
}

// RemovePrefix of each separated string elements
func (stringsNS) RemovePrefix(s, prefix, sep string) string {
	parts := strings.Split(s, sep)
	for i, p := range parts {
		parts[i] = strings.TrimPrefix(p, prefix)
	}

	return strings.Join(parts, sep)
}

// AddSuffix to each separated string elements
func (stringsNS) AddSuffix(s, suffix, sep string) string {
	parts := strings.Split(s, sep)
	for i, p := range parts {
		if len(p) == 0 && i == len(parts)-1 {
			continue
		}

		parts[i] = p + suffix
	}

	return strings.Join(parts, sep)
}

func (stringsNS) RemoveSuffix(s, suffix, sep string) string {
	parts := strings.Split(s, sep)
	for i, p := range parts {
		parts[i] = strings.TrimSuffix(p, suffix)
	}

	return strings.Join(parts, sep)
}
