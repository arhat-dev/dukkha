package unionfs

import (
	"io/fs"
	"time"

	"arhat.dev/unionfs/trie"
)

func New() *UnionFS {
	return &UnionFS{
		t: trie.New(PathKeyFunc),
	}
}

var (
	_ fs.ReadDirFS = (*UnionFS)(nil)
	// _ fs.ReadFileFS = (*UnionFS)(nil)
	// TODO
	// _ fs.GlobFS = (*UnionFS)(nil)
	// _ fs.StatFS = (*UnionFS)(nil)
	// _ fs.SubFS  = (*UnionFS)(nil)
)

type UnionFS struct {
	t *trie.Trie
}

// Map existing vfs with new path prefix inside the ufs, old prefix is required to access
// actual vfs data
func (ufs *UnionFS) Map(newPrefix, oldPrefix string, vfs fs.FS) {
	_ = ufs.t.Add(newPrefix, &fsEntry{
		newPrefix: newPrefix,
		oldPrefix: oldPrefix,
		vfs:       vfs,
	})
}

func (ufs *UnionFS) find(name string) (node *trie.Node, isExact bool, err error) {
	node, isExact = ufs.t.Get(name)
	if isExact || node != nil {
		return
	}

	return nil, false, &fs.PathError{
		Op:   "open",
		Path: name,
		Err:  fs.ErrNotExist,
	}
}

func (ufs *UnionFS) Open(name string) (fs.File, error) {
	node, isExact, err := ufs.find(name)
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

func (ufs *UnionFS) ReadDir(name string) ([]fs.DirEntry, error) {
	node, _, err := ufs.find(name)
	if err != nil {
		return nil, err
	}

	ent, ok := node.Value().(*fsEntry)
	if ok {
		return ent.ReadDir(name)
	}

	return newDirEntries(node.Children()), nil
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
