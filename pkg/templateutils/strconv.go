package templateutils

import "strconv"

type strconvNS struct{}

func (strconvNS) Unquote(s String) (string, error) {
	return strconv.Unquote(toString(s))
}
