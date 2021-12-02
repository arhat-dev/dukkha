package fshelper

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"arhat.dev/pkg/pathhelper"
	"arhat.dev/pkg/wellknownerrors"
	"github.com/bmatcuk/doublestar/v4"
)

var (
	_ fs.FS         = (*OSFS)(nil)
	_ fs.ReadDirFS  = (*OSFS)(nil)
	_ fs.ReadFileFS = (*OSFS)(nil)
	_ fs.GlobFS     = (*OSFS)(nil)
	_ fs.StatFS     = (*OSFS)(nil)
	_ fs.SubFS      = (*OSFS)(nil)
)

// NewOSFS creates a new filesystem abstraction for real filesystem
func NewOSFS(
	strictIOFS bool,
	getCwd func() (string, error),
) *OSFS {
	if getCwd == nil {
		getCwd = os.Getwd
	}

	return &OSFS{
		strict: strictIOFS,
		getCwd: getCwd,
	}
}

// OSFS is a context aware filesystem abstration for afero.FS and io/fs.FS
type OSFS struct {
	strict bool
	getCwd func() (string, error)
}

func (ofs *OSFS) SetStrict(s bool) {
	ofs.strict = s
}

// getRealPath of name by joining current working dir when name is relative path
// name MUST be valid fs path value
func (ofs *OSFS) getRealPath(name string) (cwd, rpath string, _ error) {
	if !fs.ValidPath(name) && ofs.strict {
		return "", "", &fs.PathError{
			Op:   "",
			Err:  fs.ErrInvalid,
			Path: name,
		}
	}

	if path.IsAbs(name) {
		if runtime.GOOS == "windows" {
			cwd, err := ofs.getCwd()
			if err != nil {
				return "", "", err
			}

			rpath, err = pathhelper.AbsWindowsPath(cwd, name, func() (string, error) {
				// TODO: lookup fhs root
				return "", wellknownerrors.ErrNotSupported
			})

			return "", rpath, err
		}

		return "", name, nil
	}

	cwd, err := ofs.getCwd()
	if err != nil {
		return "", "", err
	}

	return cwd, filepath.Join(cwd, name), nil
}

func (ofs *OSFS) ReadDir(name string) ([]fs.DirEntry, error) {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return nil, err
	}

	return os.ReadDir(path)
}

func (ofs *OSFS) ReadFile(name string) ([]byte, error) {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path)
}

func (ofs *OSFS) Glob(pattern string) ([]string, error) {
	return doublestar.Glob(ofs, pattern)
}

func (ofs *OSFS) Sub(dir string) (fs.FS, error) {
	_, path, err := ofs.getRealPath(dir)
	if err != nil {
		return nil, err
	}

	return NewOSFS(ofs.strict, func() (string, error) { return path, nil }), nil
}

func (ofs *OSFS) Create(name string) (fs.File, error) {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return nil, err
	}

	return os.Create(path)
}

func (ofs *OSFS) Mkdir(name string, perm fs.FileMode) error {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return err
	}

	return os.Mkdir(path, perm)
}

func (ofs *OSFS) MkdirAll(name string, perm fs.FileMode) error {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return err
	}

	return os.MkdirAll(path, perm)
}

func (ofs *OSFS) Open(name string) (fs.File, error) {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return nil, err
	}

	return os.Open(path)
}

func (ofs *OSFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return nil, err
	}

	return os.OpenFile(path, flag, perm)
}

func (ofs *OSFS) Remove(name string) error {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return err
	}

	return os.Remove(path)
}

func (ofs *OSFS) RemoveAll(name string) error {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return err
	}

	return os.RemoveAll(path)
}

func (ofs *OSFS) Rename(oldname, newname string) error {
	_, oldPath, err := ofs.getRealPath(oldname)
	if err != nil {
		return err
	}

	_, newPath, err := ofs.getRealPath(newname)
	if err != nil {
		return err
	}

	return os.Rename(oldPath, newPath)
}

func (ofs *OSFS) Stat(name string) (fs.FileInfo, error) {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return nil, err
	}

	return os.Stat(path)
}

func (ofs *OSFS) Chmod(name string, mode fs.FileMode) error {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return err
	}

	return os.Chmod(path, mode)
}

func (ofs *OSFS) Chown(name string, uid, gid int) error {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return err
	}

	return os.Chown(path, uid, gid)
}

func (ofs *OSFS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return err
	}

	return os.Chtimes(path, atime, mtime)
}

func (ofs *OSFS) Lstat(name string) (fs.FileInfo, error) {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return nil, err
	}

	return os.Lstat(path)
}

func (ofs *OSFS) Symlink(oldname, newname string) error {
	if path.IsAbs(oldname) && path.IsAbs(newname) {
		if ofs.strict {
			return &fs.PathError{
				Op:   "",
				Path: oldname,
				Err:  fs.ErrInvalid,
			}
		}

		// nothing to do with cwd
		return os.Symlink(oldname, newname)
	}

	if ofs.strict {
		if !fs.ValidPath(oldname) {
			return &fs.PathError{
				Op:   "",
				Path: oldname,
				Err:  fs.ErrInvalid,
			}
		}

		if !fs.ValidPath(newname) {
			return &fs.PathError{
				Op:   "",
				Path: newname,
				Err:  fs.ErrInvalid,
			}
		}
	}

	cwd, err := ofs.getCwd()
	if err != nil {
		return err
	}

	// either oldname or newname is relative path
	// so we need to create symlink based on the current working dir
	return Symlinkat(cwd, oldname, newname)
}

func (ofs *OSFS) Readlink(name string) (string, error) {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return "", err
	}

	return os.Readlink(path)
}

func (ofs *OSFS) Abs(name string) (string, error) {
	_, path, err := ofs.getRealPath(name)
	if err != nil {
		return "", err
	}

	return path, nil
}
