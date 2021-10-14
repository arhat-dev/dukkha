package templateutils

import "strconv"

var (
	strconvNS = &_strconvNS{}
)

type _strconvNS struct{}

func (ns *_strconvNS) Unquote(s string) (string, error) {
	return strconv.Unquote(s)
}
