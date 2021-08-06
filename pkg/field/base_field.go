package field

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

type (
	unresolvedFieldKey struct {
		// NOTE: put `suffix` and `yamlKey` in key is to support fields with
		// 		 `dukkha:"other"` field tag, each item should be able
		// 		 to have its own renderer
		yamlKey string
		suffix  string
	}

	unresolvedFieldValue struct {
		fieldName  string
		fieldValue reflect.Value
		rawData    []*alterInterface
		renderers  []string

		isCatchOtherField bool
	}
)

// alterInterface is a direct `interface{}`` replacement for data unmarshaling
// with no support of rendering suffix
type alterInterface struct {
	mapData map[string]*alterInterface
	arrData []*alterInterface

	scalarData interface{}
}

func (f *alterInterface) Value() interface{} {
	if f.mapData != nil {
		return f.mapData
	}

	if f.arrData != nil {
		return f.arrData
	}

	return f.scalarData
}

func (f *alterInterface) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		switch n.ShortTag() {
		case "!!str":
			f.scalarData = n.Value
		case "!!binary":
			f.scalarData = n.Value
		default:
			return n.Decode(&f.scalarData)
		}

		return nil
	case yaml.MappingNode:
		f.mapData = make(map[string]*alterInterface)
		return n.Decode(&f.mapData)
	case yaml.SequenceNode:
		f.arrData = make([]*alterInterface, 0)
		return n.Decode(&f.arrData)
	default:
		return fmt.Errorf("unexpected node tag %q", n.ShortTag())
	}
}

func (f *alterInterface) MarshalYAML() (interface{}, error) {
	return f.Value(), nil
}

type BaseField struct {
	_initialized uint32

	// _parentValue is always a pointer type with .Elem() to the struct
	// when initialized
	_parentValue reflect.Value

	unresolvedFields map[unresolvedFieldKey]*unresolvedFieldValue

	ifaceTypeHandler InterfaceTypeHandler

	// yamlKey -> map value
	// TODO: separate static data and runtime generated data
	catchOtherCache map[string]reflect.Value

	catchOtherFields map[string]struct{}
}

func (f *BaseField) Inherit(b *BaseField) error {
	if len(b.unresolvedFields) == 0 {
		return nil
	}

	if f.unresolvedFields == nil {
		f.unresolvedFields = make(map[unresolvedFieldKey]*unresolvedFieldValue)
	}

	for k, v := range b.unresolvedFields {
		existingV, ok := f.unresolvedFields[k]
		if !ok {
			f.unresolvedFields[k] = &unresolvedFieldValue{
				fieldName:         v.fieldName,
				fieldValue:        f._parentValue.Elem().FieldByName(v.fieldName),
				isCatchOtherField: v.isCatchOtherField,
				rawData:           v.rawData,
				renderers:         v.renderers,
			}

			continue
		}

		switch {
		case existingV.fieldName != v.fieldName,
			existingV.isCatchOtherField != v.isCatchOtherField:
			return fmt.Errorf(
				"invalid not match for same target, want %q, got %q",
				existingV.fieldName, v.fieldName,
			)
		}

		existingV.rawData = append(existingV.rawData, v.rawData...)
	}

	if len(b.catchOtherCache) != 0 {
		if f.catchOtherCache == nil {
			f.catchOtherCache = make(map[string]reflect.Value)
		}

		for k, v := range b.catchOtherCache {
			f.catchOtherCache[k] = v
		}
	}

	if len(b.catchOtherFields) != 0 {
		if f.catchOtherFields == nil {
			f.catchOtherFields = make(map[string]struct{})
		}

		for k, v := range b.catchOtherFields {
			f.catchOtherFields[k] = v
		}
	}

	return nil
}

