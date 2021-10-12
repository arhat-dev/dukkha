package textquery

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v3"
)

// JQ runs query over json data
func YQ(query, data string) (string, error) {
	return YQBytes(query, []byte(data))
}

// JQ runs query over yaml data bytes
func YQBytes(query string, dataBytes []byte) (string, error) {
	return Query(query, NewYAMLIterator(bytes.NewReader(dataBytes)), yaml.Marshal)
}

func NewYAMLIterator(r io.Reader) func() (interface{}, bool) {
	dec := yaml.NewDecoder(r)

	return func() (interface{}, bool) {
		var data interface{}
		err := dec.Decode(&data)
		if err != nil {
			// return plain text on unexpected error
			return nil, false
		}

		return data, true
	}
}
