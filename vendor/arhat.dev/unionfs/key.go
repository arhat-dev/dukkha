package unionfs

import (
	"strings"
)

func PathKeyFunc(key string, pos int) (string, int) {
	size := len(key)
	if pos >= size {
		return "", -1
	}

	switch {
	case pos == 0:
		next := strings.IndexByte(key, '/')
		switch next {
		case 0:
			// `/`
			return "/", 1
		case -1:
			return key, size
		default:
			return key[:next], next
		}
	case pos == 1 && key[0] == '/':
		next := strings.IndexByte(key[1:], '/')
		if next == -1 {
			return key, size
		}

		return key[:pos+next], pos + next
	default:
	}

	next := strings.IndexByte(key[pos+1:], '/')
	switch {
	case next == -1:
		return key[pos:], size
	default:
		return key[pos : pos+next+1], pos + next + 1
	}
}
