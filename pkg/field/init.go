package field

import (
	"reflect"
	"sync/atomic"
)

var (
	baseFieldPtrType    = reflect.TypeOf(&BaseField{})
	baseFieldStructType = baseFieldPtrType.Elem()
)

// Init the BaseField embedded in your struct, the BaseField must be the first field
//
// 		type Foo struct {
// 			field.BaseField // or *field.BaseField
// 		}
//
// if the arg `in` doesn't contain BaseField or the BaseField is not the first element
// it does nothing and will return `in` as is.
func Init(in Field, h InterfaceTypeHandler) Field {
	v := reflect.ValueOf(in)
	switch v.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		// no pointer to pointer support
		v = v.Elem()

		if v.Kind() != reflect.Struct {
			// the target is not a struct, not using BaseField
			return in
		}
	default:
		return in
	}

	if !v.CanAddr() {
		panic("invalid non addressable value")
	}

	if v.NumField() == 0 {
		// empty struct, no BaseField
		return in
	}

	firstField := v.Field(0)

	var baseField *BaseField
	switch firstField.Type() {
	case baseFieldStructType:
		// using BaseField

		baseField = firstField.Addr().Interface().(*BaseField)
	case baseFieldPtrType:
		// using *BaseField

		if firstField.IsZero() {
			// not initialized
			baseField = new(BaseField)
			firstField.Set(reflect.ValueOf(baseField))
		} else {
			baseField = firstField.Interface().(*BaseField)
		}
	default:
		// BaseField is not the first field
		return in
	}

	if !atomic.CompareAndSwapUint32(&baseField._initialized, 0, 1) {
		// already initialized
		return in
	}

	baseField._parentValue = v.Addr()
	baseField.ifaceTypeHandler = h

	return in
}
