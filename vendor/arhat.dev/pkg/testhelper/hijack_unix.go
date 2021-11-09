//go:build !windows
// +build !windows

package testhelper

func getSyscallFD(fd uintptr) int {
	return int(fd)
}
