package archive

import (
	"archive/zip"
	"compress/flate"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"arhat.dev/dukkha/pkg/constant"
	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
)

func createZip(
	w io.Writer, files []*entry,
	compressionMethod *string,
	compressionLevel string,
) error {
	zw := zip.NewWriter(w)
	defer func() { _ = zw.Close() }()

	var cm uint16
	if compressionMethod != nil {
		switch method := *compressionMethod; method {
		case "deflate":
			level, err := parseFlateCompressionLevel(compressionLevel)
			if err != nil {
				return err
			}

			cm = zip.Deflate
			zw.RegisterCompressor(cm, func(w io.Writer) (io.WriteCloser, error) {
				return flate.NewWriter(w, level)
			})
		case "bzip2":
			level, err := parseBzip2CompresssionLevel(compressionLevel)
			if err != nil {
				return err
			}

			cm = uint16(constant.ZipCompressionMethod_BZIP2)
			zw.RegisterCompressor(cm, func(w io.Writer) (io.WriteCloser, error) {
				return bzip2.NewWriter(w, &bzip2.WriterConfig{
					Level: level,
				})
			})
		case "lzma":
			cm = uint16(constant.ZipCompressionMethod_LZMA)
			zw.RegisterCompressor(cm, func(w io.Writer) (io.WriteCloser, error) {
				return lzma.WriterConfig{}.NewWriter(w)
			})
		case "zstd":
			level, err := parseZstdCompressionLevel(compressionLevel)
			if err != nil {
				return err
			}

			cm = uint16(constant.ZipCompressionMethod_ZSTD)
			zw.RegisterCompressor(cm, func(w io.Writer) (io.WriteCloser, error) {
				return zstd.NewWriter(w, zstd.WithEncoderLevel(level))
			})
		case "xz":
			cm = uint16(constant.ZipCompressionMethod_XZ)
			zw.RegisterCompressor(cm, func(w io.Writer) (io.WriteCloser, error) {
				return xz.WriterConfig{}.NewWriter(w)
			})
		default:
			return fmt.Errorf("unsupported compression method %q", method)
		}
	}

	for _, f := range files {
		hdr, err := zip.FileInfoHeader(f.info)
		if err != nil {
			return err
		}

		hdr.Name = f.to
		hdr.Method = cm

		mode := f.info.Mode()
		if mode.IsDir() && !strings.HasSuffix(hdr.Name, "/") {
			hdr.Name += "/"
		}

		wr, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}

		switch {
		case mode&fs.ModeSymlink != 0:
			_, err = wr.Write([]byte(f.link))
			if err != nil {
				return err
			}
		case mode.IsRegular():
			err = copyFileContent(wr, f.from)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
