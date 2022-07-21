package templateutils

import (
	"io"
	"strings"
	"time"
	"unsafe"
)

func parseArgs_MaybeOUTPUT_OBJ(args []any) (obj any, outWriter io.Writer, actualArgs []any) {
	var ok bool

	switch len(args) {
	case 0:
		return
	case 1:
	default:
		outWriter, ok = args[len(args)-2].(io.Writer)
		if ok {
			actualArgs = args[:len(args)-2]
		} else {
			outWriter = nil
			actualArgs = args[:len(args)-1]
		}
	}

	obj = args[len(args)-1]
	return
}

func parseArgs_MaybeOUTPUT_DATA(args []any) (
	inData []byte,
	inReader io.Reader,
	outWriter io.Writer,
	actualArgs []any,
	err error,
) {
	var ok bool

	switch len(args) {
	case 0:
		return
	case 1:
	default:
		outWriter, ok = args[len(args)-2].(io.Writer)
		if ok {
			actualArgs = args[:len(args)-2]
		} else {
			outWriter = nil
			actualArgs = args[:len(args)-1]
		}
	}

	inData, inReader, ok, err = toBytesOrReader(args[len(args)-1])
	if !ok {
		inReader = nil
	}

	return
}

// TODO: move all of functions below to arhat.dev/pkg

// RemoveMatchedRunesInPlace removes all matched runes in p, return new size of p, which is always <= len(p)
func RemoveMatchedRunesInPlace(p []byte, match func(r rune) bool) (n int) {
	var (
		matchStart  = int(-1)
		startRuneSZ int

		i int
		j int
		r rune
	)

	n = len(p)
	str := *(*string)(unsafe.Pointer(&p))

loop:
	if i >= n {
		return
	}

	for j, r = range str[i:n] {
		switch {
		case match(r):
			if matchStart == -1 {
				matchStart = j
			} else if startRuneSZ == 0 {
				startRuneSZ = j - matchStart
			}
		case matchStart != -1:
			if startRuneSZ == 0 {
				startRuneSZ = j - matchStart
			}

			_ = copy(p[i+matchStart:], p[i+j:n])
			n -= j - matchStart
			i += matchStart + startRuneSZ

			startRuneSZ = 0
			matchStart = -1
			goto loop
		case matchStart != -1 && startRuneSZ == 0:
			startRuneSZ = j - matchStart
		}
	}

	if matchStart != -1 {
		n -= n - (i + matchStart)
	}

	return
}

func RemoveMatchedRunesCopy(str string, match func(r rune) bool) string {
	var sb strings.Builder

	lastIdx := 0
	wasMatched := true
	for i, c := range str {
		if match(c) {
			if !wasMatched {
				sb.WriteString(str[lastIdx:i])
				wasMatched = true
			}
		} else {
			if wasMatched {
				lastIdx = i
				wasMatched = false
			}
		}
	}

	if !wasMatched {
		sb.WriteString(str[lastIdx:])
	}

	return sb.String()
}

func ReplaceMatchedRunesCopy(
	str string,
	match func(r rune) bool,
	replace func(matched string) string,
) string {
	var sb strings.Builder

	lastMismatchedIdx := len(str)
	lastMatchedIdx := 0
	wasMatched := true
	for i, c := range str {
		if match(c) {
			if !wasMatched {
				// this matches, prev mismatched,
				//
				// i is the end of latest match
				sb.WriteString(str[lastMatchedIdx:i])
				wasMatched = true
				lastMismatchedIdx = i
			}
		} else {
			if wasMatched {
				// this mismatch, prev matched
				//
				// i is the end of latest mismatch
				if lastMismatchedIdx < i {
					sb.WriteString(replace(str[lastMismatchedIdx:i]))
				}

				lastMatchedIdx = i
				wasMatched = false
			}
		}
	}

	if wasMatched {
		sb.WriteString(replace(str[lastMismatchedIdx:]))
	} else {
		sb.WriteString(str[lastMatchedIdx:])
	}

	return sb.String()
}

var _ io.Reader = (*FilterReader)(nil)

func NewFilterReader(
	underlay io.Reader,
	doFilter func(p []byte) int, // return new size of p (MUST be less than len(p))
) FilterReader {
	return FilterReader{
		underlay: underlay,
		filter:   doFilter,
	}
}

type FilterReader struct {
	underlay io.Reader

	filter func(p []byte) int
}

func (fr FilterReader) Read(p []byte) (n int, err error) {
	n, err = fr.underlay.Read(p)
	n = fr.filter(p[:n])
	return
}

func NewChunkedWriter(
	chunkSize int,
	underlay io.Writer,
	doBeforeEachChunkWriting func() error,
	doAfterEachChunkWrote func() error,
) ChunkedWriter {
	if chunkSize <= 0 {
		panic("invalid non-positive chunk size")
	}

	return ChunkedWriter{
		chunkSize: chunkSize,
		remainder: 0,

		doBeforeEachChunkWriting: doBeforeEachChunkWriting,
		doAfterEachChunkWrote:    doAfterEachChunkWrote,

		underlay: underlay,
	}
}

