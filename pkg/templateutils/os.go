package templateutils

import (
	"fmt"
	"os"
)

var osNS = &_osNS{}

type _osNS struct{}

func (ns *_osNS) ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (ns *_osNS) WriteFile(filename string, d interface{}) error {
	var data []byte
	switch dt := d.(type) {
	case string:
		data = []byte(dt)
	case []byte:
		data = dt
	default:
		return fmt.Errorf("invalid non string nor bytes data: %T", d)
	}

	return os.WriteFile(filename, data, 0640)
}
