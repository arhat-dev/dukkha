package rs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
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

	oneLevelMap, err := unmarshalYamlMap(n)
	if err != nil {
		return fmt.Errorf(
			"rs: data unmarshal failed for %s: %w",
			f._parentType.String(), err,
		)
	}

	visited := make(map[string]struct{})
	allowUnknown := f._opts != nil && f._opts.AllowUnknownFields

	// set values
	for _, kv := range oneLevelMap {
		rawYamlKey := kv[0].Value
		suffixStart := strings.LastIndexByte(rawYamlKey, '@')
		yamlKey := rawYamlKey

		// try to find a field with rendering suffix unstripped
		// cause you can have your tag `yaml:"foo@http"`, then your yaml key
		// is `foo@http`, not `foo` with rendering suffix `@http`

		field := f.getField(yamlKey)
		if suffixStart == -1 || field != nil {
			// no rendering suffix or matched some field, fill value directly

			if _, ok := visited[yamlKey]; ok {
				return fmt.Errorf(
					"rs: duplicate yaml field %q for %s",
					yamlKey, f._parentType.String(),
				)
			}

			visited[yamlKey] = struct{}{}

			v := kv[1]
			if field == nil {
				if f.inlineMap == nil {
					if allowUnknown {
						continue
					}

					return fmt.Errorf(
						"rs: unknown yaml field %q for %s",
						yamlKey, f._parentType.String(),
					)
				}

				field = f.inlineMap
				v = fakeMap(kv[0], kv[1])
			}

			err = unmarshal(
				yamlKey,
				v,
				field.fieldValue,
				&field.isInlineMapItem,
				f._opts,
			)
			if err != nil {
				return fmt.Errorf(
					"rs: failed to unmarshal yaml field %q to %s.%s: %w",
					yamlKey, f._parentType.String(), field.fieldName, err,
				)
			}

			continue
		}

		// has rendering suffix

		yamlKey, suffix := rawYamlKey[:suffixStart], rawYamlKey[suffixStart+1:]

		v := kv[1]
		field = f.getField(yamlKey)
		if field == nil {
			if f.inlineMap == nil {
				if allowUnknown {
					continue
				}

				return fmt.Errorf(
					"rs: unknown yaml field %q for %s",
					yamlKey, f._parentType.String(),
				)
			}

			field = f.inlineMap
			v = fakeMap(cloneYamlNode(kv[0], strTag, yamlKey), kv[1])
		}

		if _, ok := visited[yamlKey]; ok && !field.isInlineMapItem {
			return fmt.Errorf(
				"rs: duplicate yaml data field %q for %s, please note"+
					" rendering suffix won't change the field name",
				yamlKey, f._parentType.String(),
			)
		}

		// do not unmarshal now, we only render the value
		// when calling ResolveFields

		visited[yamlKey] = struct{}{}

		field.base.addUnresolvedField(yamlKey, suffix, nil,
			field.fieldName, field.fieldValue, field.isInlineMapItem,
			v,
		)
	}

	return nil
}

