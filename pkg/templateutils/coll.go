package templateutils

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"

	"arhat.dev/pkg/sorthelper"
	"arhat.dev/pkg/stringhelper"
	"arhat.dev/rs"
	"github.com/mitchellh/copystructure"
)

// collNS for collections (map, slice)
type collNS struct{}

func (collNS) List(v ...any) []any                   { return v }
func (collNS) Bools(v ...Bool) ([]bool, error)       { return toBools[bool](v) }
func (collNS) Uints(v ...Number) ([]uint64, error)   { return toIntegers[uint64](v) }
func (collNS) Ints(v ...Number) ([]int64, error)     { return toIntegers[int64](v) }
func (collNS) Floats(v ...Number) ([]float64, error) { return toFloats[float64](v) }
func (collNS) Strings(v ...String) ([]string, error) { return toStrings(v) }

// Slice operation on slice/array (the last argument)
//
// Slice(s Slice): s[:]
//
// Slice(start Number, s Slice): s[start:]
//
// Slice(start, end Number, s Slice): s[start:end]
//
// Slice(start, end, newCap Number, s Slice): s[start:end:cap]
func (collNS) Slice(args ...any) (_ Slice, err error) {
	n := len(args)
	if n == 0 {
		err = errAtLeastOneArgGotZero
		return
	}

	var (
		low, high, max int
	)

	switch n {
	default:
		max, err = parseInteger[int](args[2])
		if err != nil {
			return
		}
		fallthrough
	case 3:
		high, err = parseInteger[int](args[1])
		if err != nil {
			return
		}
		fallthrough
	case 2:
		low, err = parseInteger[int](args[0])
		if err != nil {
			return
		}
	case 1: // do nothing
	}

	if n < 4 { // max not set
		max = math.MaxInt
	}

	if n < 3 { // high not set
		high = math.MaxInt
	}

	switch t := args[n-1].(type) {
	case string:
		low, high, _ = validSliceArgs(low, high, max, len(t), len(t))
		return t[low:high], nil

	case []string:
		return slice3(t, low, high, max)

	case []int:
		return slice3(t, low, high, max)
	case []uint:
		return slice3(t, low, high, max)

	case []int8:
		return slice3(t, low, high, max)
	case []uint8:
		return slice3(t, low, high, max)

	case []int16:
		return slice3(t, low, high, max)
	case []uint16:
		return slice3(t, low, high, max)

	case []int32:
		return slice3(t, low, high, max)
	case []uint32:
		return slice3(t, low, high, max)

	case []int64:
		return slice3(t, low, high, max)
	case []uint64:
		return slice3(t, low, high, max)

	case []uintptr:
		return slice3(t, low, high, max)

	case []float32:
		return slice3(t, low, high, max)
	case []float64:
		return slice3(t, low, high, max)

	case []any:
		return slice3(t, low, high, max)

	default:
		switch val := reflect.Indirect(reflect.ValueOf(t)); val.Kind() {
		case reflect.String:
			n := val.Len()
			low, high, _ = validSliceArgs(low, high, max, n, n)
			return val.Slice(low, high).Interface(), nil
		case reflect.Slice, reflect.Array:
			n, k := val.Len(), val.Cap()
			return val.Slice3(validSliceArgs(low, high, max, n, k)).Interface(), nil
		default:
			err = fmt.Errorf("invalid slice target %T", t)
			return
		}
	}
}

// validSliceArgs adjusts low, high, max for sz and kap
func validSliceArgs(low, high, max int, sz, kap int) (l, h, m int) {
	// ensure 0 <= max <= cap(x)
	if max < 0 { // index from end
		if max+kap+1 < 0 {
			max = 0
		} else {
			max += kap + 1
		}
	}

	if max > kap || max < 0 /* overflow */ {
		max = kap
	}

	// ensure 0 <= high <= len(x)
	if high < 0 { // index from end
		if high+sz+1 < 0 {
			high = 0
		} else {
			high += sz + 1
		}
	}

	if high > sz || high < 0 /* overflow */ {
		high = sz
	}

	// ensure low >= 0
	if low < 0 { // index from end
		if low+sz+1 < 0 {
			low = 0
		} else {
			low += sz + 1
		}
	}

	// ensure high <= max
	if high > max {
		high = max
	}

	// 0 <= high <= len(t) <= max <= cap(t)

	// ensure: low <= high
	if low > high {
		low = high
	}

	return low, high, max
}

// slice3 return s[low:high:max]
//
// low, high, max are adjusted to valid values for s
//
// ref: https://go.dev/ref/spec#Slice_expressions
func slice3[T any](s []T, low, high, max int) ([]T, error) {
	low, high, max = validSliceArgs(low, high, max, len(s), cap(s))
	return s[low:high:max], nil
}

// Reverse string or slice
func (collNS) Reverse(stringOrSlice any) (_ any, err error) {
	switch t := stringOrSlice.(type) {
	case string:
		return stringhelper.Reverse[byte](t), nil
	case []string:
		return reverseSlice(t), nil
	case []any:
		return reverseSlice(t), nil

	case []int:
		return reverseSlice(t), nil
	case []uint:
		return reverseSlice(t), nil

	case []int8:
		return reverseSlice(t), nil
	case []uint8:
		return reverseSlice(t), nil

	case []int16:
		return reverseSlice(t), nil
	case []uint16:
		return reverseSlice(t), nil

	case []int32:
		return reverseSlice(t), nil
	case []uint32:
		return reverseSlice(t), nil

	case []int64:
		return reverseSlice(t), nil
	case []uint64:
		return reverseSlice(t), nil

	case []uintptr:
		return reverseSlice(t), nil

	case []float32:
		return reverseSlice(t), nil
	case []float64:
		return reverseSlice(t), nil
	default:
		switch val := reflect.Indirect(reflect.ValueOf(stringOrSlice)); val.Kind() {
		case reflect.String:
			return reflect.ValueOf(stringhelper.Reverse[byte](val.String())).Convert(val.Type()).Interface(), nil
		case reflect.Slice, reflect.Array:
			n := val.Len()
			for i := 0; i < (n+1)/2; i++ {
				vi := val.Index(i)
				tmp := val.Index(n - i - 1).Interface()

				val.Index(n - i - 1).Set(vi)
				vi.Set(reflect.ValueOf(tmp))
			}
			return t, nil
		default:
			err = fmt.Errorf("unsupported operation on %T", t)
			return
		}
	}
}

