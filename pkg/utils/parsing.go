package utils

import "fmt"

// ParseBrackets `()`
func ParseBrackets(s string) (string, error) {
	leftBrackets := 0
	for i := range s {
		switch s[i] {
		case '(':
			leftBrackets++
		case ')':
			if leftBrackets == 0 {
				return s[:i], nil
			}
			leftBrackets--
		}
	}

	// invalid data
	return "", fmt.Errorf("unexpected non-terminated brackets")
}
