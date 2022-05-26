package byteshelper

import "unsafe"

// ToString returns string view of byte slice s (zero allocation)
// this is different from `string(*s)`
func ToString[B ~byte](s []B) string {
	return *(*string)(unsafe.Pointer(&s))
}
