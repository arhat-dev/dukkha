package af

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type alwaysFailSeeker struct {
	io.Reader
	err error
}

func (s *alwaysFailSeeker) Seek(_ int64, _ int) (int64, error) {
	return -1, s.err
}

func NewAlwaysFailSeeker(r io.Reader, err error) io.ReadSeeker {
	if err == nil {
		err = io.ErrUnexpectedEOF
	}

	return &alwaysFailSeeker{
		Reader: r,
		err:    err,
	}
}

func TestPrepareSeekRestore(t *testing.T) {
	t.Parallel()

	t.Run("Reader Error", func(t *testing.T) {
		_, err := prepareSeekRestore(NewAlwaysFailSeeker(nil, io.ErrClosedPipe))
		assert.ErrorIs(t, err, io.ErrClosedPipe)
	})

	t.Run("Restore Ok", func(t *testing.T) {
		const (
			testData = "test-data"
			original = 3
		)
		r := io.NewSectionReader(bytes.NewReader([]byte(testData)), original, 1024)
		do := func() {
			restore, err := prepareSeekRestore(r)
			assert.NoError(t, err)

			data, err := io.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, len(testData)-original, len(data))
			assert.EqualValues(t, testData[original:], string(data))

			assert.NoError(t, restore())
		}

		do()
		do()
		do()
	})

}

func TestBufferedReaderAt_ReadAt(t *testing.T) {
	t.Parallel()

	const (
		testData = "test-data"
	)

	t.Run("Offset 0 - Full", func(t *testing.T) {
		const offset = 0
		r := newBufferedReaderAt(strings.NewReader(testData)).(*bufferredReaderAt)

		data := make([]byte, len(testData))
		n, err := r.ReadAt(data, offset)
		assert.NoError(t, err)
		assert.Equal(t, len(testData), n)

		assert.EqualValues(t, testData[offset:], string(data))
		assert.EqualValues(t, testData, string(r.buf))
	})

	t.Run("Offset 5 - Full", func(t *testing.T) {
		const offset = 5
		r := newBufferedReaderAt(strings.NewReader(testData)).(*bufferredReaderAt)

		data := make([]byte, len(testData))
		n, err := r.ReadAt(data, offset)
		assert.NoError(t, err)
		assert.Equal(t, len(testData)-offset, n)

		data = data[:n]
		assert.EqualValues(t, testData[offset:], string(data))
		assert.Len(t, r.buf, len(testData))
		assert.EqualValues(t, testData, string(r.buf[:len(r.buf)]))
	})

	t.Run("Read Buffered", func(t *testing.T) {
		const offset = 5
		r := newBufferedReaderAt(strings.NewReader(testData)).(*bufferredReaderAt)

		data := make([]byte, len(testData))
		n, err := r.ReadAt(data, offset)
		assert.NoError(t, err)
		assert.Equal(t, len(testData)-offset, n)

		data = data[:n]
		assert.EqualValues(t, testData[offset:], string(data))
		assert.EqualValues(t, testData, string(r.buf))

		// read from 0, should read from buf directly and read zero byte from upstream reader
		data = make([]byte, len(testData))
		n, err = r.ReadAt(data, 0)
		assert.NoError(t, err)
		assert.Equal(t, len(testData), n)

		data = data[:n]
		assert.EqualValues(t, testData, string(data))
		assert.EqualValues(t, testData, string(r.buf))

		// read from 1, should read from buf directly and read 1 byte from upstream reader
		// and should not error
		data = make([]byte, len(testData))
		n, err = r.ReadAt(data, 1)
		assert.NoError(t, err)
		assert.Equal(t, len(testData)-1, n)

		data = data[:n]
		assert.EqualValues(t, testData[1:], string(data))
		assert.EqualValues(t, testData, string(r.buf))
	})
}
