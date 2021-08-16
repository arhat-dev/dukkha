package textquery

import (
	"gopkg.in/yaml.v3"
)

// JQ runs query over json data
func YQ(query, data string) (string, error) {
	return YQBytes(query, []byte(data))
}

// JQ runs query over yaml data bytes
func YQBytes(query string, dataBytes []byte) (string, error) {
	return Query(query, dataBytes, yaml.Unmarshal, yaml.Marshal)
}