func reverseSlice[T any](s []T) []T {
	n := len(s)
	if n < 2 {
		return s
	}

	for i := 0; i < (n+1)/2; i++ {
		s[i], s[n-i-1] = s[n-i-1], s[i]
	}

	return s
}

// Unique return a slice of unique items in s
func (collNS) Unique(s Slice) (_ any, err error) {
	switch t := s.(type) {
	case []string:
		return typedSliceUnique(t), nil

	case []int:
		return typedSliceUnique(t), nil
	case []uint:
		return typedSliceUnique(t), nil

	case []int8:
		return typedSliceUnique(t), nil
	case []uint8:
		return typedSliceUnique(t), nil

	case []int16:
		return typedSliceUnique(t), nil
	case []uint16:
		return typedSliceUnique(t), nil

	case []int32:
		return typedSliceUnique(t), nil
	case []uint32:
		return typedSliceUnique(t), nil

	case []int64:
		return typedSliceUnique(t), nil
	case []uint64:
		return typedSliceUnique(t), nil

	case []uintptr:
		return typedSliceUnique(t), nil

	case []float32:
		return typedSliceUnique(t), nil
	case []float64:
		return typedSliceUnique(t), nil

	case []any:
		return anySliceUnique(t), nil

	case map[string]any:
		return t, nil
	case map[string]string:
		return t, nil
	case map[string]struct{}:
		return t, nil
	case map[any]any:
		return t, nil

	default:
		switch val := reflect.Indirect(reflect.ValueOf(s)); val.Kind() {
		case reflect.Map:
			// TODO: TBD map keys are unique, shall we support this?
			return s, nil
		case reflect.Slice, reflect.Array:
			var (
				ret = reflect.Zero(val.Type())
				m   = make(map[any]None)
				ok  bool
			)

			n := val.Len()
			for i := 0; i < n; i++ {
				v := val.Index(i)
				if !v.IsValid() {
					continue
				}

				obj := v.Interface()
				_, ok = m[obj]
				if ok {
					continue
				}

				m[obj] = None{}
				ret = reflect.Append(ret, v)
			}

			return ret.Interface(), nil
		default:
			err = fmt.Errorf("invalid unique operation on %T", s)
			return
		}
	}
}

func anySliceUnique(list []any) (ret []any) {
	var (
		m  = make(map[any]None)
		ok bool
	)

	for _, v := range list {
		_, ok = m[v]
		if ok {
			continue
		}

		m[v] = None{}
		ret = append(ret, v)
	}

	return
}

func typedSliceUnique[T comparable](list []T) (ret []T) {
	var (
		m  = make(map[T]None)
		ok bool
	)

	for _, v := range list {
		_, ok = m[v]
		if ok {
			continue
		}

		m[v] = None{}
		ret = append(ret, v)
	}

	return
}

// Merge maps into a new map, type of the returned map is determined by the last argument
// it merges in last to first order
//
// nolint:gocyclo
func (collNS) Merge(maps ...Map) (_ Map, err error) {
	n := len(maps)
	if n < 2 {
		err = fmt.Errorf("at least 2 args expected, got %d", n)
		return
	}

	switch maps[n-1].(type) {
	case map[string]any:
		ret := make(map[string]any)

		for i := n - 1; i >= 0; i-- {
			switch t := maps[i].(type) {
			case map[string]any:
				ret, err = rs.MergeMap(ret, t, true, false)
				if err != nil {
					return
				}
			case map[string]string:
				for k, v := range t {
					ret[k] = v
				}
			case map[string]struct{}:
				for k, v := range t {
					ret[k] = v
				}
			case map[any]any:
				var key string

				in := make(map[string]any, len(t))
				for k, v := range t {
					key, err = toString(k)
					if err != nil {
						return
					}

					in[key] = v
				}

				ret, err = rs.MergeMap(ret, in, true, false)
				if err != nil {
					return
				}
			default:
				err = fmt.Errorf("incompatible map %T with target %T", t, ret)
				return
			}
		}

		return ret, nil
	case map[string]string:
		ret := make(map[string]string)
		for i := n - 1; i >= 0; i-- {
			switch t := maps[i].(type) {
			case map[string]any:
				for k, v := range t {
					ret[k], err = toString(v)
					if err != nil {
						return
					}
				}
			case map[string]string:
				for k, v := range t {
					ret[k] = v
				}
			case map[string]struct{}:
				for k := range t {
					ret[k] = ""
				}
			case map[any]any:
				var key string
				for k, v := range t {
					key, err = toString(k)
					if err != nil {
						return
					}

					ret[key], err = toString(v)
					if err != nil {
						return
					}
				}
			default:
				err = fmt.Errorf("incompatible map %T with target %T", t, ret)
				return
			}
		}

		return ret, nil
	case map[string]struct{}:
		ret := make(map[string]struct{})
		for i := n - 1; i >= 0; i-- {
			switch t := maps[i].(type) {
			case map[string]any:
				for k := range t {
					ret[k] = struct{}{}
				}
			case map[string]string:
				for k := range t {
					ret[k] = struct{}{}
				}
			case map[string]struct{}:
				for k := range t {
					ret[k] = struct{}{}
				}
			case map[any]any:
				var key string
				for k := range t {
					key, err = toString(k)
					if err != nil {
						return
					}

					ret[key] = struct{}{}
				}
			default:
				err = fmt.Errorf("incompatible map %T with target %T", t, ret)
				return
			}
		}

		return ret, nil
	case map[any]any:
		ret := make(map[any]any)
		for i := n - 1; i >= 0; i-- {
			switch t := maps[i].(type) {
			case map[string]any:
				ret, err = mergeIntoAnyMap(ret, t, true, false)
				if err != nil {
					return
				}
			case map[string]string:
				for k, v := range t {
					ret[k] = v
				}
			case map[string]struct{}:
				for k := range t {
					ret[k] = struct{}{}
				}
			case map[any]any:
				ret, err = mergeBothAnyMap(ret, t, true, false)
				if err != nil {
					return
				}
			default:
				err = fmt.Errorf("incompatible map %T with target %T", t, ret)
				return
			}
		}

		return ret, nil
	default:
		err = fmt.Errorf("unsupported merge target type %T", maps[n-1])
		return
	}
}

