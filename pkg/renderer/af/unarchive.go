package af

import (
	"compress/bzip2"
	"compress/flate"
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
	"github.com/ulikunitz/xz/lzma"
)

var (
	lzmaType = filetype.NewType("lzma", "application/vnd.lzma")
	lz4Type  = filetype.NewType("lz4", "application/vnd.lz4")
)

// magic number ref: https://www.kernel.org/doc/html/latest/x86/boot.html
func init() {
	filetype.AddMatcher(lzmaType, func(b []byte) bool {
		return len(b) >= 2 &&
			b[0] == 0x5d && b[1] == 0x00
	})

	filetype.AddMatcher(lz4Type, func(b []byte) bool {
		return len(b) >= 2 &&
			b[0] == 0x02 && b[1] == 0x21
	})
}

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
	case matchers.TypeXz:
		r, err := xz.ReaderConfig{}.NewReader(src)
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
				r.Close()
				return src.Close()
			}), nil
		}

		return unarchiveNext(src, r, inArchivePath, password)
	case lz4Type:
		r := lz4.NewReader(src)
		if len(inArchivePath) == 0 {
			return iohelper.CustomReadCloser(r, func() error {
				return src.Close()
			}), nil
		}

		return unarchiveNext(src, r, inArchivePath, password)
	case lzmaType:
		r, err := lzma.ReaderConfig{}.NewReader(src)
		if err != nil {
			return nil, err
		}

		if len(inArchivePath) == 0 {
			return iohelper.CustomReadCloser(r, func() error {
				return src.Close()
			}), nil
		}

		return unarchiveNext(src, r, inArchivePath, password)
	case matchers.Type7z:
		// TODO
	case matchers.TypePdf:
		// TODO
	case matchers.TypeDeb:
		// TODO
	case matchers.TypeAr:
		// unix archive
		// TODO
	case matchers.TypeRpm:
		// TODO
	default:
		// assume deflate
		r := flate.NewReader(src)

		if len(inArchivePath) == 0 {
			return iohelper.CustomReadCloser(r, func() error {
				_ = r.Close()
				return src.Close()
			}), nil
		}

		return unarchiveNext(src, r, inArchivePath, password)
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
