package templateutils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"arhat.dev/pkg/stringhelper"
)

type Integer any

func toUint64(i Integer) uint64 {
	switch t := i.(type) {
	case string:
		return strToUint64(t)
	case []byte:
		return strToUint64(stringhelper.Convert[string, byte](t))

	case int8:
		return uint64(t)
	case int16:
		return uint64(t)
	case int32:
		return uint64(t)
	case int64:
		return uint64(t)
	case int:
		return uint64(t)

	case uint8:
		return uint64(t)
	case uint16:
		return uint64(t)
	case uint32:
		return uint64(t)
	case uint64:
		return uint64(t)
	case uint:
		return uint64(t)
	case uintptr:
		return uint64(t)

	case float32:
		return uint64(t)
	case float64:
		return uint64(t)

	case bool:
		if t {
			return 1
		}

		return 0
	default:
		switch val := reflect.Indirect(reflect.ValueOf(i)); val.Kind() {
		case reflect.String:
			return strToUint64(val.String())
		case reflect.Slice, reflect.Array:
			switch typ := val.Elem().Type(); typ.Kind() {
			case reflect.Uint8:
				return strToUint64(stringhelper.Convert[string, byte](unsafe.Slice((*byte)(unsafe.Pointer(val.UnsafeAddr())), val.Len())))
			default:
				return 0
			}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			return uint64(val.Int())
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			return val.Uint()
		case reflect.Float32, reflect.Float64:
			return uint64(val.Float())
		case reflect.Bool:
			if val.Bool() {
				return 1
			}
			return 0
		default:
			return 0
		}
	}
}

func strToUint64(str string) uint64 {
	if strings.Contains(str, ",") {
		str = strings.ReplaceAll(str, ",", "")
	}

	iv, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		// maybe it's a float?
		var fv float64
		fv, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return 0
		}
		return uint64(fv)
	}
	return uint64(iv)
}

type String any

func toString(s String) string {
	switch t := s.(type) {
	case []byte:
		return *(*string)(unsafe.Pointer(&t))
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	default:
		return fmt.Sprint(t)
	}
}

func toStrings(strArr []String) []string {
	out := make([]string, len(strArr))
	for i, v := range strArr {
		out[i] = toString(v)
	}
	return out
}

type Bytes any

func toBytes(data Bytes) []byte {
	switch t := data.(type) {
	case []byte:
		return t
	case string:
		if n := len(t); n == 0 {
			return []byte{}
		} else {
			return unsafe.Slice((*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&t)).Data)), n)
		}
	default:
		return []byte(fmt.Sprint(t))
	}
}
