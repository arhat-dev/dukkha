package rs

import (
	"fmt"
	"reflect"
	"strconv"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

func (f *BaseField) HasUnresolvedField() bool {
	return len(f.unresolvedFields) != 0
}

func (f *BaseField) ResolveFields(rc RenderingHandler, depth int, fieldNames ...string) error {
	if atomic.LoadUint32(&f._initialized) == 0 {
		return fmt.Errorf("rs: struct not intialized before resolving")
	}

	if depth == 0 {
		return nil
	}

	parentStructType := f._parentValue.Type()

	if len(fieldNames) == 0 {
		// resolve all

		for i := 1; i < f._parentValue.NumField(); i++ {
			sf := parentStructType.Field(i)
			if len(sf.PkgPath) != 0 {
				// not exported
				continue
			}

			err := f.resolveSingleField(
				rc, depth, sf.Name, f._parentValue.Field(i),
			)
			if err != nil {
				return fmt.Errorf(
					"rs: failed to resolve field %s.%s: %w",
					parentStructType.String(), sf.Name, err,
				)
			}
		}

		return nil
	}

	for _, name := range fieldNames {
		fv := f._parentValue.FieldByName(name)
		if !fv.IsValid() {
			return fmt.Errorf("rs: no such field %q in struct %q", name, parentStructType.String())
		}

		err := f.resolveSingleField(rc, depth, name, fv)
		if err != nil {
			return fmt.Errorf(
				"rs: failed to resolve requested field %s.%s: %w",
				parentStructType.String(), name, err,
			)
		}
	}

	return nil
}

func (f *BaseField) resolveSingleField(
	rc RenderingHandler,
	depth int,
	fieldName string,
	targetField reflect.Value,
) error {
	keepOld := false

	for k, v := range f.unresolvedFields {
		if v.fieldName == fieldName {
			err := f.handleUnResolvedField(
				rc, depth, k, v, keepOld,
			)
			if err != nil {
				return err
			}

			keepOld = true
		}
	}

	return f.handleResolvedField(rc, depth, targetField)
}

// nolint:gocyclo
func (f *BaseField) handleResolvedField(
	rc RenderingHandler,
	depth int,
	field reflect.Value,
) error {
	if depth == 0 {
		return nil
	}

	switch fk := field.Kind(); fk {
	case reflect.Map:
		if field.IsNil() {
			return nil
		}

		iter := field.MapRange()
		for iter.Next() {
			err := f.handleResolvedField(rc, depth-1, iter.Value())
			if err != nil {
				return err
			}
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if fk == reflect.Slice && field.IsNil() {
			// this is a resolved field, slice/array empty means no value
			return nil
		}

		for i := 0; i < field.Len(); i++ {
			tt := field.Index(i)
			err := f.handleResolvedField(rc, depth-1, tt)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		// handled after switch
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if !field.IsValid() || field.IsZero() || field.IsNil() {
			return nil
		}
		// handled after switch
	default:
		// scalar types, no action required
		return nil
	}

	return tryResolve(rc, depth, field)
}

func tryResolve(rc RenderingHandler, depth int, targetField reflect.Value) error {
	if targetField.CanInterface() {
		fVal, canCallResolve := targetField.Interface().(Field)
		if canCallResolve {
			return fVal.ResolveFields(rc, depth)
		}
	}

	if targetField.CanAddr() && targetField.Addr().CanInterface() {
		fVal, canCallResolve := targetField.Addr().Interface().(Field)
		if canCallResolve {
			return fVal.ResolveFields(rc, depth)
		}
	}

	// no more field to resolve
	return nil
}

func (f *BaseField) handleUnResolvedField(
	rc RenderingHandler,
	depth int,
	key unresolvedFieldKey,
	v *unresolvedFieldValue,
	keepOld bool,
) error {
	target := v.fieldValue

	for i, rawData := range v.rawDataList {
		toResolve := rawData
		if v.isCatchOtherField {
			// unwrap map data for resolving
			toResolve = rawData.mapData[key.yamlKey]
		}

		var (
			err error
		)

		for _, renderer := range v.renderers {

			// a patch is implied when the renderer has a `!` suffix
			var patchSpec *renderingPatchSpec
			if renderer.patchSpec {
				patchSpec, toResolve, err = f.resolvePatchSpec(rc, toResolve)
				if err != nil {
					return fmt.Errorf(
						"failed to resolve patch spec for renderer %q: %w",
						renderer.name, err,
					)
				}
			}

			if len(renderer.name) != 0 {
				var tmp interface{}

				// toResolve can be nil when no patch value is set
				if toResolve != nil {
					tmp = toResolve.NormalizedValue()
				}

				tmp, err = rc.RenderYaml(renderer.name, tmp)
				if err != nil {
					return fmt.Errorf(
						"renderer %q failed to render value: %w",
						renderer.name, err,
					)
				}

				toResolve = &alterInterface{
					scalarData: tmp,
				}
			}

			// apply patch if set
			if patchSpec != nil {
				// apply type hint before patching to make sure value
				// being patched is correctly typed
				hint := renderer.typeHint
				toResolve, err = applyTypeHint(hint, toResolve)
				if err != nil {
					return fmt.Errorf(
						"failed to ensure type hint %q on yaml key %q: %w",
						hint, key.yamlKey, err,
					)
				}

				var tmp interface{}
				tmp, err = patchSpec.ApplyTo(toResolve.NormalizedValue())
				if err != nil {
					return fmt.Errorf(
						"failed to apply patches: %w",
						err,
					)
				}

				toResolve = &alterInterface{
					scalarData: tmp,
				}
			}

			// apply hint after resolving (rendering)
			hint := renderer.typeHint
			toResolve, err = applyTypeHint(hint, toResolve)
			if err != nil {
				return fmt.Errorf(
					"failed to ensure type hint %q on yaml key %q: %w",
					hint, key.yamlKey, err,
				)
			}
		}

		resolved := toResolve.NormalizedValue()
		if v.isCatchOtherField {
			// wrap back for catch other filed
			resolved = map[string]interface{}{
				key.yamlKey: resolved,
			}
		}

		// TODO: currently we always keepOld when the field has tag
		// 		 `rs:"other"`, need to ensure this behavior won't
		// 	     leave inconsistent data

		actualKeepOld := keepOld || v.isCatchOtherField || i != 0
		err = f.unmarshal(key.yamlKey, reflect.ValueOf(resolved), target, actualKeepOld)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal resolved value of yaml key %q to field %q: %w",
				key.yamlKey, v.fieldName, err,
			)
		}
	}

	return tryResolve(rc, depth-1, target)
}

// resolve user provided data as patch spec
func (f *BaseField) resolvePatchSpec(
	rc RenderingHandler,
	toResolve interface{},
) (
	patchSpec *renderingPatchSpec,
	value *alterInterface,
	err error,
) {
	var patchSpecBytes []byte
	switch t := toResolve.(type) {
	case string:
		patchSpecBytes = []byte(t)
	case []byte:
		patchSpecBytes = t
	default:
		patchSpecBytes, err = yaml.Marshal(toResolve)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to marshal renderer data for patch spec: %w",
				err,
			)
		}
	}

	patchSpec = Init(&renderingPatchSpec{}, f.ifaceTypeHandler).(*renderingPatchSpec)
	err = yaml.Unmarshal(patchSpecBytes, patchSpec)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to unmarshal patch spec %q: %w",
			string(patchSpecBytes), err,
		)
	}

	err = patchSpec.ResolveFields(rc, -1)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to resolve patch spec %q: %w",
			string(patchSpecBytes), err,
		)
	}

	value = convertAnyObjectToAlterInterface(&patchSpec.Value)
	return patchSpec, value, nil
}

