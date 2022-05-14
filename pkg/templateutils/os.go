package templateutils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"
	"strconv"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"

	"arhat.dev/pkg/stringhelper"
	"mvdan.cc/sh/v3/interp"
)

func createOSNS(rc dukkha.RenderingContext) osNS { return osNS{rc: rc} }

type osNS struct{ rc dukkha.RenderingContext }

func (osNS) Stdin() *os.File  { return os.Stdin }
func (osNS) Stdout() *os.File { return os.Stdout }
func (osNS) Stderr() *os.File { return os.Stderr }

func (ns osNS) UserHomeDir() (ret string) {
	env := "HOME"
	goos, _ := constant.GetGolangOS(ns.rc.HostKernel())

	switch goos {
	case "windows":
		env = "USERPROFILE"
	case "plan9":
		env = "home"
	}

	if ret = ns.rc.Get(env).String(); ret != "" {
		return
	}

	switch goos {
	case "android":
		return "/sdcard"
	case "ios":
		return "/"
	}

	return
}

func (ns osNS) UserConfigDir() string {
	goos, _ := constant.GetGolangOS(ns.rc.HostKernel())

	var dir string

	switch goos {
	case "windows":
		dir = ns.rc.Get("AppData").String()
		if dir == "" {
			return ""
		}

	case "darwin", "ios":
		dir = ns.rc.Get("HOME").String()
		if dir == "" {
			return ""
		}
		dir += "/Library/Application Support"

	case "plan9":
		dir = ns.rc.Get("home").String()
		if dir == "" {
			return ""
		}
		dir += "/lib"

	default: // Unix
		dir = ns.rc.Get("XDG_CONFIG_HOME").String()
		if dir == "" {
			dir = ns.rc.Get("HOME").String()
			if dir == "" {
				return ""
			}
			dir += "/.config"
		}
	}

	return dir
}

func (ns osNS) UserCacheDir() (ret string) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = ns.rc.Get("LocalAppData").String()
		if dir == "" {
			return ""
		}

	case "darwin", "ios":
		dir = ns.rc.Get("HOME").String()
		if dir == "" {
			return ""
		}
		dir += "/Library/Caches"

	case "plan9":
		dir = ns.rc.Get("home").String()
		if dir == "" {
			return ""
		}
		dir += "/lib/cache"

	default: // Unix
		dir = ns.rc.Get("XDG_CACHE_HOME").String()
		if dir == "" {
			dir = ns.rc.Get("HOME").String()
			if dir == "" {
				return ""
			}
			dir += "/.cache"
		}
	}

	return dir
}

// Lookup looks up executable by name in PATH list, return empty string if not found
//
// NOTE: it will try extra suffices (e.g. `.exe`) on windows
func (ns osNS) Lookup(args ...String) string {
	var (
		rc   = ns.rc
		exec string
	)
	switch len(args) {
	case 0:
		return ""
	case 1:
		// Lookup(file)
		exec = toString(args[0])
	case 2:
		// Lookup(PATH, file)
		rc = ns.rc.(dukkha.Context).DeriveNew()
		rc.AddEnv(true, &dukkha.EnvEntry{
			Name:  "PATH",
			Value: toString(args[0]),
		})

		exec = toString(args[1])
	default:
		// Lookup(PATH, PATHEXT, ... file)

		rc = ns.rc.(dukkha.Context).DeriveNew()
		rc.AddEnv(true, &dukkha.EnvEntry{
			Name:  "PATH",
			Value: toString(args[0]),
		}, &dukkha.EnvEntry{
			Name:  "PATHEXT",
			Value: toString(args[1]),
		})

		exec = toString(args[len(args)-1])
	}

	goos, ok := constant.GetGolangOS(rc.HostKernel())
	if !ok {
		goos = runtime.GOOS
	}

	ret, err := interp.DukkhaLookPathDir(goos, rc.WorkDir(), exec, rc, interp.DukkhaFindExecutable)
	if err != nil {
		return ""
	}

	return ret
}

// LookupFile looks up file by name in PATH list, return empty string if not found
//
// it's like Lookup but doesn't require file to be executable
//
// NOTE: it will not try extra suffices (e.g. `.com`) on windows
func (ns osNS) LookupFile(args ...String) string {
	var (
		rc   = ns.rc
		file string
	)
	switch len(args) {
	case 0:
		return ""
	case 1:
		// LookupFile(file)
		file = toString(args[0])
	default:
		// LookupFile(PATH, ... file)

		rc = ns.rc.(dukkha.Context).DeriveNew()
		rc.AddEnv(true, &dukkha.EnvEntry{
			Name:  "PATH",
			Value: toString(args[0]),
		})

		file = toString(args[len(args)-1])
	}

	goos, ok := constant.GetGolangOS(rc.HostKernel())
	if !ok {
		goos = runtime.GOOS
	}

	ret, err := interp.DukkhaLookPathDir(goos, rc.WorkDir(), file, rc, interp.DukkhaFindFile)
	if err != nil {
		return ""
	}

	return ret
}

func (ns osNS) ReadFile(file String) (string, error) {
	data, err := ns.rc.FS().ReadFile(toString(file))
	if err != nil {
		return "", err
	}

	return stringhelper.Convert[string, byte](data), nil
}

func (ns osNS) WriteFile(file String, d Bytes, args ...interface{}) error {
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

	return ns.rc.FS().WriteFile(toString(file), toBytes(d), perm)
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
	defer f.Close()

	_, err = f.(io.Writer).Write(toBytes(data))
	return err
}

func (ns osNS) MkdirAll(path String, args ...interface{}) error {
	perm := fs.FileMode(0755)
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
