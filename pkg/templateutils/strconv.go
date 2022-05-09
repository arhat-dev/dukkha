package templateutils

import "strconv"

type strconvNS struct{}

func (ns strconvNS) Unquote(s string) (string, error) {
	return strconv.Unquote(s)
}
