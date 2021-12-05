package iohelper

import (
	"io"
)

func NewDevNull() *DevNull {
	return &DevNull{
		closeSig: make(chan struct{}),
	}
}

var (
	_ io.ReadWriteCloser = (*DevNull)(nil)
	_ io.StringWriter    = (*DevNull)(nil)
)

type DevNull struct {
	closeSig chan struct{}
}

func (d *DevNull) Write(p []byte) (int, error) {
	return len(p), nil
}

func (d *DevNull) WriteString(s string) (int, error) {
	return len(s), nil
}

func (d *DevNull) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	// mock actual read from /dev/null, always block
	<-d.closeSig
	return 0, io.EOF
}

func (r *DevNull) Close() error {
	select {
	case <-r.closeSig:
		// already closed
	default:
		close(r.closeSig)
	}

	return nil
}
