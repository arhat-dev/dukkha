package templateutils

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"arhat.dev/pkg/stringhelper"
)

func must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}

	return x
}

type None struct{}

// Bool type is applicable to
// - ~string (convert "true", "yes", "y", "1" to true, "false", "no", "n", "0", "" to false)
// - ~[]~byte (like string)
// - integers (true if non-zero)
// - floats (true if non-zero)
// - bool value
type Bool any

func toBoolOrFalse(v Bool) (ret bool) {
	ret, err := parseBool(v)
	if err != nil {
		return false
	}

	return
}

func parseBool(v Bool) (bool, error) {
	switch v := v.(type) {
	case bool:
		return v, nil

	case string:
		return strToBool(v)
	case []byte:
		return strToBool(stringhelper.Convert[string, byte](v))

	case int:
		return v != 0, nil
	case uint:
		return v != 0, nil

	case int8:
		return v != 0, nil
	case uint8:
		return v != 0, nil

	case int16:
		return v != 0, nil
	case uint16:
		return v != 0, nil

	case int32:
		return v != 0, nil
	case uint32:
		return v != 0, nil

	case int64:
		return v != 0, nil
	case uint64:
		return v != 0, nil

	case uintptr:
		return v != 0, nil

	case float32:
		return v != 0, nil
	case float64:
		return v != 0, nil

	case nil:
		return false, nil

	default:
		switch val := reflect.Indirect(reflect.ValueOf(v)); val.Kind() {
		case reflect.String:
			return strToBool(val.String())
		case reflect.Slice, reflect.Array:
			switch typ := val.Elem().Type(); typ.Kind() {
			case reflect.Uint8:
				return strToBool(*(*string)(val.Addr().UnsafePointer()))
			default:
				return false, fmt.Errorf("unsupport bool source type %T", v)
			}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			return val.Int() != 0, nil
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			return val.Uint() != 0, nil
		case reflect.Float32, reflect.Float64:
			return val.Float() != 0, nil
		case reflect.Bool:
			return val.Bool(), nil
		default:
			return false, fmt.Errorf("unsupport bool source type %T", v)
		}
	}
}

func strToBool(s string) (bool, error) {
	switch v := *(*string)(unsafe.Pointer(&s)); v {
	case "true", "yes", "y", "1":
		return true, nil
	case "false", "no", "n", "0", "":
		return false, nil
	default:
		return false, fmt.Errorf("unknown string bool value %q", v)
	}
}

func toBools[T ~bool](b []Bool) (ret []T, err error) {
	ret = make([]T, len(b))
	var v bool
	for i := range b {
		v, err = parseBool(b[i])
		ret[i] = T(v)
	}

	return
}

// Number type is applicable to
// - ~string
// - ~[]~byte
// - integers
// - floats
// - bool value (true as 1, false as 0)
type Number any

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type float interface {
	~float32 | ~float64
}

// toIntegerOrPanic panics on error condition
func toIntegerOrPanic[T integer](i Number) (ret T) {
	ret, err := parseInteger[T](i)
	if err != nil {
		panic(err)
	}

	return
}

// toFloatOrPanic panics on error condition
func toFloatOrPanic[T float](i Number) (ret T) {
	ret, err := parseFloat[T](i)
	if err != nil {
		panic(err)
	}

	return
}

func toIntegers[T integer, V Number](arr []V) (ret []T, err error) {
	if len(arr) == 0 {
		return
	}

	ret = make([]T, len(arr))
	for i := range arr {
		ret[i], err = parseInteger[T](arr[i])
		if err != nil {
			return
		}
	}

	return
}

func toFloats[T float, V Number](arr []V) (ret []T, err error) {
	if len(arr) == 0 {
		return
	}

	ret = make([]T, len(arr))
	for i := range arr {
		ret[i], err = parseFloat[T](arr[i])
		if err != nil {
			return
		}
	}

	return
}

func parseFloat[T float](i Number) (T, error) {
	iv, isFloat, err := parseNumber(i)
	if err != nil {
		return 0, err
	}

	if isFloat {
		return T(math.Float64frombits(iv)), nil
	}

	return T(iv), nil
}

func parseInteger[T integer](i Number) (T, error) {
	ret, isFloat, err := parseNumber(i)
	if err != nil {
		return 0, err
	}

	if isFloat {
		return T(math.Float64frombits(ret)), nil
	}

	return T(ret), nil
}

