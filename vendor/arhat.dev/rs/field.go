package rs

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
)

type BaseField struct {
	// fields with `_` prefix are supposed to be initialized
	// after init()

	_initialized uint32

	// _parentType and _parentValue represents the parent struct type and value
	// they are set in Init function call
	_parentType  reflect.Type
	_parentValue reflect.Value

	// key: yamlKey
	unresolvedFields map[string]*unresolvedFieldSpec

	_opts *Options

	// yamlKey -> map value
	// TODO: separate static data and runtime generated data
	catchOtherCache map[string]reflect.Value

	catchOtherFields map[string]struct{}

	// normalFields are those without `rs:"other"` field tag
	// key: yaml tag name or lower case exported field name
	normalFields map[string]*fieldRef

	catchOtherField *fieldRef
}

func (f *BaseField) initialized() bool {
	return atomic.LoadUint32(&f._initialized) == 1
}

type tagSpec struct {
	yamlKey string

	inline     bool
	omitempty  bool
	catchOther bool
}

// parseFieldTags
// if the field is unexported or ignored by its data tag (e.g. `yaml:"-"`)
// the return value will be nil
func (f *BaseField) parseFieldTags(sf *reflect.StructField, dataTagNS string) (*tagSpec, error) {
	if len(sf.PkgPath) != 0 {
		// unexported
		return nil, nil
	}

	yTags := strings.Split(sf.Tag.Get(dataTagNS), ",")
	yamlKey := yTags[0]

	if yamlKey == "-" {
		// ignored explicitly
		return nil, nil
	}

	if len(yamlKey) == 0 {
		yamlKey = strings.ToLower(sf.Name)
	}

	ret := &tagSpec{
		yamlKey: yamlKey,
	}

	for _, t := range yTags[1:] {
		switch t {
		case "inline":
			if sf.Type.Kind() != reflect.Map {
				ret.inline = true
				continue
			}

			// inline map, MUST have string key
			if sf.Type.Key() != stringType {
				return nil, fmt.Errorf(
					"inline option not applicable to %s.%s: "+
						"inline map MUST have string key",
					f._parentType.String(), sf.Name,
				)
			}

			// inline map is equivalent to `rs:"other"`
			ret.catchOther = true
		case "omitempty":
			ret.omitempty = true
		default:
			// TBD: currently we just ignore unknown tag, shall we panic?
		}
	}

	// rs tag is used to extend yaml tag
	for _, t := range strings.Split(sf.Tag.Get(TagNameRS), ",") {
		switch t {
		case "other":
			// other is used to match unhandled values
			// only supports map[string]Any
			ret.catchOther = true
		case "":
		default:
			return nil, fmt.Errorf(
				"unknown rs tag value %q for %s.%s",
				t, f._parentType.String(), sf.Name,
			)
		}
	}

	return ret, nil
}

