package templateutils

import (
	"io"
	"runtime"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

func createOSNS(rc dukkha.RenderingContext) osNS { return osNS{rc: rc} }

type osNS struct{ rc dukkha.RenderingContext }

func (ns osNS) Stdin() io.Reader  { return ns.rc.Stdin() }
func (ns osNS) Stdout() io.Writer { return ns.rc.Stdout() }
func (ns osNS) Stderr() io.Writer { return ns.rc.Stderr() }

func (ns fsNS) UserHomeDir() (ret string) {
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

func (ns fsNS) UserConfigDir() string {
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

func (ns fsNS) UserCacheDir() (ret string) {
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