// UnmarshalYAML handles parsing of rendering suffix and normal yaml
// unmarshaling
// nolint:gocyclo
func (f *BaseField) UnmarshalYAML(n *yaml.Node) error {
	if atomic.LoadUint32(&f._initialized) == 0 {
		return fmt.Errorf("field unmarshal: struct not intialized with Init()")
	}

	type fieldKey struct {
		yamlKey string
	}

	type fieldSpec struct {
		fieldName  string
		fieldValue reflect.Value
		base       *BaseField

		isCatchOther bool
	}

	fields := make(map[fieldKey]*fieldSpec)
	pt := f._parentValue.Type().Elem()

	addField := func(
		yamlKey, fieldName string,
		fieldValue reflect.Value,
		base *BaseField,
	) bool {
		key := fieldKey{yamlKey: yamlKey}
		if _, exists := fields[key]; exists {
			return false
		}

		fields[key] = &fieldSpec{
			fieldName: fieldName,

			fieldValue: fieldValue,
			base:       base,
		}
		return true
	}

	getField := func(yamlKey string) *fieldSpec {
		return fields[fieldKey{
			yamlKey: yamlKey,
		}]
	}

	var catchOtherField *fieldSpec

	// get expected fields first, skip the first field (the BaseField itself)
fieldLoop:
	for i := 1; i < pt.NumField(); i++ {
		fieldType := pt.Field(i)
		fieldValue := f._parentValue.Elem().Field(i)

		// initialize struct fields accepted by Init(), in case being used later
		InitRecursively(fieldValue, f.ifaceTypeHandler)

		yTags := strings.Split(fieldType.Tag.Get("yaml"), ",")

		// check if ignored
		for _, t := range yTags {
			if t == "-" {
				// ignored
				continue fieldLoop
			}
		}

		// get yaml field name
		yamlKey := yTags[0]
		if len(yamlKey) != 0 {
			if !addField(yamlKey, fieldType.Name, fieldValue, f) {
				return fmt.Errorf(
					"field: duplicate yaml key %q in %s",
					yamlKey, pt.String(),
				)
			}
		}

		// process yaml tag flags
		for _, t := range yTags[1:] {
			switch t {
			case "inline":
				kind := fieldType.Type.Kind()
				switch {
				case kind == reflect.Struct:
				case kind == reflect.Ptr && fieldType.Type.Elem().Kind() == reflect.Struct:
				default:
					return fmt.Errorf(
						"field: inline tag applied to non struct nor pointer field %s.%s",
						pt.String(), fieldType.Name,
					)
				}

				inlineFv := fieldValue
				inlineFt := f._parentValue.Type().Elem().Field(i).Type

				var iface interface{}
				switch inlineFv.Kind() {
				case reflect.Ptr:
					iface = inlineFv.Interface()
				default:
					iface = inlineFv.Addr().Interface()
				}

				base := f
				fVal, canCallInit := iface.(Field)
				if canCallInit {
					innerBaseF := reflect.ValueOf(
						Init(fVal, base.ifaceTypeHandler),
					).Elem().Field(0)

					if innerBaseF.Kind() == reflect.Struct {
						if innerBaseF.Addr().Type() == baseFieldPtrType {
							base = innerBaseF.Addr().Interface().(*BaseField)
						}
					} else {
						if innerBaseF.Type() == baseFieldPtrType {
							base = innerBaseF.Interface().(*BaseField)
						}
					}
				}

				for j := 0; j < inlineFv.NumField(); j++ {
					innerFv := inlineFv.Field(j)
					innerFt := inlineFt.Field(j)

					innerYamlKey := strings.Split(innerFt.Tag.Get("yaml"), ",")[0]
					if innerYamlKey == "-" {
						continue
					}

					if len(innerYamlKey) == 0 {
						// already in a inline field, do not check inline anymore
						continue
					}

					if !addField(innerYamlKey, innerFt.Name, innerFv, base) {
						return fmt.Errorf(
							"field: duplicate yaml key %q in inline field %s of %s",
							innerYamlKey, innerFt.Name, pt.String(),
						)
					}
				}
			default:
				// TODO: handle other yaml tag flags
			}
		}

		// dukkha tag is used to extend yaml tag
		dTags := strings.Split(fieldType.Tag.Get("dukkha"), ",")
		for _, t := range dTags {
			switch t {
			case "other":
				// other is used to match unhandled values
				// only supports map[string]Any

				if catchOtherField != nil {
					return fmt.Errorf(
						"field: bad field tags in %s: only one map in a struct can have `dukkha:\"other\"` tag",
						pt.String(),
					)
				}

				catchOtherField = &fieldSpec{
					fieldName:  fieldType.Name,
					fieldValue: fieldValue,
					base:       f,

					isCatchOther: true,
				}
			case "":
			default:
				return fmt.Errorf("field: unknown dukkha tag value %q", t)
			}
		}
	}

	switch n.Kind {
	case yaml.MappingNode:
	default:
		return fmt.Errorf("field: unsupported yaml tag %q when handling %s", n.Tag, pt.String())
	}

	m := make(map[string]*alterInterface)
	err := n.Decode(&m)
	if err != nil {
		return fmt.Errorf("field: data unmarshal failed for %s: %w", pt.String(), err)
	}

	handledYamlValues := make(map[string]struct{})
	// handle rendering suffix
	for rawYamlKey, v := range m {
		yamlKey := rawYamlKey

		parts := strings.SplitN(rawYamlKey, "@", 2)
		if len(parts) == 1 {
			// no rendering suffix, fill value

			if _, ok := handledYamlValues[yamlKey]; ok {
				return fmt.Errorf(
					"field: duplicate yaml field name %q",
					yamlKey,
				)
			}

			handledYamlValues[yamlKey] = struct{}{}

			fSpec := getField(yamlKey)
			if fSpec == nil {
				if catchOtherField == nil {
					return fmt.Errorf("field: unknown yaml field %q for %s", yamlKey, pt.String())
				}

				fSpec = catchOtherField
				v = &alterInterface{
					mapData: map[string]*alterInterface{
						yamlKey: v,
					},
				}
			}

			err = f.unmarshal(
				yamlKey, v, fSpec.fieldValue, fSpec.isCatchOther,
			)
			if err != nil {
				return fmt.Errorf(
					"field: failed to unmarshal yaml field %q to struct field %q: %w",
					yamlKey, fSpec.fieldName, err,
				)
			}

			continue
		}

		// has rendering suffix

		yamlKey, suffix := parts[0], parts[1]

		if _, ok := handledYamlValues[yamlKey]; ok {
			return fmt.Errorf(
				"field: duplicate yaml field name %q, please note"+
					" rendering suffix won't change the field name",
				yamlKey,
			)
		}

		fSpec := getField(yamlKey)
		if fSpec == nil {
			if catchOtherField == nil {
				return fmt.Errorf("field: unknown yaml field %q for %s", yamlKey, pt.String())
			}

			fSpec = catchOtherField
			v = &alterInterface{
				mapData: map[string]*alterInterface{
					yamlKey: v,
				},
			}
		}

		// do not unmarshal now, we will render the value
		// and unmarshal at runtime

		handledYamlValues[yamlKey] = struct{}{}
		// don't forget the raw name with rendering suffix
		handledYamlValues[rawYamlKey] = struct{}{}

		err = fSpec.base.addUnresolvedField(
			yamlKey, suffix,
			fSpec.fieldName, fSpec.fieldValue, fSpec.isCatchOther,
			v,
		)
		if err != nil {
			return fmt.Errorf("field: failed to add unresolved field: %w", err)
		}
	}

	for k := range handledYamlValues {
		delete(m, k)
	}

	if len(m) == 0 {
		// all values consumed
		return nil
	}

	var unknownFields []string
	for k := range m {
		unknownFields = append(unknownFields, k)
	}
	sort.Strings(unknownFields)

	return fmt.Errorf(
		"field: unknown yaml fields for %s: %s",
		pt.String(), strings.Join(unknownFields, ", "),
	)
}

