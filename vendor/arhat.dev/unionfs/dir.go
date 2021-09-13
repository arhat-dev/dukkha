package unionfs

import (
	"io"
	"io/fs"
	"strings"
	"time"

	"arhat.dev/unionfs/trie"
)

func newDir(name string, files []*trie.Node) fs.File {
	return &dir{path: name, entry: newDirEntries(files)}
}

func newDirEntries(files []*trie.Node) []fs.DirEntry {
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

var _ fs.File = (*dir)(nil)

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
