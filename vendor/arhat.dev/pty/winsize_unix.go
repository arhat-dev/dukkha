//go:build !windows
// +build !windows

package pty

import (
	"syscall"
	"unsafe"
)

// Setsize resizes t to s.
func Setsize(t FdHolder, ws Winsize) error {
	//nolint:gosec // Expected unsafe pointer for Syscall call.
	return ioctl(t.Fd(), syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(&ws)))
}

// GetsizeFull returns the full terminal size description.
func GetsizeFull(t FdHolder) (size Winsize, err error) {
	//nolint:gosec // Expected unsafe pointer for Syscall call.
	err = ioctl(t.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&size)))
	return
}
