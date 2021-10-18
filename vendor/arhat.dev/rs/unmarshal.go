package rs

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	fieldInterfaceType = reflect.TypeOf((*Field)(nil)).Elem()
	rawInterfaceType   = reflect.TypeOf((*interface{})(nil)).Elem()
)

// UnmarshalYAML handles parsing of rendering suffix and normal yaml
// unmarshaling
func (f *BaseField) UnmarshalYAML(n *yaml.Node) error {
	if !f.initialized() {
		return fmt.Errorf("rs: struct not intialized before unmarshaling")
	}

	if n.Kind != yaml.MappingNode {
		return fmt.Errorf(
			"rs: unexpected non map data %q for struct %q unmarshaling",
			n.Tag, f._parentType.String(),
		)
	}

	m := make(map[string]*alterInterface)
	err := n.Decode(&m)
	if err != nil {
		return fmt.Errorf(
			"rs: data unmarshal failed for %s: %w",
			f._parentType.String(), err,
		)
	}

	handledYamlValues := make(map[string]struct{})
	allowUnknown := f.opts != nil && f.opts.AllowUnknownFields

	// set values
	for rawYamlKey, v := range m {
		suffixAt := strings.LastIndexByte(rawYamlKey, '@')
		yamlKey := rawYamlKey
		if suffixAt == -1 {
			// no rendering suffix, fill value directly

			if _, ok := handledYamlValues[yamlKey]; ok {
				return fmt.Errorf(
					"rs: duplicate yaml field %q",
					yamlKey,
				)
			}

			handledYamlValues[yamlKey] = struct{}{}

			fSpec := f.getField(yamlKey)
			if fSpec == nil {
				if f.catchOtherField == nil {
					if allowUnknown {
						continue
					}

					return fmt.Errorf(
						"rs: unknown yaml field %q in %s",
						yamlKey, f._parentType.String(),
					)
				}

				fSpec = f.catchOtherField
				v = &alterInterface{
					mapData: map[string]*alterInterface{
						yamlKey: v,
					},
				}
			}

			err = f.unmarshal(
				yamlKey,
				reflect.ValueOf(v.NormalizedValue()),
				fSpec.fieldValue,
				fSpec.isCatchOther,
			)
			if err != nil {
				return fmt.Errorf(
					"rs: failed to unmarshal yaml field %q to struct field %s.%s: %w",
					yamlKey, f._parentType.String(), fSpec.fieldName, err,
				)
			}

			continue
		}

		// has rendering suffix

		yamlKey, suffix := rawYamlKey[:suffixAt], rawYamlKey[suffixAt+1:]

		if _, ok := handledYamlValues[yamlKey]; ok {
			return fmt.Errorf(
				"rs: duplicate yaml field name %q for %s, please note"+
					" rendering suffix won't change the field name",
				yamlKey, f._parentType.String(),
			)
		}

		fSpec := f.getField(yamlKey)
		if fSpec == nil {
			if f.catchOtherField == nil {
				if allowUnknown {
					continue
				}

				return fmt.Errorf(
					"rs: unknown yaml field %q for %s",
					yamlKey, f._parentType.String(),
				)
			}

			fSpec = f.catchOtherField
			v = &alterInterface{
				mapData: map[string]*alterInterface{
					yamlKey: v,
				},
			}
		}

		// do not unmarshal now, we will render the value
		// and unmarshal when calling ResolveFields

		handledYamlValues[yamlKey] = struct{}{}
		// don't forget the raw name with rendering suffix
		handledYamlValues[rawYamlKey] = struct{}{}

		fSpec.base.addUnresolvedField(
			yamlKey, suffix,
			fSpec.fieldName, fSpec.fieldValue, fSpec.isCatchOther,
			v,
		)
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
		"rs: unknown yaml fields to %s: %s",
		f._parentType.String(), strings.Join(unknownFields, ", "),
	)
}