func (f *BaseField) init(
	parentType reflect.Type,
	parentVal reflect.Value,
	opts *Options,
) error {
	if !atomic.CompareAndSwapUint32(&f._initialized, 0, 1) {
		// already initialized
		return nil
	}

	f._parentValue = parentVal
	f._parentType = parentType
	f._opts = opts

	var dataTagNS string
	if opts != nil && len(opts.DataTagNamespace) != 0 {
		dataTagNS = opts.DataTagNamespace
	} else {
		dataTagNS = "yaml"
	}

	// get known fields for unmarshaling, skip the first field (the BaseField itself)
	for i := 1; i < f._parentType.NumField(); i++ {
		sf := f._parentType.Field(i)

		ts, err := f.parseFieldTags(&sf, dataTagNS)
		if err != nil {
			return err
		}

		if ts == nil {
			continue
		}

		fieldValue := f._parentValue.Field(i)

		// initialize struct fields accepted by Init(), in case being used later
		// DO NOT USE tryInit, that will only init current field, which will cause
		// error when user try to resolve data not unmarshaled from yaml
		InitRecursively(fieldValue, opts)

		if !ts.inline {
			if !f.addField(sf.Name, fieldValue, f, ts) {
				return fmt.Errorf(
					"duplicate yaml key %q in struct %s.%s",
					ts.yamlKey, parentType.String(), sf.Name,
				)
			}

			continue
		}

		// handle inline fields

		err = f.collectInlineFields(&sf, fieldValue, dataTagNS)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *BaseField) collectInlineFields(
	sf *reflect.StructField,
	fieldValue reflect.Value,
	dataTagNS string,
) error {
	// try to find BaseField in inline field
	// if exists, let it manage dynamnic fields
	var innerBaseF reflect.Value
	switch kind := sf.Type.Kind(); {
	case kind == reflect.Struct:
		// embedded struct
		if reflect.PtrTo(sf.Type).Implements(fieldInterfaceType) || sf.Type.Implements(fieldInterfaceType) {
			innerBaseF = fieldValue.Field(0)
		}
	case kind == reflect.Ptr:
		typ := sf.Type
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		if typ.Kind() != reflect.Struct {
			return fmt.Errorf(
				"invalid inline tag applied to pointer of non struct %s.%s",
				f._parentType.String(), sf.Name,
			)
		}

		// TODO: fix BaseField lookup
		// embedded ptr to struct
		// if sf.Type.Implements(fieldInterfaceType) || typ.Implements(fieldInterfaceType) {
		// 	innerBaseF = fieldValue.Elem().Field(0)
		// }
	default:
		// TBD: support interface type inline?
		return fmt.Errorf(
			"invalid inline tag applied to non struct nor struct pointer field %s.%s",
			f._parentType.String(), sf.Name,
		)
	}

	base := f
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

	// using fieldValue.NumField() instead of sf.Type.NumField()
	// to support already initialized interface value
	for i := start; i < fieldValue.NumField(); i++ {
		sf := sf.Type.Field(i)

		ts, err := f.parseFieldTags(&sf, dataTagNS)
		if err != nil {
			return err
		}

		if ts == nil {
			continue
		}

		if !ts.inline {
			if !f.addField(sf.Name, fieldValue.Field(i), base, ts) {
				return fmt.Errorf(
					"duplicate yaml key %q in inline field %s.%s",
					ts.yamlKey, f._parentType.String(), sf.Name,
				)
			}

			continue
		}

		// handle inline fields

		err = f.collectInlineFields(&sf, fieldValue.Field(i), dataTagNS)
		if err != nil {
			return err
		}
	}

	return nil
}

type fieldRef struct {
	fieldName  string
	fieldValue reflect.Value
	base       *BaseField

	omitempty bool

	// this field is only set to true for fields with
	// `rs:"other"` struct field tag
	isCatchOther bool
}

// addField adds one field identified by its yamlKey
// it may be a catch-other field
func (f *BaseField) addField(
	fieldName string,
	fieldValue reflect.Value,
	base *BaseField,
	ts *tagSpec,
) bool {
	if ts.catchOther {
		if f.catchOtherField != nil {
			panic(fmt.Errorf(
				"bad field tags in %s: only one map in the struct can have `rs:\"other\"` tag",
				f._parentType.String(),
			))
		}

		f.catchOtherField = &fieldRef{
			fieldName:  fieldName,
			fieldValue: fieldValue,
			base:       base,

			isCatchOther: true,
		}

		return true
	}

	if f.normalFields == nil {
		f.normalFields = make(map[string]*fieldRef)
	}

	// handle normal field

	if _, exists := f.normalFields[ts.yamlKey]; exists {
		return false
	}

	f.normalFields[ts.yamlKey] = &fieldRef{
		fieldName: fieldName,

		fieldValue: fieldValue,
		base:       base,

		omitempty:    ts.omitempty,
		isCatchOther: false,
	}

	return true
}

func (f *BaseField) getField(yamlKey string) *fieldRef {
	return f.normalFields[yamlKey]
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
