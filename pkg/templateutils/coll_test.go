package templateutils

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func test_1AnyIn_Error(
	t assert.TestingT,
	do func(any) (any, error),
	in any,
) {
	_, err := do(in)
	assert.Error(t, err)
}

func test_1AnyIn_Ok[T any](
	t assert.TestingT,
	do func(any) (any, error),
	expected T,
	in any,
) {
	ret, err := do(in)
	assert.NoError(t, err)
	assert.IsType(t, expected, ret)
	assert.EqualValues(t, expected, ret)
}

func test_2AnyIn_Error(
	t assert.TestingT,
	do func(a, b any) (any, error),
	arg0, arg1 any,
) {
	_, err := do(arg0, arg1)
	assert.Error(t, err)
}

func test_2AnyIn_Ok[T any](
	t assert.TestingT,
	do func(a, b any) (any, error),
	expected T,
	arg0, arg1 any,
) {
	ret, err := do(arg0, arg1)
	assert.NoError(t, err)
	assert.IsType(t, expected, ret)
	assert.EqualValues(t, expected, ret)
}

func test_VarAnyIn_Ok[T any](
	t assert.TestingT,
	do func(...any) (any, error),
	expected T,
	in ...any,
) {
	ret, err := do(in...)
	assert.NoError(t, err)
	assert.IsType(t, expected, ret)
	assert.EqualValues(t, expected, ret)
}

func fromAnySlice[T any](s []any) (ret []T) {
	ret = make([]T, len(s))
	for i := range s {
		ret[i] = *(*T)(unsafe.Pointer(&s[i]))
	}
	return
}

type TypedInt int
type TypedString string
type TypedIntSlice []TypedInt

func TestValidSliceArgs(t *testing.T) {
	for _, test := range []struct {
		// args
		low, high, max, sz, kap int
		// ret
		l, h, m int
	}{
		{}, // all zero
		{1, 1, 1, 1, 1, 1, 1, 1},
	} {
		l, h, m := validSliceArgs(test.low, test.high, test.max, test.sz, test.kap)
		assert.Equal(t, test.l, l)
		assert.Equal(t, test.h, h)
		assert.Equal(t, test.m, m)
	}
}

func TestCollNS_Slice(t *testing.T) {
	// TODO
}

func TestCollNS_Reverse(t *testing.T) {
	reverseFn := collNS{}.Reverse

	test_1AnyIn_Error(t, reverseFn, map[string]string{"1": "1"})

	test_1AnyIn_Ok(t, reverseFn, "12345", "54321")
	test_1AnyIn_Ok(t, reverseFn, TypedString("一二三四"), TypedString("四三二一"))

	test_1AnyIn_Ok(t, reverseFn, []string{"1"}, []string{"1"})
	test_1AnyIn_Ok(t, reverseFn, []any{3.3, 2, "1"}, []any{"1", 2, 3.3})

	test_1AnyIn_Ok(t, reverseFn, []TypedInt{4, 3, 2, 1}, []TypedInt{1, 2, 3, 4})

	test_1AnyIn_Ok(t, reverseFn, TypedIntSlice{}, TypedIntSlice{})
	test_1AnyIn_Ok(t, reverseFn, []uint{4, 3, 2, 1}, []uint{1, 2, 3, 4})

	test_1AnyIn_Ok(t, reverseFn, []int8{4, 3, 2, 1}, []int8{1, 2, 3, 4})
	test_1AnyIn_Ok(t, reverseFn, []uint8{4, 3, 2, 1}, []uint8{1, 2, 3, 4})

	test_1AnyIn_Ok(t, reverseFn, []int16{4, 3, 2, 1}, []int16{1, 2, 3, 4})
	test_1AnyIn_Ok(t, reverseFn, []uint16{4, 3, 2, 1}, []uint16{1, 2, 3, 4})

	test_1AnyIn_Ok(t, reverseFn, []int32{4, 3, 2, 1}, []int32{1, 2, 3, 4})
	test_1AnyIn_Ok(t, reverseFn, []uint32{4, 3, 2, 1}, []uint32{1, 2, 3, 4})

	test_1AnyIn_Ok(t, reverseFn, []int64{4, 3, 2, 1}, []int64{1, 2, 3, 4})
	test_1AnyIn_Ok(t, reverseFn, []uint64{4, 3, 2, 1}, []uint64{1, 2, 3, 4})

	test_1AnyIn_Ok(t, reverseFn, []uintptr{4, 3, 2, 1}, []uintptr{1, 2, 3, 4})

	test_1AnyIn_Ok(t, reverseFn, []float32{3, 2, 1}, []float32{1, 2, 3})
	test_1AnyIn_Ok(t, reverseFn, []float64{3, 2, 1}, []float64{1, 2, 3})
}

