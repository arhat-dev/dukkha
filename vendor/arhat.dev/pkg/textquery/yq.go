package textquery

import (
	"gopkg.in/yaml.v3"
)

func YQ(query, data string) (string, error) {
	return YQBytes(query, []byte(data))
}

func YQBytes(query string, input []byte) (string, error) {
	return Query(query, input, yaml.Unmarshal, yaml.Marshal)
}
