package input

import (
	"io"
	"runtime"
)

// copied from golang.org/x/term@v0.0.0-20210615171337-6886f2dfbf5b/terminal.go

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// readline reads from reader until it finds \n or io.EOF.
// The slice returned does not include the \n.
// readline also ignores any \r it finds.
// Windows uses \r as end of line. So, on Windows, readline
// reads until it finds \r and ignores any \n it finds during processing.
func readline(reader io.Reader) ([]byte, error) {
	var buf [1]byte
	var ret []byte

	for {
		n, err := reader.Read(buf[:])
		if n > 0 {
			switch buf[0] {
			case '\b':
				if len(ret) > 0 {
					ret = ret[:len(ret)-1]
				}
			case '\n':
				if runtime.GOOS != "windows" {
					return ret, nil
				}
				// otherwise ignore \n
			case '\r':
				if runtime.GOOS == "windows" {
					return ret, nil
				}
				// otherwise ignore \r
			default:
				ret = append(ret, buf[0])
			}
			continue
		}
		if err != nil {
			if err == io.EOF && len(ret) > 0 {
				return ret, nil
			}
			return ret, err
		}
	}
}
