package utils

import "gopkg.in/yaml.v3"

type Size int64

var _ yaml.Unmarshaler = (*Size)(nil)

func (s *Size) UnmarshalYAML(node *yaml.Node) error {
	// TODO: implement
	// KB, MB, GB
	return nil
}
