package rs

import (
	"reflect"
	"strings"
)

type BaseField struct {
	_initialized uint32

	// _parentType and _parentValue represents the parent struct type and value
	// they are set in Init function call
	_parentType  reflect.Type
	_parentValue reflect.Value

	unresolvedFields map[string]*unresolvedFieldSpec

	opts *Options

	// yamlKey -> map value
	// TODO: separate static data and runtime generated data
	catchOtherCache map[string]reflect.Value

	catchOtherFields map[string]struct{}

	// key: yaml tag name or lower case exported field name
	fields map[string]*fieldRef

	catchOtherField *fieldRef
}

type fieldRef struct {
	fieldName  string
	fieldValue reflect.Value
	base       *BaseField

	omitempty bool

	isCatchOther bool
}

func (f *BaseField) addField(
	yamlKey, fieldName string,
	fieldValue reflect.Value,
	base *BaseField,
	omitempty bool,
) bool {
	if _, exists := f.fields[yamlKey]; exists {
		return false
	}

	f.fields[yamlKey] = &fieldRef{
		fieldName: fieldName,

		fieldValue: fieldValue,
		base:       base,

		omitempty: omitempty,
	}
	return true
}

func (f *BaseField) getField(yamlKey string) *fieldRef {
	return f.fields[yamlKey]
}

type unresolvedFieldSpec struct {
	fieldName   string
	fieldValue  reflect.Value
	rawDataList []*alterInterface
	renderers   []*suffixSpec

	isCatchOtherField bool
}

func (f *BaseField) addUnresolvedField(
	// key part
	yamlKey string,
	suffix string,

	// value part
	fieldName string,
	fieldValue reflect.Value,
	isCatchOtherField bool,
	rawData *alterInterface,
) {
	if f.unresolvedFields == nil {
		f.unresolvedFields = make(map[string]*unresolvedFieldSpec)
	}

	if isCatchOtherField {
		if f.catchOtherFields == nil {
			f.catchOtherFields = make(map[string]struct{})
		}

		f.catchOtherFields[yamlKey] = struct{}{}
	}

	if old, exists := f.unresolvedFields[yamlKey]; exists {
		// TODO: no idea how can this happen, the key suggests this can only
		// 	     happen when there are duplicate yaml keys, which is invalid yaml
		//       go-yaml should errored before we add this
		// 		 so this is considered as unreachable code

		// unreachable
		old.rawDataList = append(old.rawDataList, rawData)
		return
	}

	f.unresolvedFields[yamlKey] = &unresolvedFieldSpec{
		fieldName:   fieldName,
		fieldValue:  fieldValue,
		rawDataList: []*alterInterface{rawData},
		renderers:   parseRenderingSuffix(suffix),

		isCatchOtherField: isCatchOtherField,
	}
}

type suffixSpec struct {
	name string

	patchSpec bool
	typeHint  TypeHint
}

func parseRenderingSuffix(rs string) []*suffixSpec {
	var ret []*suffixSpec
	for _, part := range strings.Split(rs, "|") {
		size := len(part)
		if size == 0 {
			continue
		}

		spec := &suffixSpec{
			patchSpec: part[size-1] == '!',
		}

		if spec.patchSpec {
			part = part[:size-1]
			// size-- // size not used any more
		}

		if idx := strings.LastIndexByte(part, '?'); idx >= 0 {
			// TODO: do we really want to panic when type hint is not valid?
			var err error
			spec.typeHint, err = ParseTypeHint(part[idx+1:])
			if err != nil {
				panic(err)
			}

			part = part[:idx]
		}

		spec.name = part

		ret = append(ret, spec)
	}

	return ret
}

// TODO: shall we generate type hint for those without one?
// func generateTypeHintForType(typ reflect.Type) TypeHint {
// 	switch typ.Kind() {
// 	case reflect.Int, reflect.Int8, reflect.Int16,
// 		reflect.Int32, reflect.Int64:
// 		return TypeHintInt
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16,
// 		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
// 		return TypeHintInt
// 	case reflect.Float32, reflect.Float64:
// 		return TypeHintFloat
// 	case reflect.String:
// 		return TypeHintStr
// 	case reflect.Array, reflect.Slice:
// 		switch typ.Elem().Kind() {
// 		case reflect.Uint8:
// 			return TypeHintBytes
// 		default:
// 			return TypeHintObjects
// 		}
// 	case reflect.Map, reflect.Struct:
// 		return TypeHintObject
// 	default:
// 		return TypeHintNone
// 	}
// }
