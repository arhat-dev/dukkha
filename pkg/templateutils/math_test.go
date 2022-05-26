package templateutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMath_Seq(t *testing.T) {
	for _, test := range []struct {
		args     []Number
		expected []int64
	}{
		{
			args:     []Number{0},
			expected: []int64{},
		},
		{
			args:     []Number{1},
			expected: []int64{0},
		},
		{
			args:     []Number{2},
			expected: []int64{0, 1},
		},
		{
			args:     []Number{1, 4},
			expected: []int64{1, 2, 3},
		},
		{
			args:     []Number{1, 6, 2},
			expected: []int64{1, 3, 5},
		},
		{
			args:     []Number{6, 1, 2},
			expected: nil,
		},
		{
			args:     []Number{6, 1, -2},
			expected: []int64{6, 4, 2},
		},
		{
			args:     []Number{-2},
			expected: []int64{0, -1},
		},
		{
			args:     []Number{-1, -6, -2},
			expected: []int64{-1, -3, -5},
		},
	} {
		ret, err := mathNS{}.Seq(test.args...)
		assert.NoError(t, err)
		assert.EqualValues(t, test.expected, ret)
	}
}

func TestMathNS_Add(t *testing.T) {
	for _, test := range []struct {
		args     []Number
		expected Number
	}{
		{
			args:     []Number{1, 2, 3},
			expected: int64(6),
		},
		{
			args:     []Number{-2, 2, uint(3)},
			expected: uint64(3),
		},
		{
			args:     []Number{1, 2, 3.0},
			expected: float64(6),
		},
		{
			args:     []Number{1, "2.1", "3"},
			expected: "6.1",
		},
	} {
		ret, err := mathNS{}.Add(test.args...)
		assert.NoError(t, err)
		assert.IsType(t, test.expected, ret)
		assert.EqualValues(t, test.expected, ret)
	}
}

func TestMathNS_Sub(t *testing.T) {
	for _, test := range []struct {
		args     []Number
		expected Number
	}{
		{
			args:     []Number{1, 2, 3},
			expected: int64(0),
		},
		{
			args:     []Number{-2, 2, uint(3)},
			expected: uint64(3),
		},
		{
			args:     []Number{1, 2, 3.0},
			expected: float64(0),
		},
		{
			args:     []Number{1, "2", "3"},
			expected: "0",
		},
	} {
		ret, err := mathNS{}.Sub(test.args...)
		assert.NoError(t, err)
		assert.IsType(t, test.expected, ret)
		assert.EqualValues(t, test.expected, ret)
	}
}

func TestMathNS_Mul(t *testing.T) {
	for _, test := range []struct {
		args     []Number
		expected Number
	}{
		{
			args:     []Number{1, 2, 3},
			expected: int64(6),
		},
		{
			args:     []Number{2, 2, uint(3)},
			expected: uint64(12),
		},
		{
			args:     []Number{1, 2, 3.0},
			expected: float64(6),
		},
		{
			args:     []Number{1, "2", "3"},
			expected: "6",
		},
	} {
		ret, err := mathNS{}.Mul(test.args...)
		assert.NoError(t, err)
		assert.IsType(t, test.expected, ret)
		assert.EqualValues(t, test.expected, ret)
	}
}

func TestMathNS_Div(t *testing.T) {
	for _, test := range []struct {
		args     []Number
		expected Number
	}{
		{
			args:     []Number{2, 2, 99},
			expected: int64(24),
		},
		{
			args:     []Number{2, 2, uint(20)},
			expected: uint64(5),
		},
		{
			args:     []Number{1, 2, 3.0},
			expected: float64(1.5),
		},
		{
			args:     []Number{1, "2", "3"},
			expected: "1.5",
		},
	} {
		ret, err := mathNS{}.Div(test.args...)
		assert.NoError(t, err)
		assert.IsType(t, test.expected, ret)
		assert.EqualValues(t, test.expected, ret)
	}
}

func TestMathNS_Mod(t *testing.T) {
	for _, test := range []struct {
		args     []Number
		expected Number
	}{
		{
			args:     []Number{2, 4, 99},
			expected: int64(1),
		},
		{
			args:     []Number{2, 2, uint(20)},
			expected: uint64(0),
		},
		{
			args:     []Number{1, 2, 3.0},
			expected: float64(0),
		},
		// TODO: support big.Float mod
		// {
		// 	args:     []Number{1, "2", "3"},
		// 	expected: "1.5",
		// },
	} {
		ret, err := mathNS{}.Mod(test.args...)
		assert.NoError(t, err)
		assert.IsType(t, test.expected, ret)
		assert.EqualValues(t, test.expected, ret)
	}
}
