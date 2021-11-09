package testhelper

import (
	"syscall"
)

func getSyscallFD(fd uintptr) syscall.Handle {
	return syscall.Handle(fd)
}
