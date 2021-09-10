package unionfs

import (
	"io"
	"io/fs"
	"path"
	"strings"
	"time"
)

type fsEntry struct {
	oldPrefix string
	newPrefix string
	_fs       fs.FS
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
	f, err := fe._fs.Open(actualName)
	if err != nil {
		return nil, err
	}

	if actualName != name {
		return &fileEntry{File: f, nameOverride: name}, nil
	}

	return f, nil
}

func (fe *fsEntry) ReadDir(name string) ([]fs.DirEntry, error) {
	rdFS, ok := fe._fs.(fs.ReadDirFS)
	if !ok {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrInvalid}
	}

	dirs, err := rdFS.ReadDir(fe.getActualName(name))
	if err != nil {
		return nil, err
	}

	return rewriteDirEntries(fe.oldPrefix, dirs), nil
}

var _ fs.File = (*fileEntry)(nil)

type fileEntry struct {
	fs.File

	nameOverride string
}

func (fe *fileEntry) Stat() (fs.FileInfo, error) {
	info, err := fe.File.Stat()
	if err != nil {
		return nil, err
	}

	return &fileInfo{FileInfo: info, name: fe.nameOverride}, nil
}

type fileInfo struct {
	fs.FileInfo
	name string
}

func (fi *fileInfo) Name() string {
	return fi.name
}

func New() *UnionFS {
	return &UnionFS{
		t: NewTrie(PathReverseKeysFunc),
	}
}

type UnionFS struct {
	t *Trie
}

func (f *UnionFS) Remap(newPrefix, oldPrefix string, _fs fs.FS) {
	_ = f.t.Add(newPrefix, &fsEntry{
		newPrefix: newPrefix,
		oldPrefix: oldPrefix,
		_fs:       _fs,
	})
}

func (f *UnionFS) find(name string) (_ *Node, isExact bool, err error) {
	exact, nearest, ok := f.t.Get(name)
	if ok {
		return exact, true, nil
	}

	if nearest != nil {
		return nearest, false, err
	}

	return nil, false, &fs.PathError{
		Op:   "open",
		Path: name,
		Err:  fs.ErrNotExist,
	}
}

func (f *UnionFS) Open(name string) (fs.File, error) {
	node, isExact, err := f.find(name)
	if err != nil {
		return nil, err
	}

	ent, ok := node.Value().(*fsEntry)
	if ok {
		return ent.Open(name)
	}

	if isExact {
		return newDir(name, node.Children()), nil
	}

	return nil, &fs.PathError{
		Op:   "open",
		Path: name,
		Err:  fs.ErrNotExist,
	}
}

func (f *UnionFS) ReadDir(name string) ([]fs.DirEntry, error) {
	node, _, err := f.find(name)
	if err != nil {
		return nil, err
	}

	ent, ok := node.Value().(*fsEntry)
	if ok {
		return ent.ReadDir(name)
	}

	return newDirEntries(node.Children()), nil
}

func newDir(name string, files []Node) fs.File {
	return &dir{path: name, entry: newDirEntries(files)}
}

func rewriteDirEntries(trimNamePrefix string, original []fs.DirEntry) []fs.DirEntry {
	ret := make([]fs.DirEntry, len(original))
	for i, d := range original {
		ret[i] = &dirEntry{
			DirEntry: original[i],
			name:     strings.TrimPrefix(strings.TrimPrefix(d.Name(), trimNamePrefix), "/"),
		}
	}
	return ret
}

type dirEntry struct {
	fs.DirEntry
	name string
}

func (de *dirEntry) Name() string {
	return de.name
}

func newDirEntries(files []Node) []fs.DirEntry {
	var entries []fs.DirEntry
	for _, f := range files {
		entries = append(entries, &mapFileInfo{
			name:    f.ElementKey(),
			mode:    fs.ModeDir,
			modTime: time.Time{},
			sys:     nil,
		})
	}

	return entries
}

var (
	_ fs.File = (*dir)(nil)
)

type dir struct {
	path string

	entry  []fs.DirEntry
	offset int
}

func (d *dir) Stat() (fs.FileInfo, error) {
	return &mapFileInfo{
		name:    d.path,
		mode:    fs.ModeDir,
		modTime: time.Time{},
		sys:     nil,
	}, nil
}

func (d *dir) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.path, Err: fs.ErrInvalid}
}

func (d *dir) Close() error {
	return nil
}

func (d *dir) ReadDir(count int) ([]fs.DirEntry, error) {
	n := len(d.entry) - d.offset
	if n == 0 && count > 0 {
		return nil, io.EOF
	}
	if count > 0 && n > count {
		n = count
	}
	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = d.entry[d.offset+i]
	}
	d.offset += n
	return list, nil
}

// A mapFileInfo implements fs.FileInfo and fs.DirEntry for a given map file.
type mapFileInfo struct {
	name string

	data    []byte      // file content
	mode    fs.FileMode // FileInfo.Mode
	modTime time.Time   // FileInfo.ModTime
	sys     interface{} // FileInfo.Sys
}

func (i *mapFileInfo) Name() string               { return i.name }
func (i *mapFileInfo) Size() int64                { return int64(len(i.data)) }
func (i *mapFileInfo) Mode() fs.FileMode          { return i.mode }
func (i *mapFileInfo) Type() fs.FileMode          { return i.mode.Type() }
func (i *mapFileInfo) ModTime() time.Time         { return i.modTime }
func (i *mapFileInfo) IsDir() bool                { return i.mode&fs.ModeDir != 0 }
func (i *mapFileInfo) Sys() interface{}           { return i.sys }
func (i *mapFileInfo) Info() (fs.FileInfo, error) { return i, nil }
