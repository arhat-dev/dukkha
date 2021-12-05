package cache

import (
	"io"
)

type (
	IdentifiableObject interface {
		ScopeUniqueID() string
		Ext() string
	}

	RemoteCacheRefreshFunc func(obj IdentifiableObject) (io.ReadCloser, error)
	LocalCacheRefreshFunc  func(obj IdentifiableObject) ([]byte, error)
)

type IdentifiableString string

func (s IdentifiableString) ScopeUniqueID() string {
	return string(s)
}

// Ext returns file extension
func (s IdentifiableString) Ext() string {
	data := []rune(s)
	for i := len(data) - 1; i >= 0; i-- {
		switch {
		case data[i] == '/':
			return ""
		case data[i] == '.':
			return string(data[i:])
		}
	}

	return ""
}
