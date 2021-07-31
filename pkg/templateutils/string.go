package templateutils

import "strings"

// AddPrefix to each seperated string elements
func AddPrefix(s, prefix, sep string) string {
	parts := strings.Split(s, sep)
	for i, p := range parts {
		if len(p) == 0 && i == len(parts)-1 {
			continue
		}

		parts[i] = prefix + p
	}

	return strings.Join(parts, sep)
}

// RemovePrefix of each seperated string elements
func RemovePrefix(s, prefix, sep string) string {
	parts := strings.Split(s, sep)
	for i, p := range parts {
		parts[i] = strings.TrimPrefix(p, prefix)
	}

	return strings.Join(parts, sep)
}

// AddSuffix to each seperated string elements
func AddSuffix(s, suffix, sep string) string {
	parts := strings.Split(s, sep)
	for i, p := range parts {
		if len(p) == 0 && i == len(parts)-1 {
			continue
		}

		parts[i] = p + suffix
	}

	return strings.Join(parts, sep)
}

func RemoveSuffix(s, suffix, sep string) string {
	parts := strings.Split(s, sep)
	for i, p := range parts {
		parts[i] = strings.TrimSuffix(p, suffix)
	}

	return strings.Join(parts, sep)
}
