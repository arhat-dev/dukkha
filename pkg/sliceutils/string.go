package sliceutils

func NewStringSlice(base []string, other ...string) []string {
	return append(append([]string(nil), base...), other...)
}
