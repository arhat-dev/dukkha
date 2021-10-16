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

type Options struct {
	// InterfaceTypeHandler handles interface value creation for any inner field
	// with a interface{...} type (including []interface{...} and map[string]interface{...})
	//
	// defaults to `nil`
	InterfaceTypeHandler InterfaceTypeHandler

	// TagNamespace used for struct field tag parsing
	// you can set it to json to use json tag as name source
	//
	// supported tag values are:
	// - <first tag value as data field name>
	// - inline
	// - omitempty
	//
	// unsupported tag values are ignored
	//
	// defaults to `yaml`
	TagNamespace string

	// AllowUnknownFields whether restrict unmarshaling to known fields
	//
	// NOTE: if there is a map field in struct with field type `rs:"other"`
	// even when AllowUnknownFields is set to false, it still gets these unknown
	// fields marshaled to that map field
	//
	// defaults to `false`
	AllowUnknownFields bool
}

// Init the BaseField embedded in your struct, the BaseField must be the first field
//
// 		type Foo struct {
// 			rs.BaseField // or *rs.BaseField
// 		}
//
// if the arg `in` doesn't contain BaseField or the BaseField is not the first element
// it does nothing and will return `in` as is.
// nolint:gocyclo
func Init(in Field, opts *Options) Field {
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
	baseField.opts = opts

	var tagNamespace string
	if opts != nil && len(opts.TagNamespace) != 0 {
		tagNamespace = opts.TagNamespace
	} else {
		tagNamespace = "yaml"
	}

	// get known fields for unmarshaling, skip the first field (the BaseField itself)
	baseField.fields = make(map[string]*fieldRef)
	for i := 1; i < parentType.NumField(); i++ {
		sf := parentType.Field(i)
		if len(sf.PkgPath) != 0 {
			// unexported
			continue
		}

		yTags := strings.Split(sf.Tag.Get(tagNamespace), ",")
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
		InitRecursively(fieldValue, opts)

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

				innerYamlTags := strings.Split(innerFt.Tag.Get(tagNamespace), ",")
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
func InitRecursively(fv reflect.Value, opts *Options) {
	switch fv.Type() {
	case baseFieldPtrType, baseFieldStructType:
		return
	}

	target := fv
findStruct:
	switch target.Kind() {
	case reflect.Struct:
		_ = tryInit(fv, opts)
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
		InitRecursively(target.Field(i), opts)
	}
}

// nolint:unparam
func tryInit(fieldValue reflect.Value, opts *Options) bool {
	if fieldValue.CanInterface() {
		fVal, canCallInit := fieldValue.Interface().(Field)
		if canCallInit {
			_ = Init(fVal, opts)
			return true
		}
	}

	if fieldValue.CanAddr() && fieldValue.Addr().CanInterface() {
		fVal, canCallInit := fieldValue.Addr().Interface().(Field)
		if canCallInit {
			_ = Init(fVal, opts)
			return true
		}
	}

	return false
}