func (f *BaseField) unmarshal(
	yamlKey string,
	in reflect.Value,
	outVal reflect.Value,
	keepOld bool,
) error {
	outKind := outVal.Kind()

	if !outVal.IsValid() {
		// no way to know what value we can set
		// NOTE: this should not happen unless user called yaml.Unmarshal with
		// 	     something nil as out
		return fmt.Errorf("invalid nil unmarshal target")
	}

	// reset to zero value if already set
	if !in.IsValid() {
		if outVal.CanSet() {
			outVal.Set(reflect.Zero(outVal.Type()))
		}
		return nil
	}

	// we are trying to set value of it, so initialize the pointer when not set before
	for outKind == reflect.Ptr {
		if outVal.IsZero() {
			outVal.Set(reflect.New(outVal.Type().Elem()))
		}

		outVal = outVal.Elem()
		outKind = outVal.Kind()
	}

	switch outKind {
	case reflect.Invalid:
		// TODO: this should not happen since we have already checked outVal.IsValid before

		// unreachable code
		return fmt.Errorf("unexpected nil out value for yaml key %q", yamlKey)
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
		// val := reflect.ValueOf(in.NormalizedValue())
		// inKind := val.Kind()

		inKind := in.Kind()
		if inKind == outKind {
			// same kind means same scalar type
			outVal.Set(in)
			return nil
		}

		// not same kind, try to convert form input value
		// currently we only
		switch vt := in.Interface().(type) {
		case string:
			// input is string, while output is not, maybe it's []byte ([]uint8)
			switch outVal.Interface().(type) {
			case []byte:
				outVal.SetBytes([]byte(vt))
				return nil
			default:
				return f.unmarshalRaw(in, outVal)
			}
		case []byte:
			// input is bytes, while output is not, maybe it's string
			switch outVal.Interface().(type) {
			case string:
				outVal.SetString(string(vt))
				return nil
			default:
				return f.unmarshalRaw(in, outVal)
			}
		default:
			// not compatible, check if assignable
			if err := checkAssignable(yamlKey, in, outVal); err != nil {
				return err
			}

			outVal.Set(in)
			return nil
		}
	}
}

// unmarshalRaw unmarshals interface{} type value to outVal
func (f *BaseField) unmarshalRaw(in, outVal reflect.Value) error {
	tryInit(outVal, f.opts)

	dataBytes, err := yaml.Marshal(in.Interface())
	if err != nil {
		return fmt.Errorf(
			"rs: failed to marshal back yaml for %q: %w",
			outVal.String(), err,
		)
	}

	var out interface{}
	if outVal.Kind() != reflect.Ptr {
		out = outVal.Addr().Interface()
	} else {
		out = outVal.Interface()
	}

	err = yaml.Unmarshal(dataBytes, out)
	if err != nil {
		return err
	}

	return nil
}

// unmarshalInterface handles interface type creation
func (f *BaseField) unmarshalInterface(
	yamlKey string,
	in, outVal reflect.Value,
	keepOld bool,
) (bool, error) {
	if f.opts == nil || f.opts.InterfaceTypeHandler == nil {
		// there is no interface type handler
		// use default behavior for interface{} types
		return false, nil
	}

	fVal, err := f.opts.InterfaceTypeHandler.Create(outVal.Type(), yamlKey)
	if err != nil {
		if errors.Is(err, ErrInterfaceTypeNotHandled) && outVal.Type() == rawInterfaceType {
			// no type information provided, decode using go-yaml directly
			return false, nil
		}

		return true, fmt.Errorf(
			"failed to create interface field: %w",
			err,
		)
	}

	val := reflect.ValueOf(fVal)

	if err := checkAssignable(yamlKey, val, outVal); err != nil {
		return true, err
	}

	if outVal.CanSet() {
		outVal.Set(val)
	} else {
		outVal.Elem().Set(val)
	}

	// DO NOT use outVal directly, which will always match reflect.Interface
	return true, f.unmarshal(yamlKey, in, val, keepOld)
}

func (f *BaseField) unmarshalArray(yamlKey string, in, outVal reflect.Value) error {
	if ik := in.Kind(); ik != reflect.Slice && ik != reflect.Array {
		var err error
		switch it := in.Interface().(type) {
		case []byte:
			if len(it) == 0 {
				// nil value
				return f.unmarshal(yamlKey, reflect.Value{}, outVal, false)
			}

			var inVal []interface{}
			err = yaml.Unmarshal(it, &inVal)
			if err != nil {
				break
			}

			in = reflect.ValueOf(inVal)
		case string:
			if len(it) == 0 {
				// nil value
				return f.unmarshal(yamlKey, reflect.Value{}, outVal, false)
			}

			var inVal []interface{}
			err = yaml.Unmarshal([]byte(it), &inVal)
			if err != nil {
				break
			}

			in = reflect.ValueOf(inVal)
		default:
			err = fmt.Errorf("incompatible value %q with type %T", it, it)
		}

		if err != nil {
			return fmt.Errorf(
				"unexpected non array data of yaml field %q for %s: %w",
				yamlKey, outVal.Type().String(), err,
			)
		}
	}

	size := in.Len()
	expectedSize := outVal.Len()
	if size < expectedSize {
		return fmt.Errorf(
			"array size not match for %s: want %d got %d",
			outVal.Type().String(), expectedSize, size,
		)
	}

	for i := 0; i < expectedSize; i++ {
		err := f.unmarshal(
			yamlKey, in.Index(i), outVal.Index(i),
			// always drop existing inner data
			// (actually doesn't matter since it's new)
			false,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal #%d array item of yaml field %q for %s: %w",
				i, yamlKey, outVal.Type().String(), err,
			)
		}
	}

	return nil
}