func mergeIntoAnyMap[K comparable](
	original map[any]any, additional map[K]any,

	// options
	appendList bool,
	uniqueInListItems bool,
) (map[any]any, error) {
	out := make(map[any]any, len(original))
	for k, v := range original {
		out[k] = v
	}

	var err error
	for k, v := range additional {
		switch newVal := v.(type) {
		case map[K]any:
			if originalVal, ok := out[k]; ok {
				if orignalMap, ok := originalVal.(map[any]any); ok {
					out[k], err = mergeIntoAnyMap(orignalMap, newVal, appendList, uniqueInListItems)
					if err != nil {
						return nil, err
					}

					continue
				} else {
					return nil, fmt.Errorf("unexpected non map data %v: %v", k, orignalMap)
				}
			} else {
				out[k] = newVal
			}
		case []any:
			if originalVal, ok := out[k]; ok {
				if originalList, ok := originalVal.([]any); ok {
					if appendList {
						originalList = append(originalList, newVal...)
					} else {
						originalList = newVal
					}

					if uniqueInListItems {
						originalList = rs.UniqueList(originalList)
					}

					out[k] = originalList

					continue
				} else {
					return nil, fmt.Errorf("unexpected non list data %v: %v", k, originalList)
				}
			} else {
				out[k] = newVal
			}
		default:
			out[k] = newVal
		}
	}

	return out, nil
}

func mergeBothAnyMap(
	original, additional map[any]any,

	// options
	appendList bool,
	uniqueInListItems bool,
) (map[any]any, error) {
	out := make(map[any]any, len(original))
	for k, v := range original {
		out[k] = v
	}

	var err error
	for k, v := range additional {
		switch newVal := v.(type) {
		case map[any]any:
			if originalVal, ok := out[k]; ok {
				if orignalMap, ok := originalVal.(map[any]any); ok {
					out[k], err = mergeBothAnyMap(orignalMap, newVal, appendList, uniqueInListItems)
					if err != nil {
						return nil, err
					}

					continue
				} else {
					return nil, fmt.Errorf("unexpected non map data %q: %v", k, orignalMap)
				}
			} else {
				out[k] = newVal
			}
		case []any:
			if originalVal, ok := out[k]; ok {
				if originalList, ok := originalVal.([]any); ok {
					if appendList {
						originalList = append(originalList, newVal...)
					} else {
						originalList = newVal
					}

					if uniqueInListItems {
						originalList = rs.UniqueList(originalList)
					}

					out[k] = originalList

					continue
				} else {
					return nil, fmt.Errorf("unexpected non list data %q: %v", k, originalList)
				}
			} else {
				out[k] = newVal
			}
		default:
			out[k] = newVal
		}
	}

	return out, nil
}

// Index slice or map, return indexed element
//
// for slices, idx may be < 0 to index from end
//
// nolint:gocyclo
func (collNS) Index(idxOrKey any, sliceOrMap any) (_ any, err error) {
	switch t := sliceOrMap.(type) {
	case nil:
		return "", nil
	case map[any]any:
		return t[idxOrKey], nil
	}

	var (
		val reflect.Value
		n   int
		get func(int) any
	)
	i, err := parseInteger[int](idxOrKey)
	if err != nil {
		goto indexMap
	}

	switch t := sliceOrMap.(type) {
	case string:
		n, get = len(t), func(i int) any { return t[i] }
	case []string:
		n, get = len(t), func(i int) any { return t[i] }

	case []int:
		n, get = len(t), func(i int) any { return t[i] }
	case []uint:
		n, get = len(t), func(i int) any { return t[i] }

	case []int8:
		n, get = len(t), func(i int) any { return t[i] }
	case []uint8:
		n, get = len(t), func(i int) any { return t[i] }

	case []int16:
		n, get = len(t), func(i int) any { return t[i] }
	case []uint16:
		n, get = len(t), func(i int) any { return t[i] }

	case []int32:
		n, get = len(t), func(i int) any { return t[i] }
	case []uint32:
		n, get = len(t), func(i int) any { return t[i] }

	case []int64:
		n, get = len(t), func(i int) any { return t[i] }
	case []uint64:
		n, get = len(t), func(i int) any { return t[i] }

	case []uintptr:
		n, get = len(t), func(i int) any { return t[i] }

	case []float32:
		n, get = len(t), func(i int) any { return t[i] }
	case []float64:
		n, get = len(t), func(i int) any { return t[i] }

	case []any:
		n, get = len(t), func(i int) any { return t[i] }

	default:
		switch val = reflect.Indirect(reflect.ValueOf(sliceOrMap)); val.Kind() {
		case reflect.String:
			str := val.String()
			n, get = len(str), func(i int) any { return str[i] }
		case reflect.Slice, reflect.Array:
			n, get = val.Len(), func(i int) any { return val.Index(i).Interface() }
		case reflect.Map:
			goto indexMap
		default:
			err = fmt.Errorf("unsupported numeric index for %T", sliceOrMap)
			return
		}
	}

	if i < -n || i >= n {
		err = fmt.Errorf(
			"invalid index out of range: expected in range [-%d,%d), got %d", n, n, i)
		return
	}

	if i < 0 {
		i = n + i
	}

	return get(i), nil

indexMap:
	key, err := toString(idxOrKey)
	if err != nil {
		return
	}

	switch t := sliceOrMap.(type) {
	case map[string]string:
		return t[key], nil
	case map[string]any:
		return t[key], nil
	case map[string]struct{}:
		return "", nil
	default:
		if !val.IsValid() {
			val = reflect.Indirect(reflect.ValueOf(sliceOrMap))
		}

		switch val.Kind() {
		case reflect.Map:
			ret := val.MapIndex(reflect.ValueOf(idxOrKey))
			if ret.IsValid() {
				return ret.Interface(), nil
			}

			return nil, nil
		default:
			err = fmt.Errorf("indexing of %T not supported", sliceOrMap)
			return
		}
	}
}

