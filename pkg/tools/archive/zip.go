package archive

import (
	"archive/zip"
	"compress/flate"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"

	"arhat.dev/pkg/fshelper"

	"arhat.dev/dukkha/pkg/constant"
)

func createZip(
	ofs *fshelper.OSFS,
	w io.Writer, files []*entry,
	compressionMethod *string,
	compressionLevel string,
) error {
	zw := zip.NewWriter(w)
	defer func() { _ = zw.Close() }()

	var cm uint16
	if compressionMethod != nil {
		switch method := *compressionMethod; method {
		case constant.CompressionMethod_DEFLATE:
			level, err := parseFlateCompressionLevel(compressionLevel)
			if err != nil {
				return err
			}

			cm = zip.Deflate
			zw.RegisterCompressor(cm, func(w io.Writer) (io.WriteCloser, error) {
				return flate.NewWriter(w, level)
			})
		case constant.CompressionMethod_Bzip2:
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
		case constant.CompressionMethod_LZMA:
			cm = uint16(constant.ZipCompressionMethod_LZMA)
			zw.RegisterCompressor(cm, func(w io.Writer) (io.WriteCloser, error) {
				return lzma.WriterConfig{}.NewWriter(w)
			})
		case constant.CompressionMethod_ZSTD:
			level, err := parseZstdCompressionLevel(compressionLevel)
			if err != nil {
				return err
			}

			cm = uint16(constant.ZipCompressionMethod_ZSTD)
			zw.RegisterCompressor(cm, func(w io.Writer) (io.WriteCloser, error) {
				return zstd.NewWriter(w, zstd.WithEncoderLevel(level))
			})
		case constant.CompressionMethod_XZ:
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
			err = copyFileContent(ofs, wr, f.from)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