func unmarshal(
	yamlKey string,
	in *yaml.Node,
	outVal reflect.Value,
	keepOld *bool,
	opts *Options,
) error {
	if !outVal.IsValid() {
		// no way to know what value we can set
		// NOTE: this should not happen unless user called yaml.Unmarshal with
		// 	     something nil as out
		return fmt.Errorf(
			"invalid nil unmarshal target for yaml key %q",
			yamlKey,
		)
	}

	outKind := outVal.Kind()

	in = prepareYamlNode(in)
	// reset to zero value if already set
	if isEmpty(in) {
		// only clear out ptr value to keep same behavior as
		// yaml.Unmarshal
		if outVal.CanSet() {
			switch outKind {
			case reflect.Ptr, reflect.Map, reflect.Slice:
				outVal.Set(reflect.Zero(outVal.Type()))
			}
		}
		return nil
	}

	// we are trying to set value of it, so initialize the pointer when not set before
	for outKind == reflect.Ptr {
		if outVal.IsNil() {
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
	case reflect.Chan, reflect.Func:
		return fmt.Errorf("invalid out value is not data type for yaml key %q", yamlKey)
	case reflect.Array:
		return unmarshalArray(yamlKey, in, outVal, opts)
	case reflect.Slice:
		return unmarshalSlice(yamlKey, in, outVal, keepOld, opts)
	case reflect.Map:
		_, err := unmarshalMap(yamlKey, in, outVal, nil, keepOld, opts)
		return err
	case reflect.Struct:
		return unmarshalStruct(yamlKey, in, outVal, opts)
	case reflect.Interface:
		handled, err := unmarshalInterface(yamlKey, in, outVal, keepOld, opts)
		if !handled {
			// fallback to go-yaml behavior
			return in.Decode(outVal.Addr().Interface())
		}

		return err
	default:
		return in.Decode(outVal.Addr().Interface())
	}
}

func unmarshalStruct(
	yamlKey string,
	in *yaml.Node,
	outVal reflect.Value,
	opts *Options,
) error {
	tryInit(outVal, opts)

	var (
		err error
		out = outVal.Addr().Interface()
	)

	switch ot := out.(type) {
	case yaml.Unmarshaler:
		err = ot.UnmarshalYAML(in)
	default:
		err = in.Decode(ot)
	}
	if in.Kind == yaml.MappingNode {
		return err
	}

	if isEmpty(in) {
		return nil
	}

	in, err = TypeHintNone{}.apply(in)
	if err != nil {
		return fmt.Errorf(
			"unexpected input for struct %q %s: %w",
			yamlKey, outVal.Type().String(), err,
		)
	}

	switch ot := out.(type) {
	case yaml.Unmarshaler:
		return ot.UnmarshalYAML(in)
	default:
		return in.Decode(ot)
	}
}

// unmarshalInterface handles interface type creation
func unmarshalInterface(
	yamlKey string,
	in *yaml.Node,
	outVal reflect.Value,
	keepOld *bool,
	opts *Options,
) (bool, error) {
	if opts == nil || opts.InterfaceTypeHandler == nil {
		// there is no interface type handler
		// use default behavior for interface{} types
		return false, nil
	}

	fVal, err := opts.InterfaceTypeHandler.Create(outVal.Type(), yamlKey)
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
	return true, unmarshal(yamlKey, in, val, keepOld, opts)
}

func unmarshalArray(
	yamlKey string,
	in *yaml.Node,
	outVal reflect.Value,
	opts *Options,
) error {
	if in.Kind != yaml.SequenceNode {
		var err error
		in, err = applyObjectsHint(in)
		if err != nil {
			return fmt.Errorf(
				"unexpected input for array %q (%s): %w",
				yamlKey, outVal.Type().String(), err,
			)
		}
	}

	size := len(in.Content)
	expectedSize := outVal.Len()
	if size < expectedSize {
		return fmt.Errorf(
			"array size not match for %s: want %d got %d",
			outVal.Type().String(), expectedSize, size,
		)
	}

	for i := 0; i < expectedSize; i++ {
		err := unmarshal(
			yamlKey, in.Content[i], outVal.Index(i),
			// always drop existing inner data
			// (actually doesn't matter since it's new)
			nil, opts,
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

func unmarshalSlice(
	yamlKey string,
	in *yaml.Node,
	outVal reflect.Value,
	keepOld *bool,
	opts *Options,
) error {
	if in.Kind != yaml.SequenceNode {
		if isEmpty(in) {
			return nil
		}

		var err error
		in, err = applyObjectsHint(in)
		if err != nil {
			return fmt.Errorf(
				"unexpected input for slice %q (%s): %w",
				yamlKey, outVal.Type().String(), err,
			)
		}
	}

	size := len(in.Content)
	tmpVal := reflect.MakeSlice(outVal.Type(), size, size)

	for i := 0; i < size; i++ {
		err := unmarshal(
			yamlKey, in.Content[i], tmpVal.Index(i),
			// always drop existing inner data
			// (actually doesn't matter since it's new)
			nil, opts,
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

	if outVal.IsZero() || keepOld == nil || !*keepOld {
		outVal.Set(tmpVal)
	} else {
		outVal.Set(reflect.AppendSlice(outVal, tmpVal))
	}

	return nil
}

// map key MUST be string
// the return value is only meaningful when keepOld is set (resolving inline map pair)
func unmarshalMap(
	yamlKey string,
	in *yaml.Node,
	outVal reflect.Value,
	inlineMapItemCache *reflect.Value,
	keepOld *bool,
	opts *Options,
) (ret *reflect.Value, _ error) {
	if in.Kind != yaml.MappingNode {
		var err error
		in, err = applyObjectHint(in)
		if err != nil {
			return nil, fmt.Errorf(
				"unexpected input for map %q (%s): %w",
				yamlKey, outVal.Type().String(), err,
			)
		}
	}

	// keepOld can only be set when resolving inline map items
	// then the `in` arg only contain a pair of key value
	//
	// when keepOld is set, there can be some data already unmarshaled
	// inside the map, so we should not create a new map in that case
	if outVal.IsNil() || keepOld == nil {
		outVal.Set(reflect.MakeMap(outVal.Type()))
	}

	m, err := unmarshalYamlMap(in)
	if err != nil {
		return nil, err
	}

	valType := outVal.Type().Elem()
	for i, kv := range m {
		if i == 0 && keepOld != nil && *keepOld && inlineMapItemCache != nil {
			ret = inlineMapItemCache
		} else {
			val := reflect.New(valType).Elem()
			ret = &val
		}

		k := kv[0].Value
		err := unmarshal(
			// use k rather than `yamlKey`
			k,
			kv[1], *ret, keepOld, opts,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal map value %s for key %q: %w",
				valType.String(), k, err,
			)
		}

		outVal.SetMapIndex(reflect.ValueOf(k), *ret)
	}

	return ret, nil
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
