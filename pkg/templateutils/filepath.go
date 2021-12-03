package templateutils

import (
	"path/filepath"
	"runtime"
	"strings"

	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/dukkha/third_party/gomplate/conv"
	"arhat.dev/pkg/pathhelper"
	"arhat.dev/pkg/wellknownerrors"
)

type filepathNS struct {
	rc dukkha.RenderingContext
}

func createFilePathNS(rc dukkha.RenderingContext) *filepathNS {
	return &filepathNS{
		rc: rc,
	}
}

func (f *filepathNS) Base(in interface{}) string {
	return filepath.Base(conv.ToString(in))
}

func (f *filepathNS) Clean(in interface{}) string {
	return filepath.Clean(conv.ToString(in))
}

func (f *filepathNS) Dir(in interface{}) string {
	return filepath.Dir(conv.ToString(in))
}

func (f *filepathNS) Ext(in interface{}) string {
	return filepath.Ext(conv.ToString(in))
}

func (f *filepathNS) FromSlash(in interface{}) string {
	return filepath.FromSlash(conv.ToString(in))
}

func (f *filepathNS) IsAbs(in interface{}) bool {
	return filepath.IsAbs(conv.ToString(in))
}

func (f *filepathNS) Join(elem ...interface{}) string {
	s := conv.ToStrings(elem...)
	return filepath.Join(s...)
}

func (f *filepathNS) Match(pattern, name interface{}) (matched bool, err error) {
	return filepath.Match(conv.ToString(pattern), conv.ToString(name))
}

func (f *filepathNS) Rel(basepath, targpath interface{}) (string, error) {
	return filepath.Rel(conv.ToString(basepath), conv.ToString(targpath))
}

func (f *filepathNS) Split(in interface{}) []string {
	dir, file := filepath.Split(conv.ToString(in))
	return []string{dir, file}
}

func (f *filepathNS) ToSlash(in interface{}) string {
	return filepath.ToSlash(conv.ToString(in))
}

func (f *filepathNS) VolumeName(in interface{}) string {
	return filepath.VolumeName(conv.ToString(in))
}

func (f *filepathNS) Glob(in interface{}) ([]string, error) {
	return f.rc.FS().Glob(conv.ToString(in))
}

func (f *filepathNS) Abs(in interface{}) (string, error) {
	path := conv.ToString(in)
	if filepath.IsAbs(path) {
		return path, nil
	}

	if runtime.GOOS != "windows" {
		return filepath.Join(f.rc.WorkDir(), path), nil
	}

	return pathhelper.AbsWindowsPath(f.rc.WorkDir(), path, func() (string, error) {
		// TODO: get fhs root
		return "", wellknownerrors.ErrNotSupported
	})
}

func absPath(workdir, path string) string {
	if len(path) == 0 {
		return workdir
	}

	if runtime.GOOS != "windows" {
		if filepath.IsAbs(path) {
			return path
		}

		return filepath.Join(workdir, path)
	}

	// windows

	if filepath.IsAbs(path) {
		// e.g. c:/some-path or \\nethost\share
		return path
	}

	// non absolute path

	if path[0] != '/' {
		// non potential cygpath
		return filepath.Join(workdir, path)
	}

	// cygpath or driver relative path

	switch {
	case strings.HasPrefix(path, "/cygdrive/"):
		// provide empty default volume name since /cygdrive/ is followed
		// by the volume name
		return pathhelper.ConvertFSPathToWindowsPath("", strings.TrimPrefix(path, "/cygdrive"))
	default:
		// filter common prefixes used in hfs
		// ref: https://en.wikipedia.org/wiki/Filesystem_Hierarchy_Standard

		isHFS := false
		for _, posixPrefix := range []string{
			"/bin", "/boot", "/dev", "/etc", "/home",
			"/lib", "/lib64", "/lib32", "/media", "/mnt",
			"/opt", "/proc", "/root", "/run", "/sbin",
			"/srv", "/sys", "/tmp", "/usr", "/var",
		} {
			if strings.HasPrefix(path, posixPrefix) {
				isHFS = true
				break
			}
		}

		if isHFS {
			// TODO: lookup root path of cygwin/mingw64/msys2
			return path
		}

		return pathhelper.ConvertFSPathToWindowsPath(filepath.VolumeName(workdir), path)
	}
}