// Dup returns a deep copy of obj
func (collNS) Dup(obj any) (_ any, err error) {
	switch t := obj.(type) {
	case []string:
		return cloneScalarSlice(t), nil

	case []int:
		return cloneScalarSlice(t), nil
	case []uint:
		return cloneScalarSlice(t), nil

	case []int8:
		return cloneScalarSlice(t), nil
	case []uint8:
		return cloneScalarSlice(t), nil

	case []int16:
		return cloneScalarSlice(t), nil
	case []uint16:
		return cloneScalarSlice(t), nil

	case []int32:
		return cloneScalarSlice(t), nil
	case []uint32:
		return cloneScalarSlice(t), nil

	case []int64:
		return cloneScalarSlice(t), nil
	case []uint64:
		return cloneScalarSlice(t), nil

	case []uintptr:
		return cloneScalarSlice(t), nil

	case []float32:
		return cloneScalarSlice(t), nil
	case []float64:
		return cloneScalarSlice(t), nil

	case map[string]string:
		return cloneScalarMap(t), nil
	case map[string]struct{}:
		return cloneScalarMap(t), nil

	default:
		return copystructure.Copy(t)
	}
}

func cloneScalarMap[K comparable, V integer | float | string | struct{}](m map[K]V) (ret map[K]V) {
	ret = make(map[K]V)
	for k, v := range m {
		ret[k] = v
	}
	return
}

func cloneScalarSlice[T integer | float | string](s []T) (ret []T) {
	ret = make([]T, len(s))
	_ = copy(ret, s)
	return
}

// Sort a slice in increasing order
func (collNS) Sort(s Slice) (_ Slice, err error) {
	switch s := s.(type) {
	case []string:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil

	case []int:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil
	case []uint:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil

	case []int8:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil
	case []uint8:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil

	case []int16:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil
	case []uint16:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil

	case []int32:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil
	case []uint32:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil

	case []int64:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil
	case []uint64:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil

	case []uintptr:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil

	case []float32:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil
	case []float64:
		sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
		return s, nil

	case []any:
		sort.SliceStable(s, func(i, j int) bool {
			if err != nil {
				return true
			}

			var less bool
			less, err = lessThan(s[i], s[j])
			if err != nil {
				return true
			}

			return less
		})

		return s, err

	default:
		switch val := reflect.Indirect(reflect.ValueOf(s)); val.Kind() {
		case reflect.Slice:
			var err2 error
			sort.SliceStable(s, func(i, j int) bool {
				var less bool
				less, err2 = reflectValueLessThan(val.Index(i), val.Index(j))
				if err2 != nil {
					err = err2
				}
				return less
			})

			return s, err
		case reflect.Array:
			sz := val.Len()
			ret := reflect.MakeSlice(reflect.SliceOf(val.Elem().Type()), sz, sz)
			for i := 0; i < sz; i++ {
				ret.Index(i).Set(val.Index(i))
			}

			var err2 error
			sort.SliceStable(ret.Interface(), func(i, j int) bool {
				var less bool
				less, err2 = reflectValueLessThan(ret.Index(i), ret.Index(j))
				if err2 != nil {
					err = err2
				}
				return less
			})

			return ret.Interface(), err
		default:
			return s, fmt.Errorf("unsupported sort operation on %T", s)
		}
	}
}

func reflectValueLessThan(a, b reflect.Value) (bool, error) {
	switch a.Kind() {
	case reflect.String:
		bstr, err := toString(b.Interface())
		if err != nil {
			return false, err
		}

		return a.String() < bstr, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if b.CanInt() {
			return a.Int() < b.Int(), nil
		}

		bi, err := parseInteger[int64](b.Interface())
		if err != nil {
			return false, err
		}

		return a.Int() < bi, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if b.CanUint() {
			return a.Uint() < b.Uint(), nil
		}

		bi, err := parseInteger[uint64](b.Interface())
		if err != nil {
			return false, err
		}

		return a.Uint() < bi, nil
	case reflect.Float32, reflect.Float64:
		if b.CanFloat() {
			return a.Float() < b.Float(), nil
		}

		bi, err := parseFloat[float64](b.Interface())
		if err != nil {
			return false, err
		}

		return a.Float() < bi, nil
	default:
		return false, fmt.Errorf("unsupported operation on %T", a.Interface())
	}
}

// return true when a < b
func lessThan(a, b any) (_ bool, err error) {
	switch t := a.(type) {
	case string:
		var str string
		str, err = toString(b)
		return t < str, err

	case int:
		var x int64
		x, err = parseInteger[int64](b)
		return int64(t) < x, err
	case uint:
		var x uint64
		x, err = parseInteger[uint64](b)
		return uint64(t) < x, err

	case int8:
		var x int64
		x, err = parseInteger[int64](b)
		return int64(t) < x, err
	case uint8:
		var x uint64
		x, err = parseInteger[uint64](b)
		return uint64(t) < x, err

	case int16:
		var x int64
		x, err = parseInteger[int64](b)
		return int64(t) < x, err
	case uint16:
		var x uint64
		x, err = parseInteger[uint64](b)
		return uint64(t) < x, err

	case int32:
		var x int64
		x, err = parseInteger[int64](b)
		return int64(t) < x, err
	case uint32:
		var x uint64
		x, err = parseInteger[uint64](b)
		return uint64(t) < x, err

	case int64:
		var x int64
		x, err = parseInteger[int64](b)
		return t < x, err
	case uint64:
		var x uint64
		x, err = parseInteger[uint64](b)
		return t < x, err

	case uintptr:
		var x uint64
		x, err = parseInteger[uint64](b)
		return uint64(t) < x, err

	case float32:
		var f float64
		f, err = parseFloat[float64](b)
		return float64(t) < f, err
	case float64:
		var f float64
		f, err = parseFloat[float64](b)
		return t < f, err
	default:
		return reflectValueLessThan(reflect.Indirect(reflect.ValueOf(a)), reflect.Indirect(reflect.ValueOf(b)))
	}
}

