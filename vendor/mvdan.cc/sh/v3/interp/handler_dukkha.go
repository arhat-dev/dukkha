package interp

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
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

	// paths like `./foo`, `../foo`, `/foo`, `foo/bar`
	if strings.ContainsAny(target, chars) {
		return find(cwd, target, exts)
	}

	for _, dir := range splitPathList(goos, cwd, target, env) {
		var path string
		switch dir {
		case "", ".":
			// otherwise "foo" won't be "./foo"
			path = "." + string(filepath.Separator) + target
		default:
			path = filepath.Join(dir, target)
		}

		if f, err := find(cwd, path, exts); err == nil {
			return f, nil
		}
	}

	return "", fmt.Errorf("%q: executable file not found in $PATH", target)
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
		if len(e) == 0 {
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

// splitPathList normalize PATH list as absolute paths
// parameter target is the file we are looking for
//
// both ; and : are treated as path separator on windows
func splitPathList(goos, cwd, target string, env expand.Environ) []string {
	isWindows := goos == "windows"

	isSlash := pathhelper.IsUnixSlash
	if isWindows {
		isSlash = pathhelper.IsWindowsSlash
	}

	// split both colon and semi-colon on windows
	dirList := pathhelper.SplitList(env.Get("PATH").String(), true, isWindows /* semi colon sep */, isSlash)

	const BAD = "\n" // just a invalid char for path string (both windows and posix)
	var (
		// absolute path of `uname`, `cygpath`, `winepath`
		uname, cygpath, winepath string
		// output of `uname -s`
		uname_s string
	)

	for i, thisDir := range dirList {
		if !isWindows {
			// unix path

			if path.IsAbs(thisDir) {
				continue
			}

			dirList[i] = pathhelper.JoinUnixPath(cwd, thisDir)
			continue
		}

		// windows can be one of following conditions:
		// - native windows, thus windows path
		// - cygwin with unix path on windows
		// - wine with window path on darwin/linux
		// - cygwin inside wine

		if pathhelper.IsWindowsAbs(thisDir) {
			continue
		}

		const (
			unameBin    = "uname"
			cygpathBin  = "cygpath"
			winepathBin = "winepath"
		)

		// special case `uname`, `cygpath` and `winepath` to avoid infinite loop when there is
		// unix path in PATH list
		switch target {
		case cygpathBin, winepathBin, unameBin:
			dirList[i], _ = pathhelper.AbsWindowsPath(cwd, thisDir, func(path string) (string, error) {
				return "", nil
			})

			continue
		}

		// here we ignore the returned error since it's always nil
		dirList[i], _ = pathhelper.AbsWindowsPath(cwd, thisDir, func(path string) (ret string, err2 error) {
		UNAME:
			switch {
			case uname_s == "": // uname -s not set
				switch uname {
				case "":
					uname, err2 = lookPathDir(goos, cwd, unameBin, env, findExecutable)
					if err2 != nil {
						uname, uname_s = BAD, BAD
						goto WINE
					}
				case BAD: // there is no `uname`
					uname_s = BAD
					goto WINE
				}

				// uname bin exists
				uname_s, err2 = run_uname_s(cwd, uname, env)
				if err2 != nil {
					uname_s = BAD
					goto WINE
				}

				// check uname_s again
				goto UNAME
			case uname_s == BAD:
				goto WINE
			case strings.HasPrefix(uname_s, "CYGWIN_NT") ||
				strings.HasPrefix(uname_s, "MINGW32_NT") ||
				strings.HasPrefix(uname_s, "MINGW64_NT"):
				goto CYGWIN
			default:
				// TODO: this case seems impossible: detected uname -s on windows and is not cygwin
				//       merge with `case uname_s == BAD` ?
				goto WINE
			}

		CYGWIN:
			// lookup and cache cygpath for this split
			switch cygpath {
			case "":
				cygpath, err2 = lookPathDir(goos, cwd, cygpathBin, env, findExecutable)
				if err2 != nil {
					// TODO: fallback to default cygpath for some environments:
					// 	 github actions windows virtual environment:
					//		(bash) `cygpath` at C:\Program Files\Git\usr\bin\cygpath.exe
					cygpath = BAD
					goto WINE
				}
			case BAD:
				goto WINE
			}

			// there is `cygpath`
			ret, err2 = run_cygpath_w(cwd, cygpath, path, env)
			if err2 == nil {
				return
			}

			// errored, try winepath anyway

		WINE:
			// lookup and cache winepath for this split
			switch winepath {
			case "":
				winepath, err2 = lookPathDir(goos, cwd, winepathBin, env, findExecutable)
				if err2 != nil {
					winepath = BAD
					goto FALLBACK
				}
			case BAD:
				goto FALLBACK
			}

			// there is `winepath`
			ret, err2 = run_winepath_w(cwd, winepath, path, env)
			if err2 == nil {
				return
			}

		FALLBACK:

			// both cygpath and winepath are missing or not working

			switch {
			case os.Getenv("GITHUB_ACTIONS") == "true":
				// github action has msys2 installed
				return pathhelper.JoinWindowsPath(`C:\msys64`, path), nil
			default:
				return pathhelper.ConvertFSPathToWindowsPath(filepath.VolumeName(cwd), path), nil
			}
		})
	}

	return dirList
}

func run_cygpath_w(cwd, bin, target string, env expand.Environ) (string, error) {
	var buf strings.Builder

	cmd := exec.Cmd{
		Path:   bin,
		Args:   []string{bin, "-w", target},
		Env:    execEnv(env),
		Dir:    cwd,
		Stdin:  nil,
		Stdout: &buf,
		Stderr: &buf,
	}

	err := exechelper.StartNoLookPath(&cmd)
	if err == nil {
		err = cmd.Wait()
	}

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func run_winepath_w(cwd, bin, target string, env expand.Environ) (string, error) {
	var buf strings.Builder

	cmd := exec.Cmd{
		Path:   bin,
		Args:   []string{bin, "-w", target},
		Env:    execEnv(env),
		Dir:    cwd,
		Stdin:  nil,
		Stdout: &buf,
		Stderr: &buf,
	}

	err := exechelper.StartNoLookPath(&cmd)
	if err == nil {
		err = cmd.Wait()
	}

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func run_uname_s(cwd, bin string, env expand.Environ) (string, error) {
	var buf strings.Builder
	cmd := exec.Cmd{
		Path:   bin,
		Args:   []string{bin, "-s"},
		Env:    execEnv(env),
		Dir:    cwd,
		Stdin:  nil,
		Stdout: &buf,
		Stderr: &buf,
	}

	err := exechelper.StartNoLookPath(&cmd)
	if err == nil {
		err = cmd.Wait()
	}

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}
