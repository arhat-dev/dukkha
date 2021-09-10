package sorthelper

import "sort"

// NewCustomSortable creates a sortable from provided functions
func NewCustomSortable(
	swap func(i, j int),
	less func(i, j int) bool,
	len func() int,
) sort.Interface {
	return &customSortable{
		doSwap: swap,
		isLess: less,
		getLen: len,
	}
}

var _ sort.Interface = (*customSortable)(nil)

type customSortable struct {
	doSwap func(i, j int)
	isLess func(i, j int) bool
	getLen func() int
}

func (s *customSortable) Len() int           { return s.getLen() }
func (s *customSortable) Less(i, j int) bool { return s.isLess(i, j) }
func (s *customSortable) Swap(i, j int)      { s.doSwap(i, j) }