// Flatten a slice to some depth
//
// Flatten(s Slice): flatten all slice
//
// Flatten(depth Integer, s Slice): flatten slice with specified depth (-1 to flatten all)
func (collNS) Flatten(args ...any) (ret []any, err error) {
	n := len(args)
	if n == 0 {
		err = errAtLeastOneArgGotZero
		return
	}

	depth := -1
	if n > 1 {
		depth, err = parseInteger[int](args[0])
		if err != nil {
			return
		}
	}

	var src []any
	switch t := args[n-1].(type) {
	case []any:
		src = t
	default:
		switch val := reflect.Indirect(reflect.ValueOf(t)); val.Kind() {
		case reflect.Slice, reflect.Array:
			src = make([]any, val.Len())
			for i := range src {
				src[i] = val.Index(i).Interface()
			}
		default:
			src = []any{t}
		}
	}

	return flattenAnySlice(depth, src), nil
}

func flattenAnySlice(depth int, s []any) (ret []any) {
	if depth == 0 {
		return s
	}

	for _, elem := range s {
		switch t := elem.(type) {
		case []any:
			ret = append(ret, flattenAnySlice(depth-1, t)...)
		default:
			switch val := reflect.Indirect(reflect.ValueOf(elem)); val.Kind() {
			case reflect.Slice, reflect.Array:
				src := make([]any, val.Len())
				for i := range src {
					src[i] = val.Index(i).Interface()
				}
				ret = append(ret, flattenAnySlice(depth-1, src)...)
			default:
				ret = append(ret, elem)
			}
		}
	}

	return
}

// Pick multiple indices (all but the last argument) in container (the last argument)
// retrun same type as container
func (collNS) Pick(args ...any) (_ any, err error) { return pickOrMmit(args, true) }

// Omit multiple indices (all but the last argument) in container (the last argument)
// retrun same type as container
func (collNS) Omit(args ...any) (_ any, err error) { return pickOrMmit(args, false) }

// nolint:gocyclo
func pickOrMmit(args []any, pick bool) (_ any, err error) {
	n := len(args)
	if n < 2 {
		err = fmt.Errorf("at least 2 args expected, got %d", n)
		return
	}

	indices := args[:n-1]
	switch t := args[n-1].(type) {
	case map[any]any:
		ret := make(map[any]any)
		if pick {
			for _, k := range indices {
				v, ok := t[k]
				if ok {
					ret[k] = v
				}
			}
		} else { // omit
			var buf [1]any
			for k, v := range t {
				buf[0] = k
				if !anySliceContainsAll(indices, buf[:]) {
					ret[k] = v
				}
			}
		}
	case nil:
		return
	}

	numIndices, err := toIntegers[int](indices)
	if err != nil {
		goto handleMap
	}

	switch t := args[n-1].(type) {
	case []string:
		return pickOrOmitSlice(t, numIndices, pick)

	case []int:
		return pickOrOmitSlice(t, numIndices, pick)
	case []uint:
		return pickOrOmitSlice(t, numIndices, pick)

	case []int8:
		return pickOrOmitSlice(t, numIndices, pick)
	case []uint8:
		return pickOrOmitSlice(t, numIndices, pick)

	case []int16:
		return pickOrOmitSlice(t, numIndices, pick)
	case []uint16:
		return pickOrOmitSlice(t, numIndices, pick)

	case []int32:
		return pickOrOmitSlice(t, numIndices, pick)
	case []uint32:
		return pickOrOmitSlice(t, numIndices, pick)

	case []int64:
		return pickOrOmitSlice(t, numIndices, pick)
	case []uint64:
		return pickOrOmitSlice(t, numIndices, pick)

	case []uintptr:
		return pickOrOmitSlice(t, numIndices, pick)

	case []float32:
		return pickOrOmitSlice(t, numIndices, pick)
	case []float64:
		return pickOrOmitSlice(t, numIndices, pick)

	case []any:
		return pickOrOmitSlice(t, numIndices, pick)

	default:
		switch val := reflect.Indirect(reflect.ValueOf(t)); val.Kind() {
		case reflect.Map:
			goto handleMap
		case reflect.Array, reflect.Slice:
			// TODO
			return
		default:
			err = fmt.Errorf("unsupported operation on %T", t)
			return
		}
	}

handleMap:
	keys, err := toStrings(indices)
	if err != nil {
		return
	}

	switch t := args[n-1].(type) {
	case map[string]any:
		return pickOrOmitMap(t, keys, pick), nil
	case map[string]string:
		return pickOrOmitMap(t, keys, pick), nil
	case map[string]struct{}:
		return pickOrOmitMap(t, keys, pick), nil
	default:
		err = fmt.Errorf("unsupported operation on %T", t)
		return
	}
}

func pickOrOmitSlice[T any](s []T, indices []int, pick bool) (ret []T, err error) {
	n := len(s)

	var filterOut map[int]None
	if !pick {
		filterOut = make(map[int]None)
	}

	for _, i := range indices {
		if i < -n || i >= n {
			err = fmt.Errorf(
				"invalid index out of range: %d not in range [-%d,%d)", i, n, n)
			return
		}

		if i < 0 {
			i = n + i
		}

		if pick {
			ret = append(ret, s[i])
		} else {
			filterOut[i] = None{}
		}
	}

	if pick {
		return
	}

	at := 0
	ret = make([]T, n-len(filterOut))
	for i := range s {
		_, nopick := filterOut[i]
		if nopick {
			continue
		}

		ret[at] = s[i]
		at++
	}

	return
}

func pickOrOmitMap[K comparable, V any](t map[K]V, indices []K, pick bool) (ret map[K]V) {
	ret = make(map[K]V)
	if pick {
		for _, k := range indices {
			v, ok := t[k]
			if ok {
				ret[k] = v
			}
		}
	} else { // omit
		var buf [1]K
		for k, v := range t {
			buf[0] = k

			found := false
			for _, tk := range indices {
				if k == tk {
					found = true
					break
				}
			}

			if found {
				continue
			}

			ret[k] = v
		}
	}

	return
}

// Push is an alias of Append as stack operation push
func (ns collNS) Push(args ...any) (any, error) { return ns.Append(args...) }

