package templateutils

import (
	"io"
	"reflect"
)

// typeNS for scalar type conversion
type typeNS struct{}

func (typeNS) Close(file any) (_ None, err error) {
	c, ok := file.(io.Closer)
	if ok {
		err = c.Close()
		return
	}

	return
}

func (typeNS) Default(def, v any) any {
	if IsZero(v) {
		return def
	}

	return v
}

func (typeNS) FirstNoneZero(v ...any) any {
	for _, elem := range v {
		if !IsZero(elem) {
			return elem
		}
	}

	return nil
}

func (typeNS) AllTrue(v ...any) bool {
	for _, elem := range v {
		if IsZero(elem) {
			return false
		}
	}

	return true
}

func (typeNS) AnyTrue(v ...any) bool {
	for _, elem := range v {
		if !IsZero(elem) {
			return true
		}
	}

	return len(v) == 0
}

func (typeNS) ToBool(v any) (bool, error)        { return parseBool(v) }
func (typeNS) ToUint(v any) (uint64, error)      { return parseInteger[uint64](v) }
func (typeNS) ToInt(v any) (int64, error)        { return parseInteger[int64](v) }
func (typeNS) ToFloat(v any) (float64, error)    { return parseFloat[float64](v) }
func (typeNS) ToString(v any) (string, error)    { return toString(v) }
func (typeNS) ToStrings(v any) ([]string, error) { return anyToStrings(v) }

// TODO: fix IsXXX logic

func (typeNS) IsBool(v any) bool {
	_, err := parseBool(v)
	return err == nil
}

func (typeNS) IsInt(v any) bool {
	_, err := parseInteger[int](v)
	return err == nil
}

func (typeNS) IsFloat(v any) bool {
	_, err := parseFloat[float64](v)
	return err == nil
}

func (typeNS) IsNum(v any) bool {
	_, _, err := parseNumber(v)
	return err == nil
}

func (typeNS) IsZero(v any) bool { return IsZero(v) }

func IsZero(v any) bool {
	switch t := v.(type) {
	case string:
		return len(t) == 0

	case bool:
		return !t
	case *bool:
		return t == nil

	case int:
		return t == 0
	case uint:
		return t == 0

	case int8:
		return t == 0
	case uint8:
		return t == 0

	case int16:
		return t == 0
	case uint16:
		return t == 0

	case int32:
		return t == 0
	case uint32:
		return t == 0

	case int64:
		return t == 0
	case uint64:
		return t == 0

	case uintptr:
		return t == 0

	case float32:
		return t == 0
	case float64:
		return t == 0

	case []int:
		return len(t) == 0

	case []uint:
		return len(t) == 0

	case []int8:
		return len(t) == 0

	case []uint8:
		return len(t) == 0

	case []int16:
		return len(t) == 0

	case []uint16:
		return len(t) == 0

	case []int32:
		return len(t) == 0

	case []uint32:
		return len(t) == 0

	case []int64:
		return len(t) == 0

	case []uint64:
		return len(t) == 0

	case []uintptr:
		return len(t) == 0

	case []float32:
		return len(t) == 0

	case []float64:
		return len(t) == 0

	case map[string]struct{}:
		return len(t) == 0

	case map[any]any:
		return len(t) == 0

	case map[string]string:
		return len(t) == 0

	case map[string]any:
		return len(t) == 0

	case *struct{}:
		return true
	case struct{}:
		return true

	case *None:
		return true
	case None:
		return true

	default:
		switch val := reflect.ValueOf(t); val.Kind() {
		case reflect.Invalid:
			return true
		case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
			return val.Len() == 0
		case reflect.Bool:
			return !val.Bool()
		case reflect.Complex64, reflect.Complex128:
			return val.Complex() == 0
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return val.Int() == 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return val.Uint() == 0
		case reflect.Float32, reflect.Float64:
			return val.Float() == 0
		case reflect.Struct:
			return false
		default:
			return val.IsNil()
		}
	}
}
