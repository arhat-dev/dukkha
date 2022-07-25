package templateutils

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"arhat.dev/pkg/clihelper"
	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/stringhelper"
	"github.com/spf13/pflag"
	"mvdan.cc/sh/v3/interp"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/dukkha/pkg/dukkha"
)

func createFSNS(rc dukkha.RenderingContext) fsNS { return fsNS{rc: rc} }

type fsNS struct{ rc dukkha.RenderingContext }

func (fsNS) Base(p String) string              { return filepath.Base(must(toString(p))) }
func (fsNS) Clean(p String) string             { return filepath.Clean(must(toString(p))) }
func (fsNS) Dir(p String) string               { return filepath.Dir(must(toString(p))) }
func (fsNS) Ext(p String) string               { return filepath.Ext(must(toString(p))) }
func (fsNS) IsAbs(p String) bool               { return filepath.IsAbs(must(toString(p))) }
func (fsNS) Split(p String) (dir, file string) { return filepath.Split(must(toString(p))) }
func (fsNS) FromSlash(p String) string         { return filepath.FromSlash(must(toString(p))) }
func (fsNS) ToSlash(p String) string           { return filepath.ToSlash(must(toString(p))) }
func (fsNS) VolumeName(p String) string        { return filepath.VolumeName(must(toString(p))) }

func (fsNS) Join(elem ...String) (_ string, err error) {
	elements, err := toStrings(elem)
	if err != nil {
		return
	}

	return filepath.Join(elements...), nil
}

func (fsNS) Match(pattern, name String) (matched bool, err error) {
	ptn, err := toString(pattern)
	if err != nil {
		return
	}

	tgt, err := toString(name)
	if err != nil {
		return
	}

	return filepath.Match(ptn, tgt)
}

func (fsNS) Rel(basepath, targpath String) (_ string, err error) {
	bp, err := toString(basepath)
	if err != nil {
		return
	}

	tp, err := toString(targpath)
	if err != nil {
		return
	}

	return filepath.Rel(bp, tp)
}

func (ns fsNS) Glob(pattern String) (_ []string, err error) {
	ptn, err := toString(pattern)
	if err != nil {
		return
	}

	return ns.rc.FS().Glob(ptn)
}

func (ns fsNS) Abs(path String) (_ string, err error) {
	p, err := toString(path)
	if err != nil {
		return
	}

	return ns.rc.FS().Abs(p)
}

func (ns fsNS) Exists(path String) (_ bool, err error) {
	p, err := toString(path)
	if err != nil {
		return
	}

	_, err = ns.rc.FS().Lstat(p)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return err == nil, err
}

func (ns fsNS) IsDir(path String) bool        { return ns.fileTypeIs(fs.ModeDir, path) }
func (ns fsNS) IsSymlink(path String) bool    { return ns.fileTypeIs(fs.ModeSymlink, path) }
func (ns fsNS) IsDevice(path String) bool     { return ns.fileTypeIs(fs.ModeDevice, path) }
func (ns fsNS) IsCharDevice(path String) bool { return ns.fileTypeIs(fs.ModeCharDevice, path) }
func (ns fsNS) IsFIFO(path String) bool       { return ns.fileTypeIs(fs.ModeNamedPipe, path) }
func (ns fsNS) IsSocket(path String) bool     { return ns.fileTypeIs(fs.ModeSocket, path) }
func (ns fsNS) IsOther(path String) bool      { return ns.fileTypeIs(fs.ModeIrregular, path) }

func (ns fsNS) fileTypeIs(expected fs.FileMode, file String) bool {
	f, err := toString(file)
	if err != nil {
		return false
	}

	info, err := ns.rc.FS().Lstat(f)
	if err != nil {
		return false
	}

	return info.Mode()&expected == expected
}

// ReadDir read a list of entry names (files & dirs) in directory dir
func (ns fsNS) ReadDir(dir String) (ret []string, err error) {
	d, err := toString(dir)
	if err != nil {
		return
	}

	osfs := ns.rc.FS()
	entries, err := osfs.ReadDir(d)
	if err != nil {
		return
	}

	ret = make([]string, len(entries))
	for i, ent := range entries {
		ret[i] = ent.Name()
	}

	return
}

// Find works like unix cli find (without -exec supprot), return a list of matched entries by walking every
// filesystem entry from the start path
func (ns fsNS) Find(path String, args ...String) (ret []string, err error) {
	p, err := toString(path)
	if err != nil {
		return
	}

	flags, err := toStrings(args)
	if err != nil {
		return
	}

	return clihelper.FindCli(ns.rc.FS(), p, flags...)
}

