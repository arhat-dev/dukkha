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

// nolint:revive
func (self *BaseField) Inherit(b *BaseField) error {
	if len(b.unresolvedFields) == 0 {
		return nil
	}

	if self.unresolvedFields == nil {
		self.unresolvedFields = make(map[unresolvedFieldKey]*unresolvedFieldValue)
	}

	for k, v := range b.unresolvedFields {
		existingV, ok := self.unresolvedFields[k]
		if !ok {
			self.unresolvedFields[k] = &unresolvedFieldValue{
				fieldName:         v.fieldName,
				fieldValue:        self._parentValue.Elem().FieldByName(v.fieldName),
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
		if self.catchOtherCache == nil {
			self.catchOtherCache = make(map[string]reflect.Value)
		}

		for k, v := range b.catchOtherCache {
			self.catchOtherCache[k] = v
		}
	}

	if len(b.catchOtherFields) != 0 {
		if self.catchOtherFields == nil {
			self.catchOtherFields = make(map[string]struct{})
		}

		for k, v := range b.catchOtherFields {
			self.catchOtherFields[k] = v
		}
	}

	return nil
}

// UnmarshalYAML handles parsing of rendering suffix and normal yaml
// unmarshaling
// nolint:gocyclo,revive
func (self *BaseField) UnmarshalYAML(n *yaml.Node) error {
	if atomic.LoadUint32(&self._initialized) == 0 {
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
	pt := self._parentValue.Type().Elem()

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
		fieldValue := self._parentValue.Elem().Field(i)

		// initialize struct fields accepted by Init(), in case being used later
		InitRecursively(fieldValue, self.ifaceTypeHandler)

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
			if !addField(yamlKey, fieldType.Name, fieldValue, self) {
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
				inlineFt := self._parentValue.Type().Elem().Field(i).Type

				var iface interface{}
				switch inlineFv.Kind() {
				case reflect.Ptr:
					iface = inlineFv.Interface()
				default:
					iface = inlineFv.Addr().Interface()
				}

				base := self
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
					base:       self,

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

			err = self.unmarshal(
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

// nolint:revive,gocyclo
func (self *BaseField) unmarshal(
	yamlKey string,
	in *alterInterface,
	outVal reflect.Value,
	keepOld bool,
) error {
	for {
		switch outVal.Kind() {
		case reflect.Slice:
			if in.arrData == nil && in.Value() != nil {
				return fmt.Errorf("unexpected non slice data for %q", outVal.String())
			}

			inSlice := in.arrData
			size := len(inSlice)

			sliceVal := reflect.MakeSlice(outVal.Type(), size, size)

			for i := 0; i < size; i++ {
				itemVal := sliceVal.Index(i)

				err := self.unmarshal(
					yamlKey, inSlice[i], itemVal,
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
		case reflect.Map:
			if in.mapData == nil && in.Value() != nil {
				return fmt.Errorf("unexpected non map value for %q", outVal.String())
			}

			// map key MUST be string
			if outVal.IsNil() || !keepOld {
				outVal.Set(reflect.MakeMap(outVal.Type()))
			}

			valType := outVal.Type().Elem()

			isCatchOtherField := self.isCatchOtherField(yamlKey)
			if isCatchOtherField {
				if self.catchOtherCache == nil {
					self.catchOtherCache = make(map[string]reflect.Value)
				}
			}

			for k, v := range in.mapData {
				// since indexed map value is not addressable
				// we have to keep the original data in BaseField cache

				// TODO: test keepOld behavior with cached data

				var val reflect.Value

				if isCatchOtherField {
					cachedData, ok := self.catchOtherCache[k]
					if ok {
						val = cachedData
					} else {
						val = reflect.New(valType)
						self.catchOtherCache[k] = val
					}
				} else {
					val = reflect.New(valType)
				}

				err := self.unmarshal(
					// use iter.Key() rather than `yamlKey`
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
		case reflect.Interface:
			if self.ifaceTypeHandler == nil {
				// use default behavior for interface{} types
				break
			}

			fVal, err := self.ifaceTypeHandler.Create(outVal.Type(), yamlKey)
			if err != nil {
				if errors.Is(err, ErrInterfaceTypeNotHandled) && outVal.Type() == rawInterfaceType {
					// no type information proviede, decode using go-yaml directly
					break
				}

				return fmt.Errorf("failed to create interface field: %w", err)
			}

			val := reflect.ValueOf(fVal)
			if outVal.CanSet() {
				outVal.Set(val)
			} else {
				outVal.Elem().Set(val)
			}

			// DO NOT use outVal directly, which will always match reflect.Interface
			return self.unmarshal(yamlKey, in, val, keepOld)
		case reflect.Ptr:
			// handled after switch
		default:
			// scalar types or struct/array/func/chan/unsafe.Pointer
			// hand it to go-yaml
		}

		if outVal.Kind() != reflect.Ptr {
			break
		}

		if outVal.IsZero() {
			outVal.Set(reflect.New(outVal.Type().Elem()))
		}

		outVal = outVal.Elem()
	}

	var out interface{}
	if outVal.Kind() != reflect.Ptr {
		out = outVal.Addr().Interface()
	} else {
		out = outVal.Interface()
	}

	fVal, canCallInit := out.(Field)
	if canCallInit {
		_ = Init(fVal, self.ifaceTypeHandler)
	}

	dataBytes, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("field: failed to marshal back yaml field %q: %w", yamlKey, err)
	}

	return yaml.Unmarshal(dataBytes, out)
}
