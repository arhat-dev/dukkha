package unionfs

import "io/fs"

var (
	_ fs.File        = (*fileOverride)(nil)
	_ fs.ReadDirFile = (*fileOverride)(nil)
)

type fileOverride struct {
	fs.File

	nameOverride string
}

func (fe *fileOverride) Stat() (fs.FileInfo, error) {
	info, err := fe.File.Stat()
	if err != nil {
		return nil, err
	}

	return &fileInfoOverride{FileInfo: info, nameOverride: fe.nameOverride}, nil
}

func (fe *fileOverride) ReadDir(n int) ([]fs.DirEntry, error) {
	rdf, ok := fe.File.(fs.ReadDirFile)
	if !ok {
		return nil, &fs.PathError{Op: "readdir", Path: fe.nameOverride, Err: fs.ErrInvalid}
	}

	return rdf.ReadDir(n)
}

type fileInfoOverride struct {
	fs.FileInfo

	nameOverride string
}

func (fi *fileInfoOverride) Name() string {
	return fi.nameOverride
}
