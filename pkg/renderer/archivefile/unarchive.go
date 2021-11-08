package archivefile

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"

	"arhat.dev/pkg/iohelper"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"
)

func unarchive(src archiveSource, typ types.Type, inArchivePath, password string) (io.ReadCloser, error) {
	switch typ {
	case matchers.TypeZip:
		rd, err := unzip(src, inArchivePath, password)
		if err != nil {
			return nil, err
		}

		return iohelper.CustomReadCloser(rd, src.Close), nil
	case matchers.TypeTar:
		r, err := untar(src, inArchivePath)
		if err != nil {
			return nil, err
		}

		return iohelper.CustomReadCloser(r, src.Close), nil
	case matchers.TypeRar:
		rd, err := unrar(src, inArchivePath, password)
		if err != nil {
			return nil, err
		}

		return iohelper.CustomReadCloser(rd, src.Close), nil
	case matchers.TypeGz:
		r, err := gzip.NewReader(src)
		if err != nil {
			return nil, err
		}

		if len(inArchivePath) == 0 {
			return iohelper.CustomReadCloser(r, func() error {
				_ = r.Close()
				return src.Close()
			}), nil
		}

		return unarchiveNext(src, r, inArchivePath, password)
	case matchers.TypeBz2:
		r := bzip2.NewReader(src)
		if len(inArchivePath) == 0 {
			return iohelper.CustomReadCloser(r, func() error {
				return src.Close()
			}), nil
		}

		return unarchiveNext(src, r, inArchivePath, password)
	case matchers.Type7z:
		// TODO
	case matchers.TypeXz:
		r, err := xz.NewReader(src)
		if err != nil {
			return nil, err
		}

		if len(inArchivePath) == 0 {
			return iohelper.CustomReadCloser(r, func() error {
				return src.Close()
			}), nil
		}

		return unarchiveNext(src, r, inArchivePath, password)
	case matchers.TypeZstd:
		r, err := zstd.NewReader(src)
		if err != nil {
			return nil, err
		}

		if len(inArchivePath) == 0 {
			return iohelper.CustomReadCloser(r, func() error {
				return src.Close()
			}), nil
		}

		return unarchiveNext(src, r, inArchivePath, password)
	case matchers.TypePdf:
		// TODO
	case matchers.TypeDeb:
		// TODO
	case matchers.TypeAr:
		// unix archive
		// TODO
	case matchers.TypeLz:
		r := lz4.NewReader(src)
		if len(inArchivePath) == 0 {
			return iohelper.CustomReadCloser(r, func() error {
				return src.Close()
			}), nil
		}

		return unarchiveNext(src, r, inArchivePath, password)
	case matchers.TypeRpm:
		// TODO
	}

	return nil, fmt.Errorf("no implementation")
}

func unarchiveNext(src archiveSource, r io.Reader, inArchivePath, password string) (io.ReadCloser, error) {
	src = newArchiveSource(src, r)

	restore, err := prepareSeekRestore(src)
	if err != nil {
		return nil, err
	}

	typ, err := filetype.MatchReader(src)
	if err != nil {
		return nil, err
	}

	err = restore()
	if err != nil {
		return nil, err
	}

	return unarchive(src, typ, inArchivePath, password)
}