func TestCollNS_Unique(t *testing.T) {
	uniqueFn := func(in any) (any, error) {
		return collNS{}.Unique(in)
	}

	test_1AnyIn_Error(t, uniqueFn, "12345")
	test_1AnyIn_Error(t, uniqueFn, TypedString("一二三四"))

	test_1AnyIn_Ok(t, uniqueFn, map[string]string{"1": "1"}, map[string]string{"1": "1"})
	test_1AnyIn_Ok(t, uniqueFn, map[string]any{"1": "1"}, map[string]any{"1": "1"})
	test_1AnyIn_Ok(t, uniqueFn, map[string]struct{}{"1": {}}, map[string]struct{}{"1": {}})
	test_1AnyIn_Ok(t, uniqueFn, map[any]any{"1": "1"}, map[any]any{"1": "1"})
	test_1AnyIn_Ok(t, uniqueFn, map[int]string{1: "1"}, map[int]string{1: "1"})

	test_1AnyIn_Ok(t, uniqueFn, []string{"1"}, []string{"1"})
	test_1AnyIn_Ok(t, uniqueFn, []any{"1", 2}, []any{"1", 2, 2})

	test_1AnyIn_Ok(t, uniqueFn, []TypedInt{1, 3, 2}, []TypedInt{1, 3, 1, 2, 2})

	test_1AnyIn_Ok[TypedIntSlice](t, uniqueFn, nil, TypedIntSlice{})
	test_1AnyIn_Ok(t, uniqueFn, []uint{1}, []uint{1, 1})

	test_1AnyIn_Ok[[]int8](t, uniqueFn, nil, []int8{})
	test_1AnyIn_Ok(t, uniqueFn, []uint8{1, 2, 3, 4}, []uint8{1, 2, 3, 4})

	test_1AnyIn_Ok(t, uniqueFn, []int16{1, 2, 3}, []int16{1, 2, 3, 3})
	test_1AnyIn_Ok(t, uniqueFn, []uint16{1, 2, 3}, []uint16{1, 2, 3, 3})

	test_1AnyIn_Ok(t, uniqueFn, []int32{1, 2, 3}, []int32{1, 2, 3, 3})
	test_1AnyIn_Ok(t, uniqueFn, []uint32{1, 2, 3}, []uint32{1, 2, 3, 3})

	test_1AnyIn_Ok(t, uniqueFn, []int64{1, 2, 3}, []int64{1, 2, 3, 3})
	test_1AnyIn_Ok(t, uniqueFn, []uint64{1, 2, 3}, []uint64{1, 2, 3, 3})

	test_1AnyIn_Ok(t, uniqueFn, []uintptr{1, 2, 3}, []uintptr{1, 2, 3, 3})

	test_1AnyIn_Ok(t, uniqueFn, []float32{1, 2, 3.1}, []float32{1, 2, 3.1, 3.1})
	test_1AnyIn_Ok(t, uniqueFn, []float64{1, 2, 3.1}, []float64{1, 2, 3.1, 3.1})
}

func TestCollNS_Merge(t *testing.T) {
	mergeFn := func(args ...any) (any, error) {
		return collNS{}.Merge(fromAnySlice[Map](args)...)
	}

	test_VarAnyIn_Ok(t, mergeFn,
		map[string]string{"a": "1", "b": "", "c": "3"},
		// in
		map[string]string{"a": "1"}, map[string]struct{}{"b": {}}, map[string]string{"c": "3"},
	)

	test_VarAnyIn_Ok(t, mergeFn,
		map[string]any{"a": "1", "b": struct{}{}, "c": map[string]any{"d": []any{4, 1, 2, 3}}},
		// in
		map[string]string{"a": "1"},
		map[string]struct{}{"b": {}},
		map[string]any{
			"c": map[string]any{
				"d": []any{1, 2, 3},
			},
		},
		map[string]any{
			"c": map[string]any{
				"d": []any{4},
			},
		},
	)

	test_VarAnyIn_Ok(t, mergeFn,
		map[string]struct{}{"a": {}, "b": {}, "c": {}},
		// in
		map[string]string{"a": "1"}, map[string]any{"b": 2}, map[string]struct{}{"c": {}},
	)

	test_VarAnyIn_Ok(t, mergeFn,
		map[any]any{"a": "1", "b": struct{}{}, "c": 3},
		// in
		map[string]string{"a": "1"}, map[string]struct{}{"b": {}}, map[any]any{"c": 3},
	)
}