// Lookup lookup executable by name in PATH list, return empty string if not found
//
// Lookup(bin String): lookup bin according to PATH env, take PATHEXT into consideration on windows
//
// Lookup(path, file String): like the one above, but use specified path as PATH env
//
// Lookup(path, pathext, file String): specify both PATH and PATHEXT manually,
// 										but if path is empty, fallback to PATH env
//
// NOTE: it will try extra suffices (e.g. `.exe`) on windows
//
// only return err when input is problematic
func (ns fsNS) Lookup(args ...String) (ret string, err error) {
	var (
		rc   = ns.rc
		exec String
	)

	switch n := len(args); n {
	case 0:
		err = errAtLeastOneArgGotZero
		return
	case 1:
		// Lookup(file)
		exec = args[0]
	case 2:
		// Lookup(PATH, file)
		var pathEnv string
		pathEnv, err = toString(args[0])
		if err != nil {
			return
		}

		rc = ns.rc.(dukkha.Context).DeriveNew()
		rc.AddEnv(true, &dukkha.EnvEntry{
			Name:  "PATH",
			Value: pathEnv,
		})

		exec = args[1]
	default:
		// Lookup(PATH, PATHEXT, ... file)

		var (
			pathEnv, pathextEnv string
		)

		pathEnv, err = toString(args[0])
		if err != nil {
			return
		}

		pathextEnv, err = toString(args[1])
		if err != nil {
			return
		}

		exec = args[n-1]

		if len(pathEnv) == 0 {
			pathEnv = rc.Get("PATH").String()
		}

		rc = ns.rc.(dukkha.Context).DeriveNew()
		rc.AddEnv(true, &dukkha.EnvEntry{
			Name:  "PATH",
			Value: pathEnv,
		}, &dukkha.EnvEntry{
			Name:  "PATHEXT",
			Value: pathextEnv,
		})
	}

	execFile, err := toString(exec)
	if err != nil {
		return
	}

	goos, ok := constant.GetGolangOS(rc.HostKernel())
	if !ok {
		goos = runtime.GOOS
	}

	ret, _ = interp.DukkhaLookPathDir(goos, rc.WorkDir(), execFile, rc, interp.DukkhaFindExecutable)
	return ret, nil
}

// LookupFile lookup file by name in PATH list, return empty string if not found
//
// it's like Lookup but doesn't require file to be executable
//
// NOTE: it will not try extra suffices on all platforms
//
// only return err when input is problematic
func (ns fsNS) LookupFile(args ...String) (ret string, err error) {
	var (
		rc   = ns.rc
		file String
	)
	switch n := len(args); n {
	case 0:
		err = errAtLeastOneArgGotZero
		return
	case 1:
		// LookupFile(file)
		file = args[0]
	default:
		// LookupFile(PATH, ... file)
		var pathEnv string
		pathEnv, err = toString(args[0])
		if err != nil {
			return
		}

		rc = ns.rc.(dukkha.Context).DeriveNew()
		rc.AddEnv(true, &dukkha.EnvEntry{
			Name:  "PATH",
			Value: pathEnv,
		})

		file = args[n-1]
	}

	goos, ok := constant.GetGolangOS(rc.HostKernel())
	if !ok {
		goos = runtime.GOOS
	}

	f, err := toString(file)
	if err != nil {
		return
	}

	return interp.DukkhaLookPathDir(goos, rc.WorkDir(), f, rc, interp.DukkhaFindFile)
}

// ReadFile reads all content from local file
func (ns fsNS) ReadFile(path String) (_ string, err error) {
	f, err := toString(path)
	if err != nil {
		return
	}

	data, err := ns.rc.FS().ReadFile(f)
	if err != nil {
		return
	}

	return stringhelper.Convert[string, byte](data), nil
}

// OpenFile opens a local file
//
// OpenFile(file String): open a local file as read-only to read from start
//
// OpenFile(...<options>, file String): open a local file with options
// where options are:
// - `--mode` or `-m` mode: permission bits
// - `--flags` or `-f` flags: file open flags
func (ns fsNS) OpenFile(args ...String) (_ *os.File, err error) {
	n := len(args)
	if n == 0 {
		err = errAtLeastOneArgGotZero
		return
	}

	path, err := toString(args[0])
	if err != nil {
		return
	}

	if n == 1 {
		var f fs.File
		f, err = ns.rc.FS().Open(path)
		if err != nil {
			return
		}

		return f.(*os.File), nil
	}

	flags, err := toStrings(args[:n-1])
	if err != nil {
		return
	}

	var (
		pfs pflag.FlagSet

		mode          uint32
		fflags        string
		fileOpenFlags int
		read, write   bool
	)

	clihelper.InitFlagSet(&pfs, "open-file")

	pfs.Uint32VarP(&mode, "mode", "m", 0, "")
	pfs.StringVarP(&fflags, "flags", "f", "r", "")

	err = pfs.Parse(flags)
	if err != nil {
		return
	}

	for _, c := range fflags {
		switch c {
		case 'r':
			read = true
		case 'w':
			write = true
		case 'a':
			fileOpenFlags |= os.O_APPEND
		case 'x':
			fileOpenFlags |= os.O_EXCL
		case 's':
			fileOpenFlags |= os.O_SYNC
		}
	}

	// nolint:gocritic
	if read && write {
		fileOpenFlags |= os.O_RDWR
	} else if read {
		fileOpenFlags |= os.O_RDONLY
	} else if write {
		fileOpenFlags |= os.O_WRONLY
	}

	var f fs.File
	f, err = ns.rc.FS().OpenFile(path, fileOpenFlags, fs.FileMode(mode))
	if err != nil {
		return
	}

	return f.(*os.File), nil
}

