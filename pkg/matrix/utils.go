package matrix

import "sort"

var _ sort.Interface = (*swapper)(nil)

type swapper struct {
	doSwap func(i, j int)
	isLess func(i, j int) bool
	getLen func() int
}

func (s *swapper) Len() int           { return s.getLen() }
func (s *swapper) Less(i, j int) bool { return s.isLess(i, j) }
func (s *swapper) Swap(i, j int)      { s.doSwap(i, j) }

func CartesianProduct(m map[string][]string) []map[string]string {
	names := make([]string, 0)
	mat := make([][]string, 0)
	for k, v := range m {
		if len(v) == 0 {
			// ignore empty list
			continue
		}

		names = append(names, k)
		mat = append(mat, v)
	}

	// sort names and mat
	sort.Sort(&swapper{
		getLen: func() int { return len(names) },
		isLess: func(i, j int) bool { return names[i] < names[j] },
		doSwap: func(i, j int) {
			names[i], names[j] = names[j], names[i]
			mat[i], mat[j] = mat[j], mat[i]
		},
	})

	listCart := cartNext(mat)
	if len(listCart) == 0 {
		return nil
	}

	result := make([]map[string]string, 0, len(listCart))
	for _, list := range listCart {
		vMap := make(map[string]string)
		for i, v := range list {
			vMap[names[i]] = v
		}
		result = append(result, vMap)
	}

	return result
}

func cartNext(mat [][]string) [][]string {
	if len(mat) == 0 {
		return nil
	}

	tupleCount := 1
	for _, list := range mat {
		if len(list) == 0 {
			// ignore empty list
			continue
		}

		tupleCount *= len(list)
	}

	result := make([][]string, tupleCount)

	buf := make([]string, tupleCount*len(mat))
	indexPerList := make([]int, len(mat))

	start := 0
	for i := range result {
		end := start + len(mat)

		tuple := buf[start:end]
		result[i] = tuple

		start = end

		for j, idx := range indexPerList {
			// mat[j] is the list

			tuple[j] = mat[j][idx]
		}

		for j := len(indexPerList) - 1; j >= 0; j-- {
			indexPerList[j]++
			if indexPerList[j] < len(mat[j]) {
				break
			}

			// reset for next tuple
			indexPerList[j] = 0
		}
	}

	return result
}
