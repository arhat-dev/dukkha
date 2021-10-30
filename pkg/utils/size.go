package utils

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Size is human-readable data size in KB, MB, GB, TB, TB
//
// Schema type is `string`
type Size uint64

var (
	_ yaml.Unmarshaler = (*Size)(nil)
)

func (s *Size) UnmarshalYAML(n *yaml.Node) error {
	val := strings.TrimSuffix(n.Value, "B")
	base := uint64(1)
	switch {
	case strings.HasSuffix(val, "P"):
		base *= 1024
		fallthrough
	case strings.HasSuffix(val, "T"):
		base *= 1024
		fallthrough
	case strings.HasSuffix(val, "G"):
		base *= 1024
		fallthrough
	case strings.HasSuffix(val, "M"):
		base *= 1024
		fallthrough
	case strings.HasSuffix(val, "K"):
		base *= 1024
		fallthrough
	default:
		val = strings.TrimRight(val, "PTGMK")
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid size value %q: %w", n.Value, err)
		}

		*s = Size(v * base)
	}

	return nil
}