// Touch is an alias of WriteFile(file)
func (ns fsNS) Touch(file String) (None, error) { return ns.WriteFile(file) }

// WriteFile write data to file in O_TRUNC mode
//
// WriteFile(path String): touch file at path, if didn't exist, create one with permission 0640
//
// WriteFile(path String, data BytesOrReader): write data to file at path with permission 0640
//
// WriteFile(path String, perm Integer, data BytesOrReader): write data to file at path with specified permission
func (ns fsNS) WriteFile(file String, args ...any) (_ None, err error) {
	f, err := toString(file)
	if err != nil {
		return
	}

	err = handleFileWrite(args, ns.rc.FS(), f, false)
	return
}

// AppendFile writes data to file in O_APPEND mode
//
// AppendFile(path String): touch file at path, if didn't exist, create one with permission 0640
//
// AppendFile(path String, data BytesOrReader): append data to file at path with permission 0640
//
// AppendFile(path String, perm Integer, data BytesOrReader): append data to file at path with specified permission
func (ns fsNS) AppendFile(file String, args ...any) (_ None, err error) {
	f, err := toString(file)
	if err != nil {
		return
	}

	err = handleFileWrite(args, ns.rc.FS(), f, true)
	return
}

func handleFileWrite(args []any, osfs *fshelper.OSFS, file string, append bool) (err error) {
	path, err := toString(file)
	if err != nil {
		return
	}

	perm := fs.FileMode(0640)
	n := len(args)

	if n == 0 {
		// touch file (DO NOT use State)
		var f fs.File
		f, err = osfs.Open(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return osfs.WriteFile(path, []byte{}, perm)
			}

			return
		}

		return f.Close()
	}

	if n > 1 {
		perm = toIntegerOrPanic[fs.FileMode](args[0])
	}

	mode := os.O_TRUNC
	if append {
		mode = os.O_APPEND
	}

	inData, inReader, isReader, err := toBytesOrReader(args[n-1])
	if err != nil {
		return
	}

	if append || isReader {
		f, err := osfs.OpenFile(path, os.O_WRONLY|os.O_CREATE|mode, perm)
		if err != nil {
			return err
		}

		if isReader {
			_, err = f.(*os.File).ReadFrom(inReader)
		} else {
			_, err = f.(io.Writer).Write(inData)
		}

		err2 := f.Close()
		if err == nil {
			return err2
		}

		return err
	}

	return osfs.WriteFile(path, inData, perm)
}

// Mkdir works like unix cli mkdir
//
// Mkdir(path String): create an empty dir at path with permission 0755, return error if failed or path exists
//
// Mkdir(...<options>, path String)
// where options are:
// - `--mode` or `-m` mode: set permission bits
// - `--parents`, `-p`: create intermediate directories
func (ns fsNS) Mkdir(args ...String) (_ None, err error) {
	n := len(args)
	if n == 0 {
		err = errAtLeastOneArgGotZero
		return
	}

	if n == 1 {
		var path string
		path, err = toString(args[n-1])
		if err != nil {
			return
		}

		err = ns.rc.FS().Mkdir(path, 0755)
		return
	}

	flags, err := toStrings(args)
	if err != nil {
		return
	}

	var (
		pfs pflag.FlagSet

		mode    uint32
		parents bool
	)

	clihelper.InitFlagSet(&pfs, "mkdir")

	pfs.Uint32VarP(&mode, "mode", "m", 0755, "")
	pfs.BoolVarP(&parents, "parents", "p", false, "")
	err = pfs.Parse(flags)
	if err != nil {
		return
	}

	var mkdirFn func(string, fs.FileMode) error
	if parents {
		mkdirFn = ns.rc.FS().MkdirAll
	} else {
		mkdirFn = ns.rc.FS().Mkdir
	}

	for _, dir := range pfs.Args() {
		err = mkdirFn(dir, fs.FileMode(mode))
		if err != nil {
			return
		}
	}

	return
}
