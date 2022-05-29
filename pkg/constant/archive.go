package constant

import "fmt"

const (
	ArchiveFormat_Zip = "zip"
	ArchiveFormat_Tar = "tar"
)

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
		return CompressionMethod_DEFLATE
	case ZipCompressionMethod_BZIP2:
		return CompressionMethod_Bzip2
	case ZipCompressionMethod_LZMA:
		return CompressionMethod_LZMA
	case ZipCompressionMethod_ZSTD:
		return CompressionMethod_ZSTD
	case ZipCompressionMethod_XZ:
		return CompressionMethod_XZ
	default:
		return fmt.Sprintf("<unknown %d>", m)
	}
}
