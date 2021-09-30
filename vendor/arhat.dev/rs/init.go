package rs

import (
	"reflect"
	"sync/atomic"
)

var (
	baseFieldPtrType    = reflect.TypeOf((*BaseField)(nil))
	baseFieldStructType = baseFieldPtrType.Elem()
)

// Init the BaseField embedded in your struct, the BaseField must be the first field
//
// 		type Foo struct {
// 			rs.BaseField // or *rs.BaseField
// 		}
//
// if the arg `in` doesn't contain BaseField or the BaseField is not the first element
// it does nothing and will return `in` as is.
func Init(in Field, h InterfaceTypeHandler) Field {
	parentVal := reflect.ValueOf(in)
	switch parentVal.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		// no pointer to pointer support
		parentVal = parentVal.Elem()

		if parentVal.Kind() != reflect.Struct {
			// the target is not a struct, not using BaseField
			return in
		}
	default:
		return in
	}

	if !parentVal.CanAddr() {
		panic("invalid non addressable value")
	}

	if parentVal.NumField() == 0 {
		// empty struct, no BaseField
		return in
	}

	firstField := parentVal.Field(0)

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

	baseField._parentValue = parentVal
	baseField.ifaceTypeHandler = h

	return in
}

func InitRecursively(fv reflect.Value, h InterfaceTypeHandler) {
	switch fv.Type() {
	case baseFieldPtrType, baseFieldStructType:
		return
	}

	target := fv
findStruct:
	switch target.Kind() {
	case reflect.Struct:
		_ = tryInit(fv, h)
	case reflect.Ptr:
		if !target.IsValid() || target.IsZero() || target.IsNil() {
			return
		}

		target = target.Elem()
		goto findStruct
	default:
		return
	}

	for i := 0; i < target.NumField(); i++ {
		InitRecursively(target.Field(i), h)
	}
}

// nolint:unparam
func tryInit(fieldValue reflect.Value, h InterfaceTypeHandler) bool {
	if fieldValue.CanInterface() {
		fVal, canCallInit := fieldValue.Interface().(Field)
		if canCallInit {
			_ = Init(fVal, h)
			return true
		}
	}

	if fieldValue.CanAddr() && fieldValue.Addr().CanInterface() {
		fVal, canCallInit := fieldValue.Addr().Interface().(Field)
		if canCallInit {
			_ = Init(fVal, h)
			return true
		}
	}

	return false
}
