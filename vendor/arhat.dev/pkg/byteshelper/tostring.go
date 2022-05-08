package byteshelper

import "unsafe"

// ToString returns string view of byte slice s (zero allocation)
// this is different from `string(*s)`
func ToString[B ~byte, T ~[]B](s *T) string {
	return *(*string)(unsafe.Pointer(s))
}