// parseNumber converts i to uint64 (with sign kept)
//
// if i is a float number (indicated by return value isFloat), return IEEE 754 bits of i
func parseNumber(i Number) (_ uint64, isFloat bool, _ error) {
	switch i := i.(type) {
	case string:
		return strToInteger(i)
	case []byte:
		return strToInteger(stringhelper.Convert[string, byte](i))

	case int:
		return uint64(i), false, nil
	case uint:
		return uint64(i), false, nil

	case int8:
		return uint64(i), false, nil
	case uint8:
		return uint64(i), false, nil

	case int16:
		return uint64(i), false, nil
	case uint16:
		return uint64(i), false, nil

	case int32:
		return uint64(i), false, nil
	case uint32:
		return uint64(i), false, nil

	case int64:
		return uint64(i), false, nil
	case uint64:
		return uint64(i), false, nil

	case uintptr:
		return uint64(i), false, nil

	case float32:
		return uint64(math.Float32bits(i)), true, nil
	case float64:
		return math.Float64bits(i), true, nil

	case bool:
		if i {
			return 1, false, nil
		}

		return 0, false, nil

	case nil:
		return 0, false, nil

	default:
		switch val := reflect.Indirect(reflect.ValueOf(i)); val.Kind() {
		case reflect.String:
			return strToInteger(val.String())
		case reflect.Slice, reflect.Array:
			switch typ := val.Elem().Type(); typ.Kind() {
			case reflect.Uint8:
				return strToInteger(*(*string)(val.Addr().UnsafePointer()))
			default:
				return 0, false, fmt.Errorf("unhandled array type %q", typ.String())
			}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			return uint64(val.Int()), false, nil
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			return val.Uint(), false, nil
		case reflect.Float32, reflect.Float64:
			return math.Float64bits(val.Float()), true, nil
		case reflect.Bool:
			if val.Bool() {
				return 1, false, nil
			}
			return 0, false, nil
		default:
			return 0, false, fmt.Errorf("unhandled value %T", i)
		}
	}
}

func strToInteger(str string) (uint64, bool, error) {
	iv, err := strconv.ParseInt(str, 0, 64)
	if err == nil {
		return uint64(iv), false, nil
	}

	// maybe it's a float?
	fv, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return math.Float64bits(fv), true, nil
	}

	return math.Float64bits(math.NaN()), true, strconv.ErrSyntax
}

// String type applicable to
// - ~string
// - ~[]~byte
// - number (formated as decimal string)
// - bool value (converted to "true" or "false")
// - fmt.Stringer
type String any

func toString(s String) (_ string, err error) {
	switch t := s.(type) {
	case []byte:
		return stringhelper.Convert[string, byte](t), nil
	case string:
		return t, nil
	case []rune:
		return string(t), nil

	case int:
		return strconv.FormatInt(int64(t), 10), nil
	case uint:
		return strconv.FormatUint(uint64(t), 10), nil

	case int8:
		return strconv.FormatInt(int64(t), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(t), 10), nil

	case int16:
		return strconv.FormatInt(int64(t), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(t), 10), nil

	case int32:
		return strconv.FormatInt(int64(t), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(t), 10), nil

	case int64:
		return strconv.FormatInt(int64(t), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(t), 10), nil

	case uintptr:
		return strconv.FormatUint(uint64(t), 10), nil

	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(float64(t), 'f', -1, 64), nil

	case bool:
		return strconv.FormatBool(t), nil

	case fmt.Stringer:
		return t.String(), nil
	case nil:
		return
	default:
		switch val := reflect.Indirect(reflect.ValueOf(s)); val.Kind() {
		case reflect.String:
			return val.String(), nil
		case reflect.Slice, reflect.Array:
			switch typ := val.Elem().Type(); typ.Kind() {
			case reflect.Uint8: // []byte
				return *(*string)(val.Addr().UnsafePointer()), nil
			case reflect.Int32: // []rune
				data, ok := val.Interface().([]rune)
				if ok {
					return string(data), nil
				}

				fallthrough
			default:
				err = fmt.Errorf("unsupported converting %T to string", s)
				return
			}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			return strconv.FormatInt(val.Int(), 10), nil
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			return strconv.FormatUint(val.Uint(), 10), nil
		case reflect.Float32, reflect.Float64:
			return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
		case reflect.Bool:
			return strconv.FormatBool(val.Bool()), nil
		default:
			err = fmt.Errorf("unknown source string type %T", s)
			return
		}
	}
}

func anyToStrings(s any) (ret []string, err error) {
	switch t := s.(type) {
	case []string:
		return t, nil
	case []any:
		return toStrings(t)
	case []String:
		return toStrings(t)

	case []int:
		return toStrings(t)
	case []uint:
		return toStrings(t)

	case []int8:
		return toStrings(t)
	case []uint8:
		return toStrings(t)

	case []int16:
		return toStrings(t)
	case []uint16:
		return toStrings(t)

	case []int32:
		return toStrings(t)
	case []uint32:
		return toStrings(t)

	case []int64:
		return toStrings(t)
	case []uint64:
		return toStrings(t)

	case []uintptr:
		return toStrings(t)

	case []float32:
		return toStrings(t)
	case []float64:
		return toStrings(t)

	default:
		switch val := reflect.Indirect(reflect.ValueOf(s)); val.Kind() {
		case reflect.Slice, reflect.Array:
			n := val.Len()
			tmp := make([]any, n)
			for i := 0; i < n; i++ {
				tmp[i] = val.Index(i).Interface()
			}

			return toStrings(tmp)
		default:
			err = fmt.Errorf("invalid non slice type %T", t)
			return
		}
	}
}