var (
	rawInterfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
)

func (f *BaseField) unmarshal(
	yamlKey string,
	in *alterInterface,
	outVal reflect.Value,
	keepOld bool,
) error {
	for outVal.Kind() == reflect.Ptr {
		// we are trying to set value of it, so init it when not set
		if outVal.IsZero() {
			outVal.Set(reflect.New(outVal.Type().Elem()))
		}

		outVal = outVal.Elem()
	}

	switch kind := outVal.Kind(); kind {
	case reflect.Array:
		return f.unmarshalArray(yamlKey, in, outVal)
	case reflect.Slice:
		return f.unmarshalSlice(yamlKey, in, outVal, keepOld)
	case reflect.Map:
		return f.unmarshalMap(yamlKey, in, outVal, keepOld)
	case reflect.Struct:
		return f.unmarshalRaw(in, outVal)
	case reflect.Interface:
		handled, err := f.unmarshalInterface(yamlKey, in, outVal, keepOld)
		if !handled {
			// using default go-yaml behavior
			return f.unmarshalRaw(in, outVal)
		}

		return err
	default:
		// scalar types
		switch t := in.Value().(type) {
		case string:
			if kind == reflect.String {
				outVal.SetString(t)
				return nil
			}

			return f.unmarshalRaw(in, outVal)
		default:
			outVal.Set(reflect.ValueOf(in.Value()))
		}
		return nil
	}
}

func (f *BaseField) unmarshalRaw(in *alterInterface, outVal reflect.Value) error {
	var out interface{}
	if outVal.Kind() != reflect.Ptr {
		out = outVal.Addr().Interface()
	} else {
		out = outVal.Interface()
	}

	fVal, canCallInit := out.(Field)
	if canCallInit {
		_ = Init(fVal, f.ifaceTypeHandler)
	}

	dataBytes, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("field: failed to marshal back yaml for %q: %w", outVal.String(), err)
	}

	err = yaml.Unmarshal(dataBytes, out)
	if err != nil {
		return err
	}

	_ = tryInit(outVal, f.ifaceTypeHandler)
	return nil
}

