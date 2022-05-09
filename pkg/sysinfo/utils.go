package sysinfo

import (
	"arhat.dev/pkg/cpuhelper"
)

func Arch() string {
	return string(cpuhelper.Arch(nil))
}
