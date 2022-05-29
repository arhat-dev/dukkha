package templateutils

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"unsafe"

	"arhat.dev/pkg/stringhelper"
)

type mathNS struct{}

func (mathNS) Abs(v Number) (_ Number, err error) {
	return
}

// Seq generates a sequence like unix cli `seq` but with step as the third argument
//
// Seq(end Number): generate sequence 0..end, step is set to 1 if end > 0 otherwise step is set to -1
//
// Seq(start, end Number): generate sequence start...end, step is set to 1 if start < end, otherwise step is set to -1
//
// Seq(start, end, step Number): generate sequence start...end with specified step
//
// start is inclusive, end is not inclusive,
// so
//  Seq(0) will generate an empty sequence
// 	Seq(1) will generate sequence [0]
//
// when step == 0, return an empty sequence (not nil)
// when start <= end && step < 0, return a nil sequence
// when start > end && step > 0, return a nil sequence
func (mathNS) Seq(args ...Number) (ret []int64, err error) {
	n := len(args)
	if n == 0 {
		err = errAtLeastOneArgGotZero
		return
	}

	var (
		start, end, step int64
	)

	params, err := toIntegers[int64](args)
	if err != nil {
		return
	}

	switch n {
	case 1:
		start, end = 0, params[0]
		if end < 0 {
			step = -1
		} else {
			step = 1
		}
	case 2:
		start, end = params[0], params[1]

		if start > end {
			step = -1
		} else {
			step = 1
		}
	default:
		start, end, step = params[0], params[1], params[2]
		if step == 0 {
			return []int64{}, nil
		}
	}

	var (
		low, high, factor int64
	)
	if start > end {
		if step > 0 {
			return nil, nil
		}

		low, high, factor = end, start, -step
	} else {
		if step < 0 {
			return nil, nil
		}

		low, high, factor = start, end, step
	}

	sz := (high - low) / factor
	if (high-low)%factor != 0 {
		sz++
	}

	ret = make([]int64, sz)

	i := int64(0)
	for val := start; i < sz; val += step {
		ret[i] = val
		i++
	}

	return
}

func (mathNS) Min(args ...Number) (_ any, err error) {
	n := len(args)
	if n == 0 {
		err = errAtLeastOneArgGotZero
		return
	}

	var (
		less   bool
		minIdx int
	)

	for i := range args {
		less, err = lessThan(args[i], args[minIdx])
		if err != nil {
			return
		}

		if less {
			minIdx = i
		}
	}

	return args[minIdx], nil
}

func (mathNS) Max(args ...Number) (_ any, err error) {
	n := len(args)
	if n == 0 {
		err = errAtLeastOneArgGotZero
		return
	}

	var (
		less   bool
		maxIdx int
	)

	for i := range args {
		less, err = lessThan(args[maxIdx], args[i])
		if err != nil {
			return
		}

		if less {
			maxIdx = i
		}
	}

	return args[maxIdx], nil
}

func (mathNS) Floor(v Number) (_ float64, err error) {
	f, err := parseFloat[float64](v)
	if err != nil {
		return
	}

	return math.Floor(f), nil
}

func (mathNS) Ceil(v Number) (_ float64, err error) {
	f, err := parseFloat[float64](v)
	if err != nil {
		return
	}

	return math.Ceil(f), nil
}

func (mathNS) Round(v Number) (_ float64, err error) {
	f, err := parseFloat[float64](v)
	if err != nil {
		return
	}

	return math.Round(f), nil
}

func (ns mathNS) Add1(v Number) (Number, error)   { return ns.Add(1, v) }
func (ns mathNS) Sub1(v Number) (Number, error)   { return ns.Sub(1, v) }
func (ns mathNS) Half(v Number) (Number, error)   { return ns.Div(2, v) }
func (ns mathNS) Double(v Number) (Number, error) { return ns.Mul(2, v) }

// Add sums all arguments
//
// return type is determined by the last argument
// that is, when it's a float number, return float64
// when it's a integer number, return int64
// when it's a string, return string
//
// when the last argument is string, precision is preseved
func (mathNS) Add(v ...Number) (_ Number, err error) {
	var ret any

	err = forEachNumber(v,
		func(last any) { ret = last },
		func(i int64) { ret = ret.(int64) + i },
		func(u uint64) { ret = ret.(uint64) + u },
		func(f float64) { ret = ret.(float64) + f },
		func(f *big.Float) { ret = ret.(*big.Float).Add(ret.(*big.Float), f) },
	)
	if err != nil {
		return
	}

	switch t := ret.(type) {
	case *big.Float:
		return t.Text('f', -1), nil
	default:
		return t, nil
	}
}

