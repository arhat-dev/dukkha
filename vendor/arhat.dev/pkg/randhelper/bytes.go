package randhelper

import (
	"crypto/rand"
	"io"
)

// Bytes fill the buf with random bytes
func Bytes(buf []byte) ([]byte, error) {
	_, err := io.ReadFull(rand.Reader, buf)
	return buf, err
}
