package templateutils

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"arhat.dev/pkg/clihelper"
	"arhat.dev/pkg/regexphelper"
	"arhat.dev/pkg/stringhelper"
	"github.com/spf13/pflag"
)

// regexpNS for regular expression
//
// all functions' arguments starts with a expr, followed by an optional list of options for
// that expr
// XXX(expr String, ...<options for expr>, ...<function specific options>, data String)
//
// where `options for regexpr` are:
// - `--case-insensitive` or `-i`: run in case insensitive mode
// - `--multi-line` or `-m`: run in multi-line mode
// - `--dot-newline` or `-s`: match newline with dots(`.`)
// - `--ungreedy` or `-U`: reverse meaning of `.*?` and `.*`
// - `--in-place`: operate on input directly, do not copy
//    only applicable when the size of result is less than or equal to the input data
type regexpNS struct{}

// FindFirst returns first matched string
//
// FindFirst(expr string, args ...string, data string)
func (regexpNS) FindFirst(args ...String) (ret string, err error) {
	err = handleRE_DATA(args, func(re *regexp.Regexp, data []byte) error {
		ret = re.FindString(stringhelper.Convert[string, byte](data))
		return nil
	})

	return
}

// FindN returns n matched (non-overlapped) strings
//
// FindN(expr string, args ...string, data string): this is FindAll
// FindN(expr string, args ...string, count int, data string): when count < 0, it's FindAll
func (regexpNS) FindN(args ...String) (ret []string, err error) {
	err = handleRE_OptionalN_DATA(args, func(re *regexp.Regexp, data []byte, c int) error {
		ret = re.FindAllString(stringhelper.Convert[string, byte](data), c)
		return nil
	})

	return
}

// Find returns all matched (non-overlapped) string
//
// FindAll(expr string, args ...string, data string)
func (regexpNS) FindAll(args ...String) (ret []string, err error) {
	err = handleRE_OptionalN_DATA(args, func(re *regexp.Regexp, data []byte, c int) error {
		ret = re.FindAllString(stringhelper.Convert[string, byte](data), -1)
		return nil
	})
	return
}

// Match returns true when data matched expr
// Match(expr string, args ...string, data string)
func (regexpNS) Match(args ...String) (matched bool, err error) {
	err = handleRE_DATA(args, func(re *regexp.Regexp, data []byte) error {
		matched = re.Match(data)
		return nil
	})

	return
}

// Replace first matched string, it will expand capture group references (e.g. `$1`) in expr
//
// Replace(expr string, args ...string, data string)
// TODO: support writer as the second last argument
func (regexpNS) ReplaceFirst(args ...String) (ret string, err error) {
	err = handleRE_OPDATA_DATA(args, func(re *regexp.Regexp, inplace bool, opData, data []byte) error {
		ret0 := RegexpReplaceExpandN(re, data, opData, 1)
		ret = stringhelper.Convert[string, byte](ret0)
		return nil
	})

	return
}

// ReplaceAll replace all matched string (non-overlapped)
// TODO: support writer as the second last argument
func (regexpNS) ReplaceAll(args ...String) (ret string, err error) {
	err = handleRE_OPDATA_DATA(args, func(re *regexp.Regexp, inplace bool, repl, src []byte) error {
		ret0 := RegexpReplaceExpandN(re, src, repl, -1)
		ret = stringhelper.Convert[string, byte](ret0)
		return nil
	})

	return
}

// Replace first matched string, it will expand capture group references (e.g. `$1`) in expr
//
// Replace(expr string, args ...string, data string)
// TODO: support writer as the second last argument
func (regexpNS) ReplaceFirstNoExpand(args ...String) (ret string, err error) {
	err = handleRE_OPDATA_DATA(args, func(re *regexp.Regexp, inplace bool, opData, data []byte) error {
		ret0 := RegexpReplaceNoExpandN(re, data, opData, 1)
		ret = stringhelper.Convert[string, byte](ret0)
		return nil
	})

	return
}

// ReplaceAll replace all matched string (non-overlapped)
// TODO: support writer as the second last argument
func (regexpNS) ReplaceAllNoExpand(args ...String) (ret string, err error) {
	err = handleRE_OPDATA_DATA(args, func(re *regexp.Regexp, inplace bool, opData, data []byte) error {
		ret0 := RegexpReplaceNoExpandN(re, data, opData, -1)
		ret = stringhelper.Convert[string, byte](ret0)
		return nil
	})

	return
}

func RegexpReplaceExpandN(re *regexp.Regexp, src, repl []byte, n int) []byte {
	template := stringhelper.Convert[string, byte](repl)
	data := stringhelper.Convert[string, byte](src)
	return RegexpReplaceN(re, src, n, func(dest []byte, match []int) []byte {
		// here we use re.ExpandString because it doesn't do type conversion
		// and re.Expand converts []byte to string
		return re.ExpandString(dest, template, data, match)
	})
}

func RegexpReplaceNoExpandN(re *regexp.Regexp, src, repl []byte, n int) []byte {
	return RegexpReplaceN(re, src, n, func(dest []byte, match []int) []byte {
		return append(dest, repl...)
	})
}