func (f *BaseField) unmarshalInterface(
	yamlKey string,
	in *alterInterface,
	outVal reflect.Value,
	keepOld bool,
) (bool, error) {
	if f.ifaceTypeHandler == nil {
		// use default behavior for interface{} types
		return false, nil
	}

	fVal, err := f.ifaceTypeHandler.Create(outVal.Type(), yamlKey)
	if err != nil {
		if errors.Is(err, ErrInterfaceTypeNotHandled) && outVal.Type() == rawInterfaceType {
			// no type information provided, decode using go-yaml directly
			return false, nil
		}

		return true, fmt.Errorf("failed to create interface field: %w", err)
	}

	val := reflect.ValueOf(fVal)
	if outVal.CanSet() {
		outVal.Set(val)
	} else {
		outVal.Elem().Set(val)
	}

	// DO NOT use outVal directly, which will always match reflect.Interface
	return true, f.unmarshal(yamlKey, in, val, keepOld)
}

func (f *BaseField) unmarshalArray(yamlKey string, in *alterInterface, outVal reflect.Value) error {
	if in.arrData == nil && in.Value() != nil {
		return fmt.Errorf("unexpected non array data for %q", outVal.String())
	}

	size := len(in.arrData)
	if size != outVal.Len() {
		return fmt.Errorf(
			"array size not match for %q: want %d got %d",
			outVal.String(), outVal.Len(), size,
		)
	}

	for i := 0; i < size; i++ {
		itemVal := outVal.Index(i)

		err := f.unmarshal(
			yamlKey, in.arrData[i], itemVal,
			// always drop existing inner data
			// (actually doesn't matter since it's new)
			false,
		)
		if err != nil {
			return fmt.Errorf("failed to unmarshal slice item %s: %w", itemVal.Type().String(), err)
		}
	}

	return nil
}

func (f *BaseField) unmarshalSlice(yamlKey string, in *alterInterface, outVal reflect.Value, keepOld bool) error {
	if in.arrData == nil && in.Value() != nil {
		return fmt.Errorf("unexpected non slice data for %q", outVal.String())
	}

	size := len(in.arrData)
	sliceVal := reflect.MakeSlice(outVal.Type(), size, size)

	for i := 0; i < size; i++ {
		itemVal := sliceVal.Index(i)

		err := f.unmarshal(
			yamlKey, in.arrData[i], itemVal,
			// always drop existing inner data
			// (actually doesn't matter since it's new)
			false,
		)
		if err != nil {
			return fmt.Errorf("failed to unmarshal slice item %s: %w", itemVal.Type().String(), err)
		}
	}

	if outVal.IsZero() || !keepOld {
		outVal.Set(sliceVal)
	} else {
		outVal.Set(reflect.AppendSlice(outVal, sliceVal))
	}

	return nil
}

func (f *BaseField) unmarshalMap(yamlKey string, in *alterInterface, outVal reflect.Value, keepOld bool) error {
	if in.mapData == nil && in.Value() != nil {
		return fmt.Errorf("unexpected non map value for %q", outVal.String())
	}

	// map key MUST be string
	if outVal.IsNil() || !keepOld {
		outVal.Set(reflect.MakeMap(outVal.Type()))
	}

	valType := outVal.Type().Elem()

	isCatchOtherField := f.isCatchOtherField(yamlKey)
	if isCatchOtherField {
		if f.catchOtherCache == nil {
			f.catchOtherCache = make(map[string]reflect.Value)
		}
	}

	for k, v := range in.mapData {
		// since indexed map value is not addressable
		// we have to keep the original data in BaseField cache

		// TODO: test keepOld behavior with cached data

		var val reflect.Value

		if isCatchOtherField {
			cachedData, ok := f.catchOtherCache[k]
			if ok {
				val = cachedData
			} else {
				val = reflect.New(valType)
				f.catchOtherCache[k] = val
			}
		} else {
			val = reflect.New(valType)
		}

		err := f.unmarshal(
			// use k rather than `yamlKey`
			// because it can be the field catching other
			// (field tag: `dukkha:"other"`)
			k,
			v,
			val,
			keepOld,
		)
		if err != nil {
			return fmt.Errorf("failed to unmarshal map value %s for key %q: %w",
				valType.String(), k, err,
			)
		}

		outVal.SetMapIndex(reflect.ValueOf(k), val.Elem())
	}

	return nil
}
