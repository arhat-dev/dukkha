/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package exechelper

import (
	"io"
	"os/exec"

	"arhat.dev/pty"
)

func startCmdWithTty(cmd *exec.Cmd) (
	doResize resizeFunc,
	close func(),
	stdin io.WriteCloser,
	stdout io.ReadCloser,
	err error,
) {
	var f pty.Pty
	f, err = pty.Start(cmd)
	if err != nil {
		return
	}

	doResize = func(cols, rows uint16) error {
		return pty.Setsize(f, pty.Winsize{
			Cols: uint16(cols), Rows: uint16(rows),
		})
	}

	close = func() { _ = f.Close() }

	switch t := f.(type) {
	case *pty.WindowsPty:
		stdin = t.InputPipe()
		stdout = t.OutputPipe()
	default:
		// unreachable, defensive
		stdin = f
		stdout = f
	}

	return
}
