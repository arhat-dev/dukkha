package utils

import (
	"io"

	"gopkg.in/yaml.v3"
)

func UnmarshalStrict(r io.Reader, out interface{}) error {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	return dec.Decode(out)
}