func (f *BaseField) isCatchOtherField(yamlKey string) bool {
	if f.catchOtherFields == nil {
		return false
	}

	_, ok := f.catchOtherFields[yamlKey]
	return ok
}

type TypeHint int8

const (
	TypeHintNone TypeHint = iota
	TypeHintStr
	TypeHintBytes
	TypeHintObject
	TypeHintObjects
	TypeHintInt
	TypeHintUint
	TypeHintFloat
)

func (h TypeHint) String() string {
	switch h {
	case TypeHintNone:
		return ""
	case TypeHintStr:
		return "str"
	case TypeHintBytes:
		return "[]byte"
	case TypeHintObjects:
		return "[]obj"
	case TypeHintObject:
		return "obj"
	case TypeHintInt:
		return "int"
	case TypeHintUint:
		return "uint"
	case TypeHintFloat:
		return "float"
	default:
		return "<unknown>"
	}
}

func ParseTypeHint(h string) (TypeHint, error) {
	switch h {
	case "":
		return TypeHintNone, nil
	case "str":
		return TypeHintStr, nil
	case "[]byte":
		return TypeHintBytes, nil
	case "[]obj":
		return TypeHintObjects, nil
	case "obj":
		return TypeHintObject, nil
	case "int":
		return TypeHintInt, nil
	case "uint":
		return TypeHintUint, nil
	case "float":
		return TypeHintFloat, nil
	default:
		return -1, fmt.Errorf("unknown type hint %q", h)
	}
}