type ChunkedWriter struct {
	chunkSize int
	// last time wrote for next chunk
	remainder int // 0 <= remainder < chunk size

	doBeforeEachChunkWriting func() error
	doAfterEachChunkWrote    func() error

	underlay io.Writer
}

func (cw *ChunkedWriter) Remainder() int { return cw.remainder }

func (cw *ChunkedWriter) Write(p []byte) (wrote int, err error) {
	i, n, chunksize := 0, len(p), cw.chunkSize

	// handle remainder of last wrote
	if cw.remainder != 0 {
		// write i bytes fills a chunk
		i = chunksize - cw.remainder

		if i > n {
			wrote, err = cw.underlay.Write(p)
			cw.remainder += wrote
			return
		}

		wrote, err = cw.underlay.Write(p[0:i])
		if err != nil {
			cw.remainder += wrote
			if cw.remainder == chunksize {
				cw.remainder = 0
			}

			return
		}

		cw.remainder = 0
		if i == n {
			err = cw.doAfterEachChunkWrote()
			if err != nil {
				return
			}
		}
	}

	var err2 error
	wn, end := 0, i+chunksize
	for ; end < n; i, end = end, end+chunksize {
		err = cw.doBeforeEachChunkWriting()
		if err != nil {
			return
		}

		wn, err2 = cw.underlay.Write(p[i:end])
		wrote += wn

		// do after action regardless of error
		if wn == chunksize {
			err = cw.doAfterEachChunkWrote()
			if err != nil {
				// in this case, remainder is 0
				return
			}
		}

		if err2 != nil {
			cw.remainder = chunksize - wn

			err = err2
			return
		}
	}

	if i == n {
		return
	}

	// n - i < chunksize

	err = cw.doBeforeEachChunkWriting()
	if err != nil {
		return
	}

	wn, err = cw.underlay.Write(p[i:])
	wrote += wn
	cw.remainder = wn
	if wn == chunksize {
		err2 := cw.doAfterEachChunkWrote()
		if err == nil {
			err = err2
		}
	}

	return
}

type modAction uint8

const (
	mod_Round modAction = iota
	mod_Floor
	mod_Ceil
)

// ModDuration return a new duration of x that is multiple of d
func ModDuration(d, x time.Duration, action modAction) time.Duration {
	switch action {
	case mod_Floor:
		return x.Truncate(d)
	case mod_Round:
		return x.Round(d)
	case mod_Ceil:
		ret := x.Truncate(d)
		if ret != x {
			return ret + d
		}

		return ret
	default:
		panic("unreachable")
	}
}

// ModTime returns a new value base on t in timezone of the time t
//
// time.{Round, Truncat} operates on UTC time
//
// e.g. 01:01:31 in TZ -11:30 (12:31:31 in UTC)
//
// timezone offset = -11:30
//
// thus will have default behavior when dur is 1hr
//
// time.Round(12:31:31):	13:00:00 UTC -> 01:30:00 in TZ -11:30
// time.Truncat(12:31:31):	12:00:00 UTC -> 00:30:00 in TZ -11:30
//
// what actually expected to local time is
//
// floor:	01:00:00
// round:	01:00:00
// ceil:	02:00:00
func ModTime(d time.Duration, t time.Time, action modAction) (ret time.Time) {
	// offset is seconds east of UTC
	_, offSec := t.Zone()

	offset := time.Duration(offSec) * time.Second

	ret = t.Add(offset)

	switch action {
	case mod_Floor:
		return ret.Truncate(d).Add(-offset)
	case mod_Round:
		return ret.Round(d).Add(-offset)
	case mod_Ceil:
		ret = ret.Truncate(d).Add(-offset)
		if t.After(ret) {
			ret = ret.Add(d)
		}

		return
	default:
		panic("unreachable")
	}
}

func NewLazyAutoCloseReader(open func() (io.Reader, error)) LazyAutoCloseReader {
	return LazyAutoCloseReader{
		open: open,
	}
}

type readerState uint8

const (
	readerIdle readerState = iota
	readerOpened
	readerClosed
)

// LazyAutoCloseReader open underlay stream on first read call, close on error (if implements io.Closer)
type LazyAutoCloseReader struct {
	cur     io.Reader
	lastErr error

	open  func() (io.Reader, error)
	state readerState
}

func (l *LazyAutoCloseReader) Read(p []byte) (ret int, err error) {
	switch l.state {
	case readerIdle:
		l.cur, err = l.open()
		if err != nil {
			l.state = readerClosed
			return
		}

		l.state = readerOpened
		fallthrough
	case readerOpened:
		ret, err = l.cur.Read(p)
		if err != nil {
			l.lastErr = err
			l.state = readerClosed
			clo, ok := l.cur.(io.Closer)
			if ok {
				_ = clo.Close()
			}
		}
		return
	default:
		err = l.lastErr
		if err == nil {
			err = io.EOF
		}
		return
	}
}