// Div divides the last argument by all but the last arguments
//
// see Add for details about the return type
func (mathNS) Div(v ...Number) (_ Number, err error) {
	var ret any

	err = forEachNumber(v,
		func(last any) { ret = last },
		func(i int64) { ret = ret.(int64) / i },
		func(u uint64) { ret = ret.(uint64) / u },
		func(f float64) { ret = ret.(float64) / f },
		func(f *big.Float) { ret = ret.(*big.Float).Quo(ret.(*big.Float), f) },
	)
	if err != nil {
		return
	}

	switch t := ret.(type) {
	case *big.Float:
		return t.Text('f', -1), nil
	default:
		return t, nil
	}
}

// Sub subtracts the last argument by all but the last arguments
//
// see Add for details about the return type
func (mathNS) Sub(v ...Number) (_ Number, err error) {
	var ret any

	err = forEachNumber(v,
		func(last any) { ret = last },
		func(i int64) { ret = ret.(int64) - i },
		func(u uint64) { ret = ret.(uint64) - u },
		func(f float64) { ret = ret.(float64) - f },
		func(f *big.Float) { ret = ret.(*big.Float).Sub(ret.(*big.Float), f) },
	)
	if err != nil {
		return
	}

	switch t := ret.(type) {
	case *big.Float:
		return t.Text('f', -1), nil
	default:
		return t, nil
	}
}

// Mod does the modulus operation to the last argument by all but the last arguments
//
// see Add for details about the return type
func (mathNS) Mod(v ...Number) (_ Number, err error) {
	var ret any

	err = forEachNumber(v,
		func(last any) { ret = last },
		func(i int64) { ret = ret.(int64) % i },
		func(u uint64) { ret = ret.(uint64) % u },
		func(f float64) { ret = math.Mod(ret.(float64), f) },
		// TODO: implement math.Mod for big.Float
		func(f *big.Float) {
			a, _ := ret.(*big.Float).Float64()
			b, _ := f.Float64()
			ret.(*big.Float).SetFloat64(math.Mod(a, b))
		},
	)
	if err != nil {
		return
	}

	switch t := ret.(type) {
	case *big.Float:
		return t.Text('f', -1), nil
	default:
		return t, nil
	}
}

// Mul multiplies the last argument by all but the last arguments
//
// see Add for details about the return type
func (mathNS) Mul(v ...Number) (_ Number, err error) {
	var ret any

	err = forEachNumber(v,
		func(last any) { ret = last },
		func(i int64) { ret = ret.(int64) * i },
		func(u uint64) { ret = ret.(uint64) * u },
		func(f float64) { ret = ret.(float64) * f },
		func(f *big.Float) { ret = ret.(*big.Float).Mul(ret.(*big.Float), f) },
	)
	if err != nil {
		return
	}

	switch t := ret.(type) {
	case *big.Float:
		return t.Text('f', -1), nil
	default:
		return t, nil
	}
}

func (mathNS) Pow(v ...Number) (_ Number, err error) {
	var ret any

	err = forEachNumber(v,
		func(last any) { ret = last },
		func(i int64) { ret = int64(math.Pow(float64(ret.(int64)), float64(i))) },
		func(u uint64) { ret = uint64(math.Pow(float64(ret.(uint64)), float64(u))) },
		func(f float64) { ret = math.Pow(ret.(float64), f) },
		func(f *big.Float) {
			a, _ := ret.(*big.Float).Float64()
			b, _ := f.Float64()
			ret.(*big.Float).SetFloat64(math.Pow(a, b))
		},
	)

	if err != nil {
		return
	}

	switch t := ret.(type) {
	case *big.Float:
		return t.Text('f', -1), nil
	default:
		return t, nil
	}
}

func (mathNS) LogE(v Number) (_ float64, err error) {
	f, err := parseFloat[float64](v)
	if err != nil {
		return
	}

	return math.Log(f), nil
}

func (mathNS) Log10(v Number) (_ float64, err error) {
	f, err := parseFloat[float64](v)
	if err != nil {
		return
	}

	return math.Log10(f), nil
}

func (mathNS) Log2(v Number) (_ float64, err error) {
	f, err := parseFloat[float64](v)
	if err != nil {
		return
	}

	return math.Log2(f), nil
}