func TestCollNS_Index(t *testing.T) {
	indexFn := collNS{}.Index

	test_2AnyIn_Error(t, indexFn, -1, []int{})
	test_2AnyIn_Error(t, indexFn, "a", []int{})

	test_2AnyIn_Ok[any](t, indexFn, nil, "x", map[string]any{})

	test_2AnyIn_Ok[int](t, indexFn, 2, "b", map[string]any{"a": 1, "b": 2, "c": 3})
	test_2AnyIn_Ok[string](t, indexFn, "b", 2, map[string]string{"1": "a", "2": "b", "3": "c"})
	// reflect
	test_2AnyIn_Ok[int](t, indexFn, 2, 2, map[int]int{1: 1, 2: 2, 3: 3})

	test_2AnyIn_Ok[string](t, indexFn, "2", 1, []string{"1", "2", "3"})

	// reflect
	test_2AnyIn_Ok[TypedInt](t, indexFn, 2, 1, TypedIntSlice{1, 2, 3})

	test_2AnyIn_Ok[int](t, indexFn, 3, -1, []int{1, 2, 3})
	test_2AnyIn_Ok[uint](t, indexFn, 2, -2, []uint{1, 2, 3})

	test_2AnyIn_Ok[int8](t, indexFn, 1, -3, []int8{1, 2, 3})
	test_2AnyIn_Ok[uint8](t, indexFn, 2, 1, []uint8{1, 2, 3})

	test_2AnyIn_Ok[int16](t, indexFn, 2, 1, []int16{1, 2, 3})
	test_2AnyIn_Ok[uint16](t, indexFn, 2, 1, []uint16{1, 2, 3})

	test_2AnyIn_Ok[int32](t, indexFn, 2, 1, []int32{1, 2, 3})
	test_2AnyIn_Ok[uint32](t, indexFn, 2, 1, []uint32{1, 2, 3})

	test_2AnyIn_Ok[int64](t, indexFn, 2, 1, []int64{1, 2, 3})
	test_2AnyIn_Ok[uint64](t, indexFn, 2, 1, []uint64{1, 2, 3})

	test_2AnyIn_Ok[uintptr](t, indexFn, 2, 1, []uintptr{1, 2, 3})

	test_2AnyIn_Ok[float32](t, indexFn, 2.2, 1, []float32{1, 2.2, 3})
	test_2AnyIn_Ok[float64](t, indexFn, 2.2, 1, []float64{1, 2.2, 3})
}

func TestCollNS_Clone(t *testing.T) {
	// TODO
}

func TestCollNS_Sort(t *testing.T) {
	sortFn := func(in any) (any, error) {
		return collNS{}.Sort(in)
	}

	test_1AnyIn_Error(t, sortFn, map[string]any{"1": 1})

	test_1AnyIn_Ok(t, sortFn, TypedIntSlice{1, 2, 3}, TypedIntSlice{2, 1, 3})

	test_1AnyIn_Ok(t, sortFn, []string{"1", "2", "3"}, []string{"2", "1", "3"})

	test_1AnyIn_Ok(t, sortFn, []int{1, 2, 3}, []int{2, 1, 3})
	test_1AnyIn_Ok(t, sortFn, []uint{1, 2, 3}, []uint{2, 1, 3})

	test_1AnyIn_Ok(t, sortFn, []int8{1, 2, 3}, []int8{2, 1, 3})
	test_1AnyIn_Ok(t, sortFn, []uint8{1, 2, 3}, []uint8{2, 1, 3})

	test_1AnyIn_Ok(t, sortFn, []int16{1, 2, 3}, []int16{2, 1, 3})
	test_1AnyIn_Ok(t, sortFn, []uint16{1, 2, 3}, []uint16{2, 1, 3})

	test_1AnyIn_Ok(t, sortFn, []int32{1, 2, 3}, []int32{2, 1, 3})
	test_1AnyIn_Ok(t, sortFn, []uint32{1, 2, 3}, []uint32{2, 1, 3})

	test_1AnyIn_Ok(t, sortFn, []int64{1, 2, 3}, []int64{2, 1, 3})
	test_1AnyIn_Ok(t, sortFn, []uint64{1, 2, 3}, []uint64{2, 1, 3})

	test_1AnyIn_Ok(t, sortFn, []uintptr{1, 2, 3}, []uintptr{2, 1, 3})

	test_1AnyIn_Ok(t, sortFn, []float32{1.1, 1.2, 1.3}, []float32{1.2, 1.1, 1.3})
	test_1AnyIn_Ok(t, sortFn, []float64{1.1, 1.2, 1.3}, []float64{1.2, 1.1, 1.3})
}

func TestCollNS_Flatten(t *testing.T) {
	// TODO
}

func TestCollNS_Pick(t *testing.T) {
	// TODO
}

func TestCollNS_Omit(t *testing.T) {
	// TODO
}

func TestCollNS_Append(t *testing.T) {
	// TODO
}

func TestCollNS_Prepend(t *testing.T) {
	// TODO
}

func TestCollNS_MapStringAny(t *testing.T) {
	// TODO
}

func TestCollNS_MapAnyAny(t *testing.T) {
	// TODO
}

func TestCollNS_Keys(t *testing.T) {
	// TODO
}

func TestCollNS_Values(t *testing.T) {
	// TODO
}

func TestCollNS_HasAny(t *testing.T) {
	// TODO
}

func TestCollNS_HasAll(t *testing.T) {
	// TODO
}
