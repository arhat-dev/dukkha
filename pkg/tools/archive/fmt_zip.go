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
	"arhat.dev/pkg/stringhelper"

	"arhat.dev/dukkha/pkg/constant"
)

func createZip(
	ofs *fshelper.OSFS,
	w io.Writer,
	files []*entry,
	enableCompression bool,
	compressionMethod string,
	compressionLevel string,
) (err error) {
	zw := zip.NewWriter(w)
	defer func() { _ = zw.Close() }()

	var method uint16
	if enableCompression {
		switch compressionMethod {
		case constant.CompressionMethod_DEFLATE:
			var level int
			level, err = parseFlateCompressionLevel(compressionLevel)
			if err != nil {
				return
			}

			method = zip.Deflate
			zw.RegisterCompressor(method, func(w io.Writer) (io.WriteCloser, error) {
				return flate.NewWriter(w, level)
			})
		case constant.CompressionMethod_Bzip2:
			var level int
			level, err = parseBzip2CompresssionLevel(compressionLevel)
			if err != nil {
				return
			}

			method = uint16(constant.ZipCompressionMethod_BZIP2)
			zw.RegisterCompressor(method, func(w io.Writer) (io.WriteCloser, error) {
				return bzip2.NewWriter(w, &bzip2.WriterConfig{
					Level: level,
				})
			})
		case constant.CompressionMethod_LZMA:
			method = uint16(constant.ZipCompressionMethod_LZMA)
			zw.RegisterCompressor(method, func(w io.Writer) (io.WriteCloser, error) {
				return lzma.WriterConfig{}.NewWriter(w)
			})
		case constant.CompressionMethod_ZSTD:
			var level zstd.EncoderLevel
			level, err = parseZstdCompressionLevel(compressionLevel)
			if err != nil {
				return err
			}

			method = uint16(constant.ZipCompressionMethod_ZSTD)
			zw.RegisterCompressor(method, func(w io.Writer) (io.WriteCloser, error) {
				return zstd.NewWriter(w, zstd.WithEncoderLevel(level))
			})
		case constant.CompressionMethod_XZ:
			method = uint16(constant.ZipCompressionMethod_XZ)
			zw.RegisterCompressor(method, func(w io.Writer) (io.WriteCloser, error) {
				return xz.WriterConfig{}.NewWriter(w)
			})
		default:
			err = fmt.Errorf("unsupported compression method %q", compressionMethod)
			return
		}
	}

	var (
		hdr  *zip.FileHeader
		mode fs.FileMode

		entryWriter io.Writer
	)

	for _, f := range files {
		hdr, err = zip.FileInfoHeader(f.info)
		if err != nil {
			return
		}

		hdr.Name = f.to
		hdr.Method = method

		mode = f.info.Mode()
		if mode.IsDir() && !strings.HasSuffix(hdr.Name, "/") {
			hdr.Name += "/"
		}

		entryWriter, err = zw.CreateHeader(hdr)
		if err != nil {
			return err
		}

		switch {
		case mode&fs.ModeSymlink != 0:
			_, err = entryWriter.Write(stringhelper.ToBytes[byte, byte](f.link))
			if err != nil {
				return
			}
		case mode.IsRegular():
			err = copyFileContent(ofs, entryWriter, f.from)
			if err != nil {
				return
			}
		}
	}

	return nil
}
