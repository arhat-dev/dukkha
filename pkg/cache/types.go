package cache

import "io"

type (
	IdentifiableObject interface {
		ScopeUniqueID() string
	}
	IdentifiableString string

	RemoteCacheRefreshFunc func(obj IdentifiableObject) (io.ReadCloser, error)
	LocalCacheRefreshFunc  func(obj IdentifiableObject) ([]byte, error)
)

func (s IdentifiableString) ScopeUniqueID() string {
	return string(s)
}
