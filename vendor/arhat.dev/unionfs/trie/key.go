package trie

type KeyFunc func(key string, pos int) (string, int)

func RuneKeyFunc(key string, pos int) (k string, next int) {
	data := []rune(key)
	size := len(data)

	if pos >= size {
		return "", -1
	}

	return string(data[pos]), pos + 1
}
