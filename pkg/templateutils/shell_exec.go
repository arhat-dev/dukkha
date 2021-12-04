package templateutils

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

	"arhat.dev/pkg/log"
	"arhat.dev/pkg/pathhelper"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

func newExecHandler(rc dukkha.RenderingContext, stdin io.Reader) interp.ExecHandlerFunc {
	defaultCmdExecHandler := sysExecHandler(0)

	return func(
		ctx context.Context,
		args []string,
	) error {
		hc := interp.HandlerCtx(ctx)

		if !strings.HasPrefix(args[0], "tpl:") {
			err := defaultCmdExecHandler(ctx, args)
			if err != nil {
				return fmt.Errorf("%q: %w", strings.Join(args, " "), err)
			}

			return nil
		}

		var pipeReader io.Reader
		if hc.Stdin != stdin {
			// piped context
			pipeReader = hc.Stdin
		}

		return ExecCmdAsTemplateFuncCall(
			rc,
			pipeReader,
			hc.Stdout,
			append(
				[]string{strings.TrimPrefix(args[0], "tpl:")},
				args[1:]...,
			),
		)
	}
}

// sysExecHandler returns an ExecHandlerFunc used by default.
// It finds binaries in PATH and executes them.
// When context is canceled, interrupt signal is sent to running processes.
// KillTimeout is a duration to wait before sending kill signal.
// A negative value means that a kill signal will be sent immediately.
// On Windows, the kill signal is always sent immediately,
// because Go doesn't currently support sending Interrupt on Windows.
// Runner.New sets killTimeout to 2 seconds by default.
func sysExecHandler(killTimeout time.Duration) interp.ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		hc := interp.HandlerCtx(ctx)
		path, err := lookPathDir(hc.Dir, hc.Env, args[0], findExecutable)
		if err != nil {
			fmt.Fprintln(hc.Stderr, err)
			return interp.NewExitStatus(127)
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

					if killTimeout <= 0 || runtime.GOOS == constant.KERNEL_WINDOWS {
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
					return interp.NewExitStatus(uint8(128 + status.Signal()))
				}
				return interp.NewExitStatus(uint8(status.ExitStatus()))
			}
			return interp.NewExitStatus(1)
		case *exec.Error:
			// did not start
			fmt.Fprintf(hc.Stderr, "%v\n", err)
			return interp.NewExitStatus(127)
		default:
			return err
		}
	}
}

// findExecutable returns the path to an existing executable file.
func findExecutable(dir, file string, exts []string) (string, error) {
	if len(exts) == 0 {
		// non-windows
		return checkStat(dir, file, true)
	}
	if winHasExt(file) {
		if file2, err := checkStat(dir, file, true); err == nil {
			return file2, nil
		}
	}

	for _, e := range exts {
		if f, err := checkStat(dir, file+e, true); err == nil {
			return f, nil
		}
	}
	return "", fmt.Errorf("not found")
}

func winHasExt(file string) bool {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		return false
	}
	return strings.LastIndexAny(file, `:\/`) < i
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
	if checkExec && runtime.GOOS != constant.KERNEL_WINDOWS && m&0o111 == 0 {
		return "", fmt.Errorf("permission denied")
	}
	return file, nil
}

type findAny = func(dir string, file string, exts []string) (string, error)

func lookPathDir(cwd string, env expand.Environ, file string, find findAny) (string, error) {
	if find == nil {
		panic("no find function found")
	}

	chars := `/`
	if runtime.GOOS == constant.KERNEL_WINDOWS {
		chars = `:\/`
	}
	exts := pathExts(env)

	if strings.ContainsAny(file, chars) {
		return find(cwd, file, exts)
	}

	pathList := splitPathList(cwd, env.Get("PATH").String())

	log.Log.V("lookup path list",
		log.String("PATH", env.Get("PATH").String()),
		log.Strings("path_list", pathList),
		log.Strings("exts", exts),
	)

	for _, elem := range pathList {
		var path string
		switch elem {
		case "", ".":
			// otherwise "foo" won't be "./foo"
			path = "." + string(filepath.Separator) + file
		default:
			path = filepath.Join(elem, file)
		}

		if f, err := find(cwd, path, exts); err == nil {
			return f, nil
		}
	}

	return "", fmt.Errorf("%q: executable file not found in $PATH", file)
}

// splitPathList normalize PATH list as absolute paths
// both ; and : are treated as path separator
func splitPathList(cwd, path string) []string {
	isWindows := runtime.GOOS == constant.KERNEL_WINDOWS

	isSlash := pathhelper.IsUnixSlash
	if isWindows {
		isSlash = pathhelper.IsWindowsSlash
	}

	// only split semi-colon on windows
	list := pathhelper.SplitList(path, true, isWindows, isSlash)

	for i, v := range list {
		if filepath.IsAbs(v) {
			continue
		}

		if !isWindows {
			list[i] = pathhelper.JoinUnixPath(cwd, v)
			continue
		}

		var err error
		list[i], err = pathhelper.AbsWindowsPath(cwd, v, func() (string, error) {
			// find root path of the fhs root using cygpath
			output, err2 := exec.Command("cygpath", "-w", "/").CombinedOutput()
			if err2 == nil {
				return strings.TrimSpace(string(output)), nil
			}

			// github action has msys2 installed without PATH added
			// so we cannot find `cygpath` executable
			switch {
			case os.Getenv("GITHUB_ACTIONS") == "true":
				return `C:\msys64`, nil
			default:
				// TODO: other defaults?
				return "", err2
			}
		})

		// error can only happen when looking up fhs root
		_ = err
	}

	return list
}

func execEnv(env expand.Environ) []string {
	list := make([]string, 0, 64)
	env.Each(func(name string, vr expand.Variable) bool {
		if !vr.IsSet() {
			// If a variable is set globally but unset in the
			// runner, we need to ensure it's not part of the final
			// list. Seems like zeroing the element is enough.
			// This is a linear search, but this scenario should be
			// rare, and the number of variables shouldn't be large.
			for i, kv := range list {
				if strings.HasPrefix(kv, name+"=") {
					list[i] = ""
				}
			}
		}
		if vr.Exported && vr.Kind == expand.String {
			list = append(list, name+"="+vr.String())
		}
		return true
	})
	return list
}