// nolint:gocyclo
func applyTypeHint(hint TypeHint, v *alterInterface) (*alterInterface, error) {
	switch hint {
	case TypeHintNone:
		// no hint, use default behavior:
		// 	return string value if unmarshal failed or result
		//  is string or []byte

		var rawBytes []byte
		switch vt := v.scalarData.(type) {
		case string:
			rawBytes = []byte(vt)
		case []byte:
			rawBytes = vt
		default:
			return v, nil
		}

		tmp := new(alterInterface)
		err := yaml.Unmarshal(rawBytes, tmp)
		if err != nil {
			// couldn't unmarshal, return original value in string format
			return &alterInterface{
				scalarData: string(rawBytes),
			}, nil
		}

		switch tmp.Value().(type) {
		case string, []byte, nil:
			// yaml.Unmarshal will do some transformation on plaintext value
			// when it's not valid yaml, so return the original value
			return &alterInterface{
				scalarData: string(rawBytes),
			}, nil
		default:
			return tmp, nil
		}
	case TypeHintStr:
		if v.originalNode != nil && v.originalNode.Kind == yaml.ScalarNode {
			return &alterInterface{
				scalarData: v.originalNode.Value,
			}, nil
		}

		switch vt := v.Value().(type) {
		case []byte:
			return &alterInterface{
				scalarData: string(vt),
			}, nil
		case string:
			return &alterInterface{
				scalarData: vt,
			}, nil
		default:
			bytesV, err := yaml.Marshal(vt)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to marshal %T value as str",
					vt,
				)
			}
			return &alterInterface{
				scalarData: string(bytesV),
			}, nil
		}
	case TypeHintBytes:
		if v.originalNode != nil && v.originalNode.Kind == yaml.ScalarNode {
			return &alterInterface{
				scalarData: []byte(v.originalNode.Value),
			}, nil
		}

		switch vt := v.Value().(type) {
		case []byte:
			return &alterInterface{
				scalarData: vt,
			}, nil
		case string:
			return &alterInterface{
				scalarData: []byte(vt),
			}, nil
		default:
			bytesV, err := yaml.Marshal(vt)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to marshal %T value as str",
					vt,
				)
			}
			return &alterInterface{
				scalarData: bytesV,
			}, nil
		}
	case TypeHintObjects:
		switch vt := v.scalarData.(type) {
		case []byte:
			ret := new(alterInterface)
			err := yaml.Unmarshal(vt, &ret.sliceData)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal bytes %q as object array: %w",
					string(vt), err,
				)
			}
			return ret, nil
		case string:
			ret := new(alterInterface)
			err := yaml.Unmarshal([]byte(vt), &ret.sliceData)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal string %q as object array: %w",
					vt, err,
				)
			}
			return ret, nil
		default:
			switch vk := reflect.ValueOf(v.Value()).Kind(); vk {
			case reflect.Array, reflect.Slice, reflect.Interface, reflect.Ptr:
				return v, nil
			default:
				return nil, fmt.Errorf(
					"incompatible type %T as object array", vt,
				)
			}
		}
	case TypeHintObject:
		switch vt := v.scalarData.(type) {
		case []byte:
			ret := new(alterInterface)
			err := yaml.Unmarshal(vt, &ret.mapData)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal bytes %q as object: %w",
					string(vt), err,
				)
			}
			return ret, nil
		case string:
			ret := new(alterInterface)
			err := yaml.Unmarshal([]byte(vt), &ret.mapData)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal string %q as object: %w",
					vt, err,
				)
			}
			return ret, nil
		default:
			switch vk := reflect.ValueOf(v.Value()).Kind(); vk {
			case reflect.Map, reflect.Struct, reflect.Interface, reflect.Ptr:
				return v, nil
			default:
				return nil, fmt.Errorf(
					"incompatible type %T as object", vt,
				)
			}
		}
	case TypeHintInt:
		switch vt := v.scalarData.(type) {
		case []byte:
			strV, err := strconv.Unquote(string(vt))
			if err != nil {
				strV = string(vt)
			}

			intV, err := strconv.ParseInt(strV, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal bytes %q as int: %w",
					string(vt), err,
				)
			}

			return &alterInterface{
				scalarData: int(intV),
			}, nil
		case string:
			strV, err := strconv.Unquote(vt)
			if err != nil {
				strV = vt
			}

			intV, err := strconv.ParseInt(strV, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal string %q as int: %w",
					vt, err,
				)
			}

			return &alterInterface{
				scalarData: int(intV),
			}, nil
		case int, int8, int16, int32, int64:
			return v, nil
		default:
			rv := reflect.ValueOf(v.NormalizedValue())
			switch vk := rv.Kind(); vk {
			case reflect.Float32, reflect.Float64:
				return &alterInterface{
					scalarData: int(rv.Float()),
				}, nil
			case reflect.Uint, reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return &alterInterface{
					scalarData: int(rv.Uint()),
				}, nil
			default:
				return nil, fmt.Errorf(
					"incompatible type %T as number", vt,
				)
			}
		}
	case TypeHintUint:
		switch vt := v.scalarData.(type) {
		case []byte:
			strV, err := strconv.Unquote(string(vt))
			if err != nil {
				strV = string(vt)
			}

			uintV, err := strconv.ParseUint(strV, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal bytes %q as uint: %w",
					string(vt), err,
				)
			}

			return &alterInterface{
				scalarData: uint(uintV),
			}, nil
		case string:
			strV, err := strconv.Unquote(vt)
			if err != nil {
				strV = vt
			}

			uintV, err := strconv.ParseUint(strV, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal string %q as uint: %w",
					vt, err,
				)
			}

			return &alterInterface{
				scalarData: uint(uintV),
			}, nil
		case uint, uint8, uint16, uint32, uint64, uintptr:
			return v, nil
		default:
			rv := reflect.ValueOf(v.NormalizedValue())
			switch vk := rv.Kind(); vk {
			case reflect.Float32, reflect.Float64:
				return &alterInterface{
					scalarData: uint(rv.Float()),
				}, nil
			case reflect.Int, reflect.Int8, reflect.Int16,
				reflect.Int32, reflect.Int64:
				return &alterInterface{
					scalarData: uint(rv.Int()),
				}, nil
			default:
				return nil, fmt.Errorf(
					"incompatible type %T as number", vt,
				)
			}
		}
	case TypeHintFloat:
		switch vt := v.scalarData.(type) {
		case []byte:
			strV, err := strconv.Unquote(string(vt))
			if err != nil {
				strV = string(vt)
			}

			f64v, err := strconv.ParseFloat(strV, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal bytes %q as float: %w",
					string(vt), err,
				)
			}

			return &alterInterface{
				scalarData: f64v,
			}, nil
		case string:
			strV, err := strconv.Unquote(vt)
			if err != nil {
				strV = vt
			}

			f64v, err := strconv.ParseFloat(strV, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal string %q as float: %w",
					vt, err,
				)
			}

			return &alterInterface{
				scalarData: f64v,
			}, nil
		case float32, float64:
			return v, nil
		default:
			rv := reflect.ValueOf(v.Value())
			switch vk := rv.Kind(); vk {
			case reflect.Int, reflect.Int8, reflect.Int16,
				reflect.Int32, reflect.Int64:
				return &alterInterface{
					scalarData: float64(rv.Int()),
				}, nil
			case reflect.Uint, reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return &alterInterface{
					scalarData: float64(rv.Uint()),
				}, nil
			default:
				return nil, fmt.Errorf(
					"incompatible type %T as number", vt,
				)
			}
		}
	default:
		return nil, fmt.Errorf("unknown type hint %d", hint)
	}
}
