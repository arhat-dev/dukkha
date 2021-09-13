package unionfs

import (
	"io/fs"
	"path"
	"strings"
)

var (
	_ fs.ReadDirFS  = (*fsEntry)(nil)
	_ fs.ReadFileFS = (*fsEntry)(nil)
	// TODO
	// _ fs.GlobFS = (*fsEntry)(nil)
	// _ fs.StatFS = (*fsEntry)(nil)
	// _ fs.SubFS  = (*fsEntry)(nil)
)

type fsEntry struct {
	oldPrefix string
	newPrefix string
	vfs       fs.FS
}

func (fe *fsEntry) getActualName(name string) string {
	if fe.oldPrefix == fe.newPrefix {
		return name
	}

	actualName := strings.TrimPrefix(name, fe.newPrefix)
	actualName = path.Join(fe.oldPrefix, strings.TrimPrefix(actualName, "/"))
	return actualName
}

func (fe *fsEntry) Open(name string) (fs.File, error) {
	actualName := fe.getActualName(name)
	f, err := fe.vfs.Open(actualName)
	if err != nil {
		return nil, err
	}

	if actualName != name {
		return &fileOverride{File: f, nameOverride: name}, nil
	}

	return f, nil
}

func (fe *fsEntry) ReadDir(name string) ([]fs.DirEntry, error) {
	rdFS, ok := fe.vfs.(fs.ReadDirFS)
	if !ok {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrInvalid}
	}

	dirs, err := rdFS.ReadDir(fe.getActualName(name))
	if err != nil {
		return nil, err
	}

	return rewriteDirEntries(fe.oldPrefix, dirs), nil
}

func (fe *fsEntry) ReadFile(name string) ([]byte, error) {
	return nil, nil
}
