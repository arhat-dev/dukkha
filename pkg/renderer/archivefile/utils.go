package archivefile

import (
	"io"
	"math"
	"path"
	"strings"

	"arhat.dev/pkg/iohelper"
)

type sizeIface interface {
	Size() int64
}

type archiveSource interface {
	io.ReadSeekCloser
	io.ReaderAt
	sizeIface
}

func newArchiveSource(parent archiveSource, r io.Reader) archiveSource {
	readerAt := newBufferedReaderAt(r)

	closeImpl := parent.Close
	if s, ok := r.(io.Closer); ok {
		closeImpl = func() error {
			_ = s.Close()
			return parent.Close()
		}
	}

	sR := io.NewSectionReader(readerAt, 0, math.MaxInt64)
	type result struct {
		io.Seeker
		io.ReaderAt
		io.ReadCloser
		sizeIface
	}

	return &result{
		sizeIface:  readerAt,
		Seeker:     sR,
		ReaderAt:   sR,
		ReadCloser: iohelper.CustomReadCloser(sR, closeImpl),
	}
}

func prepareSeekRestore(r io.Seeker) (func() error, error) {
	currentAt, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	return func() error {
		_, err := r.Seek(currentAt, io.SeekStart)
		return err
	}, nil
}

func newBufferedReaderAt(r io.Reader) SizedReaderAt {
	return &bufferredReaderAt{
		upstream: r,
		buf:      make([]byte, 0, 2048),
	}
}

var (
	_ SizedReaderAt = (*bufferredReaderAt)(nil)
)

type bufferredReaderAt struct {
	upstream io.Reader
	// offset currently at
	buf []byte
}

func (s *bufferredReaderAt) Size() int64 {
	return int64(len(s.buf))
}

func (s *bufferredReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	toRead := len(p)

	switch size := int64(len(s.buf)); {
	case off < size:
		// we have alreay read some of the requested data
		n = copy(p, s.buf[off:])
		toRead -= n
		// DO NOT return early if toRead is zero, we should read zero
		// byte from upstream to check its liveness
	case off > size:
		// we need to make up some progress
		if maxCap := int64(cap(s.buf)); maxCap < off {
			// grow the buff
			s.buf = append(s.buf, make([]byte, off-maxCap)...)
		}

		_, err = io.ReadFull(s.upstream, s.buf[size:off])
		if err != nil {
			return 0, io.ErrUnexpectedEOF
		}
		s.buf = s.buf[:off]
	}

	nUpstream, err := s.upstream.Read(p[n : n+toRead])
	if ret := n + nUpstream; ret > 0 {
		s.buf = append(s.buf, p[n:ret]...)
		return ret, nil
	}

	return 0, err
}

func cleanLink(currentFile, linkName string) string {
	linkName = path.Clean(linkName)
	if path.IsAbs(linkName) {
		return linkName
	}

	var (
		upperDir  string
		remainder = linkName
		actualDir = path.Join(path.Dir(currentFile), ".")
	)

	idx := strings.IndexByte(remainder, '/')
	for idx != -1 && actualDir != "." {
		upperDir, remainder = remainder[:idx], remainder[idx+1:]

		if upperDir != ".." {
			break
		}

		idx = strings.IndexByte(remainder, '/')
		actualDir = path.Dir(actualDir)
	}

	if actualDir == "." {
		return remainder
	}

	if strings.HasPrefix(linkName, "..") {
		return path.Join(currentFile, linkName)
	}

	return path.Join(path.Dir(currentFile), linkName)
}
