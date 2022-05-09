package stringhelper

func HasPrefix[S ~string](s, prefix S) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func HasSuffix[S ~string](s, suffix S) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func TrimPrefix[S ~string](s, prefix S) S {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func TrimSuffix[S ~string](s, suffix S) S {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}