func toStrings[T String](strArr []T) (ret []string, err error) {
	if len(strArr) == 0 {
		return
	}

	ret = make([]string, len(strArr))
	for i, v := range strArr {
		ret[i], err = toString(v)
	}

	return
}

// Bytes type applicable to
// - ~string
// - ~[]~byte
// - integers/floats (formated as decimal string bytes)
// - bool value (converted to bytes of "true" or "false")
// - io.Reader
// - fmt.Stringer
type Bytes any

// toBytes is toBytesOrReader but reads all data in io.Reader
func toBytes(data Bytes) ([]byte, error) {
	b, r, isReader, err := toBytesOrReader(data)
	if err != nil {
		return nil, err
	}

	if isReader {
		b, err = io.ReadAll(r)
		if err != nil {
			return b, err
		}
	}

	return b, nil
}

// toBytesOrReader
func toBytesOrReader(data Bytes) (b []byte, r io.Reader, isReader bool, err error) {
	switch t := data.(type) {
	case []byte:
		b = t
		return
	case string:
		b = strToBytes(t)
		return

	case io.Reader:
		r, isReader = t, true
		return

	case int:
		b = strconv.AppendInt(b, int64(t), 10)
		return
	case uint:
		b = strconv.AppendUint(b, uint64(t), 10)
		return

	case int8:
		b = strconv.AppendInt(b, int64(t), 10)
		return
	case uint8:
		b = strconv.AppendUint(b, uint64(t), 10)
		return

	case int16:
		b = strconv.AppendInt(b, int64(t), 10)
		return
	case uint16:
		b = strconv.AppendUint(b, uint64(t), 10)
		return

	case int32:
		b = strconv.AppendInt(b, int64(t), 10)
		return
	case uint32:
		b = strconv.AppendUint(b, uint64(t), 10)
		return

	case int64:
		b = strconv.AppendInt(b, int64(t), 10)
		return
	case uint64:
		b = strconv.AppendUint(b, uint64(t), 10)
		return

	case uintptr:
		b = strconv.AppendUint(b, uint64(t), 10)
		return

	case float32:
		b = strconv.AppendFloat(b, float64(t), 'f', -1, 64)
		return
	case float64:
		b = strconv.AppendFloat(b, float64(t), 'f', -1, 64)
		return

	case bool:
		b = strconv.AppendBool(b, t)
		return

	case fmt.Stringer:
		b = strToBytes(t.String())
		return

	case nil:
		return

	default:
		switch val := reflect.Indirect(reflect.ValueOf(data)); val.Kind() {
		case reflect.String:
			b = strToBytes(val.String())
			return
		case reflect.Slice, reflect.Array:
			switch typ := val.Elem().Type(); typ.Kind() {
			case reflect.Uint8:
				b = *(*[]byte)(val.Addr().UnsafePointer())
				return
			default:
				err = fmt.Errorf("unsupported bytes source type %T", data)
				return
			}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			b = strconv.AppendInt(b, val.Int(), 10)
			return
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			b = strconv.AppendUint(b, val.Uint(), 10)
			return
		case reflect.Float32, reflect.Float64:
			b = strconv.AppendFloat(b, val.Float(), 'f', -1, 64)
			return
		case reflect.Bool:
			b = strconv.AppendBool(b, val.Bool())
			return
		default:
			err = fmt.Errorf("unsupported bytes source type %T", data)
			return
		}
	}
}

// strToBytes casts string directly to []byte, return value is not expected to be modified
func strToBytes(s string) []byte {
	n := len(s)
	if n == 0 {
		return nil
	}

	return unsafe.Slice(
		(*byte)(unsafe.Pointer(
			(*reflect.StringHeader)(
				unsafe.Pointer(&s),
			).Data),
		),
		n,
	)
}

// Time type applicable to
// - ~string
// - ~[]~byte
// - integers/floats (as seconds since unix epoch)
type Time any

// toTimeDefault tries RFC3339Nano format on string like values
func toTimeDefault(t Time) (time.Time, error) {
	return parseTime(time.RFC3339Nano, t, nil)
}

