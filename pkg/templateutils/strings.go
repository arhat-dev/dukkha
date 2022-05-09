package templateutils

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/goutils"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"

	"arhat.dev/dukkha/third_party/gomplate/conv"
	gompstrings "arhat.dev/dukkha/third_party/gomplate/strings"
)

type stringsNS struct{}

func (stringsNS) Abbrev(args ...interface{}) (string, error) {
	str := ""
	offset := 0
	maxWidth := 0
	if len(args) < 2 {
		return "", errors.Errorf("abbrev requires a 'maxWidth' and 'input' argument")
	}
	if len(args) == 2 {
		maxWidth = conv.ToInt(args[0])
		str = conv.ToString(args[1])
	}
	if len(args) == 3 {
		offset = conv.ToInt(args[0])
		maxWidth = conv.ToInt(args[1])
		str = conv.ToString(args[2])
	}
	if len(str) <= maxWidth {
		return str, nil
	}
	return goutils.AbbreviateFull(str, offset, maxWidth)
}

func (stringsNS) ReplaceAll(old, new string, s interface{}) string {
	return strings.ReplaceAll(conv.ToString(s), old, new)
}

func (stringsNS) Contains(substr string, s interface{}) bool {
	return strings.Contains(conv.ToString(s), substr)
}

func (stringsNS) HasPrefix(prefix string, s interface{}) bool {
	return strings.HasPrefix(conv.ToString(s), prefix)
}

func (stringsNS) HasSuffix(suffix string, s interface{}) bool {
	return strings.HasSuffix(conv.ToString(s), suffix)
}

func (stringsNS) Repeat(count int, s interface{}) (string, error) {
	if count < 0 {
		return "", errors.Errorf("negative count %d", count)
	}
	str := conv.ToString(s)
	if count > 0 && len(str)*count/count != len(str) {
		return "", errors.Errorf("count %d too long: causes overflow", count)
	}
	return strings.Repeat(str, count), nil
}

func (stringsNS) Split(sep string, s interface{}) []string {
	return strings.Split(conv.ToString(s), sep)
}

func (stringsNS) SplitN(sep string, n int, s interface{}) []string {
	return strings.SplitN(conv.ToString(s), sep, n)
}

func (stringsNS) Trim(cutset string, s interface{}) string {
	return strings.Trim(conv.ToString(s), cutset)
}

func (stringsNS) TrimPrefix(cutset string, s interface{}) string {
	return strings.TrimPrefix(conv.ToString(s), cutset)
}

func (stringsNS) TrimSuffix(cutset string, s interface{}) string {
	return strings.TrimSuffix(conv.ToString(s), cutset)
}

func (stringsNS) Title(s interface{}) string {
	return strings.Title(conv.ToString(s))
}

func (stringsNS) ToUpper(s interface{}) string {
	return strings.ToUpper(conv.ToString(s))
}

func (stringsNS) ToLower(s interface{}) string {
	return strings.ToLower(conv.ToString(s))
}

func (stringsNS) TrimSpace(s interface{}) string {
	return strings.TrimSpace(conv.ToString(s))
}

func (stringsNS) Trunc(length int, s interface{}) string {
	return gompstrings.Trunc(length, conv.ToString(s))
}

func (stringsNS) Indent(args ...interface{}) (string, error) {
	input := conv.ToString(args[len(args)-1])
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

func (stringsNS) Slug(in interface{}) string {
	return slug.Make(conv.ToString(in))
}

func (stringsNS) Quote(in interface{}) string {
	return fmt.Sprintf("%q", conv.ToString(in))
}

func (stringsNS) ShellQuote(in interface{}) string {
	val := reflect.ValueOf(in)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		var sb strings.Builder
		max := val.Len()
		for n := 0; n < max; n++ {
			sb.WriteString(gompstrings.ShellQuote(conv.ToString(val.Index(n))))
			if n+1 != max {
				sb.WriteRune(' ')
			}
		}
		return sb.String()
	}
	return gompstrings.ShellQuote(conv.ToString(in))
}

func (stringsNS) Squote(in interface{}) string {
	s := conv.ToString(in)
	s = strings.ReplaceAll(s, `'`, `''`)
	return fmt.Sprintf("'%s'", s)
}

func (stringsNS) SnakeCase(in interface{}) (string, error) {
	return gompstrings.SnakeCase(conv.ToString(in)), nil
}

func (stringsNS) CamelCase(in interface{}) (string, error) {
	return gompstrings.CamelCase(conv.ToString(in)), nil
}

func (stringsNS) KebabCase(in interface{}) (string, error) {
	str := conv.ToString(in)
	if len(str) == 0 {
		return str, nil
	}

	return gompstrings.KebabCase(str), nil
}

func (stringsNS) WordWrap(args ...interface{}) (string, error) {
	if len(args) == 0 || len(args) > 3 {
		return "", errors.Errorf("expected 1, 2, or 3 args, got %d", len(args))
	}
	in := conv.ToString(args[len(args)-1])

	opts := gompstrings.WordWrapOpts{}
	if len(args) == 2 {
		switch a := (args[0]).(type) {
		case string:
			opts.LBSeq = a
		default:
			opts.Width = uint(conv.ToInt(a))
		}
	}
	if len(args) == 3 {
		opts.Width = uint(conv.ToInt(args[0]))
		opts.LBSeq = conv.ToString(args[1])
	}
	return gompstrings.WordWrap(in, opts), nil
}

// RuneCount - like len(s), but for runes
func (stringsNS) RuneCount(args ...interface{}) (int, error) {
	s := ""
	for _, arg := range args {
		s += conv.ToString(arg)
	}
	return utf8.RuneCountInString(s), nil
}
