package tlang

import (
	"reflect"
	"runtime"
	"sync/atomic"
)

// GetLazyValue returns return value of x.GetLazyValue() when x implements that method,
// otherwise return x directly
func GetLazyValue(x reflect.Value) reflect.Value {
	if !x.IsValid() || x.IsZero() {
		return x
	}

	methodGet := x.MethodByName("GetLazyValue")
	if !methodGet.IsValid() {
		return x
	}

	ret := methodGet.Call(nil)
	if len(ret) == 0 {
		return x
	}

	return ret[0]
}

type LazyValueType[T any] interface {
	GetLazyValue() T
}

var _ LazyValueType[struct{}] = (*LazyValue[struct{}])(nil)

type LazyValue[T any] struct {
	initialized int32
	writing     int32

	Create func() T
	value  T
}

func (v *LazyValue[T]) GetLazyValue() T {
	_ = atomic.AddInt32(&v.writing, 1)

	if atomic.CompareAndSwapInt32(&v.initialized, 0, 1) {
		// I'm a writer
		// set the value
		v.value = v.Create()

		_ = atomic.AddInt32(&v.writing, -1)
	} else {
		_ = atomic.AddInt32(&v.writing, -1)

		// I'm just a reader, wait until there is no writer
		for atomic.LoadInt32(&v.writing) != 0 {
			runtime.Gosched()
		}
	}

	return v.value
}

var _ LazyValueType[string] = ImmediateString("")

type ImmediateString string

func (s ImmediateString) GetLazyValue() string { return string(s) }
