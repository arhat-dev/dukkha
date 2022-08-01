package exechelper

import (
	"os"
	"os/exec"
	"unsafe"

	"arhat.dev/pkg/errhelper"
)

// StartNoLookPath is an alternative to os/exec.Cmd.Start()
// it starts cmd without looking up path with extensions on windows
func StartNoLookPath(cmd *exec.Cmd) error {
	c := (*execCmd)(unsafe.Pointer(cmd))

	if c.lookPathErr != nil {
		exec_cmd_closeDescriptors(c, c.closeAfterStart)
		exec_cmd_closeDescriptors(c, c.closeAfterWait)
		return c.lookPathErr
	}

	if c.Process != nil {
		return errhelper.ErrString("exec: already started")
	}
	if c.ctx != nil {
		select {
		case <-c.ctx.Done():
			exec_cmd_closeDescriptors(c, c.closeAfterStart)
			exec_cmd_closeDescriptors(c, c.closeAfterWait)
			return c.ctx.Err()
		default:
		}
	}

	c.childFiles = make([]*os.File, 0, 3+len(c.ExtraFiles))
	type F func(*execCmd) (*os.File, error)
	for _, setupFd := range []F{exec_cmd_stdin, exec_cmd_stdout, exec_cmd_stderr} {
		fd, err := setupFd(c)
		if err != nil {
			exec_cmd_closeDescriptors(c, c.closeAfterStart)
			exec_cmd_closeDescriptors(c, c.closeAfterWait)
			return err
		}
		c.childFiles = append(c.childFiles, fd)
	}
	c.childFiles = append(c.childFiles, c.ExtraFiles...)

	envv, err := exec_cmd_envv(c)
	if err != nil {
		return err
	}

	c.Process, err = os.StartProcess(c.Path, exec_cmd_argv(c), &os.ProcAttr{
		Dir:   c.Dir,
		Files: c.childFiles,
		Env:   exec_addCriticalEnv(exec_dedupEnv(envv)),
		Sys:   c.SysProcAttr,
	})
	if err != nil {
		exec_cmd_closeDescriptors(c, c.closeAfterStart)
		exec_cmd_closeDescriptors(c, c.closeAfterWait)
		return err
	}

	exec_cmd_closeDescriptors(c, c.closeAfterStart)

	// Don't allocate the channel unless there are goroutines to fire.
	if len(c.goroutine) > 0 {
		c.errch = make(chan error, len(c.goroutine))
		for _, fn := range c.goroutine {
			go func(fn func() error) {
				c.errch <- fn()
			}(fn)
		}
	}

	if c.ctx != nil {
		c.waitDone = make(chan struct{})
		go func() {
			select {
			case <-c.ctx.Done():
				_ = c.Process.Kill()
			case <-c.waitDone:
			}
		}()
	}

	return nil
}
