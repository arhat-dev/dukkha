package templateutils

import (
	"fmt"
	"regexp"

	"arhat.dev/pkg/clihelper"
	"arhat.dev/pkg/regexphelper"
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
type regexpNS struct{}

func (regexpNS) Find(args ...String) (ret string, err error) {
	err = handleReTemplateFunc_DATA(args, func(re *regexp.Regexp, data string) error {
		ret = re.FindString(data)
		return nil
	})

	return
}

func (regexpNS) FindAll(args ...String) (ret []string, err error) {
	err = handleRETemplateFunc_OptionalN_DATA(args, func(re *regexp.Regexp, data string, c int) error {
		ret = re.FindAllString(data, c)
		return nil
	})

	return
}

func (regexpNS) Match(args ...String) (matched bool, err error) {
	err = handleReTemplateFunc_DATA(args, func(re *regexp.Regexp, data string) error {
		matched = re.MatchString(data)
		return nil
	})

	return
}

// TODO: support writer as the second last argument
func (regexpNS) Replace(args ...String) (ret string, err error) {
	err = handleReTemplateFunc_OPDATA_DATA(args, func(re *regexp.Regexp, opData, data string) error {
		ret = re.ReplaceAllString(data, opData)
		return nil
	})

	return
}

// TODO: support writer as the second last argument
func (regexpNS) ReplaceLiteral(args ...String) (ret string, err error) {
	err = handleReTemplateFunc_OPDATA_DATA(args, func(re *regexp.Regexp, opData, data string) error {
		ret = re.ReplaceAllLiteralString(data, opData)
		return nil
	})

	return
}

func (regexpNS) Split(args ...String) (ret []string, err error) {
	err = handleRETemplateFunc_OptionalN_DATA(args, func(re *regexp.Regexp, data string, c int) error {
		ret = re.Split(data, c)
		return nil
	})

	return
}

func (regexpNS) QuoteMeta(expr String) string { return regexp.QuoteMeta(must(toString(expr))) }

// last arg is the input data
func handleReTemplateFunc_DATA(args []String, doRE func(re *regexp.Regexp, data string) error) (err error) {
	n := len(args)
	if n < 2 {
		return fmt.Errorf("at least 2 args expected: got %d", n)
	}

	re, err := parseRegexp(args[:n-1])
	if err != nil {
		return err
	}

	data, err := toString(args[n-1])
	if err != nil {
		return
	}

	return doRE(re, data)
}

// last two arga are: opration action data (e.g. replacement text) and data
func handleReTemplateFunc_OPDATA_DATA(
	args []String,
	doRE func(re *regexp.Regexp, opData, data string) error,
) (err error) {
	n := len(args)
	if n < 3 {
		return fmt.Errorf("at least 3 args expected: got %d", n)
	}

	re, err := parseRegexp(args[:n-2])
	if err != nil {
		return err
	}

	opData, err := toString(args[n-2])
	if err != nil {
		return
	}

	data, err := toString(args[n-1])
	if err != nil {
		return
	}

	return doRE(re, opData, data)
}

// last two args are: an optional number, input data
func handleRETemplateFunc_OptionalN_DATA(
	args []String,
	doRE func(re *regexp.Regexp, data string, c int) error,
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

	re, err := parseRegexp(args[:n-1])
	if err != nil {
		return err
	}

	data, err := toString(args[n-1])
	if err != nil {
		return
	}

	return doRE(re, data, c)
}

// parseRegexp args[0] is expected to be expr
// args[1:] is a list of regexp flags
func parseRegexp(args []String) (_ *regexp.Regexp, err error) {
	n := len(args)

	var expr string
	if len(args) == 0 {
		err = errAtLeastOneArgGotZero
		return
	}

	expr, err = toString(args[0])
	if err != nil {
		return
	}

	if n == 1 {
		return regexp.Compile(expr)
	}

	var (
		fs pflag.FlagSet

		flags []string
		opts  regexphelper.Options
	)

	clihelper.InitFlagSet(&fs, "regexp")

	fs.BoolVarP(&opts.CaseInsensitive, "case-insensitive", "i", false, "")
	fs.BoolVarP(&opts.Multiline, "multi-line", "m", false, "")
	fs.BoolVarP(&opts.DotNewLine, "dot-newline", "s", false, "")
	fs.BoolVarP(&opts.Ungreedy, "ungreedy", "U", false, "")

	flags, err = toStrings(args[1:])
	if err != nil {
		return
	}

	err = fs.Parse(flags)
	if err != nil {
		return nil, err
	}

	return regexp.Compile(opts.Wrap(expr))
}
