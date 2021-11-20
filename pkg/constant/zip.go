package constant

import "fmt"

// https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT#4.4.5
type ZipCompressionMethod uint16

const (
	ZipCompressionMethod_Store   ZipCompressionMethod = 0
	ZipCompressionMethod_Deflate ZipCompressionMethod = 8
	ZipCompressionMethod_BZIP2   ZipCompressionMethod = 12
	ZipCompressionMethod_LZMA    ZipCompressionMethod = 14
	ZipCompressionMethod_ZSTD    ZipCompressionMethod = 93
	ZipCompressionMethod_XZ      ZipCompressionMethod = 95
)

func (m ZipCompressionMethod) String() string {
	switch m {
	case ZipCompressionMethod_Store:
		return "store"
	case ZipCompressionMethod_Deflate:
		return "deflate"
	case ZipCompressionMethod_BZIP2:
		return "bzip2"
	case ZipCompressionMethod_LZMA:
		return "lzma"
	case ZipCompressionMethod_ZSTD:
		return "zstd"
	case ZipCompressionMethod_XZ:
		return "xz"
	default:
		return fmt.Sprintf("<unknown %d>", m)
	}
}
