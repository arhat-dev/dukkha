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
	"mvdan.cc/sh/v3/expand"
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

func lookPathDir(goos, cwd, target string, env expand.Environ, find findAny) (string, error) {
	if find == nil {
		panic("no find function found")
	}

	chars := `/`
	if goos == "windows" {
		chars = `:\/`
	}
	exts := pathExts(goos, env.Get("PATHEXT").String())

	// paths like `./`, `../`, `/`
	if strings.ContainsAny(target, chars) {
		return find(cwd, target, exts)
	}

	pathList := splitPathList(goos, cwd, target, env)
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
func splitPathList(goos, cwd, target string, env expand.Environ) []string {
	isWindows := goos == "windows"

	isSlash := pathhelper.IsUnixSlash
	if isWindows {
		isSlash = pathhelper.IsWindowsSlash
	}

	// split both colon and semi-colon on windows
	list := pathhelper.SplitList(env.Get("PATH").String(), true, isWindows, isSlash)

	for i, v := range list {
		if filepath.IsAbs(v) {
			continue
		}

		if !isWindows {
			list[i] = pathhelper.JoinUnixPath(cwd, v)
			continue
		}

		const (
			cygpathBin  = "cygpath"
			winepathBin = "winepath"
		)

		if target == cygpathBin {
			list[i], _ = pathhelper.AbsWindowsPath(cwd, v, func(path string) (string, error) {
				return "", nil
			})
			continue
		}

		// here we ignore the returned error since it's always nil
		list[i], _ = pathhelper.AbsWindowsPath(cwd, v, func(path string) (string, error) {
			// find root path of the fhs root using cygpath
			// but first lookup cygpath itself
			cygpath, err2 := lookPathDir(goos, cwd, cygpathBin, env, findExecutable)
			if err2 == nil {
				// NOTE for some environments:
				// 	 github action windows:
				//		`cygpath` at C:\Program Files\Git\usr\bin\cygpath

				var buf strings.Builder
				cmd := exec.Cmd{
					Path:   cygpath,
					Args:   []string{"-w", path},
					Env:    execEnv(env),
					Dir:    cwd,
					Stdin:  nil,
					Stdout: &buf,
					Stderr: &buf,
				}

				err2 = exechelper.StartNoLookPath(&cmd)
				if err2 == nil {
					err2 = cmd.Wait()
				}

				if err2 == nil {
					return strings.TrimSpace(buf.String()), nil
				}
			}

			// cygpath missing not working

			// try winepath
			winepath, err2 := lookPathDir(goos, cwd, winepathBin, env, findExecutable)
			if err2 == nil {
				var buf strings.Builder
				cmd := exec.Cmd{
					Path:   winepath,
					Args:   []string{"-w", path},
					Env:    execEnv(env),
					Dir:    cwd,
					Stdin:  nil,
					Stdout: &buf,
					Stderr: &buf,
				}

				err2 = exechelper.StartNoLookPath(&cmd)
				if err2 == nil {
					err2 = cmd.Wait()
				}

				if err2 == nil {
					return strings.TrimSpace(buf.String()), nil
				}
			}

			// winepath missing or not working

			switch {
			case os.Getenv("GITHUB_ACTIONS") == "true":
				// github action has msys2 installed
				return pathhelper.JoinWindowsPath(`C:\msys64`, path), nil
			default:
				return pathhelper.ConvertFSPathToWindowsPath(filepath.VolumeName(cwd), path), nil
			}
		})
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
