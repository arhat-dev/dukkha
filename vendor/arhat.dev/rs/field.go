package rs

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

type BaseField struct {
	// fields with `_` prefix are supposed to be initialized
	// after init()

	_initialized uint32

	// _parentType and _parentValue represents the parent struct type and value
	// they are set in Init function call
	_parentType  reflect.Type
	_parentValue reflect.Value

	_opts *Options

	// normalFields are those without `rs:"other"` field tag
	// key: first value of the specified field tag namespace
	// 		or lower case version of the exported field name (go-yaml's default behavior)
	normalFields map[string]*fieldRef
	inlineMap    *fieldRef

	// key: yamlKey
	unresolvedNormalFields   map[string]*unresolvedFieldSpec
	unresolvedInlineMapItems map[string][]*unresolvedFieldSpec
}

func (f *BaseField) initialized() bool {
	return atomic.LoadUint32(&f._initialized) == 1
}

type tagSpec struct {
	yamlKey string

	inline    bool
	omitempty bool
	inlineMap bool
}

// parseFieldTags
//
// return value will be nil if the field is unexported or ignored by its data tag (e.g. `yaml:"-"`)
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
			ret.inlineMap = true
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
			ret.inlineMap = true
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
	isInlineMapItem bool
}

// addField adds one field identified by its yamlKey
// it may be a catch-other field
func (f *BaseField) addField(
	fieldName string,
	fieldValue reflect.Value,
	base *BaseField,
	ts *tagSpec,
) bool {
	if ts.inlineMap {
		if f.inlineMap != nil {
			panic(fmt.Errorf(
				"bad field tags in %s: only one map in the struct can have `rs:\"other\"` tag",
				f._parentType.String(),
			))
		}

		f.inlineMap = &fieldRef{
			fieldName:  fieldName,
			fieldValue: fieldValue,
			base:       base,

			isInlineMapItem: true,
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

		omitempty:       ts.omitempty,
		isInlineMapItem: false,
	}

	return true
}

func (f *BaseField) getField(yamlKey string) *fieldRef {
	return f.normalFields[yamlKey]
}

type unresolvedFieldSpec struct {
	// fieldName is struct field name when isInlineMapItem = false
	// otherwise it's inline map item key
	fieldName  string
	fieldValue reflect.Value
	rawData    *yaml.Node
	renderers  []*rendererSpec

	isInlineMapItem bool
}

func (f *BaseField) addUnresolvedField(
	// key part
	yamlKey string,
	suffix string,
	resolvedSuffix []*rendererSpec,

	// value part
	fieldName string,
	fieldValue reflect.Value,
	isInlineMapItem bool,
	rawData *yaml.Node,
) {
	if resolvedSuffix == nil {
		resolvedSuffix = parseRenderingSuffix(suffix)
	}

	if !isInlineMapItem {
		f.addUnresolvedNormalField(yamlKey, resolvedSuffix, fieldName, fieldValue, rawData)
		return
	}

	if f.unresolvedInlineMapItems == nil {
		f.unresolvedInlineMapItems = make(map[string][]*unresolvedFieldSpec)
	}

	f.unresolvedInlineMapItems[yamlKey] = append(f.unresolvedInlineMapItems[yamlKey], &unresolvedFieldSpec{
		fieldName:  yamlKey,
		fieldValue: f.inlineMap.fieldValue,
		rawData:    rawData,
		renderers:  resolvedSuffix,

		isInlineMapItem: true,
	})
}

func (f *BaseField) addUnresolvedNormalField(
	// key part
	yamlKey string,
	renderers []*rendererSpec,

	// value part
	fieldName string,
	fieldValue reflect.Value,
	rawData *yaml.Node,
) {
	if f.unresolvedNormalFields == nil {
		f.unresolvedNormalFields = make(map[string]*unresolvedFieldSpec)
	}

	// it can have multiple values only when it's an inline map item
	if old, exists := f.unresolvedNormalFields[yamlKey]; exists {
		old.rawData = rawData
		return
	}

	f.unresolvedNormalFields[yamlKey] = &unresolvedFieldSpec{
		fieldName:  fieldName,
		fieldValue: fieldValue,
		rawData:    rawData,
		renderers:  renderers,

		isInlineMapItem: false,
	}
}

type rendererSpec struct {
	name string

	patchSpec bool
	typeHint  TypeHint
}

func parseRenderingSuffix(rs string) []*rendererSpec {
	var ret []*rendererSpec
	for _, part := range strings.Split(rs, "|") {
		size := len(part)
		if size == 0 {
			continue
		}

		spec := &rendererSpec{
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
