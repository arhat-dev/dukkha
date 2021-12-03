package pathhelper

import (
	"strings"
)

// IsWindowsAbs reports whether the path is absolute.
func IsWindowsAbs(path string) (b bool) {
	if IsReservedWindowsName(path) {
		return true
	}
	l := volumeNameLen(path)
	if l == 0 {
		return false
	}
	path = path[l:]
	if path == "" {
		return false
	}
	return IsWindowsSlash(path[0])
}

// AbsWindowsPath returns absolute path of path with custom cwd
// the argument path SHOULD be relative path, that is:
// 	on windows: not start with `[a-zA-Z]:/` or `\\`
// 	on unix: not start with `/`
//
// It tries to handle three different styles all at once:
// 	- windows style (`\\`, `[a-zA-Z]:/`)
// 	- cygpath style absolute path (`/cygdrive/c`)
// 	- golang io/fs style absolute path for windows (`/[a-zA-Z]/`, e.g. /c/foo)
func AbsWindowsPath(
	cwd, path string,
	getFHSRoot func() (string, error), // root of dirs like /usr, /root
) (string, error) {
	if len(path) == 0 {
		return cwd, nil
	}

	// non absolute path

	if path[0] != '/' {
		// non potential relative path for windows
		return JoinWindowsPath(cwd, path), nil
	}

	// cygpath or driver relative path

	if strings.HasPrefix(path, "/cygdrive/") {
		// provide empty default volume name since /cygdrive/ MUST be followed
		// by the driver desigantor
		return ConvertFSPathToWindowsPath("", strings.TrimPrefix(path, "/cygdrive")), nil
	}

	// filter common prefixes used in FHS
	// ref: https://en.wikipedia.org/wiki/Filesystem_Hierarchy_Standard

	isFHS := false
	for _, posixPrefix := range []string{
		"/bin", "/boot", "/dev", "/etc", "/home",
		"/lib", "/lib64", "/lib32", "/media", "/mnt",
		"/opt", "/proc", "/root", "/run", "/sbin",
		"/srv", "/sys", "/tmp", "/usr", "/var",
	} {
		if strings.HasPrefix(path, posixPrefix) {
			isFHS = true
			break
		}
	}

	if isFHS {
		// TODO: lookup root path of cygwin/mingw64/msys2
		root, err := getFHSRoot()
		if err != nil {
			return "", err
		}

		return JoinWindowsPath(root, path[1:]), nil
	}

	return ConvertFSPathToWindowsPath(cwd[:volumeNameLen(cwd)], path), nil
}
