package interp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"arhat.dev/pkg/exechelper"
	"arhat.dev/pkg/pathhelper"
)

// dukkha specific handler related functions

// DukkhaExecHandler is like DefaultExecHandler but won't lookup file extension
// on start
func DukkhaExecHandler(killTimeout time.Duration) ExecHandlerFunc {
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

		err = exechelper.StartNoLookPath(&cmd)
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

func lookPathDir(goos, cwd, target, pathListEnv, pathExtEnv string, find findAny) (string, error) {
	if find == nil {
		panic("no find function found")
	}

	chars := `/`
	if goos == "windows" {
		chars = `:\/`
	}
	exts := pathExts(goos, pathExtEnv)

	// paths like `./`, `../`, `/`
	if strings.ContainsAny(target, chars) {
		return find(cwd, target, exts)
	}

	pathList := splitPathList(goos, cwd, target, pathListEnv, pathExtEnv)
	for _, elem := range pathList {
		var path string
		switch elem {
		case "", ".":
			// otherwise "foo" won't be "./foo"
			path = "." + string(filepath.Separator) + target
		default:
			path = filepath.Join(elem, target)
		}

		if f, err := find(cwd, path, exts); err == nil {
			return f, nil
		}
	}

	return "", fmt.Errorf("%q: executable file not found in $PATH", target)
}

// splitPathList normalize PATH list as absolute paths
// both ; and : are treated as path separator on windows
func splitPathList(goos, cwd, target, pathListEnv, pathExtEnv string) []string {
	isWindows := goos == "windows"

	isSlash := pathhelper.IsUnixSlash
	if isWindows {
		isSlash = pathhelper.IsWindowsSlash
	}

	// only split semi-colon on windows
	list := pathhelper.SplitList(pathListEnv, true, isWindows, isSlash)

	for i, v := range list {
		if filepath.IsAbs(v) {
			continue
		}

		if !isWindows {
			list[i] = pathhelper.JoinUnixPath(cwd, v)
			continue
		}

		const (
			cygpathBin = "cygpath"
		)

		if target == cygpathBin {
			list[i], _ = pathhelper.AbsWindowsPath(cwd, v, func(path string) (string, error) {
				return "", nil
			})
			continue
		}

		var err error
		list[i], err = pathhelper.AbsWindowsPath(cwd, v, func(path string) (string, error) {
			// find root path of the fhs root using cygpath
			// but first lookup cygpath itself
			cygpath, err := lookPathDir(goos, cwd, cygpathBin, pathListEnv, pathExtEnv, findExecutable)
			if err != nil {
				return "", err
			}

			// NOTE for some environments:
			// 	 github action windows:
			//		there is `cygpath` at C:\Program Files\Git\usr\bin\cygpath
			// 		when running inside gitbash (set `shell: bash`)

			output, err2 := exec.Command(cygpath, "-w", path).CombinedOutput()
			if err2 == nil {
				return strings.TrimSpace(string(output)), nil
			}

			switch {
			case os.Getenv("GITHUB_ACTIONS") == "true":
				// github action has msys2 installed without PATH added
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

func pathExts(goos, pathExtEnv string) []string {
	if goos != "windows" {
		return nil
	}

	if len(pathExtEnv) == 0 {
		// include ""
		return []string{".com", ".exe", ".bat", ".cmd", ""}
	}

	var exts []string
	for _, e := range strings.Split(strings.ToLower(pathExtEnv), `;`) {
		if e == "" {
			continue
		}
		if e[0] != '.' {
			e = "." + e
		}
		exts = append(exts, e)
	}

	// allow no extension at last
	return append(exts, "")
}
