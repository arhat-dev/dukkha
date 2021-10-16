package rs

import (
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

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
