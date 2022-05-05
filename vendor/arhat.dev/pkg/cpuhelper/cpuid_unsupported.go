//go:build !amd64 && !386 && !amd64p32

package cpuhelper

func detect() CPU { return nil }
