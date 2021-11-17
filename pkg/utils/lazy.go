package utils

import (
	"reflect"
	"runtime"
	"sync/atomic"
)

func GetLazyValue(x reflect.Value) reflect.Value {
	if x.IsValid() && x.CanInterface() {
		lv, ok := x.Interface().(LazyValue)
		if ok {
			return reflect.ValueOf(lv.Get())
		}
	}

	return x
}

type LazyValue interface {
	_private()

	Get() string
}

type ImmediateString string

func (s ImmediateString) _private()   {}
func (s ImmediateString) Get() string { return string(s) }

func NewLazyValue(create func() string) LazyValue {
	return &LazyValueImpl{
		initialized: 0,
		writing:     0,

		create: create,
		value:  "",
	}
}

type LazyValueImpl struct {
	initialized int32
	writing     int32

	create func() string
	value  string
}

func (s *LazyValueImpl) _private() {}
func (v *LazyValueImpl) Get() string {
	_ = atomic.AddInt32(&v.writing, 1)

	if atomic.CompareAndSwapInt32(&v.initialized, 0, 1) {
		// I'm a writer
		// set the value
		v.value = v.create()

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
