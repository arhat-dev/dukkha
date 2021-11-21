package archivefile

import (
	"archive/zip"
	"compress/bzip2"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/pkg/iohelper"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
)

type SizedReaderAt interface {
	sizeIface
	io.ReaderAt
}

func unzip(src SizedReaderAt, target, password string) (io.Reader, error) {
	// TODO: support encrypted zip file
	_ = password

	r, err := zip.NewReader(src, src.Size())
	if err != nil {
		return nil, err
	}

	r.RegisterDecompressor(uint16(constant.ZipCompressionMethod_BZIP2), func(r io.Reader) io.ReadCloser {
		return iohelper.CustomReadCloser(bzip2.NewReader(r), func() error { return nil })
	})

	r.RegisterDecompressor(uint16(constant.ZipCompressionMethod_LZMA), func(r io.Reader) io.ReadCloser {
		rd, err := lzma.ReaderConfig{}.NewReader(r)
		if err != nil {
			return nil
		}

		return iohelper.CustomReadCloser(rd, func() error { return nil })
	})

	r.RegisterDecompressor(uint16(constant.ZipCompressionMethod_ZSTD), func(r io.Reader) io.ReadCloser {
		zr, err := zstd.NewReader(r)
		if err != nil {
			return nil
		}
		return zr.IOReadCloser()
	})

	r.RegisterDecompressor(uint16(constant.ZipCompressionMethod_XZ), func(r io.Reader) io.ReadCloser {
		xr, err := xz.NewReader(r)
		if err != nil {
			return nil
		}
		return ioutil.NopCloser(xr)
	})

	for {
		f, err := r.Open(strings.TrimPrefix(target, "/"))
		if err != nil {
			return nil, fmt.Errorf("unzip: %w", err)
		}

		info, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("unzip: %w", err)
		}

		switch m := info.Mode() & fs.ModeType; m {
		case 0:
			// file
			return f, nil
		case fs.ModeSymlink:
			// TODO: redirect to links
			targetBytes, err := io.ReadAll(f)
			if err != nil {
				return nil, err
			}

			target = cleanLink(info.Name(), string(targetBytes))
			continue
		default:
			return nil, fmt.Errorf("unzip: unsupported non regular file %q: %v", info.Name(), m)
		}
	}
}