func RegexpReplaceN(re *regexp.Regexp, src []byte, n int, repl func(dest []byte, match []int) []byte) (buf []byte) {
	lastMatchEnd := 0 // end position of the most recent match
	offset := 0       // position where we next look for a match
	end := len(src)
	if n < 0 {
		n = end + 1
	}

	for i := 0; offset <= end && i < n; {
		a := re.FindIndex(src[offset:])
		if len(a) == 0 {
			break // no more matches
		}

		a[0], a[1] = a[0]+offset, a[1]+offset

		// Copy the unmatched characters before this match.
		buf = append(buf, src[lastMatchEnd:a[0]]...)

		// Now insert a copy of the replacement string, but not for a
		// match of the empty string immediately after another match.
		// (Otherwise, we get double replacement for patterns that
		// match both empty and nonempty strings.)
		if a[1] > lastMatchEnd || a[0] == 0 {
			buf = repl(buf, a)
			i++
		}
		lastMatchEnd = a[1]

		// Advance past this match; always advance at least one character.
		var width int
		_, width = utf8.DecodeRune(src[offset:])
		switch {
		case offset+width > a[1]:
			offset += width
		case offset+1 > a[1]:
			// This clause is only needed at the end of the input
			// string. In that case, DecodeRuneInString returns width=0.
			offset++
		default:
			offset = a[1]
		}
	}

	// Copy the unmatched characters after the last match.
	buf = append(buf, src[lastMatchEnd:]...)

	return buf
}

func (regexpNS) Split(args ...String) (ret []string, err error) {
	err = handleRE_OptionalN_DATA(args, func(re *regexp.Regexp, data []byte, c int) error {
		ret = re.Split(stringhelper.Convert[string, byte](data), c)
		return nil
	})

	return
}

func (regexpNS) QuoteMeta(expr String) string { return regexp.QuoteMeta(must(toString(expr))) }

// last arg is the input data
func handleRE_DATA(
	args []String,
	doRE func(re *regexp.Regexp, data []byte) error,
) (err error) {
	n := len(args)
	if n < 2 {
		return fmt.Errorf("at least 2 args expected: got %d", n)
	}

	re, _, err := parseRegexp(args[:n-1])
	if err != nil {
		return err
	}

	data, err := toBytes(args[n-1])
	if err != nil {
		return
	}

	return doRE(re, data)
}

// last two args are: opration action data (e.g. replacement text) and data
func handleRE_OPDATA_DATA(
	args []String,
	doRE func(re *regexp.Regexp, inplace bool, opData, data []byte) error,
) (err error) {
	n := len(args)
	if n < 3 {
		return fmt.Errorf("at least 3 args expected: got %d", n)
	}

	re, inplace, err := parseRegexp(args[:n-2])
	if err != nil {
		return err
	}

	opData, err := toBytes(args[n-2])
	if err != nil {
		return
	}

	data, err := toBytes(args[n-1])
	if err != nil {
		return
	}

	return doRE(re, inplace, opData, data)
}

// last two args are: an optional number, input data
func handleRE_OptionalN_DATA(
	args []String,
	doRE func(re *regexp.Regexp, data []byte, c int) error,
) (err error) {
	n := len(args)
	if n < 2 {
		return fmt.Errorf("at least 2 args expected: got %d", n)
	}

	c := -1
	if n > 2 {
		var (
			iv      uint64
			isFloat bool
		)
		iv, isFloat, err = parseNumber(args[n-2])
		if !isFloat && err == nil {
			c = int(iv)
		}
	}

	re, _, err := parseRegexp(args[:n-1])
	if err != nil {
		return err
	}

	data, err := toBytes(args[n-1])
	if err != nil {
		return
	}

	return doRE(re, data, c)
}

// parseRegexp args[0] is expected to be the regular expression
// args[1:] is a list of regexp flags
func parseRegexp(args []String) (ret *regexp.Regexp, inplace bool, err error) {
	var expr string
	switch len(args) {
	case 0:
		err = errAtLeastOneArgGotZero
		return
	case 1:
		expr, err = toString(args[0])
		if err != nil {
			return
		}

		ret, err = regexp.Compile(expr)
		return
	default:
		expr, err = toString(args[0])
		if err != nil {
			return
		}

	}

	var (
		pfs pflag.FlagSet

		flags []string
		opts  regexphelper.Options
	)

	clihelper.InitFlagSet(&pfs, "regexp")

	pfs.BoolVarP(&opts.CaseInsensitive, "case-insensitive", "i", false, "")
	pfs.BoolVarP(&opts.Multiline, "multi-line", "m", false, "")
	pfs.BoolVarP(&opts.DotNewLine, "dot-newline", "s", false, "")
	pfs.BoolVarP(&opts.Ungreedy, "ungreedy", "U", false, "")
	pfs.BoolVar(&inplace, "in-place", false, "")

	flags, err = toStrings(args[1:])
	if err != nil {
		return
	}

	err = pfs.Parse(flags)
	if err != nil {
		return
	}

	ret, err = regexp.Compile(opts.Wrap(expr))
	return
}
