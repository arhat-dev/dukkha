package matrix

import "sort"

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
	sort.Slice(names, func(i, j int) bool {
		ok := names[i] < names[j]
		if ok {
			mat[i], mat[j] = mat[j], mat[i]
		}
		return ok
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
