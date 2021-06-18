package field

import (
	"reflect"
	"sync/atomic"
	"unsafe"
)

var (
	baseFieldStructType = reflect.TypeOf(BaseField{})
)

func New(f Interface) Interface {
	v := reflect.ValueOf(f)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.NumField() == 0 {
		panic("invalid empty field, BaseField is required")
	}

	firstField := v.Field(0)

	switch firstField.Type() {
	case baseFieldStructType:
	default:
		panic("invalid BaseField usage, must be first embedded struct")
	}

	var baseField *BaseField
	switch firstField.Kind() {
	case reflect.Struct:
		baseField = firstField.Addr().Interface().(*BaseField)
	default:
		panic("unexpected non struct")
	}

	if !atomic.CompareAndSwapUint32(&baseField._initialized, 0, 1) {
		return f
	}

	structType := v.Type()
	for structType.Kind() != reflect.Struct {
		structType = structType.Elem()
	}

	baseField._parentValue = reflect.NewAt(
		structType,
		unsafe.Pointer(firstField.UnsafeAddr()),
	)

	return f
}
