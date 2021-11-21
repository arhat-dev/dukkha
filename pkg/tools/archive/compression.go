package archive

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"strconv"

	"arhat.dev/dukkha/pkg/constant"
	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
)

func createCompressionStream(w io.Writer, method, level string) (io.WriteCloser, error) {
	switch method {
	case constant.CompressionMethod_DEFLATE:
		l, err := parseFlateCompressionLevel(level)
		if err != nil {
			return nil, err
		}

		return flate.NewWriter(w, l)
	case constant.CompressionMethod_Gzip:
		l, err := parseGzipCompressionLevel(level)
		if err != nil {
			return nil, err
		}

		return gzip.NewWriterLevel(w, l)
	case constant.CompressionMethod_Bzip2:
		l, err := parseBzip2CompresssionLevel(level)
		if err != nil {
			return nil, err
		}

		return bzip2.NewWriter(w, &bzip2.WriterConfig{
			Level: l,
		})
	case constant.CompressionMethod_XZ:
		return xz.WriterConfig{}.NewWriter(w)
	case constant.CompressionMethod_LZMA:
		return lzma.WriterConfig{}.NewWriter(w)
	case constant.CompressionMethod_ZSTD:
		l, err := parseZstdCompressionLevel(level)
		if err != nil {
			return nil, err
		}

		return zstd.NewWriter(w, zstd.WithEncoderLevel(l))
	default:
		return nil, fmt.Errorf("unsupported compression method: %q", method)
	}
}

func parseFlateCompressionLevel(level string) (int, error) {
	if len(level) == 0 {
		return flate.DefaultCompression, nil
	}

	l, err := strconv.ParseInt(level, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(l), nil
}

func parseGzipCompressionLevel(level string) (int, error) {
	if len(level) == 0 {
		return gzip.DefaultCompression, nil
	}

	l, err := strconv.ParseInt(level, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(l), nil
}

func parseZstdCompressionLevel(level string) (zstd.EncoderLevel, error) {
	if len(level) == 0 {
		return zstd.SpeedDefault, nil
	}

	l, err := strconv.ParseInt(level, 10, 64)
	if err != nil {
		return 0, err
	}

	if l < 3 {
		return zstd.SpeedFastest, nil
	}

	switch l {
	case 3, 4, 5, 6:
		return zstd.SpeedDefault, nil
	case 7, 8:
		return zstd.SpeedBetterCompression, nil
	default:
		return zstd.SpeedBestCompression, nil
	}
}

func parseBzip2CompresssionLevel(level string) (int, error) {
	if len(level) == 0 {
		return bzip2.DefaultCompression, nil
	}

	l, err := strconv.ParseInt(level, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(l), nil
}
