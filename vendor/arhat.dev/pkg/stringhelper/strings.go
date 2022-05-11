package stringhelper

import (
	"strings"
	"unsafe"
)

func HasPrefix[B1, B2 ~byte, S1 String[B1], S2 String[B2]](s S1, prefix S2) bool {
	return strings.HasPrefix(*(*string)(unsafe.Pointer(&s)), *(*string)(unsafe.Pointer(&prefix)))
}

func HasSuffix[B1, B2 ~byte, S1 String[B1], S2 String[B2]](s S1, suffix S2) bool {
	return strings.HasSuffix(*(*string)(unsafe.Pointer(&s)), *(*string)(unsafe.Pointer(&suffix)))
}

func TrimPrefix[B1, B2 ~byte, S1 String[B1], S2 String[B2]](s S1, prefix S2) S1 {
	if HasPrefix[B1, B2](s, prefix) {
		return SliceStart[B1](s, len(prefix))
	}

	return s
}

func TrimSuffix[B1, B2 ~byte, S1 String[B1], S2 String[B2]](s S1, suffix S2) S1 {
	if HasSuffix[B1, B2](s, suffix) {
		return SliceEnd[B1](s, len(s)-len(suffix))
	}

	return s
}
