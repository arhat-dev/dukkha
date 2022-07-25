package matrixhelper

func CartesianProduct[K comparable, V any](
	m map[K][]V, sort func(names []K, mat [][]V),
) (ret []map[K]V) {
	names := make([]K, 0, len(m))
	mat := make([][]V, 0, len(m))
	for k, v := range m {
		if len(v) == 0 {
			// ignore empty list
			continue
		}

		names = append(names, k)
		mat = append(mat, v)
	}

	// sort names and mat
	sort(names, mat)

	listCart := cartNext(mat)
	if len(listCart) == 0 {
		return nil
	}

	ret = make([]map[K]V, len(listCart))
	for i, list := range listCart {
		vMap := make(map[K]V)
		for j, v := range list {
			vMap[names[j]] = v
		}

		ret[i] = vMap
	}

	return ret
}

func cartNext[T any](mat [][]T) (ret [][]T) {
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

	ret = make([][]T, tupleCount)

	buf := make([]T, tupleCount*len(mat))
	indexPerList := make([]int, len(mat))

	start := 0
	for i := range ret {
		end := start + len(mat)

		tuple := buf[start:end]
		ret[i] = tuple

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

	return
}