func parseTime(layout string, t Time, loc *time.Location) (ret time.Time, err error) {
	// this funcion is called at most onece
	parseTime := time.Parse
	if loc != nil {
		parseTime = func(layout, value string) (ret time.Time, err error) {
			ret, err = time.ParseInLocation(layout, value, loc)
			return
		}
	}

	switch t := t.(type) {
	case time.Time:
		ret = t
	case *time.Time:
		ret = *t

	case string:
		ret, err = parseTime(layout, t)
		if err != nil {
			return
		}

	case []byte:
		ret, err = parseTime(layout, stringhelper.Convert[string, byte](t))
		if err != nil {
			return
		}

	case int:
		ret = time.Unix(int64(t), 0)
	case uint:
		ret = time.Unix(int64(t), 0)

	case int8:
		ret = time.Unix(int64(t), 0)
	case uint8:
		ret = time.Unix(int64(t), 0)

	case int16:
		ret = time.Unix(int64(t), 0)
	case uint16:
		ret = time.Unix(int64(t), 0)

	case int32:
		ret = time.Unix(int64(t), 0)
	case uint32:
		ret = time.Unix(int64(t), 0)

	case int64:
		ret = time.Unix(t, 0)
	case uint64:
		ret = time.Unix(int64(t), 0)

	case uintptr:
		ret = time.Unix(int64(t), 0)

	case float32:
		ret = time.Unix(0, int64(float64(t)*float64(time.Second)))
	case float64:
		ret = time.Unix(0, int64(t*float64(time.Second)))

	case nil:
		// zero time expected

	default:
		switch val := reflect.Indirect(reflect.ValueOf(t)); val.Kind() {
		case reflect.String:
			ret, err = parseTime(layout, val.String())
			if err != nil {
				return
			}
		case reflect.Slice, reflect.Array:
			switch typ := val.Elem().Type(); typ.Kind() {
			case reflect.Uint8:
				ret, err = parseTime(layout, *(*string)(val.Addr().UnsafePointer()))
				if err != nil {
					return
				}
			default:
				err = fmt.Errorf("unsupported time source type %T", t)
				return
			}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			ret = time.Unix(val.Int(), 0)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			ret = time.Unix(int64(val.Uint()), 0)
		case reflect.Float32, reflect.Float64:
			ret = time.Unix(0, int64(val.Float()*float64(time.Second)))
		default:
			err = fmt.Errorf("unsupported time source type %T", t)
			return
		}
	}

	if loc == nil {
		return
	}

	return ret.In(loc), nil
}

// Duration type applicable to
// - ~string (parsed using time.ParseDuration)
// - ~[]~byte (parsed using time.ParseDuration)
// - integers (as nanoseconds)
// - floats (as seconds)
type Duration any

func parseDuration(d Duration) (time.Duration, error) {
	switch d := d.(type) {
	case time.Duration:
		return d, nil
	case string:
		return strToDuration(d)
	case []byte:
		return time.ParseDuration(stringhelper.Convert[string, byte](d))

	case int:
		return time.Duration(d), nil
	case uint:
		return time.Duration(d), nil

	case int8:
		return time.Duration(d), nil
	case uint8:
		return time.Duration(d), nil

	case int64:
		return time.Duration(d), nil
	case uint64:
		return time.Duration(d), nil

	case int16:
		return time.Duration(d), nil
	case uint16:
		return time.Duration(d), nil

	case int32:
		return time.Duration(d), nil
	case uint32:
		return time.Duration(d), nil

	case uintptr:
		return time.Duration(d), nil

	case float32:
		return time.Duration(float64(d) * float64(time.Second)), nil
	case float64:
		return time.Duration(d * float64(time.Second)), nil

	case nil:
		return 0, nil

	default:
		switch val := reflect.Indirect(reflect.ValueOf(d)); val.Kind() {
		case reflect.String:
			return strToDuration(val.String())
		case reflect.Slice, reflect.Array:
			switch typ := val.Elem().Type(); typ.Kind() {
			case reflect.Uint8:
				return time.ParseDuration(*(*string)(val.Addr().UnsafePointer()))
			default:
				return 0, fmt.Errorf("unsupported time source type %T", d)
			}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			return time.Duration(val.Int()), nil
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			return time.Duration(val.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return time.Duration(val.Float() * float64(time.Second)), nil
		default:
			return 0, fmt.Errorf("unsupported time source type %T", d)
		}
	}
}

func strToDuration(s string) (ret time.Duration, err error) {
	ret, err = time.ParseDuration(s)
	if err != nil {
		if strings.Contains(err.Error(), "missing unit") {
			iv, isFloat, err2 := parseNumber(s)
			if err2 != nil {
				err = err2
				return
			}
			if isFloat {
				return time.Duration(math.Float64frombits(iv) * float64(time.Second)), nil
			}

			return time.Duration(iv), nil
		}

		return
	}

	return
}

// Map type is applicable to
// - map[string]string
// - map[string]any
// - map[any]any
type Map any

// Slice type is applicable to
// - []string
// - []integer
// - []float
// - []any
type Slice any