// Append values (all but the last argument) to container (the last argument)
//
// Append(values..., container Slice)
func (collNS) Append(args ...any) (_ any, err error) {
	n := len(args)
	if n < 2 {
		err = fmt.Errorf("append: at least 2 args expected, got %d", n)
		return
	}

	switch t := args[n-1].(type) {
	case string:
		var val []string
		val, err = toStrings(args[:n-1])
		if err != nil {
			return
		}

		var sb strings.Builder
		sb.WriteString(t)
		for _, v := range val {
			sb.WriteString(v)
		}

		return sb.String(), nil
	case []string:
		var val []string
		val, err = toStrings(args[:n-1])
		if err != nil {
			return t, err
		}

		return append(t, val...), nil
	case []any:
		return append(t, args[:n-1]), nil

	case []int:
		return appendIntegerSlice(t, args[:n-1])
	case []uint:
		return appendIntegerSlice(t, args[:n-1])

	case []int8:
		return appendIntegerSlice(t, args[:n-1])
	case []uint8:
		return appendIntegerSlice(t, args[:n-1])

	case []int16:
		return appendIntegerSlice(t, args[:n-1])
	case []uint16:
		return appendIntegerSlice(t, args[:n-1])

	case []int32:
		return appendIntegerSlice(t, args[:n-1])
	case []uint32:
		return appendIntegerSlice(t, args[:n-1])

	case []int64:
		return appendIntegerSlice(t, args[:n-1])
	case []uint64:
		return appendIntegerSlice(t, args[:n-1])

	case []uintptr:
		return appendIntegerSlice(t, args[:n-1])

	case []float32:
		return appendFloatSlice(t, args[:n-1])
	case []float64:
		return appendFloatSlice(t, args[:n-1])

	default:
		switch ctr := reflect.Indirect(reflect.ValueOf(t)); ctr.Kind() {
		case reflect.Slice:
			// TODO: handle type conversion
			return reflect.AppendSlice(ctr, reflect.ValueOf(args[:n-1])).Interface(), nil
		default:
			err = fmt.Errorf("append: invalid non slice type %T", t)
			return
		}
	}
}

func appendIntegerSlice[T integer](s []T, v []any) ([]T, error) {
	val, err := toIntegers[T](v)
	if err != nil {
		return s, err
	}

	return append(s, val...), nil
}

func appendFloatSlice[T float](s []T, v []any) ([]T, error) {
	val, err := toFloats[T](v)
	if err != nil {
		return s, err
	}

	return append(s, val...), nil
}

// Prepend values (all but the last argument) to container (the last argument)
//
// Prepend(values..., container Slice)
func (collNS) Prepend(args ...any) (_ any, err error) {
	n := len(args)
	if n < 2 {
		err = fmt.Errorf("prepend: at least 2 args expected, got %d", n)
		return
	}

	switch t := args[n-1].(type) {
	case string:
		var val []string
		val, err = toStrings(args[:n-1])
		if err != nil {
			return
		}

		var sb strings.Builder
		for _, v := range val {
			sb.WriteString(v)
		}
		sb.WriteString(t)

		return sb.String(), nil
	case []string:
		var val []string
		val, err = toStrings(args[:n-1])
		if err != nil {
			return t, err
		}

		return append(val, t...), nil
	case []int:
		return prependIntegerSlice(t, args[:n-1])
	case []uint:
		return prependIntegerSlice(t, args[:n-1])

	case []int8:
		return prependIntegerSlice(t, args[:n-1])
	case []uint8:
		return prependIntegerSlice(t, args[:n-1])

	case []int16:
		return prependIntegerSlice(t, args[:n-1])
	case []uint16:
		return prependIntegerSlice(t, args[:n-1])

	case []int32:
		return prependIntegerSlice(t, args[:n-1])
	case []uint32:
		return prependIntegerSlice(t, args[:n-1])

	case []int64:
		return prependIntegerSlice(t, args[:n-1])
	case []uint64:
		return prependIntegerSlice(t, args[:n-1])

	case []uintptr:
		return prependIntegerSlice(t, args[:n-1])

	case []float32:
		return prependFloatSlice(t, args[:n-1])
	case []float64:
		return prependFloatSlice(t, args[:n-1])

	case []any:
		return append(args[:n-1], t...), nil

	default:
		switch ctr := reflect.Indirect(reflect.ValueOf(t)); ctr.Kind() {
		case reflect.Slice:
			// TODO: handle type conversion
			return reflect.AppendSlice(reflect.ValueOf(args[:n-1]), ctr).Interface(), nil
		default:
			err = fmt.Errorf("prepend: invalid non slice type %T", t)
			return
		}
	}
}

func prependIntegerSlice[T integer](s []T, v []any) ([]T, error) {
	val, err := toIntegers[T](v)
	if err != nil {
		return s, err
	}

	return append(val, s...), nil
}

func prependFloatSlice[T float](s []T, v []any) ([]T, error) {
	val, err := toFloats[T](v)
	if err != nil {
		return s, err
	}

	return append(val, s...), nil
}

// MapStringAny creates a map[string]any with key value pairs
func (collNS) MapStringAny(v ...any) (dict map[string]any, err error) {
	dict = make(map[string]any)

	lenv := len(v)

	var key string
	for i := 0; i < lenv; i += 2 {
		key, err = toString(v[i])
		if err != nil {
			return
		}

		if i+1 >= lenv {
			dict[key] = None{}
			break
		}

		dict[key] = v[i+1]
	}

	return
}

// MapAnyAny creates a map[any]any with key value pairs
func (collNS) MapAnyAny(v ...any) (dict map[any]any) {
	dict = make(map[any]any)
	lenv := len(v)
	for i := 0; i < lenv; i += 2 {
		if i+1 >= lenv {
			dict[v[i]] = None{}
			break
		}

		dict[v[i]] = v[i+1]
	}

	return
}

// Keys gets sorted keys in map v
func (collNS) Keys(v Map) (ret []string, err error) {
	switch t := v.(type) {
	case map[string]struct{}:
		return collectMapKeys(t)
	case map[string]any:
		return collectMapKeys(t)
	case map[string]string:
		return collectMapKeys(t)
	case map[any]any:
		i := 0
		ret = make([]string, len(t))
		for k := range t {
			ret[i], err = toString(k)
			if err != nil {
				return
			}
			i++
		}
	default:
		switch val := reflect.Indirect(reflect.ValueOf(v)); val.Kind() {
		case reflect.Map:
			keys := val.MapKeys()
			ret = make([]string, len(keys))
			for i := range keys {
				ret[i], err = toString(keys[i].Interface())
				if err != nil {
					return
				}
			}
		default:
			err = fmt.Errorf("unsupported operation on %T", v)
			return
		}
	}

	sort.SliceStable(ret, func(i, j int) bool { return ret[i] < ret[j] })
	return
}

