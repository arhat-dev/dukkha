package rs

import (
	"fmt"
	"reflect"
	"strings"
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
// nolint:gocyclo
func Init(in Field, h InterfaceTypeHandler) Field {
	parentVal := reflect.ValueOf(in)
	parentType := reflect.TypeOf(in)

	switch parentVal.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		// no pointer to pointer support
		parentVal = parentVal.Elem()
		parentType = parentType.Elem()

		if parentType.Kind() != reflect.Struct {
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
	baseField._parentType = parentType
	baseField.ifaceTypeHandler = h

	// initialize fields
	baseField.fields = make(map[string]*fieldRef)
	for i := 1; i < parentType.NumField(); i++ {
		sf := parentType.Field(i)
		if len(sf.PkgPath) != 0 {
			// unexported
			continue
		}

		yTags := strings.Split(sf.Tag.Get("yaml"), ",")
		yamlKey := yTags[0]

		if yamlKey == "-" {
			// ignored explicitly
			continue
		}

		var (
			isInlineField bool
			omitempty     bool
		)
		for _, t := range yTags[1:] {
			switch t {
			case "inline":
				isInlineField = true
			case "omitempty":
				omitempty = true
			default:
				//
			}
		}

		fieldValue := parentVal.Field(i)

		// initialize struct fields accepted by Init(), in case being used later
		// DO NOT USE tryInit, that will only init current field, which will cause
		// error when user try to resolve data not unmarshaled from yaml
		InitRecursively(fieldValue, h)

		if !isInlineField {
			if len(yamlKey) == 0 {
				yamlKey = strings.ToLower(sf.Name)
			}

			if !baseField.addField(
				yamlKey, sf.Name, fieldValue, baseField, omitempty,
			) {
				panic(fmt.Errorf(
					"duplicate yaml key %q for %s.%s",
					yamlKey, parentType.String(), sf.Name,
				))
			}
		} else {
			// handle inline fields

			// find BaseField in inline field, if exists, let it manage dynamnic fields
			base := baseField

			var innerBaseF reflect.Value
			switch kind := sf.Type.Kind(); {
			case kind == reflect.Struct:
				// embedded struct
				if sf.Type.Implements(fieldInterfaceType) || reflect.PtrTo(sf.Type).Implements(fieldInterfaceType) {
					innerBaseF = fieldValue.Field(0)
				}
			case kind == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Struct:
				// embedded ptr to struct
				if sf.Type.Implements(fieldInterfaceType) || sf.Type.Elem().Implements(fieldInterfaceType) {
					innerBaseF = fieldValue.Elem().Field(0)
				}
			default:
				panic(fmt.Errorf(
					"invalid inline tag applied to non struct nor struct pointer field %s.%s",
					parentType.String(), sf.Name,
				))
			}

			start := 0
			switch innerBaseF.Kind() {
			case reflect.Struct:
				if innerBaseF.Addr().Type() == baseFieldPtrType {
					base = innerBaseF.Addr().Interface().(*BaseField)
					start = 1
				}
			case reflect.Ptr:
				if innerBaseF.Type() == baseFieldPtrType {
					base = innerBaseF.Interface().(*BaseField)
					start = 1
				}
			}

		inlineFieldLoop:
			for j := start; j < fieldValue.NumField(); j++ {
				innerFt := sf.Type.Field(j)

				if len(innerFt.PkgPath) != 0 {
					// unexported
					continue
				}

				innerYamlTags := strings.Split(innerFt.Tag.Get("yaml"), ",")
				innerYamlKey := innerYamlTags[0]

				if innerYamlKey == "-" {
					// ignored explicitly
					continue
				}

				for _, tag := range innerYamlTags[1:] {
					switch tag {
					case "inline":
						// already in a inline struct, do not check inline anymore
						continue inlineFieldLoop
					default:
						// TODO: handle other yaml tags
					}
				}

				if len(innerYamlKey) == 0 {
					innerYamlKey = strings.ToLower(innerFt.Name)
				}

				if !baseField.addField(
					innerYamlKey, innerFt.Name, fieldValue.Field(j),
					base, omitempty,
				) {
					panic(fmt.Errorf(
						"duplicate yaml key %q in inline field %s of %s",
						innerYamlKey, innerFt.Name, parentType.String(),
					))
				}
			}
		}

		// rs tag is used to extend yaml tag
		for _, t := range strings.Split(sf.Tag.Get(TagNameRS), ",") {
			switch t {
			case "other":
				// other is used to match unhandled values
				// only supports map[string]Any or []Any

				if baseField.catchOtherField != nil {
					panic(fmt.Errorf(
						"bad field tags in %s: only one map in a struct can have `rs:\"other\"` tag",
						parentType.String(),
					))
				}

				baseField.catchOtherField = &fieldRef{
					fieldName:  sf.Name,
					fieldValue: fieldValue,
					base:       baseField,

					isCatchOther: true,
				}
			case "":
			default:
				panic(fmt.Errorf("rs: unknown rs tag value %q", t))
			}
		}
	}

	return in
}

// InitRecursively trys to call Init on all fields implementing Field interface
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
