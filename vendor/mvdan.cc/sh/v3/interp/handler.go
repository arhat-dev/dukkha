// Copyright (c) 2017, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package interp

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"mvdan.cc/sh/v3/expand"
)

// HandlerCtx returns HandlerContext value stored in ctx.
// It panics if ctx has no HandlerContext stored.
func HandlerCtx(ctx context.Context) HandlerContext {
	hc, ok := ctx.Value(handlerCtxKey{}).(HandlerContext)
	if !ok {
		panic("interp.HandlerCtx: no HandlerContext in ctx")
	}
	return hc
}

type handlerCtxKey struct{}

// HandlerContext is the data passed to all the handler functions via a context value.
// It contains some of the current state of the Runner.
type HandlerContext struct {
	// Env is a read-only version of the interpreter's environment,
	// including environment variables, global variables, and local function
	// variables.
	Env expand.Environ

	// Dir is the interpreter's current directory.
	Dir string

	// Stdin is the interpreter's current standard input reader.
	Stdin io.Reader
	// Stdout is the interpreter's current standard output writer.
	Stdout io.Writer
	// Stderr is the interpreter's current standard error writer.
	Stderr io.Writer
}

// ExecHandlerFunc is a handler which executes simple command. It is
// called for all CallExpr nodes where the first argument is neither a
// declared function nor a builtin.
//
// Returning nil error sets commands exit status to 0. Other exit statuses
// can be set with NewExitStatus. Any other error will halt an interpreter.
type ExecHandlerFunc func(ctx context.Context, args []string) error

// DefaultExecHandler returns an ExecHandlerFunc used by default.
// It finds binaries in PATH and executes them.
// When context is cancelled, interrupt signal is sent to running processes.
// KillTimeout is a duration to wait before sending kill signal.
// A negative value means that a kill signal will be sent immediately.
// On Windows, the kill signal is always sent immediately,
// because Go doesn't currently support sending Interrupt on Windows.
// Runner.New sets killTimeout to 2 seconds by default.
func DefaultExecHandler(killTimeout time.Duration) ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		hc := HandlerCtx(ctx)
		path, err := LookPathDir(hc.Dir, hc.Env, args[0])
		if err != nil {
			fmt.Fprintln(hc.Stderr, err)
			return NewExitStatus(127)
		}
		cmd := exec.Cmd{
			Path:   path,
			Args:   args,
			Env:    execEnv(hc.Env),
			Dir:    hc.Dir,
			Stdin:  hc.Stdin,
			Stdout: hc.Stdout,
			Stderr: hc.Stderr,
		}

		err = cmd.Start()
		if err == nil {
			if done := ctx.Done(); done != nil {
				go func() {
					<-done

					if killTimeout <= 0 || runtime.GOOS == "windows" {
						_ = cmd.Process.Signal(os.Kill)
						return
					}

					// TODO: don't temporarily leak this goroutine
					// if the program stops itself with the
					// interrupt.
					go func() {
						time.Sleep(killTimeout)
						_ = cmd.Process.Signal(os.Kill)
					}()
					_ = cmd.Process.Signal(os.Interrupt)
				}()
			}

			err = cmd.Wait()
		}

		switch x := err.(type) {
		case *exec.ExitError:
			// started, but errored - default to 1 if OS
			// doesn't have exit statuses
			if status, ok := x.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() {
					if ctx.Err() != nil {
						return ctx.Err()
					}
					return NewExitStatus(uint8(128 + status.Signal()))
				}
				return NewExitStatus(uint8(status.ExitStatus()))
			}
			return NewExitStatus(1)
		case *exec.Error:
			// did not start
			fmt.Fprintf(hc.Stderr, "%v\n", err)
			return NewExitStatus(127)
		default:
			return err
		}
	}
}

func checkStat(dir, file string, checkExec bool) (string, error) {
	if !filepath.IsAbs(file) {
		file = filepath.Join(dir, file)
	}
	info, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	m := info.Mode()
	if m.IsDir() {
		return "", fmt.Errorf("is a directory")
	}
	if checkExec && runtime.GOOS != "windows" && m&0o111 == 0 {
		return "", fmt.Errorf("permission denied")
	}
	return file, nil
}

func winHasExt(file string) bool {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		return false
	}
	return strings.LastIndexAny(file, `:\/`) < i
}

// findExecutable returns the path to an existing executable file.
func findExecutable(dir, file string, exts []string) (string, error) {
	if len(exts) == 0 {
		// non-windows
		return checkStat(dir, file, true)
	}
	if winHasExt(file) {
		if file, err := checkStat(dir, file, true); err == nil {
			return file, nil
		}
	}
	for _, e := range exts {
		f := file + e
		if f, err := checkStat(dir, f, true); err == nil {
			return f, nil
		}
	}
	return "", fmt.Errorf("not found")
}

// findFile returns the path to an existing file.
func findFile(dir, file string, _ []string) (string, error) {
	return checkStat(dir, file, false)
}

// LookPath is deprecated. See LookPathDir.
func LookPath(env expand.Environ, file string) (string, error) {
	return LookPathDir(env.Get("PWD").String(), env, file)
}

// LookPathDir is similar to os/exec.LookPath, with the difference that it uses the
// provided environment. env is used to fetch relevant environment variables
// such as PWD and PATH.
//
// If no error is returned, the returned path must be valid.
func LookPathDir(cwd string, env expand.Environ, file string) (string, error) {
	return lookPathDir(runtime.GOOS, cwd, file, env.Get("PATH").String(), env.Get("PATHEXT").String(), findExecutable)
}

// findAny defines a function to pass to lookPathDir.
type findAny = func(dir string, file string, exts []string) (string, error)

// scriptFromPathDir is similar to LookPathDir, with the difference that it looks
// for both executable and non-executable files.
func scriptFromPathDir(cwd string, env expand.Environ, file string) (string, error) {
	return lookPathDir(runtime.GOOS, cwd, file, env.Get("PATH").String(), env.Get("PATHEXT").String(), findFile)
}

// OpenHandlerFunc is a handler which opens files. It is
// called for all files that are opened directly by the shell, such as
// in redirects. Files opened by executed programs are not included.
//
// The path parameter may be relative to the current directory, which can be
// fetched via HandlerCtx.
//
// Use a return error of type *os.PathError to have the error printed to
// stderr and the exit status set to 1. If the error is of any other type, the
// interpreter will come to a stop.
type OpenHandlerFunc func(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error)

// DefaultOpenHandler returns an OpenHandlerFunc used by default. It uses os.OpenFile to open files.
func DefaultOpenHandler() OpenHandlerFunc {
	return func(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
		mc := HandlerCtx(ctx)
		if !filepath.IsAbs(path) {
			path = filepath.Join(mc.Dir, path)
		}
		return os.OpenFile(path, flag, perm)
	}
}