func collectMapKeys[K comparable, V any](t map[K]V) (ret []string, err error) {
	i := 0
	ret = make([]string, len(t))
	for k := range t {
		ret[i], err = toString(k)
		if err != nil {
			return
		}

		i++
	}

	return
}

// Values gets values in map v, sorted by its key
func (collNS) Values(v Map) (ret []any, err error) {
	switch t := v.(type) {
	case map[string]struct{}:
		return nil, nil
	case map[string]any:
		return collectMapValues(t)
	case map[string]string:
		return collectMapValues(t)
	case map[any]any:
		i := 0
		sz := len(t)
		keys := make([]any, sz)
		ret = make([]any, sz)
		for k, v := range t {
			keys[i], ret[i] = k, v
			i++
		}

		sort.Stable(sorthelper.NewCustomSortable(
			func(i, j int) {
				keys[i], keys[j] = keys[j], keys[i]
				ret[i], ret[j] = ret[j], ret[i]
			},
			func(i, j int) bool {
				if err != nil {
					return true
				}

				var less bool
				less, err = lessThan(keys[i], keys[j])
				if err != nil {
					return true
				}
				return less
			},
			func() int { return len(keys) },
		))

		if err != nil {
			return
		}

		return
	default:
		switch val := reflect.Indirect(reflect.ValueOf(v)); val.Kind() {
		case reflect.Map:
			reflectKeys := val.MapKeys()

			sort.Stable(sorthelper.NewCustomSortable(
				func(i, j int) { reflectKeys[i], reflectKeys[j] = reflectKeys[j], reflectKeys[i] },
				func(i, j int) bool {
					if err != nil {
						return true
					}

					var less bool
					less, err = lessThan(reflectKeys[i].Interface(), reflectKeys[j].Interface())
					if err != nil {
						return true
					}
					return less
				},
				func() int { return len(reflectKeys) },
			))

			if err != nil {
				return
			}

			ret = make([]any, len(reflectKeys))
			for i := range reflectKeys {
				ret[i] = val.MapIndex(reflectKeys[i]).Interface()
			}

			return

		default:
			err = fmt.Errorf("unsupported operation on %T", v)
			return
		}
	}
}

func sortMapValuesByKey(keys []string, ret []any) {
	sort.Stable(sorthelper.NewCustomSortable(
		func(i, j int) {
			keys[i], keys[j] = keys[j], keys[i]
			ret[i], ret[j] = ret[j], ret[i]
		},
		func(i, j int) bool { return keys[i] < keys[j] },
		func() int { return len(keys) },
	))
}

func collectMapValues[K comparable, V any, M map[K]V](t map[K]V) (ret []any, err error) {
	i := 0
	sz := len(t)
	keys := make([]string, sz)
	ret = make([]any, sz)
	for k, v := range t {
		keys[i], err = toString(k)
		if err != nil {
			return
		}

		ret[i] = v
		i++
	}

	sortMapValuesByKey(keys, ret)
	return
}

// HasAny returns true if any target (all but the last arguments) exists in container (the last argument)
//
// when container is a map, target is used as a key
// when container is a slice, target is used as a value to compare each element
// when container is string, target is used as a substring
//
// nolint:gocyclo
func (collNS) HasAny(args ...any) (_ bool, err error) {
	n := len(args)
	if n < 2 {
		return false, fmt.Errorf("at least 2 args expected, got %d", n)
	}

	switch t := args[n-1].(type) {
	case string:
		return stringContainsAny(t, args[:n-1])
	case map[string]any:
		return mapContainsAny(t, args[:n-1])
	case map[string]string:
		return mapContainsAny(t, args[:n-1])
	case map[string]struct{}:
		return mapContainsAny(t, args[:n-1])
	case map[any]any:
		for _, target := range t {
			_, ok := t[target]
			if ok {
				return true, nil
			}
		}

		return false, nil
	case []string:
		return stringSliceContainsAny(t, args[:n-1])
	case []int:
		return integerSliceContainsAny(t, args[:n-1])
	case []uint:
		return integerSliceContainsAny(t, args[:n-1])

	case []int64:
		return integerSliceContainsAny(t, args[:n-1])
	case []uint64:
		return integerSliceContainsAny(t, args[:n-1])

	case []int8:
		return integerSliceContainsAny(t, args[:n-1])
	case []uint8:
		return integerSliceContainsAny(t, args[:n-1])

	case []int16:
		return integerSliceContainsAny(t, args[:n-1])
	case []uint16:
		return integerSliceContainsAny(t, args[:n-1])

	case []int32:
		return integerSliceContainsAny(t, args[:n-1])
	case []uint32:
		return integerSliceContainsAny(t, args[:n-1])

	case []uintptr:
		return integerSliceContainsAny(t, args[:n-1])

	case []float32:
		return floatSliceContainsAny(t, args[:n-1])
	case []float64:
		return floatSliceContainsAny(t, args[:n-1])

	case []any:
		for i := range args[:n-1] {
			for j := range t {
				if reflect.DeepEqual(t[j], args[i]) {
					return true, nil
				}
			}
		}

		return false, nil
	default:
		switch val := reflect.ValueOf(t); val.Kind() {
		case reflect.String:
			return stringContainsAny(val.String(), args[:n-1])
		case reflect.Map:
			for _, tgt := range args[:n-1] {
				if val.MapIndex(reflect.ValueOf(tgt)).IsValid() {
					return true, nil
				}
			}

			return false, nil
		case reflect.Slice, reflect.Array:
			sz := val.Len()
			for i := range args[:n-1] {
				for j := 0; j < sz; j++ {
					if reflect.DeepEqual(val.Index(j).Interface(), args[i]) {
						return true, nil
					}
				}
			}

			return false, nil
		default:
			return false, fmt.Errorf("unsupported container type %T", t)
		}
	}
}

