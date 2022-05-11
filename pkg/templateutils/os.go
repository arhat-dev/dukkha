package templateutils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"

	"arhat.dev/dukkha/pkg/dukkha"
)

func createOSNS(rc dukkha.RenderingContext) osNS { return osNS{rc: rc} }

type osNS struct{ rc dukkha.RenderingContext }

func (osNS) Stdin() *os.File  { return os.Stdin }
func (osNS) Stdout() *os.File { return os.Stdout }
func (osNS) Stderr() *os.File { return os.Stderr }

func (ns osNS) ReadFile(filename String) (string, error) {
	data, err := ns.rc.FS().ReadFile(toString(filename))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (ns osNS) WriteFile(filename String, d Bytes, args ...interface{}) error {
	perm := fs.FileMode(0640)
	if len(args) != 0 {
		if permStr := toString(args[0]); len(permStr) != 0 {
			x, err := strconv.ParseUint(permStr, 0, 8)
			if err != nil {
				return fmt.Errorf("invalid permission value: %w", err)
			}

			perm = fs.FileMode(x)
		}
	}

	return ns.rc.FS().WriteFile(toString(filename), toBytes(d), perm)
}

func (ns osNS) AppendFile(filename String, data Bytes, args ...interface{}) error {
	perm := fs.FileMode(0640)
	if len(args) != 0 {
		if permStr := toString(args[0]); len(permStr) != 0 {
			x, err := strconv.ParseUint(permStr, 0, 8)
			if err != nil {
				return fmt.Errorf("invalid permission value: %w", err)
			}

			perm = fs.FileMode(x)
		}
	}

	f, err := ns.rc.FS().OpenFile(toString(filename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}

	_, err = f.(io.Writer).Write(toBytes(data))
	return err
}

func (ns osNS) MkdirAll(path String, args ...interface{}) error {
	perm := fs.FileMode(0640)
	if len(args) != 0 {
		if permStr := toString(args[0]); len(permStr) != 0 {
			x, err := strconv.ParseUint(permStr, 0, 8)
			if err != nil {
				return fmt.Errorf("invalid permission value: %w", err)
			}

			perm = fs.FileMode(x)
		}
	}

	return ns.rc.FS().MkdirAll(toString(path), perm)
}