// nolint:gocyclo
func forEachNumber(
	args []Number,
	// type of last is one of [int64, uint64, float64, *big.Float],
	// all following operation is the same type
	handleLast func(last any),
	handleOpInt64 func(int64),
	handleOpUint64 func(uint64),
	handleOpFloat64 func(float64),
	handleOpBigFloat func(*big.Float),
) error {
	n := len(args)
	if n == 0 {
		return errAtLeastOneArgGotZero
	}

	var last any
	switch t := args[n-1].(type) {
	case string:
		var lastF big.Float
		lastF, err := toBigFloat(t)
		if err != nil {
			return err
		}

		last = &lastF

	case int:
		last = int64(t)
	case uint:
		last = uint64(t)

	case int8:
		last = int64(t)
	case uint8:
		last = uint64(t)

	case int16:
		last = int64(t)
	case uint16:
		last = uint64(t)

	case int32:
		last = int64(t)
	case uint32:
		last = uint64(t)

	case int64:
		// nolint:unconvert
		last = int64(t)
	case uint64:
		// nolint:unconvert
		last = uint64(t)

	case uintptr:
		last = uint64(t)

	case float32:
		last = float64(t)
	case float64:
		// nolint:unconvert
		last = float64(t)

	default:
		switch val := reflect.Indirect(reflect.ValueOf(t)); val.Kind() {
		case reflect.String:
			lastF, err := toBigFloat(t)
			if err != nil {
				return err
			}

			last = &lastF

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			last = val.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			last = val.Uint()
		case reflect.Float32, reflect.Float64:
			last = val.Float()
		default:
			return fmt.Errorf("unsupported arithmetic operation on %T", t)
		}
	}

	handleLast(last)

	for i := n - 2; i >= 0; i-- {
		switch last.(type) {
		case *big.Float:
			f, err := toBigFloat(args[i])
			if err != nil {
				return err
			}

			handleOpBigFloat(&f)
		case float64:
			f, err := parseFloat[float64](args[i])
			if err != nil {
				return err
			}

			handleOpFloat64(f)
		case int64:
			f, err := parseInteger[int64](args[i])
			if err != nil {
				return err
			}

			handleOpInt64(f)
		case uint64:
			f, err := parseInteger[uint64](args[i])
			if err != nil {
				return err
			}

			handleOpUint64(f)
		default:
			panic("unreachable")
		}
	}

	return nil
}

func toBigFloat(v any) (ret big.Float, err error) {
	ret.SetMode(big.ToNearestEven).SetPrec(128)

	switch t := v.(type) {
	case string:
		_, _, err = ret.Parse(t, 0)
		return

	case []byte:
		_, _, err = ret.Parse(stringhelper.Convert[string, byte](t), 0)
		return

	case int:
		ret.SetInt64(int64(t))
		return
	case uint:
		ret.SetUint64(uint64(t))
		return

	case int8:
		ret.SetInt64(int64(t))
		return
	case uint8:
		ret.SetUint64(uint64(t))
		return

	case int16:
		ret.SetInt64(int64(t))
		return
	case uint16:
		ret.SetUint64(uint64(t))
		return

	case int32:
		ret.SetInt64(int64(t))
		return
	case uint32:
		ret.SetUint64(uint64(t))
		return

	case int64:
		ret.SetInt64(t)
		return
	case uint64:
		ret.SetUint64(t)
		return

	case uintptr:
		ret.SetUint64(uint64(t))
		return

	case float32:
		ret.SetFloat64(float64(t))
		return
	case float64:
		ret.SetFloat64(t)
		return

	default:
		switch val := reflect.Indirect(reflect.ValueOf(v)); val.Kind() {
		case reflect.String:
			_, _, err = ret.Parse(val.String(), 0)
			return
		case reflect.Array, reflect.Slice:
			sz := val.Len()
			if sz == 0 {
				return
			}

			if val.Type().Elem().Kind() != reflect.Uint8 {
				err = fmt.Errorf("unsupported non bytes slice %T", t)
				return
			}

			_, _, err = ret.Parse(stringhelper.Convert[string, byte](
				unsafe.Slice((*byte)(val.Index(0).Addr().UnsafePointer()), sz)), 0,
			)

			return
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ret.SetInt64(val.Int())
			return
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			ret.SetUint64(val.Uint())
			return
		case reflect.Float32, reflect.Float64:
			ret.SetFloat64(val.Float())
			return
		default:
			err = fmt.Errorf("unsupported arithmetic operation on %T", t)
			return
		}
	}
}
