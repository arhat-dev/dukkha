package stringhelper

import (
	"reflect"
	"strings"
	"unsafe"
)

type String[T ~byte] interface {
	~string | ~[]T
}

type ByteString String[byte]

// Convert String to typed string
func Convert[R ~string, B ~byte, S String[B]](s S) R {
	return *(*R)(unsafe.Pointer(&s))
}

// Append string data to typed string
func Append[B ~byte, S1 ~string, S2 String[B]](s S1, more ...S2) S1 {
	var sb strings.Builder
	sb.WriteString(*(*string)(unsafe.Pointer(&s)))
	for _, m := range more {
		sb.WriteString(*(*string)(unsafe.Pointer(&m)))
	}

	ret := sb.String()
	return *(*S1)(unsafe.Pointer(&ret))
}

// SliceStart returns immutable s[i:]
func SliceStart[B ~byte, S String[B]](s S, i int) S {
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))

	ret := reflect.SliceHeader{
		Data: p.Data + uintptr(i),
		Len:  p.Len - i,
	}
	ret.Cap = ret.Len

	return *(*S)(unsafe.Pointer(&ret))
}

// SliceEnd returns immutable s[:i]
func SliceEnd[B ~byte, S String[B]](s S, i int) S {
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))

	ret := reflect.SliceHeader{
		Data: p.Data,
		Len:  i,
	}
	ret.Cap = ret.Len

	return *(*S)(unsafe.Pointer(&ret))
}

// SliceStartEnd returns immutable s[i:j]
func SliceStartEnd[B ~byte, S String[B]](s S, i, j int) S {
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))

	ret := reflect.SliceHeader{
		Data: p.Data + uintptr(i),
		Len:  j - i,
		Cap:  j - i,
	}

	return *(*S)(unsafe.Pointer(&ret))
}
