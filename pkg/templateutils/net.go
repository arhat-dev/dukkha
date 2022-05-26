package templateutils

import (
	"fmt"

	"github.com/hashicorp/go-sockaddr"
)

type sockaddrNS struct{}

func (sockaddrNS) AllInterfaces() (sockaddr.IfAddrs, error) {
	return sockaddr.GetAllInterfaces()
}

func (sockaddrNS) DefaultInterfaces() (sockaddr.IfAddrs, error) {
	return sockaddr.GetDefaultInterfaces()
}

func (sockaddrNS) PrivateInterfaces() (sockaddr.IfAddrs, error) {
	return sockaddr.GetPrivateInterfaces()
}

func (sockaddrNS) PublicInterfaces() (sockaddr.IfAddrs, error) {
	return sockaddr.GetPublicInterfaces()
}

func (sockaddrNS) Sort(selectorParam String, inputIfAddrs sockaddr.IfAddrs) (sockaddr.IfAddrs, error) {
	return sockaddr.SortIfBy(must(toString(selectorParam)), inputIfAddrs)
}

func (sockaddrNS) Exclude(selectorName, selectorParam String, inputIfAddrs sockaddr.IfAddrs) (sockaddr.IfAddrs, error) {
	return sockaddr.ExcludeIfs(must(toString(selectorName)), must(toString(selectorParam)), inputIfAddrs)
}

func (sockaddrNS) Include(selectorName, selectorParam String, inputIfAddrs sockaddr.IfAddrs) (sockaddr.IfAddrs, error) {
	return sockaddr.IncludeIfs(must(toString(selectorName)), must(toString(selectorParam)), inputIfAddrs)
}

func (sockaddrNS) Attr(selectorName String, ifAddrsRaw any) (string, error) {
	sname := must(toString(selectorName))
	switch v := ifAddrsRaw.(type) {
	case sockaddr.IfAddr:
		return sockaddr.IfAttr(sname, v)
	case sockaddr.IfAddrs:
		return sockaddr.IfAttrs(sname, v)
	default:
		return "", fmt.Errorf("unable to obtain attribute %s from type %T (%v)", sname, ifAddrsRaw, ifAddrsRaw)
	}
}

func (sockaddrNS) Join(selectorName, joinString String, inputIfAddrs sockaddr.IfAddrs) (string, error) {
	return sockaddr.JoinIfAddrs(must(toString(selectorName)), must(toString(joinString)), inputIfAddrs)
}

func (sockaddrNS) Limit(lim Number, in sockaddr.IfAddrs) (sockaddr.IfAddrs, error) {
	return sockaddr.LimitIfAddrs(toIntegerOrPanic[uint](lim), in)
}

func (sockaddrNS) Offset(off Number, in sockaddr.IfAddrs) (sockaddr.IfAddrs, error) {
	return sockaddr.OffsetIfAddrs(toIntegerOrPanic[int](off), in)
}

func (sockaddrNS) Unique(selectorName String, inputIfAddrs sockaddr.IfAddrs) (sockaddr.IfAddrs, error) {
	return sockaddr.UniqueIfAddrsBy(must(toString(selectorName)), inputIfAddrs)
}

func (sockaddrNS) Math(op, value String, inputIfAddrs sockaddr.IfAddrs) (sockaddr.IfAddrs, error) {
	return sockaddr.IfAddrsMath(must(toString(op)), must(toString(value)), inputIfAddrs)
}

func (sockaddrNS) PrivateIP() (string, error) { return sockaddr.GetPrivateIPs() }
func (sockaddrNS) PublicIP() (string, error)  { return sockaddr.GetPublicIPs() }
func (sockaddrNS) InterfaceIP(namedIfRE String) (string, error) {
	return sockaddr.GetInterfaceIPs(must(toString(namedIfRE)))
}