func mapContainsAny[V any](m map[string]V, target []any) (_ bool, err error) {
	var (
		key string
		ok  bool
	)

	for _, tgt := range target {
		key, err = toString(tgt)
		_, ok = m[key]
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func stringContainsAny(t string, target []any) (_ bool, err error) {
	var k string
	for _, v := range target {
		k, err = toString(v)
		if err != nil {
			return false, err
		}

		if strings.Contains(t, k) {
			return true, nil
		}
	}

	return false, nil
}

func stringSliceContainsAny(s []string, target []any) (_ bool, err error) {
	var v string
	for i := range target {
		v, err = toString(target[i])
		if err != nil {
			return false, err
		}

		for j := range s {
			if s[j] == v {
				return true, nil
			}
		}
	}

	return false, nil
}

func integerSliceContainsAny[T integer](s []T, target []any) (_ bool, err error) {
	var v T
	for i := range target {
		v, err = parseInteger[T](target[i])
		if err != nil {
			return false, err
		}

		for j := range s {
			if s[j] == v {
				return true, nil
			}
		}
	}

	return false, nil
}

func floatSliceContainsAny[T float](s []T, target []any) (_ bool, err error) {
	var v T
	for i := range target {
		v, err = parseFloat[T](target[i])
		if err != nil {
			return false, err
		}

		for j := range s {
			if s[j] == v {
				return true, nil
			}
		}
	}

	return false, nil
}

// HasAll returns true if all target (all but the last arguments) exists in container (the last argument)
//
// when container is a map, target is used as a key
// when container is a slice, target is used as a value to compare each element
// when container is string, target is used as a substring
//
// nolint:gocyclo
func (collNS) HasAll(args ...any) (_ bool, err error) {
	n := len(args)
	if n < 2 {
		return false, fmt.Errorf("at least 2 args expected, got %d", n)
	}

	switch t := args[n-1].(type) {
	case string:
		return stringContainsAll(t, args[:n-1])
	case map[string]any:
		return mapContainsAll(t, args[:n-1])
	case map[string]string:
		return mapContainsAll(t, args[:n-1])
	case map[string]struct{}:
		return mapContainsAll(t, args[:n-1])
	case map[any]any:
		if len(t) == 0 {
			return false, nil
		}

		for _, target := range t {
			_, ok := t[target]
			if !ok {
				return false, nil
			}
		}

		return true, nil
	case []string:
		return stringSliceContainsAll(t, args[:n-1])
	case []int:
		return integerSliceContainsAll(t, args[:n-1])
	case []uint:
		return integerSliceContainsAll(t, args[:n-1])

	case []int64:
		return integerSliceContainsAll(t, args[:n-1])
	case []uint64:
		return integerSliceContainsAll(t, args[:n-1])

	case []int8:
		return integerSliceContainsAll(t, args[:n-1])
	case []uint8:
		return integerSliceContainsAll(t, args[:n-1])

	case []int16:
		return integerSliceContainsAll(t, args[:n-1])
	case []uint16:
		return integerSliceContainsAll(t, args[:n-1])

	case []int32:
		return integerSliceContainsAll(t, args[:n-1])
	case []uint32:
		return integerSliceContainsAll(t, args[:n-1])

	case []uintptr:
		return integerSliceContainsAll(t, args[:n-1])

	case []float32:
		return floatSliceContainsAll(t, args[:n-1])
	case []float64:
		return floatSliceContainsAll(t, args[:n-1])

	case []any:
		return anySliceContainsAll(t, args[:n-1]), nil

	default:
		switch val := reflect.ValueOf(t); val.Kind() {
		case reflect.String:
			return stringContainsAll(val.String(), args[:n-1])
		case reflect.Map:
			if val.Len() == 0 {
				return false, nil
			}

			for _, tgt := range args[:n-1] {
				if !val.MapIndex(reflect.ValueOf(tgt)).IsValid() {
					return false, nil
				}
			}

			return true, nil
		case reflect.Slice, reflect.Array:
			sz := val.Len()
			if sz == 0 {
				return false, nil
			}

			for i := range args[:n-1] {
				for j := 0; j < sz; j++ {
					if !reflect.DeepEqual(val.Index(j).Interface(), args[i]) {
						return false, nil
					}
				}
			}

			return true, nil
		default:
			return false, fmt.Errorf("unsupported container type %T", t)
		}
	}
}

func mapContainsAll[V any](m map[string]V, target []any) (_ bool, err error) {
	if len(m) == 0 {
		return len(target) == 0, nil
	}

	var (
		key string
		ok  bool
	)

	for _, tgt := range target {
		key, err = toString(tgt)
		_, ok = m[key]
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

func stringContainsAll(t string, target []any) (_ bool, err error) {
	if len(t) == 0 {
		return len(target) == 0, nil
	}

	var k string
	for _, v := range target {
		k, err = toString(v)
		if err != nil {
			return false, err
		}

		if !strings.Contains(t, k) {
			return false, nil
		}
	}

	return true, nil
}

func anySliceContainsAll(s []any, target []any) bool {
	if len(s) == 0 {
		return len(target) == 0
	}

	for i := range target {
		for j := range s {
			if !reflect.DeepEqual(s[j], target[i]) {
				return false
			}
		}
	}

	return true
}

func stringSliceContainsAll(s []string, target []any) (_ bool, err error) {
	if len(s) == 0 {
		return len(target) == 0, nil
	}

	var v string
	for i := range target {
		v, err = toString(target[i])
		if err != nil {
			return false, err
		}

		for j := range s {
			if s[j] != v {
				return false, nil
			}
		}
	}

	return true, nil
}

func integerSliceContainsAll[T integer](s []T, target []any) (_ bool, err error) {
	if len(s) == 0 {
		return len(target) == 0, nil
	}

	var v T
	for i := range target {
		v, err = parseInteger[T](target[i])
		if err != nil {
			return false, err
		}

		for j := range s {
			if s[j] != v {
				return false, nil
			}
		}
	}

	return true, nil
}

func floatSliceContainsAll[T float](s []T, target []any) (_ bool, err error) {
	if len(s) == 0 {
		return len(target) == 0, nil
	}

	var v T
	for i := range target {
		v, err = parseFloat[T](target[i])
		if err != nil {
			return false, err
		}

		for j := range s {
			if s[j] != v {
				return false, nil
			}
		}
	}

	return true, nil
}