func (f *BaseField) unmarshalSlice(yamlKey string, in, outVal reflect.Value, keepOld bool) error {
	if ik := in.Kind(); ik != reflect.Slice && ik != reflect.Array {
		var err error
		switch it := in.Interface().(type) {
		case []byte:
			if len(it) == 0 {
				// nil value
				return f.unmarshal(yamlKey, reflect.Value{}, outVal, keepOld)
			}

			var inVal []interface{}
			err = yaml.Unmarshal(it, &inVal)
			if err != nil {
				break
			}

			in = reflect.ValueOf(inVal)
		case string:
			if len(it) == 0 {
				// nil value
				return f.unmarshal(yamlKey, reflect.Value{}, outVal, keepOld)
			}

			var inVal []interface{}
			err = yaml.Unmarshal([]byte(it), &inVal)
			if err != nil {
				break
			}

			in = reflect.ValueOf(inVal)
		default:
			err = fmt.Errorf("incompatible value %q with type %T", it, it)
		}

		if err != nil {
			return fmt.Errorf(
				"unexpected non slice data of yaml field %q for %s: %w",
				yamlKey, outVal.Type().String(), err,
			)
		}
	}

	size := in.Len()
	tmpVal := reflect.MakeSlice(outVal.Type(), size, size)

	for i := 0; i < size; i++ {
		err := f.unmarshal(
			yamlKey, in.Index(i), tmpVal.Index(i),
			// always drop existing inner data
			// (actually doesn't matter since it's new)
			false,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal #%d slice item of yaml field %q for %s: %w",
				i, yamlKey, outVal.Type().String(), err,
			)
		}
	}

	if err := checkAssignable(yamlKey, tmpVal, outVal); err != nil {
		return err
	}

	if outVal.IsZero() || !keepOld {
		outVal.Set(tmpVal)
	} else {
		outVal.Set(reflect.AppendSlice(outVal, tmpVal))
	}

	return nil
}

func (f *BaseField) unmarshalMap(yamlKey string, in, outVal reflect.Value, keepOld bool) error {
	if in.Kind() != reflect.Map {
		var err error
		switch it := in.Interface().(type) {
		case []byte:
			if len(it) == 0 {
				// nil value
				return f.unmarshal(yamlKey, reflect.Value{}, outVal, keepOld)
			}

			inVal := make(map[string]interface{})
			err = yaml.Unmarshal(it, &inVal)
			if err != nil {
				break
			}

			in = reflect.ValueOf(inVal)
		case string:
			if len(it) == 0 {
				// nil value
				return f.unmarshal(yamlKey, reflect.Value{}, outVal, keepOld)
			}

			inVal := make(map[string]interface{})
			err = yaml.Unmarshal([]byte(it), &inVal)
			if err != nil {
				break
			}

			in = reflect.ValueOf(inVal)
		default:
			err = fmt.Errorf("incompatible value %q with type %T", it, it)
		}

		if err != nil {
			return fmt.Errorf(
				"unexpected non map data of yaml field %q for %s: %w",
				yamlKey, outVal.Type().String(), err,
			)
		}
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

	iter := in.MapRange()
	for iter.Next() {
		k, v := iter.Key().String(), iter.Value()
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
			// (field tag: `rs:"other"`)
			k,
			// NOTE: if the map value type is interface{} (e.g. map[string]interface{})
			//		 return value of iter.Value() will lose its real type info (always be interface{} type)
			// 		 to fix this, we have to get its underlying data and do reflect.ValueOf again
			reflect.ValueOf(v.Interface()),
			val, keepOld,
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

func checkAssignable(yamlKey string, in, out reflect.Value) error {
	if !in.Type().AssignableTo(out.Type()) {
		return fmt.Errorf(
			"unexpected value of yaml field %q: want %q, got %q",
			yamlKey, out.Type().String(), in.Type().String(),
		)
	}

	return nil
}
